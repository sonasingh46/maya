/*
Copyright 2018 The OpenEBS Authors.
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
package spc
import (
	"fmt"
	"github.com/golang/glog"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"github.com/openebs/maya/cmd/maya-apiserver/spc-actions"
)

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the spcPoolUpdated resource
// with the current status of the resource.
func (c *Controller) syncHandler(key, operation string) error {
	//glog.Infof("at sync handler")
	spcGot, err := c.getSpcResource(key)
	if err != nil {
		return err
	}
	status, err := c.spcEventHandler(operation, spcGot, key)
	if status == "Igonre" {
		return nil
	}
	return nil
}

// spcPoolEventHandler is to handle SPC related events.
func (c *Controller) spcEventHandler(operation string, spcGot *apis.StoragePoolClaim, key string) (string, error) {
	switch operation {
	case "add":
		// TO-DO : Handle Business Logic
		// Query for this spc object
		glog.Info("Create SPC Event Handler1")
		// Pass spc object from here to the function
		err := cstorpool.CreateCstorpool(spcGot)
		if err !=nil{
			fmt.Println("Could Not Create cstor pool")
		}
		break

	case "update":
		// TO-DO : Handle Business Logic
		glog.Info("Update SPC Event Handler")
		break

	case "delete":
		err := cstorpool.DeleteCstorpool(key)
		if err !=nil{
			fmt.Println("Could Not Delete cstor pool")
		}
		break
	default:
		// Ignore
		break
	}
	return "Ignore", nil
}

// enqueueSpc takes a SPC resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than SPC.
func (c *Controller) enqueueSpc(obj interface{}, q QueueLoad) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	q.Key = key
	c.workqueue.AddRateLimited(q)
}

// getSpcResource returns object corresponding to the resource key
func (c *Controller) getSpcResource(key string) (*apis.StoragePoolClaim, error) {
	// Convert the key(namespace/name) string into a distinct name
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil, nil
	}
	spcGot, err := c.clientset.OpenebsV1alpha1().StoragePoolClaims().Get(name,metav1.GetOptions{})
	if err != nil {
		// The SPC resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("spcGot '%s' in work queue no longer exists", key))
			return nil, nil
		}

		return nil, err
	}
	return spcGot, nil
}