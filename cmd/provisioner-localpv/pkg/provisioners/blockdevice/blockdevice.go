package blockdevice

import (
	"github.com/golang/glog"
	pvController "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"github.com/openebs/maya/cmd/provisioner-localpv/pkg/provisioners"
	t "github.com/openebs/maya/cmd/provisioner-localpv/pkg/types"
	mconfig "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	mPV "github.com/openebs/maya/pkg/kubernetes/persistentvolume/v1alpha1"
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

// DeleteBlockDevice is invoked by the PVC controller to perform clean-up
//  activities before deleteing the PV object. If reclaim policy is
//  set to not-retain, then this function will delete the associated BDC
func (p *LocalProvisioner) Delete(pv *v1.PersistentVolume) (err error) {
	defer func() {
		err = errors.Wrapf(err, "failed to delete volume %v", pv.Name)
	}()

	blkDevOpts := &HelperBlockDeviceOptions{
		name: pv.Name,
	}

	//Determine if a BDC is set on the PV and save it to BlockDeviceOptions
	blkDevOpts.setBlockDeviceClaimFromPV(pv)

	//Initiate clean up only when reclaim policy is not retain.
	//TODO: this part of the code could be eliminated by setting up
	// BDC owner reference to PVC.
	glog.Infof("Release the Block Device Claim %v for PV %v", blkDevOpts.bdcName, pv.Name)

	if err := p.deleteBlockDeviceClaim(blkDevOpts); err != nil {
		glog.Infof("clean up volume %v failed: %v", pv.Name, err)
		return err
	}
	return nil
}

// ProvisionBlockDevice is invoked by the Provisioner to create a Local PV
//  with a Block Device
func (p *LocalProvisioner) Provision( /*opts pvController.VolumeOptions, volumeConfig *VolumeConfig*/ ) (*v1.PersistentVolume, error) {
	pvc := p.lp.VolumeOptions.PVC
	node := p.lp.VolumeOptions.SelectedNode
	name := p.lp.VolumeOptions.PVName
	capacity := p.lp.VolumeOptions.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	stgType := p.lp.VolumeConfig.GetStorageType()
	fsType := p.lp.VolumeConfig.GetFSType()

	//Extract the details to create a Block Device Claim
	blkDevOpts := &HelperBlockDeviceOptions{
		nodeName: node.Name,
		name:     name,
		capacity: capacity.String(),
	}

	path, blkPath, err := p.getBlockDevicePath(blkDevOpts)
	if err != nil {
		glog.Infof("Initialize volume %v failed: %v", name, err)
		return nil, err
	}
	glog.Infof("Creating volume %v on %v at %v(%v)", name, node.Name, path, blkPath)
	if path == "" {
		path = blkPath
		glog.Infof("Using block device{%v} with fs{%v}", blkPath, fsType)
	}

	// TODO
	// VolumeMode will always be specified as Filesystem for host path volume,
	// and the value passed in from the PVC spec will be ignored.
	fs := v1.PersistentVolumeFilesystem

	// It is possible that the HostPath doesn't already exist on the node.
	// Set the Local PV to create it.
	//hostPathType := v1.HostPathDirectoryOrCreate

	// TODO initialize the Labels and annotations
	// Use annotations to specify the context using which the PV was created.
	volAnnotations := make(map[string]string)
	volAnnotations[bdcStorageClassAnnotation] = blkDevOpts.bdcName
	//fstype := casVolume.Spec.FSType

	labels := make(map[string]string)
	labels[string(mconfig.CASTypeKey)] = "local-" + stgType
	//labels[string(v1alpha1.StorageClassKey)] = *className

	//TODO Change the following to a builder pattern
	pvObj, err := mPV.NewBuilder().
		WithName(name).
		WithLabels(labels).
		WithAnnotations(volAnnotations).
		WithReclaimPolicy(p.lp.VolumeOptions.PersistentVolumeReclaimPolicy).
		WithAccessModes(pvc.Spec.AccessModes).
		WithVolumeMode(fs).
		WithCapacityQty(pvc.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]).
		WithLocalHostPathFormat(path, fsType).
		WithNodeAffinity(node.Name).
		Build()

	if err != nil {
		return nil, err
	}
	return pvObj, nil
}
