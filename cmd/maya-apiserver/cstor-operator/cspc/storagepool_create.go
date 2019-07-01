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

package cspc

import (
	"github.com/golang/glog"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/maya/pkg/storagepool"
	"github.com/pkg/errors"
)

// Cas template is a custom resource which has a list of runTasks.

// runTasks are configmaps which has defined yaml templates for resources that needs
// to be created or deleted for a storagepool creation or deletion respectively.

// CreateStoragePool is a function that does following:
// 1. It receives storagepoolclaim object from the cspc watcher event handler.
// 2. After successful validation, it will call a worker function for actual storage creation
//    via the cas template specified in storagepoolclaim.
func (c *Controller) CreateStoragePool(cspcGot *apis.CStorPoolCluster) error {
	//poolconfig := c.NewPoolCreateConfig(cspcGot)
	//newCasPool, err := poolconfig.getCasPool(cspcGot)

	newCasPool, err := c.getCasPool()

	if err != nil {
		return errors.Wrapf(err, "failed to build cas pool for cspc %s", cspcGot.Name)
	}

	// Calling worker function to create storagepool
	err = poolCreateWorker(newCasPool)
	if err != nil {
		return err
	}

	return nil
}

// poolCreateWorker is a worker function which will create a storagepool.
func poolCreateWorker(pool *apis.CasPool) error {

	glog.Infof("Creating storagepool for storagepoolclaim %s via CASTemplate", pool.StoragePoolClaim)

	storagepoolOps, err := storagepool.NewCasPoolOperation(pool)
	if err != nil {
		return errors.Wrapf(err, "NewCasPoolOperation failed error")
	}
	_, err = storagepoolOps.Create()
	if err != nil {
		return errors.Wrapf(err, "failed to create cas template based storagepool")

	}
	glog.Infof("Cas template based storagepool created successfully: name '%s'", pool.StoragePoolClaim)
	return nil
}

// TODO: Complete following function ( Mock Only ) 
func (c *Controller) getCasPool() (*apis.CasPool, error) {
	return &apis.CasPool{}, nil
}
