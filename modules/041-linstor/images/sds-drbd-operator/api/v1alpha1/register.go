package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	BDKind                          = "BlockDevice"
	APIGroup                        = "storage.deckhouse.io"
	APIVersion                      = "v1alpha1" // v1alpha1
	OwnerReferencesAPIVersion       = "v1"
	OwnerReferencesKind             = "BlockDevice"
	OwnerReferencesLVMChangeRequest = "LVMChangeRequest"
	TypeMediaAPIVersion             = APIGroup + "/" + APIVersion
	Node                            = "Node"
)

// SchemeGroupVersion is group version used to register these objects
var (
	SchemeGroupVersion = schema.GroupVersion{
		Group:   APIGroup,
		Version: APIVersion,
	}
	// SchemeBuilder tbd
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme tbd
	AddToScheme = SchemeBuilder.AddToScheme
)

// ConsumableBlockDeviceGVK is group version kind for BlockDevice.
var ConsumableBlockDeviceGVK = schema.GroupVersionKind{
	Group:   SchemeGroupVersion.Group,
	Version: SchemeGroupVersion.Version,
	Kind:    BDKind,
}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ConsumableBlockDevice{},
		&ConsumableBlockDeviceList{},
		&LVMChangeRequest{},
		&LVMChangeRequestList{},
		&LinstorStorageClass{},
		&LinstorStorageClassList{},
		&LinstorStoragePool{},
		&LinstorStoragePoolList{},
		&LvmVolumeGroup{},
		&LvmVolumeGroupList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
