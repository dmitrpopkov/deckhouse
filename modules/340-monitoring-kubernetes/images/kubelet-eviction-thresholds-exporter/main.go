package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/unix"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubelet/config/v1beta1"
)

const (
	nodeFsBytesAvailableEvictionSignal  = "nodefs.available"
	nodeFsInodesFreeEvictionSignal      = "nodefs.inodesFree"
	imageFsBytesAvailableEvictionSignal = "imagefs.available"
	imageFsInodesFreeEvictionSignal     = "imagefs.inodesFree"
)

var (
	containerdConfigRootDirRegex = regexp.MustCompile(`^\s*root\s*=\s*(.+)\s*$`)
)

type KubeletConfig struct {
	KubeletConfiguration v1beta1.KubeletConfiguration `json:"kubeletConfiguration"`
}

func main() {
	err := generateMetrics()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: signal handling
	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := generateMetrics()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func generateMetrics() error {
	containerRuntimeVersion, kubeletConfig, err := getContainerRuntimeAndKubeletConfig()
	if err != nil {
		log.Fatal(err)
	}

	runtimeRootDir, err := getRuntimeRootDir(strings.Split(containerRuntimeVersion, ":")[0])
	if err != nil {
		log.Fatal(err)
	}
	kubeletRootDir, err := getKubeletRootDir()
	if err != nil {
		log.Fatal(err)
	}

	nodeFsBytesAvail, nodeFsInodesAvail, err := getBytesAndInodeStatsFromPath(kubeletRootDir)
	imageFsBytesAvail, imageFsInodesAvail, err := getBytesAndInodeStatsFromPath(runtimeRootDir)
	softEvictionMap := kubeletConfig.KubeletConfiguration.EvictionSoft
	hardEvictionMap := kubeletConfig.KubeletConfiguration.EvictionHard

	evictionHardNodeFsBytesAvailable, err := extractPercent(nodeFsBytesAvail, nodeFsBytesAvailableEvictionSignal, hardEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionHardNodeFsInodesAvailable, err := extractPercent(nodeFsInodesAvail, nodeFsInodesFreeEvictionSignal, hardEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionHardImageFsBytesAvailable, err := extractPercent(imageFsBytesAvail, imageFsBytesAvailableEvictionSignal, hardEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionHardImagesFsInodesAvailable, err := extractPercent(imageFsInodesAvail, imageFsInodesFreeEvictionSignal, hardEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionSoftNodeFsBytesAvailable, err := extractPercent(nodeFsBytesAvail, nodeFsBytesAvailableEvictionSignal, softEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionSoftNodeFsInodesAvailable, err := extractPercent(nodeFsInodesAvail, nodeFsInodesFreeEvictionSignal, softEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionSoftImageFsBytesAvailable, err := extractPercent(imageFsBytesAvail, imageFsBytesAvailableEvictionSignal, softEvictionMap)
	if err != nil {
		log.Fatal(err)
	}
	evictionSoftImagesFsInodesAvailable, err := extractPercent(imageFsInodesAvail, imageFsInodesFreeEvictionSignal, softEvictionMap)
	if err != nil {
		log.Fatal(err)
	}

	fd, err := os.OpenFile("/var/run/node-exporter-textfile/kubelet-eviction.prom", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func(fd *os.File) {
		_ = fd.Close()
	}(fd)

	_, err = fmt.Fprintf(fd, `kubelet_eviction_nodefs_bytes{mountpoint="%s", type="hard"} %d`, kubeletRootDir, evictionHardNodeFsBytesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_nodefs_inodes{mountpoint="%s", type="hard"} %d`, kubeletRootDir, evictionHardNodeFsInodesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_imagefs_bytes{mountpoint="%s", type="hard"} %d`, runtimeRootDir, evictionHardImageFsBytesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_imagefs_inodes{mountpoint="%s", type="hard"} %d`, runtimeRootDir, evictionHardImagesFsInodesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_nodefs_bytes{mountpoint="%s", type="soft"} %d`, kubeletRootDir, evictionSoftNodeFsBytesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_nodefs_inodes{mountpoint="%s", type="soft"} %d`, kubeletRootDir, evictionSoftNodeFsInodesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_imagefs_bytes{mountpoint="%s", type="soft"} %d`, runtimeRootDir, evictionSoftImageFsBytesAvailable)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(fd, `kubelet_eviction_imagefs_inodes{mountpoint="%s", type="soft"} %d`, runtimeRootDir, evictionSoftImagesFsInodesAvailable)
	if err != nil {
		return err
	}

	return nil
}

func extractPercent(realResource uint64, signal string, evictionMap map[string]string) (int, error) {
	if evictionMap == nil {
		return 0, nil
	}

	evictionSignalValue, ok := evictionMap[signal]
	if !ok {
		return 0, nil
	}

	return parseThresholdStatement(realResource, evictionSignalValue)
}

func parseThresholdStatement(realResource uint64, val string) (int, error) {
	if strings.HasSuffix(val, "%") {
		// ignore 0% and 100%
		if val == "0%" || val == "100%" {
			return 0, nil
		}
		percentage, err := parsePercentage(val)
		if err != nil {
			return 0, err
		}
		if percentage < 0 {
			return 0, fmt.Errorf("eviction percentage threshold must be >= 0%%: %s", val)
		}
		// percentage is a float and should not be greater than 1 (100%)
		if percentage > 1 {
			return 0, fmt.Errorf("eviction percentage threshold must be <= 100%%: %s", val)
		}
		return int(percentage), nil
	}

	quantity, err := resource.ParseQuantity(val)
	if err != nil {
		return 0, err
	}
	if quantity.Sign() < 0 || quantity.IsZero() {
		return 0, fmt.Errorf("eviction threshold must be positive: %s", &quantity)
	}

	decQuantity, _ := quantity.AsInt64()
	return int((realResource - uint64(decQuantity)) / uint64(decQuantity) * 100), nil
}

func parsePercentage(input string) (float32, error) {
	value, err := strconv.ParseFloat(strings.TrimRight(input, "%"), 32)
	if err != nil {
		return 0, err
	}
	return float32(value) / 100, nil
}

func getBytesAndInodeStatsFromPath(path string) (bytesAvail uint64, inodeAvail uint64, err error) {
	var stat unix.Statfs_t

	err = unix.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}

	bytesAvail = stat.Bavail * uint64(stat.Bsize)
	inodeAvail = stat.Files

	return
}

func getKubeletRootDir() (string, error) {
	procs, err := process.Processes()
	if err != nil {
		return "", err
	}

	for _, p := range procs {
		cmdLine, err := p.CmdlineSlice()
		if err != nil {
			return "", err
		}

		if len(cmdLine) == 0 {
			continue
		}

		if !strings.Contains(cmdLine[0], "kubelet") {
			continue
		}

		const rootDirPrefix = "--root-dir="
		for _, arg := range cmdLine {
			if strings.HasPrefix(arg, rootDirPrefix) {
				return strings.TrimPrefix(arg, rootDirPrefix), nil
			}
		}
	}

	return "/var/lib/kubelet", nil
}

func getRuntimeRootDir(runtime string) (string, error) {
	switch runtime {
	case "containerd":
		return getContainerdRootDir()
	case "docker":
		return getDockerRootDir()
	}

	return "", fmt.Errorf(`unknown container runtime: "%s". Known containers runtimes: "docker"", "containerd"`, runtime)
}

func getDockerRootDir() (string, error) {
	type MiniDockerConfig struct {
		DockerRootDir string `json:"DockerRootDir"`
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				unixAddr, err := net.ResolveUnixAddr("unix", "/var/run/docker.sock")
				if err != nil {
					return nil, err
				}

				return net.DialUnix("unix", nil, unixAddr)
			},
		},
	}

	resp, err := httpClient.Get("http://system/info")
	if err != nil {
		return "", err
	}

	decoder := json.NewDecoder(resp.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var miniDockerConfig MiniDockerConfig
	err = decoder.Decode(&miniDockerConfig)
	if err != nil {
		return "", err
	}

	if len(miniDockerConfig.DockerRootDir) == 0 {
		return "", errors.New("DockerRootDir is empty")
	}

	return miniDockerConfig.DockerRootDir, nil
}

func getContainerdRootDir() (string, error) {
	containerdConfig, err := os.ReadFile("/etc/containerd/config.toml")
	if err != nil {
		return "/var/lib/containerd", nil
	}

	matches := containerdConfigRootDirRegex.FindSubmatch(containerdConfig)
	if len(matches) != 2 {
		return "/var/lib/containerd", nil
	}

	return string(matches[1]), err
}

func getContainerRuntimeAndKubeletConfig() (string, *KubeletConfig, error) {
	var (
		containerRuntimeVersion string
		kubeletConfig           *KubeletConfig
	)

	myNodeName, ok := os.LookupEnv("MY_NODE_NAME")
	if !ok {
		return "", nil, errors.New("no MY_NODE_NAME env")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return "", nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	nodeObj, err := clientset.CoreV1().Nodes().Get(ctx, myNodeName, metav1.GetOptions{})
	if err != nil {
		return "", nil, err
	}

	containerRuntimeVersion = nodeObj.Status.NodeInfo.ContainerRuntimeVersion

	request := clientset.CoreV1().RESTClient().Get().Resource("nodes").Name(myNodeName).SubResource("proxy").Suffix("configz")
	responseBytes, err := request.DoRaw(context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("failed to get config from Node %q: %s", myNodeName, err)
	}

	err = json.Unmarshal(responseBytes, &kubeletConfig)
	if err != nil {
		return "", nil, fmt.Errorf("can't unmarshal kubelet config %s: %s", responseBytes, err)
	}

	return containerRuntimeVersion, kubeletConfig, nil
}
