package main

// To use in docker with OIDC, add import
//
//	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/drain"
	"nhooyr.io/websocket"
)

type appConfig struct {
	listenPort   string
	resyncPeriod time.Duration
	kubeConfig   *rest.Config
}

func main() {
	appConfig := getConfig()

	// Init factory for informers for well-known types
	clientset, err := kubernetes.NewForConfig(appConfig.kubeConfig)
	if err != nil {
		klog.Fatal(fmt.Errorf("creating clientset: %v", err.Error()))
	}
	factory := informers.NewSharedInformerFactory(clientset, appConfig.resyncPeriod)
	defer factory.Shutdown()

	// Init factory for informers for custom types
	dynClient, err := dynamic.NewForConfig(appConfig.kubeConfig)
	if err != nil {
		klog.Fatal(fmt.Errorf("creating dynamic client: %v", err.Error()))
	}
	dynFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynClient, appConfig.resyncPeriod)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := httprouter.New()
	handler, err := initHandlers(ctx, router, clientset, factory, dynClient, dynFactory)
	if err != nil {
		klog.Fatal(fmt.Errorf("initializing handlers: %v", err.Error()))
	}

	router.GET("/healthz", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) { w.WriteHeader(200) })

	var inSync atomic.Bool
	router.GET("/readyz", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if inSync.Load() {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(500)
	})

	errc := make(chan error, 1)
	go func() {
		// Start informers all at once after we have inited them in initHandlers func
		factory.Start(ctx.Done()) // Start processing these informers.
		klog.Info("Started informers.")
		// Wait for cache sync
		klog.Info("Waiting for initial sync of informers.")
		synced := factory.WaitForCacheSync(ctx.Done())
		for v, ok := range synced {
			if !ok {
				errc <- fmt.Errorf("caches failed to sync: %v", v)
			}
		}

		// Start dynamic informers all at once after we have inited them in initHandlers func
		dynFactory.Start(ctx.Done())
		klog.Info("Started dynamic informers.")
		// Wait for cache sync for dynamic informers
		klog.Info("Waiting for initial sync of dynamic informers.")
		dynSynced := dynFactory.WaitForCacheSync(ctx.Done())
		for v, ok := range dynSynced {
			if !ok {
				errc <- fmt.Errorf("dynamic caches failed to sync: %v", v)
			}
		}

		inSync.Store(true)
	}()

	klog.Info("Listening :" + appConfig.listenPort)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + appConfig.listenPort,
	}

	go func() {
		errc <- srv.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		klog.Errorf("failed: %v", err)
	case sig := <-sigs:
		klog.Infof("terminating: %v", sig)
	}

	shutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		klog.Errorf("shutdown: %v", err)
	}
}

type resourceDefinition struct {
	gvr   schema.GroupVersionResource
	ns    bool         // is it namespaced?
	subh  []subHandler // subhandlers for a named object
	check gvrCheck     // should we register this API?
}

type gvrCheck func(context.Context, schema.GroupVersionResource) (bool, error)

type subHandler struct {
	method  string
	suffix  string
	handler func(*kubernetes.Clientset, informers.GenericInformer) httprouter.Handle
}

func checkCustomResourceExistence(dynClient *dynamic.DynamicClient, timeout time.Duration) func(context.Context, schema.GroupVersionResource) (bool, error) {
	return func(ctx context.Context, gvr schema.GroupVersionResource) (bool, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		_, err := dynClient.Resource(gvr).List(ctx, metav1.ListOptions{})
		if err != nil {
			if apierrors.IsForbidden(err) || apierrors.IsNotFound(err) {
				// 403 is expected if the CRD is not present locally, 404 is expected when run in a Pod
				klog.V(5).Infof("CRD %s is not available: %v", gvr.String(), err)
				return false, nil
			}
			return false, err
		}
		return true, nil
	}
}

func initHandlers(
	ctx context.Context,
	router *httprouter.Router,
	clientset *kubernetes.Clientset,
	factory informers.SharedInformerFactory,
	dynClient *dynamic.DynamicClient,
	dynFactory dynamicinformer.DynamicSharedInformerFactory,
) (http.HandlerFunc, error) {
	reh := newResourceEventHandler()
	checkTimeout := 10 * time.Second

	definitions := []resourceDefinition{
		{
			gvr: schema.GroupVersionResource{Group: "", Resource: "nodes", Version: "v1"},
			subh: []subHandler{{
				method:  http.MethodPost,
				suffix:  "drain",
				handler: handleNodeDrain,
			}},
		},

		{
			gvr: schema.GroupVersionResource{Group: "apps", Resource: "deployments", Version: "v1"},
			ns:  true,
		},

		{gvr: schema.GroupVersionResource{Group: "deckhouse.io", Resource: "nodegroups", Version: "v1"}},
		{gvr: schema.GroupVersionResource{Group: "deckhouse.io", Resource: "deckhousereleases", Version: "v1alpha1"}},
		{gvr: schema.GroupVersionResource{Group: "deckhouse.io", Resource: "moduleconfigs", Version: "v1alpha1"}},

		{
			gvr:   schema.GroupVersionResource{Group: "deckhouse.io", Resource: "awsinstanceclasses", Version: "v1"},
			check: checkCustomResourceExistence(dynClient, checkTimeout),
		},
		{
			gvr:   schema.GroupVersionResource{Group: "deckhouse.io", Resource: "azureinstanceclasses", Version: "v1"},
			check: checkCustomResourceExistence(dynClient, checkTimeout),
		},
		{
			gvr:   schema.GroupVersionResource{Group: "deckhouse.io", Resource: "gcpinstanceclasses", Version: "v1"},
			check: checkCustomResourceExistence(dynClient, checkTimeout),
		},
		{
			gvr:   schema.GroupVersionResource{Group: "deckhouse.io", Resource: "openstackinstanceclasses", Version: "v1"},
			check: checkCustomResourceExistence(dynClient, checkTimeout),
		},
		{
			gvr:   schema.GroupVersionResource{Group: "deckhouse.io", Resource: "vsphereinstanceclasses", Version: "v1"},
			check: checkCustomResourceExistence(dynClient, checkTimeout),
		},
		{
			gvr:   schema.GroupVersionResource{Group: "deckhouse.io", Resource: "yandexinstanceclasses", Version: "v1"},
			check: checkCustomResourceExistence(dynClient, checkTimeout),
		},
	}

	discovery := newDiscoveryCollector(clientset)
	for _, def := range definitions {
		gvr, namespaced := def.gvr, def.ns

		if def.check != nil {
			ok, err := def.check(ctx, gvr)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}

		collectionPath := getPathPrefix(gvr, namespaced, "k8s")
		namedItemPath := collectionPath + "/:name"

		informer := dynFactory.ForResource(gvr)
		h := newHandler(informer, dynClient.Resource(gvr), gvr, namespaced)
		_, _ = informer.Informer().AddEventHandler(reh.Handle(gvr))

		router.GET(collectionPath, h.HandleList)
		router.GET(namedItemPath, h.HandleGet)
		router.POST(collectionPath, h.HandleCreate)
		router.PUT(namedItemPath, h.HandleUpdate)
		router.DELETE(namedItemPath, h.HandleDelete)

		discovery.AddPath(collectionPath)
		discovery.AddPath(namedItemPath)

		for _, s := range def.subh {
			path := namedItemPath + "/" + s.suffix
			router.Handle(s.method, path, s.handler(clientset, informer))
			discovery.AddPath(path)
		}

		if strings.HasSuffix(gvr.Resource, "instanceclasses") {
			cloudProviderName := strings.TrimSuffix(gvr.Resource, "instanceclasses")
			discoveryCtx, discoveryCtxCancel := context.WithTimeout(ctx, checkTimeout)
			defer discoveryCtxCancel()
			if err := discovery.AddCloudProvider(discoveryCtx, cloudProviderName); err != nil {
				return nil, err
			}
		}
	}

	// Websocket
	sc := newSubscriptionController(reh)
	go sc.Start(ctx)
	router.GET("/subscribe", handleSubscribe(sc))

	// Discovery
	router.GET("/discovery", handleDiscovery(clientset, discovery.Build()))

	var wrapper http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}

		klog.V(5).Infof("Request: %s %s", r.Method, r.URL.Path)
		router.ServeHTTP(w, r)
		// TODO: Use echo/v4. To log response status, we need to wrap the response writer or
		// use non-standard library. It still will help to handle path parameters.
	}

	return wrapper, nil
}

func handleSubscribe(sc *subscriptionController) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
			// Declaring supported protocol for frontend tooling based on ActionCable;
			// "actioncable-unsupported" is omitted because it seem to be unneeded.
			Subprotocols: []string{"actioncable-v1-json"},
		})
		if err != nil {
			klog.V(5).ErrorS(err, "failed to accept websocket connection")
			return
		}
		defer c.Close(websocket.StatusInternalError, "")

		err = sc.subscribe(r.Context(), c)
		if errors.Is(err, context.Canceled) {
			klog.V(5).InfoS("websocket connection closed", "context", "cancelled")
			return
		}
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway {
			klog.V(5).InfoS("websocket connection closed", "status", websocket.CloseStatus(err))
			return
		}
		if err != nil {
			klog.V(5).ErrorS(err, "websocket connection closed with error")
			return
		}
	}
}

func handleDiscovery(clientset *kubernetes.Clientset, discovery *discoveryData) httprouter.Handle {
	lock := sync.Mutex{}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// The version of the Kubernetes API server can change, so we need to check it every time
		kubeVersion, err := clientset.ServerVersion()
		if err != nil {
			klog.Errorf("failed to get kube version: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		v := kubeVersion.String()

		if discovery.KubernetesVersion != v {
			lock.Lock()
			discovery.KubernetesVersion = v
			lock.Unlock()
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(discovery)
	}
}

func handleNodeDrain(clientset *kubernetes.Clientset, informer informers.GenericInformer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		nodeGeneric, exists, err := informer.Informer().GetIndexer().GetByKey(name)
		if err != nil {
			klog.Errorf("error getting node %q: %v", name, err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "error getting node"})
			return
		}
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}

		node := nodeGeneric.(*v1.Node)

		var sb strings.Builder
		helper := &drain.Helper{
			Client:              clientset,
			Force:               true,
			IgnoreAllDaemonSets: true,
			DeleteEmptyDirData:  true,
			GracePeriodSeconds:  -1,
			// If a pod is not evicted in 5 minutes, delete the pod
			Timeout: 5 * time.Minute,
			Out:     ioutil.Discard,
			ErrOut:  &sb,
			Ctx:     r.Context(),
		}
		if err := drain.RunCordonOrUncordon(helper, node, true); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			klog.ErrorS(err, "failed cordoning node", "name", name, "error", sb.String())
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed cordoning node: %v", err)})
			return
		}
		if err := drain.RunNodeDrain(helper, name); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			klog.ErrorS(err, "failed draining node", "name", name, "error", sb.String())
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed draining node: %v", err)})
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
	}
}

func getPathPrefix(gvr schema.GroupVersionResource, isNamespaced bool, prefixes ...string) string {
	return "/" + strings.Join(getPathSegments(gvr, isNamespaced, prefixes...), "/")
}

func getPathSegments(gvr schema.GroupVersionResource, isNamespaced bool, prefixes ...string) []string {
	n := len(prefixes) + 1 // prefixes + resource
	if len(gvr.Group) > 0 {
		n++
	}
	if isNamespaced {
		n += 2
	}
	segments := make([]string, n)
	copy(segments, prefixes)
	i := len(prefixes)
	if len(gvr.Group) > 0 {
		segments[i] = gvr.Group
		i++
	}
	if isNamespaced {
		segments[i] = "namespaces"
		segments[i+1] = ":namespace"
		i += 2
	}
	segments[i] = gvr.Resource
	return segments
}

func getConfig() *appConfig {
	flagSet := flag.NewFlagSet("dashboard", flag.ExitOnError)
	klog.InitFlags(flagSet)

	port := flagSet.String("port", "8999", "port to listen on")
	resyncPeriod := flagSet.Duration("resyncPeriod-period", 10*time.Minute, "informers resyncPeriod period")
	// create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.V(10).Info("error getting in-cluster config, falling back to local config")
		// create local config
		if !errors.Is(err, rest.ErrNotInCluster) {
			// the only recognized error
			klog.Fatal(fmt.Errorf("getting kube client config: %v", err.Error()))
		}

		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flagSet.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flagSet.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}

		err := flagSet.Parse(os.Args[1:])
		if err != nil {
			klog.Fatal(fmt.Errorf("parsing flags: %v", err.Error()))
		}

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			klog.Fatal(fmt.Errorf("building kube client config: %v", err.Error()))
		}

	} else {
		err := flagSet.Parse(os.Args[1:])
		if err != nil {
			klog.Fatal(fmt.Errorf("parsing flags: %v", err.Error()))
		}
	}

	return &appConfig{
		listenPort:   *port,
		resyncPeriod: *resyncPeriod,
		kubeConfig:   config,
	}
}

// type CRUD interface {
// 	List(context.Context, labels.Selector, fields.Selector) (*unstructured.UnstructuredList, error)
// 	Get(context.Context, string, metav1.GetOptions) (*unstructured.Unstructured, error)
// 	Create(context.Context, *unstructured.Unstructured, metav1.CreateOptions) (*unstructured.Unstructured, error)
// 	Update(context.Context, *unstructured.Unstructured, metav1.UpdateOptions) (*unstructured.Unstructured, error)
// 	Delete(context.Context, string, metav1.DeleteOptions) error
// }

// type GroupResourceConfig struct {
// 	GroupVersionResource schema.GroupVersionResource
// 	Namespace            string
// 	Informer             cache.SharedIndexInformer
// }

// type ListQuery struct {
// 	Namespace     string
// 	Name          string
// 	LabelSelector metav1.LabelSelector
// }

// type ItemQuery struct {
// 	Namespace string
// 	Name      string
// }
// type ApiGroupHandler interface {
// 	ApiGroup() string
// 	Resource() string
// 	Namespaced() bool

// 	ListHandler() func(ListQuery) (interface{}, error)
// 	ItemHandler() func(ItemQuery) (interface{}, error)
// 	UpdateHandler() func(obj interface{}) error
// 	CreateHandler() func(obj interface{}) error
// 	DeleteHandler() func(obj interface{}) error
// 	// Notify(Notifier)
// }

// func Register(mux http.ServeMux, gh ApiGroupHandler) {
// 	prefix := gh.ApiGroup() + "/" + gh.Resource()
// 	if gh.ListHandler() != nil {
// 		mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
// 			lq := ListQuery{}
// 			list, err := gh.ListHandler()(lq)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 			}
// 		})
// 	}
// 	mux.HandleFunc("/api/"+ApiGroupHandler.GroupVersion.Group+"/"+ApiGroupHandler.GroupVersion.Version+"/", ApiGroupHandler.Handle)
// }

// ii := theInformer.Informer()
// if err := ii.AddIndexers(cache.Indexers{
// 	"byName": func(obj interface{}) ([]string, error) {
// 		if node, ok := obj.(*unstructured.Unstructured); ok {
// 			name, found, err := unstructured.NestedString(node.Object, "metadata", "name")
// 			if err != nil {
// 				return nil, err
// 			}

// 			if !found {
// 				return nil, errors.New("name not found")
// 			}
// 		}
// 		return nil, nil
// 	},
// }); err != nil {
// 	klog.Fatal(fmt.Errorf("adding indexer: %v", err.Error()))
// }

/*
Там примерно такая логика:

	Клиент подключается.

	Клиент ожидает пинги { type: "ping" } . Если их не будет, он будет считать коннекшн stale и переконнекчиваться.

	Клиент делает запрос { command: "subscribe", identifier: "{\"channel\": \"MyChannel\"}"}

	Клиент ожидает ответ { type: "confirm_subscription", identifier: "{\"channel\": \"MyChannel\"}"}

	Клиент ожидает сообщения в канал { identifier: "{\"channel\": \"MyChannel\"}", message: "SOME JSON"}

	Клиент может слать в канал  { identifier: "{\"channel\": \"MyChannel\"}", command: "message", data: "SOME JSON"}

*/
