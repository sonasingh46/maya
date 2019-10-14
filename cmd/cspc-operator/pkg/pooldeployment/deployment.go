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

package pooldeployment

import (
	provisioner "github.com/openebs/maya/cmd/cspc-operator/pkg/poolprovisioner/deployment-provisioner"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PoolDeploymentController struct{}

func NewPoolDeploymentController() *PoolDeploymentController {
	return &PoolDeploymentController{}
}

// TODO : Conplete this
func (pdc *PoolDeploymentController) GetSpec(CSPI *apis.CStorPoolInstance) (*appsv1.Deployment, error) {
	return nil, nil
}

// TODO : Conplete this
func (pdc *PoolDeploymentController) Create(poolDeploy *appsv1.Deployment) (*appsv1.Deployment, error) {
	return nil, nil
}

// TODO : Conplete this
func (pdc *PoolDeploymentController) Delete(name string, options *metav1.DeleteOptions) error {
	return nil
}

type Config struct {
	CSPC               *apis.CStorPoolCluster
	PoolDeployment     provisioner.PoolDeploymentInterface
	AlgorithmInterface provisioner.PoolDeploymentAgorithmInterface
}

func NewConfig(CSPC *apis.CStorPoolCluster) *Config {
	return &Config{
		CSPC:           CSPC,
		PoolDeployment: NewPoolDeploymentController(),
		// TODO : Pass the implementer of AlfgorithmInterface
	}
}

func (c *Config) GetDeploymentSpecForOrphanedCSPI() (*appsv1.DeploymentList, error) {
	return c.AlgorithmInterface.GetDeploymentSpecForOrphanedCSPI()
}

func (c *Config) ProvisionDeployment(CSPI *apis.CStorPoolInstance) (*appsv1.Deployment, error) {
	deploy, err := c.PoolDeployment.GetSpec(CSPI)
	if err != nil {
		return nil, err
	}
	gotDeployment, err := c.PoolDeployment.Create(deploy)
	if err != nil {
		return nil, err
	}
	return gotDeployment, err
}
// TODO : Conplete this
func (c *Config) Sync() error {
	return nil
}
