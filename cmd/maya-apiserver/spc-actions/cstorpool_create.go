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
	"strconv"
)

func CreateCstorpool(spcGot *apis.StoragePoolClaim) (error) {
	// Business logic for creation of cstor pool cr
	// Launch as many go routines as the number of cstor pool crs need to be created.
	// How to handle the cr creation failure?

	// Fetch the number of cstor pools that should be created
	maxPool := spcGot.Spec.MaxPools
	// Convert maxPool which is a string type to integer type
	maxPoolCount, err := strconv.Atoi(maxPool)
	if err != nil {
		fmt.Println("conversion of max pool count failed : ", err)
		return err
	}
	poolType := spcGot.Spec.PoolSpec.PoolType
	if(poolType==""){
		fmt.Println("aborting cstor pool create operation as no poolType specified")
		return fmt.Errorf("error")
	}
	diskList := spcGot.Spec.Disks.DiskList
	if(len(diskList)==0){
		fmt.Println("aborting cstor pool create operation as no disk specified")
		return fmt.Errorf("error")
	}
	// Handle max pool count invalid inputs
	if maxPoolCount <= 0 {
		fmt.Println("aborting cstor pool create operation as pool count is not greater then ", maxPoolCount)
		return fmt.Errorf("error")
	}
	//cstorPoolCreator(spcGot)
	// Launch as many cstorPoolCreator go routines as maxPoolCount so as to start parallel creation of cstor pool CR
	for poolCount := 0; poolCount < maxPoolCount; poolCount++ {
		go cstorPoolCreator(spcGot, poolCount)
	}
	return nil
}

// function that creates a cstorpool CR
func cstorPoolCreator(spcGot *apis.StoragePoolClaim, nodeIndex int) {
	fmt.Println("Creation of cstor pool CR initiated Now Image.1")
	fmt.Println("Creating cstorpool cr for spc %s via CASTemplate", spcGot.ObjectMeta.Name)
	// Wether business logic will add some information other then extracted from spc for cstropool cr creation?
	// Create an empty cstor pool object
	cstorPool := &v1alpha1.CStorPool{}

	//Generate name using the prefix of StoragePoolClaim name and nodename hash
	cstorPool.ObjectMeta.Name = spcGot.Name
	// cstorPool.ObjectMeta.Name = spcGot.Name + "-" + spcGot.Spec.NodeSelector[nodeIndex]

	// Add Pooltype specification
	cstorPool.Spec.PoolSpec.PoolType = spcGot.Spec.PoolSpec.PoolType

	// Fetch castemplate from spc object( we need not fetch it)
	castName := spcGot.Annotations[string(v1alpha1.CASTemplateCVK)]
	// make a map that should contain the castemplate name
	mapCastName := make(map[string]string)
	// Fill the map with castemplate name
	mapCastName[string(v1alpha1.CASTemplateCVK)] = castName
	// Push the map to cstor pool cr object
	cstorPool.Annotations = mapCastName

	mapLabels := make(map[string]string)
	// Push storage pool claim name to cstor pool cr object as a label
	mapLabels[string(v1alpha1.StoragePoolClaimCVK)] = spcGot.Name

	// Add init status
	cstorPool.Status.Phase= v1alpha1.CStorPoolStatusInit


	// Push node hostname to cstor pool cr object as a label.

	// mapLabels[string(v1alpha1.CstorPoolHostNameCVK)] = spcGot.Spec.NodeSelector[nodeIndex]
	cstorPool.Labels = mapLabels

	// TODO : Select disks from nodes and push it to cstor pool cr object

	cstorOps, err := NewCstorPoolOperation(cstorPool)
	if err != nil {
		fmt.Println("NewCstorPoolOPeration Failed with following error")
		fmt.Println(err)
	}
	cstorPoolObject, err := cstorOps.Create()
	if err != nil {
		glog.Errorf("failed to create cas template based cstorpool: error '%s'", err.Error())
		//return nil, CodedError(500, err.Error())
	} else {
		glog.Infof("cas template based cstorpool created successfully: name '%s'", cstorPoolObject.Name)
	}
}
