/*
Copyright 2019 The OpenEBS Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package types

import (
	mconfig "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (

	//KeyPVStorageType defines if the PV should be backed
	// a hostpath ( sub directory or a storage device)
	KeyPVStorageType = "StorageType"
	//KeyPVBasePath defines base directory for hostpath volumes
	// can be configured via the StorageClass annotations.
	KeyPVBasePath = "BasePath"
	//KeyPVFSType defines filesystem type to be used with devices
	// and can be configured via the StorageClass annotations.
	KeyPVFSType = "FSType"
	//KeyPVRelativePath defines the alternate folder name under the BasePath
	// By default, the pv name will be used as the folder name.
	// KeyPVBasePath can be useful for providing the same underlying folder
	// name for all replicas in a Statefulset.
	// Will be a property of the PVC annotations.
	//KeyPVRelativePath = "RelativePath"
	//KeyPVAbsolutePath specifies a complete hostpath instead of
	// auto-generating using BasePath and RelativePath. This option
	// is specified with PVC and is useful for granting shared access
	// to underlying hostpaths across multiple pods.
	//KeyPVAbsolutePath = "AbsolutePath"
	// Some of the PVCs launched with older helm charts, still
	// refer to the StorageClass via beta annotations.
	BetaStorageClassAnnotation = "volume.beta.kubernetes.io/storage-class"

	// LocalPVFinalizer represents finalizer string used by LocalPV
	LocalPVFinalizer = "local.openebs.io/finalizer"
)

//Provisioner struct has the configuration and utilities required
// across the different work-flows.
type Provisioner struct {
	StopCh      chan struct{}
	KubeClient  *clientset.Clientset
	Namespace   string
	HelperImage string
	// defaultConfig is the default configurations
	// provided from ENV or Code
	DefaultConfig []mconfig.Config
	// getVolumeConfig is a reference to a function
	GetVolumeConfig GetVolumeConfigFn
}

//olumeConfig struct contains the merged configuration of the PVC
// and the associated SC. The configuration is derived from the
// annotation `cas.openebs.io/config`. The configuration will be
// in the following json format:
// {
//   Key1:{
//	enabled: true
//	value: "string value"
//   },
//   Key2:{
//	enabled: true
//	value: "string value"
//   },
// }
type VolumeConfig struct {
	PVName  string
	PVCName string
	SCName  string
	Options map[string]interface{}
}

// GetVolumeConfigFn allows to plugin a custom function
//  and makes it easy to unit test provisioner
type GetVolumeConfigFn func(pvName string, pvc *v1.PersistentVolumeClaim) (*VolumeConfig, error)
