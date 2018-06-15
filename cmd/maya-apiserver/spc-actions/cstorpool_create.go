/*
Copyright 2017 The Kubernetes Authors.

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
	"github.com/golang/glog"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/kubernetes"
	"github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
)
/*
var (
	masterURL  string
	kubeconfig string
)
*/
func CreateCstorpool(spcGot *apis.StoragePoolClaim){
	// Business logic for creation of cstor pool cr
	// Launch as many go routines as the number of cstor pool crs need to be created.
	// How to handle the cr creation failure?
	fmt.Println("Creating cstorpool cr for spc %s via CASTemplate",spcGot.ObjectMeta.Name)
	// Wether business logic will add some information other then extracted from spc for cstropool cr creation?
	cstorPool:= &v1alpha1.CStorPool{}
	cstorPool.Spec.PoolSpec.PoolName= "Pool1"
	cstorPool.Namespace= "default"
	// Fetch castemplate from spc object
	castName := spcGot.Annotations[string(v1alpha1.CASTemplateCVK)]
	fmt.Println("Cast Name Fetched:")
	fmt.Println(castName)

	cstorOps, err := NewCstorPoolOperation(cstorPool)
	if err != nil {
		fmt.Println("NewCstorPoolOPeration Failed with following error")
		fmt.Println(err)
	}
	cstorPoolObject, err := cstorOps.Create()
	if err != nil {
		glog.Errorf("failed to create cas template based cstorpool: error '%s'", err.Error())
		//return nil, CodedError(500, err.Error())
	}else {
		glog.Infof("cas template based cstorpool created successfully: name '%s'", cstorPoolObject.Name)
	}

}
