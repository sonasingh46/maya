package hostpath

import (
	"github.com/golang/glog"
	pvController "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"github.com/openebs/maya/cmd/provisioner-localpv/pkg/provisioners"
	t "github.com/openebs/maya/cmd/provisioner-localpv/pkg/types"
	mconfig "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	persistentvolume "github.com/openebs/maya/pkg/kubernetes/persistentvolume/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
)

type LocalProvisioner struct {
	lp *provision.LocalProvisioner
}

func NewLocalProvisioner(provisioner *t.Provisioner, volumeConfig *t.VolumeConfig, volumeOptions pvController.VolumeOptions) provision.Provisioner {
	return &LocalProvisioner{
		lp: &provision.LocalProvisioner{
			Provisioner:   provisioner,
			VolumeConfig:  volumeConfig,
			VolumeOptions: volumeOptions,
		},
	}
}

func (p *LocalProvisioner) Provision() (*v1.PersistentVolume, error) {
	pvc := p.lp.VolumeOptions.PVC
	node := p.lp.VolumeOptions.SelectedNode
	name := p.lp.VolumeOptions.PVName
	stgType := p.lp.VolumeConfig.GetStorageType()

	path, err := p.lp.VolumeConfig.GetPath()
	if err != nil {
		return nil, err
	}

	glog.Infof("Creating volume %v at %v:%v", name, node.Name, path)

	//Before using the path for local PV, make sure it is created.
	initCmdsForPath := []string{"mkdir", "-m", "0777", "-p"}
	podOpts := &HelperPodOptions{
		cmdsForPath: initCmdsForPath,
		name:        name,
		path:        path,
		nodeName:    node.Name,
	}

	iErr := p.createInitPod(podOpts)
	if iErr != nil {
		glog.Infof("Initialize volume %v failed: %v", name, iErr)
		return nil, iErr
	}

	// VolumeMode will always be specified as Filesystem for host path volume,
	// and the value passed in from the PVC spec will be ignored.
	fs := v1.PersistentVolumeFilesystem

	// It is possible that the HostPath doesn't already exist on the node.
	// Set the Local PV to create it.
	//hostPathType := v1.HostPathDirectoryOrCreate

	// TODO initialize the Labels and annotations
	// Use annotations to specify the context using which the PV was created.
	//volAnnotations := make(map[string]string)
	//volAnnotations[string(v1alpha1.CASTypeKey)] = casVolume.Spec.CasType
	//fstype := casVolume.Spec.FSType

	labels := make(map[string]string)
	labels[string(mconfig.CASTypeKey)] = "local-" + stgType
	//labels[string(v1alpha1.StorageClassKey)] = *className

	//TODO Change the following to a builder pattern
	pvObj, err := persistentvolume.NewBuilder().
		WithName(name).
		WithLabels(labels).
		WithReclaimPolicy(p.lp.VolumeOptions.PersistentVolumeReclaimPolicy).
		WithAccessModes(pvc.Spec.AccessModes).
		WithVolumeMode(fs).
		WithCapacityQty(pvc.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]).
		WithLocalHostDirectory(path).
		WithNodeAffinity(node.Name).
		Build()

	if err != nil {
		return nil, err
	}

	return pvObj, nil
}

func (p *LocalProvisioner) Delete(pv *v1.PersistentVolume) (err error) {
	defer func() {
		err = errors.Wrapf(err, "failed to delete volume %v", pv.Name)
	}()

	//Determine the path and node of the Local PV.
	pvObj := persistentvolume.NewForAPIObject(pv)
	path := pvObj.GetPath()
	if path == "" {
		return errors.Errorf("no HostPath set")
	}

	node := pvObj.GetAffinitedNode()
	if node == "" {
		return errors.Errorf("cannot find affinited node")
	}

	//Initiate clean up only when reclaim policy is not retain.
	glog.Infof("Deleting volume %v at %v:%v", pv.Name, node, path)
	cleanupCmdsForPath := []string{"rm", "-rf"}
	podOpts := &HelperPodOptions{
		cmdsForPath: cleanupCmdsForPath,
		name:        pv.Name,
		path:        path,
		nodeName:    node,
	}

	if err := p.createCleanupPod(podOpts); err != nil {
		return errors.Wrapf(err, "clean up volume %v failed", pv.Name)
	}
	return nil
}
