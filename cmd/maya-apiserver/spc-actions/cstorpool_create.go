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
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"


	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/kubernetes"



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
	//cfg, err := getClusterConfig(kubeconfig)
	//if err != nil {
		//glog.Fatalf("Error building kubeconfig: %s", err.Error())
	//}
	//clientset, err := kubernetes.NewForConfig(cfg)
	//if err != nil {
	//	glog.Fatalf("Error in creating clientset: %s",err)

	//}
	//cstorPoolClientset
	//cstorPool := &v1alpha1.CStorPool{}
	//cstorPool.Spec.PoolSpec.PoolName= "Pool1"


	fmt.Println("Creating cstor pool cr "+spcGot.ObjectMeta.Name)

}
/*
// GetClusterConfig return the config for k8s.
func getClusterConfig(kubeconfig string) (*rest.Config, error) {
	var masterURL string
	cfg, err := rest.InClusterConfig()
	if err != nil {
		glog.Errorf("Failed to get k8s Incluster config. %+v", err)
		if kubeconfig == "" {
			return nil, fmt.Errorf("kubeconfig is empty: %v", err.Error())
		}
		cfg, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("Error building kubeconfig: %s", err.Error())
		}
	}
	return cfg, err
}
*/