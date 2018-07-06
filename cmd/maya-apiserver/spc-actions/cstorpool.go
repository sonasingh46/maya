/*
Copyright 2017 The OpenEBS Authors

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

package cstorpool

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	m_k8s_client "github.com/openebs/maya/pkg/client/k8s"
	mach_apis_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// cstorPool OperationOptions contains the options with respect to
// cstorPool related operations
type cstorPoolOperationOptions struct {
	// runNamespace is the namespace where cstorPool operation will happen
	//runNamespace string
	// k8sClient will make K8s API calls
	k8sClient *m_k8s_client.K8sClient
}

// cstorPoolOperation exposes methods with respect to cstorPool related operations
// e.g. read, create, delete.
type cstorPoolOperation struct {
	// cstorPoolOperationOptions has the options to various cstorPool related
	// operations
	cstorPoolOperationOptions
	// cstorPool to create or read or delete
	cstorPool *v1alpha1.CStorPool
}

// NewCstorPoolOperation returns a new instance of cstorPoolOperation
func NewCstorPoolOperation(cstorPool *v1alpha1.CStorPool) (*cstorPoolOperation, error) {
	if cstorPool == nil {
		return nil, fmt.Errorf("failed to instantiate cstorPool operation: nil cstorPool was provided")
	}

	// Clustered Scope ? Do we need a namespace ??
	//if len(cstorPool.Namespace) == 0 {
	//	return nil, fmt.Errorf("failed to instantiate cstorPool operation: missing run namespace")
	//}

	kc, err := m_k8s_client.NewK8sClient(cstorPool.Namespace)
	if err != nil {
		return nil, err
	}
	// Put cstor pool object inside cstorPoolOperation object
	return &cstorPoolOperation{
		cstorPool: cstorPool,
		cstorPoolOperationOptions: cstorPoolOperationOptions{
			k8sClient: kc,
			//runNamespace: cstorPool.Namespace,
		},
	}, nil
}

// Create provisions an OpenEBS cstorPool
func (v *cstorPoolOperation) Create() (*v1alpha1.CStorPool, error) {
	if v.k8sClient == nil {
		return nil, fmt.Errorf("unable to create cstorPool: nil k8s client")
	}

	//disks := v.cstorPool.Spec.Disks
	//if (disks) == 0 {
	//	disks = v.volume.Labels[string(v1alpha1.CapacityCVDK)]
	//}
	/*
	capacity := v.volume.Spec.Capacity
	if len(capacity) == 0 {
		capacity = v.volume.Labels[string(v1alpha1.CapacityCVDK)]
	}

	if len(capacity) == 0 {
		return nil, fmt.Errorf("unable to create volume: missing volume capacity")
	}*/

	// TODO
	//
	// UnComment below once provisioner is able to send name of PVC
	//
	// pvc name corresponding to this volume
	//pvcName := v.volume.Labels[string(v1alpha1.PersistentVolumeClaimCVK)]
	//if len(pvcName) == 0 {
	//	return nil, fmt.Errorf("unable to create volume: missing persistent volume claim")
	//}

	// fetch the pvc specifications
	//pvc, err := v.k8sClient.GetPVC(pvcName, mach_apis_meta_v1.GetOptions{})
	//if err != nil {
	//	return nil, err
	//}

	// extract the cas volume config from pvc
	//casConfigPVC := pvc.Annotations[string(v1alpha1.CASConfigCVK)]

	// TODO
	//
	// TODO :- We need not these variables in case of cstor pool cr  ( or we can have spcName )?
	// Remove below two lines once provisioner is able to send name of PVC
	//pvcName := ""
	//casConfigPVC := ""

	// get the storage class name corresponding to this volume
	/*
	scName := v.volume.Labels[string(v1alpha1.StorageClassCVK)]
	if len(scName) == 0 {
		return nil, fmt.Errorf("unable to create volume: missing storage class")
	}*/

	// fetch the storage pool claim specifications
	/*
	// Need to fetch the SPC object
	sc, err := v.k8sClient.GetStorageV1SC(scName, mach_apis_meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}*/

	// extract the cas volume config from storage class
	//casConfigSC := sc.Annotations[string(v1alpha1.CASConfigCVK)]

	// cas template to create a cas cstorPool

	//castName := sc.Annotations[string(v1alpha1.CASTemplateCVK)]
	castName := v.cstorPool.Annotations[string(v1alpha1.CASTemplateCVK)]
	if len(castName) == 0 {
		//return nil, fmt.Errorf("unable to create cstorPool: missing create cas template at '%s'", v1alpha1.CASTemplateCVK)
		return nil, fmt.Errorf("unable to create cstorPool: missing create cas template")
	}

	// fetch CASTemplate specifications
	//cast, err := v.k8sClient.GetOEV1alpha1CAST(castName, mach_apis_meta_v1.GetOptions{})
	cast, err := v.k8sClient.GetOEV1alpha1CAST(castName, mach_apis_meta_v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	// provision cas cstorPool via cas template engine
	cc, err := NewCASCreate(
		"null",
		"null",
		cast,
		map[string]string{
			// Make it cstor pool specific
			//string(v1alpha1.OwnerVTP):                 v.cstorPool.Name,
			string(v1alpha1.CstorPoolOwnerCTP):    v.cstorPool.Name,
			string(v1alpha1.StoragePoolClaimCTP):  v.cstorPool.Labels[string(v1alpha1.StoragePoolClaimCVK)],
			string(v1alpha1.CstorPoolHostNameCTP): v.cstorPool.Labels[string(v1alpha1.CstorPoolHostNameCVK)],
			string(v1alpha1.CstorPoolTypeCTP):     v.cstorPool.Spec.PoolSpec.PoolType,
			string(v1alpha1.CStorPoolPhaseCTP):    string(v.cstorPool.Status.Phase),
		},
	)
	if err != nil {
		return nil, err
	}

	// create the cstorPool
	data, err := cc.create()
	if err != nil {
		return nil, err
	}

	// unmarshall into openebs cstorPool
	cstorPool := &v1alpha1.CStorPool{}
	err = yaml.Unmarshal(data, cstorPool)
	if err != nil {
		return nil, err
	}
	return cstorPool, nil
}
