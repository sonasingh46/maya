/*
Copyright 2018 The OpenEBS Authors

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
	nodeselectbeta "github.com/openebs/maya/pkg/algorithm/nodeselect/v1beta1"
	apisv1beta1 "github.com/openebs/maya/pkg/apis/openebs.io/v1beta1"
	openebs "github.com/openebs/maya/pkg/client/generated/clientset/versioned"
	spcpackage "github.com/openebs/maya/pkg/storagepoolclaim/v1beta1"
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"time"
)

// PoolCreateConfig is config object used to create a cstor pool.
type PoolCreateConfig struct {
	*nodeselectbeta.Operations
	*Controller
}

type clientSet struct {
	oecs openebs.Interface
}

func (c *Controller) NewPoolCreateConfig(spc *apisv1beta1.StoragePoolClaim) *PoolCreateConfig {
	ops := nodeselectbeta.NewOperations(spc)
	return &PoolCreateConfig{ops, c}
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the spcPoolUpdated resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	startTime := time.Now()
	glog.V(4).Infof("Started syncing storagepoolclaim %q (%v)", key, startTime)
	defer func() {
		glog.V(4).Infof("Finished syncing storagepoolclaim %q (%v)", key, time.Since(startTime))
	}()

	// Convert the namespace/name string into a distinct namespace and name
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the spc resource with this namespace/name
	spc, err := c.spcLister.Get(name)
	if k8serror.IsNotFound(err) {
		runtime.HandleError(fmt.Errorf("spc '%s' has been deleted", key))
		return nil
	}
	if err != nil {
		return err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	spcGot := spc.DeepCopy()
	err = c.syncSpc(spcGot)
	return err
}

// enqueueSpc takes a SPC resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than SPC.
func (c *Controller) enqueueSpc(spc interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(spc); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// synSpc is the function which tries to converge to a desired state for the spc.
func (c *Controller) syncSpc(spc *apisv1beta1.StoragePoolClaim) error {
	pc := c.NewPoolCreateConfig(spc)
	err := spcpackage.Validate(spcpackage.BuilderForAPIObject(spc).Spc)
	if err != nil {
		glog.Errorf("Validation of spc failed:%s", err)
		return nil
	}
	//pendingPoolCount, err := c.getPendingPoolCount(spc)
	pendingPoolCount, err := pc.GetPendingPoolCount()
	// TODO: Add context to error
	if err != nil {
		return err
	}

	if pc.IsPoolCreationPending() {
		err = pc.create(pendingPoolCount, spc)
		if err != nil {
			return err
		}
	}
	return nil
}

// create is a wrapper function that calls the actual function to create pool as many time
// as the number of pools need to be created.
func (pc *PoolCreateConfig) create(pendingPoolCount int, spc *apisv1beta1.StoragePoolClaim) error {
	var newSpcLease Leaser
	newSpcLease = &Lease{spc, SpcLeaseKey, pc.clientset, pc.kubeclientset}
	err := newSpcLease.Hold()
	if err != nil {
		return errors.Wrapf(err, "Could not acquire lease on spc object")
	}
	glog.V(4).Infof("Lease acquired successfully on storagepoolclaim %s ", spc.Name)
	defer newSpcLease.Release()
	for poolCount := 1; poolCount <= pendingPoolCount; poolCount++ {
		glog.Infof("Provisioning pool %d/%d for storagepoolclaim %s", poolCount, pendingPoolCount, spc.Name)
		err = pc.CreateStoragePool(spc)
		if err != nil {
			runtime.HandleError(errors.Wrapf(err, "Pool provisioning failed for %d/%d for storagepoolclaim %s", poolCount, pendingPoolCount, spc.Name))
		}
	}
	return nil
}
