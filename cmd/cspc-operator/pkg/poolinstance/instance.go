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

package poolinstance

import (
	"fmt"
	"github.com/openebs/maya/cmd/cspc-operator/pkg/pooldeployment"
	pool "github.com/openebs/maya/cmd/cspc-operator/pkg/poolprovisioner"
	deploymentprovisioner "github.com/openebs/maya/cmd/cspc-operator/pkg/poolprovisioner/deployment-provisioner"
	nodeselect "github.com/openebs/maya/pkg/algorithm/nodeselect/v1alpha2"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	apicspi "github.com/openebs/maya/pkg/cstor/poolinstance/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

type PoolInstanceController struct {
	CSPC *apis.CStorPoolCluster
}

func NewPoolInstanceController(CSPC *apis.CStorPoolCluster) *PoolInstanceController {
	return &PoolInstanceController{CSPC: CSPC}
}

func (c *PoolInstanceController) Get(name string, opts metav1.GetOptions) (*apis.CStorPoolInstance, error) {
	CSPI, err := apicspi.NewKubeClient().WithNamespace(c.CSPC.Namespace).Get(name, opts)
	if err != nil {
		return nil, err
	}
	return CSPI, nil
}

func (c *PoolInstanceController) List(opts metav1.ListOptions) (*apis.CStorPoolInstanceList, error) {
	CSPI, err := apicspi.NewKubeClient().WithNamespace(c.CSPC.Namespace).List(opts)
	if err != nil {
		return nil, err
	}
	return CSPI, nil
}

func (c *PoolInstanceController) Create(CSPI *apis.CStorPoolInstance) (*apis.CStorPoolInstance, error) {
	gotCSPI, err := apicspi.NewKubeClient().WithNamespace(c.CSPC.Namespace).Create(CSPI)
	if err != nil {
		return nil, err
	}
	return gotCSPI, err
}

func (c *PoolInstanceController) Delete(name string, opts *metav1.DeleteOptions) error {
	err := apicspi.NewKubeClient().WithNamespace(c.CSPC.Namespace).Delete(name, opts)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	DeploymentProvisioner deploymentprovisioner.Provisioner
	PoolInstance          pool.InstanceInterface
	AlgorithmConfig       pool.AlgorithmInterface
}

func NewConfig(CSPC *apis.CStorPoolCluster, ac *nodeselect.Config) pool.Provisioner {
	return &Config{
		PoolInstance:          NewPoolInstanceController(CSPC),
		AlgorithmConfig:       ac,
		DeploymentProvisioner: pooldeployment.NewConfig(CSPC),
	}
}

func (c *Config) Provision() (*apis.CStorPoolInstance, error) {
	CSPI, err := c.AlgorithmConfig.GetPoolInstanceSpec()
	if err != nil {
		return nil, err
	}
	gotCSPI, err := c.PoolInstance.Create(CSPI)
	if err != nil {
		return nil, err
	}
	return gotCSPI, err
}

func (c *Config) Delete(name string, opts *metav1.DeleteOptions) error {
	err := c.PoolInstance.Delete(name, opts)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) DeleteAll() error {
	fmt.Println("Not Implemented")
	return nil
}

func (c *Config) IsPendingForCreation() bool {
	pendingPoolCount, err := c.AlgorithmConfig.GetPendingPoolCount()
	if err != nil {
		klog.Error("Error in getting pending pool count:", err)
		return false
	}
	return (pendingPoolCount > 0)
}
