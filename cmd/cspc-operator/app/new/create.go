package new

import (
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"

)

type PoolInstanceSpecGetter interface {
	Get() (*apis.CStorPoolInstance, error)
}

type PoolInstanceCreator interface {
	Create() (*apis.CStorPoolInstance, error)
}

type PoolInstanceDeleter interface {
	Delete() error
}

type PoolInstanceWorker interface {
	PoolInstanceCreator
	PoolInstanceDeleter
	PoolInstanceSpecGetter
	IsPendingForCreation() bool
}

type PoolDeploymentWorker interface {
	PoolDeploymentCreator
	PoolDeploymentDeleter
	PoolDeploymentSpecGetter
	IsPendingForCreation() bool
}

type PoolDeploymentCreator interface {
	Create() (*appsv1.Deployment,error)
}

type PoolDeploymentDeleter interface {
	Delete() error
}

type PoolDeploymentSpecGetter interface {
	Get() (*appsv1.Deployment,error)
}
