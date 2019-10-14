/*
Copyright 2019 The OpenEBS Authors

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

package app

import (
	"fmt"
	"github.com/openebs/maya/cmd/cspc-operator/pkg/pooldeployment"
	"github.com/openebs/maya/cmd/cspc-operator/pkg/poolinstance"
	pool "github.com/openebs/maya/cmd/cspc-operator/pkg/poolprovisioner"
	deploymentprovisioner "github.com/openebs/maya/cmd/cspc-operator/pkg/poolprovisioner/deployment-provisioner"
	nodeselect "github.com/openebs/maya/pkg/algorithm/nodeselect/v1alpha2"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	"time"
)

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the cspcPoolUpdated resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	startTime := time.Now()
	klog.V(4).Infof("Started syncing cstorpoolcluster %q (%v)", key, startTime)
	defer func() {
		klog.V(4).Infof("Finished syncing cstorpoolcluster %q (%v)", key, time.Since(startTime))
	}()

	// Convert the namespace/name string into a distinct namespace and name
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the cspc resource with this namespace/name
	cspc, err := c.cspcLister.CStorPoolClusters(ns).Get(name)
	if k8serror.IsNotFound(err) {
		runtime.HandleError(fmt.Errorf("cspc '%s' has been deleted", key))
		return nil
	}
	if err != nil {
		return err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	cspcGot := cspc.DeepCopy()
	CStorPool, err := NewCStorPoolProvisioner(cspcGot)
	if err != nil {
		return err
	}
	err = CStorPool.syncCSPC(cspcGot)
	return err
}

// enqueueCSPC takes a CSPC resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than CSPC.
func (c *Controller) enqueueCSPC(cspc interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(cspc); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

type CStorProvisioner struct {
	Pool           pool.Provisioner
	PoolDeployment deploymentprovisioner.Provisioner
}

func NewCStorPoolProvisioner(CSPC *apis.CStorPoolCluster) (*CStorProvisioner, error) {
	ac, err := nodeselect.
		NewBuilder().
		WithCSPC(CSPC).
		WithNameSpace(CSPC.Namespace).
		Build()
	if err != nil {
		return nil, errors.Wrap(err, "could not get algorithm config for provisioning")
	}
	return &CStorProvisioner{
		Pool: &poolinstance.Config{
			PoolInstance:    poolinstance.NewPoolInstanceController(CSPC),
			AlgorithmConfig: ac,
		},
		PoolDeployment: pooldeployment.NewConfig(CSPC),
	}, nil
}

func (CStor *CStorProvisioner) syncCSPC(cspc *apis.CStorPoolCluster) error {
	if cspc.DeletionTimestamp.IsZero() {
		err := CStor.Pool.DeleteAll()
		if err != nil {
			return err
		}
		return nil
	}

	for CStor.Pool.IsPendingForCreation() {
		_, err := CStor.Pool.Provision()
		if err != nil {
			return err
		}
	}

	err := CStor.PoolDeployment.Sync()
	if err != nil {
		return err
	}
	return nil
}
