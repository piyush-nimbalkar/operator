package storage

import (
	corev1alpha1 "github.com/libopenstorage/operator/pkg/apis/core/v1alpha1"
	"github.com/libopenstorage/operator/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
)

// Driver defines an external storage driver interface.
// Any driver that wants to be used with openstorage operator needs
// to implement these interfaces.
type Driver interface {
	// Init initializes the storage driver
	Init(interface{}) error

	// String returns the string name of the driver
	String() string

	// ClusterPluginInterface interface to manage storage cluster
	ClusterPluginInterface
}

// ClusterPluginInterface interface to manage storage cluster
type ClusterPluginInterface interface {
	// PreInstall the driver should do whatever it is needed before the pods
	// start to make sure the cluster comes up correctly. This should be
	// idempotent and subsequent calls should result in the same result.
	PreInstall(*corev1alpha1.StorageCluster) error
	// GetStoragePodSpec given the storage cluster spec it returns the pod spec
	GetStoragePodSpec(*corev1alpha1.StorageCluster) v1.PodSpec
	// GetSelectorLabels returns driver specific labels that are applied on the pods
	GetSelectorLabels() map[string]string
	// SetDefaultsOnStorageCluster sets the driver specific defaults on the storage
	// cluster spec if they are not set
	SetDefaultsOnStorageCluster(*corev1alpha1.StorageCluster)
}

var (
	storageDrivers = make(map[string]Driver)
)

// Register registers the given storage driver
func Register(name string, d Driver) error {
	logrus.Debugf("Registering storage driver: %v", name)
	storageDrivers[name] = d
	return nil
}

// Get an external storage driver to be used with the operator
func Get(name string) (Driver, error) {
	d, ok := storageDrivers[name]
	if ok {
		return d, nil
	}
	return nil, &errors.ErrNotFound{
		ID:   name,
		Type: "StorageDriver",
	}
}