package new


type PoolInstanceSpecGetter interface {
	Get() error
}

type PoolInstanceCreator interface {
	Create() error
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
	Create() error
}

type PoolDeploymentDeleter interface {
	Delete() error
}

type PoolDeploymentSpecGetter interface {
	Get() error
}

