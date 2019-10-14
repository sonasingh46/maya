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

package deploymentprovisioner

import (
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO : Provide Comments

type PoolDeploymentAgorithmInterface interface {
	GetDeploymentSpecForOrphanedCSPI() (*appsv1.DeploymentList, error)
}

type Provisioner interface {
	PoolDeploymentAgorithmInterface
	ProvisionDeployment(CSPI *apis.CStorPoolInstance) (*appsv1.Deployment, error)
	Sync() error
}

type PoolDeploymentInterface interface {
	Create(poolDeploy *appsv1.Deployment) (*appsv1.Deployment, error)
	Delete(name string, options *metav1.DeleteOptions) error
	GetSpec(CSPI *apis.CStorPoolInstance) (*appsv1.Deployment, error)
}
