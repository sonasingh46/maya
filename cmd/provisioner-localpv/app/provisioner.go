package app

import (
	"fmt"
	"github.com/golang/glog"
	pvController "github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	events "github.com/openebs/maya/cmd/provisioner-localpv/app/analytics"
	"github.com/openebs/maya/cmd/provisioner-localpv/app/env"
	"github.com/openebs/maya/cmd/provisioner-localpv/pkg/provisioners"
	"github.com/openebs/maya/cmd/provisioner-localpv/pkg/provisioners/blockdevice"
	"github.com/openebs/maya/cmd/provisioner-localpv/pkg/provisioners/hostpath"
	t "github.com/openebs/maya/cmd/provisioner-localpv/pkg/types"
	mconfig "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	cast "github.com/openebs/maya/pkg/castemplate/v1alpha1"
	analytics "github.com/openebs/maya/pkg/usage"
	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"strings"
)

// Provisioner is  OpenEBS local provisioner
type Provisioner struct {
	lp *t.Provisioner
}

// NewProvisioner returns an empty instance of Provisioner.
// Notes: The constructor should be called with proper builders( e.g. WithKubeClient etc defined below)
// else it can be abused.
func NewProvisioner() *Provisioner {
	p := &Provisioner{}
	return p
}

// WithKubeClient sets the kubeclient field of the Provisioner.
func (p *Provisioner) WithKubeClient(kubeClient *clientset.Clientset) *Provisioner {
	p.lp.KubeClient = kubeClient
	return p
}

// WithStopChannel sets the stopch field of the Provisioner.
func (p *Provisioner) WithStopChannel(ch chan struct{}) *Provisioner {
	p.lp.StopCh = ch
	return p
}

// WithNameSpace sets the namespace field of the Provisioner.
func (p *Provisioner) WithNameSpace(namespace string) *Provisioner {
	p.lp.Namespace = namespace
	return p
}

// WithOpenEBSNameSpace sets the namespace field of the Provisioner.
func (p *Provisioner) WithOpenEBSNameSpace() *Provisioner {
	p.lp.Namespace = env.GetOpenEBSNamespace()
	return p
}

// WithHelperImage sets the helperImage field of the Provisioner.
func (p *Provisioner) WithHelperImage(helperImage string) *Provisioner {
	p.lp.HelperImage = env.GetDefaultHelperImage()
	return p
}

// WithHelperImage sets the helperImage field of the Provisioner.
func (p *Provisioner) WithDefaultHelperImage() *Provisioner {
	p.lp.HelperImage = env.GetDefaultBasePath()
	return p
}

// WithDefaultConfig sets the defaultConfig of Provisioner
func (p *Provisioner) WithDefaultConfig() *Provisioner {
	defaultConfig := []mconfig.Config{
		{
			Name:  t.KeyPVBasePath,
			Value: env.GetDefaultBasePath(),
		},
	}
	p.lp.DefaultConfig = defaultConfig
	return p
}

func (p *Provisioner) WithVolumeConfigFn() *Provisioner {
	p.lp.GetVolumeConfig = p.getVolumeConfig
	return p
}

// SupportsBlock will be used by controller to determine if block mode is
//  supported by the host path provisioner. Return false.
func (p *Provisioner) SupportsBlock() bool {
	return false
}

// Provision is invoked by the PVC controller which expect the PV
//  to be provisioned and a valid PV spec returned.
func (p *Provisioner) Provision(opts pvController.VolumeOptions) (*v1.PersistentVolume, error) {
	pvc := opts.PVC
	if pvc.Spec.Selector != nil {
		return nil, fmt.Errorf("claim.Spec.Selector is not supported")
	}
	for _, accessMode := range pvc.Spec.AccessModes {
		if accessMode != v1.ReadWriteOnce {
			return nil, fmt.Errorf("Only support ReadWriteOnce access mode")
		}
	}
	//node := opts.SelectedNode
	if opts.SelectedNode == nil {
		return nil, fmt.Errorf("configuration error, no node was specified")
	}

	name := opts.PVName

	// Create a new Config instance for the PV by merging the
	// default configuration with configuration provided
	// via PVC and the associated StorageClass
	pvCASConfig, err := p.lp.GetVolumeConfig(name, pvc)
	if err != nil {
		return nil, err
	}

	//TODO: Determine if hostpath or device based Local PV should be created
	stgType := pvCASConfig.GetStorageType()
	size := resource.Quantity{}
	reqMap := pvc.Spec.Resources.Requests
	if reqMap != nil {
		size = pvc.Spec.Resources.Requests["storage"]
	}
	events.SendEventOrIgnore(name, size.String(), stgType, analytics.VolumeProvision)

	p.GetLocalProvisioner(stgType, opts, pvCASConfig).Provision()

	return nil, fmt.Errorf("PV with StorageType %v is not supported", stgType)
}

// Delete is invoked by the PVC controller to perform clean-up
//  activities before deleteing the PV object. If reclaim policy is
//  set to not-retain, then this function will create a helper pod
//  to delete the host path from the node.
func (p *Provisioner) Delete(pv *v1.PersistentVolume) (err error) {
	defer func() {
		err = errors.Wrapf(err, "failed to delete volume %v", pv.Name)
	}()
	//Initiate clean up only when reclaim policy is not retain.
	if pv.Spec.PersistentVolumeReclaimPolicy != v1.PersistentVolumeReclaimRetain {
		//TODO: Determine the type of PV
		pvType := GetLocalPVType(pv)
		size := resource.Quantity{}
		reqMap := pv.Spec.Capacity
		if reqMap != nil {
			size = pv.Spec.Capacity["storage"]
		}

		events.SendEventOrIgnore(pv.Name, size.String(), pvType, analytics.VolumeDeprovision)

		p.GetLocalProvisioner(pvType, pvController.VolumeOptions{}, nil).Delete(pv)
	}
	glog.Infof("Retained volume %v", pv.Name)
	return nil
}

func (p *Provisioner) GetLocalProvisioner(ptype string, opts pvController.VolumeOptions, config *t.VolumeConfig) provision.Provisioner {
	if ptype == "hostpath" || ptype == "local-device" {
		return hostpath.NewLocalProvisioner(p.lp, config, opts)
	}
	if ptype == "device" {
		return blockdevice.NewLocalProvisioner(p.lp, config, opts)
	}
	return &provision.LocalProvisioner{}
}

//GetVolumeConfig creates a new VolumeConfig struct by
// parsing and merging the configuration provided in the PVC
// annotation - cas.openebs.io/config with the
// default configuration of the provisioner.
func (p *Provisioner) getVolumeConfig(pvName string, pvc *v1.PersistentVolumeClaim) (*t.VolumeConfig, error) {

	pvConfig := p.lp.DefaultConfig

	//Fetch the SC
	scName := GetStorageClassName(pvc)
	sc, err := p.lp.KubeClient.StorageV1().StorageClasses().Get(*scName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get storageclass: missing sc name {%v}", scName)
	}

	// extract and merge the cas config from storageclass
	scCASConfigStr := sc.ObjectMeta.Annotations[string(mconfig.CASConfigKey)]
	glog.Infof("SC %v has config:%v", *scName, scCASConfigStr)
	if len(strings.TrimSpace(scCASConfigStr)) != 0 {
		scCASConfig, err := cast.UnMarshallToConfig(scCASConfigStr)
		if err == nil {
			pvConfig = cast.MergeConfig(scCASConfig, pvConfig)
		} else {
			return nil, errors.Wrapf(err, "failed to get config: invalid sc config {%v}", scCASConfigStr)
		}
	}
	pvConfigMap, err := cast.ConfigToMap(pvConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read volume config: pvc {%v}", pvc.ObjectMeta.Name)
	}

	c := &t.VolumeConfig{
		PVName:  pvName,
		PVCName: pvc.ObjectMeta.Name,
		SCName:  *scName,
		Options: pvConfigMap,
	}
	return c, nil
}

// GetStorageClassName extracts the StorageClass name from PVC
func GetStorageClassName(pvc *v1.PersistentVolumeClaim) *string {
	// Use beta annotation first
	class, found := pvc.Annotations[t.BetaStorageClassAnnotation]
	if found {
		return &class
	}
	return pvc.Spec.StorageClassName
}

// GetLocalPVType extracts the Local PV Type from PV
func GetLocalPVType(pv *v1.PersistentVolume) string {
	casType, found := pv.Labels[string(mconfig.CASTypeKey)]
	if found {
		return casType
	}
	return ""
}
