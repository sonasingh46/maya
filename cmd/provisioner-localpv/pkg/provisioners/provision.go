package provision

import (
	pvController "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	t "github.com/openebs/maya/cmd/provisioner-localpv/pkg/types"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
)

type LocalProvisioner struct {
	Provisioner   *t.Provisioner
	VolumeConfig  *t.VolumeConfig
	VolumeOptions pvController.VolumeOptions
}

type Provisioner interface {
	Provision() (*v1.PersistentVolume, error)
	Delete(*v1.PersistentVolume) (err error)
}

func (lp *LocalProvisioner) Provision() (*v1.PersistentVolume, error) {
	return nil, errors.New("Not Supported")
}

func (lp *LocalProvisioner) Delete(pv *v1.PersistentVolume) error {
	return errors.New("Not Supported")
}
