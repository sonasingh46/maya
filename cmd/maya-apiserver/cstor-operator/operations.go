package spc

import (
	"encoding/json"
	nodeselect "github.com/openebs/maya/pkg/algorithm/nodeselect/v1alpha1"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/maya/pkg/hash/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// PoolConfig is the config to carry out disk operations.
type PoolConfig struct {
	spc        *apis.StoragePoolClaim
	cspList    *apis.CStorPoolList
	controller *Controller
}
// handlerfunc is typed predicates to handle disk operations performed on spc.
type handlerfunc func(*PoolConfig) (*apis.StoragePoolClaim, error)

// HandlerPredicates contains a list of predicates that should be executed in order so as to
// reach desired state in response to any cahnge in disk list on spc.
var HandlerPredicates = []handlerfunc{
	HandleDiskRemoval,
	HandleDiskAddition,
}

// NewPoolConfig is the constructor for PoolConfig struct.
func (c *Controller) NewPoolConfig(spc *apis.StoragePoolClaim) (*PoolConfig, error) {
	cspList, err := c.getCsp(spc)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not list CSP for SPC %s", spc.Name)
	}
	newPoolConfig := &PoolConfig{
		spc:        spc,
		cspList:    cspList,
		controller: c,
	}
	return newPoolConfig, nil
}

// handleDiskHashChange is called if the hash of disk list changes on SPC.
func (c *Controller) handleDiskHashChange(spc *apis.StoragePoolClaim) (*apis.StoragePoolClaim, error) {
	err := c.executeHandlerPredicates(spc)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute handler predicates for disk operations for spc %s", spc.Name)
	}
	// patch spc with the new disk list hash once the execution is successful.
	spc, err = c.patchSpcWithDiskHash(spc)
	if err != nil {
		return nil, errors.Wrapf(err, "could not patch spc %s with newer disk hash for disk operations", spc.Name)
	}
	return spc, nil
}

// executeHandlerPredicates executes all the handler predicates in order.
func (c *Controller) executeHandlerPredicates(spc *apis.StoragePoolClaim) error {
	for _, p := range HandlerPredicates {
		poolConfig, err := c.NewPoolConfig(spc)
		if err != nil {
			errors.Wrapf(err, "could not initialize the disk operations object for spc %s", spc.Name)
		}
		_, err = p(poolConfig)
		if err != nil {
			errors.Wrapf(err, "disk operation was not successful for spc %s", spc.Name)
		}
	}
	return nil
}

type removeDiskOps func()

var diskRemovalOps = []removeDiskOps{}

type addDiskOps func(*PoolConfig)

var diskAdditionOps = []func() addDiskOps{
	ReattachDisk,
	ReplaceDisk,
	ExpandPool,
}

func ReattachDisk() addDiskOps {
	return func(pc *PoolConfig) {
		dettachedDisk := pc.getDettachedCspDisks()
		for disk, _ := range dettachedDisk {
			pc.reAttachDisk(disk)
		}
	}
}

func ReplaceDisk() addDiskOps {
	return func(pc *PoolConfig) {
		replacementDisks := pc.getAddedDisks()
		for _, disk := range replacementDisks {
			pc.replaceDisk(disk)
		}
	}
}

func ExpandPool() addDiskOps {
	return func(pc *PoolConfig) {
		nodeCspMap := pc.getCspNodeMap()
		newDisks := pc.getAddedDisks()
		nodeDisk := pc.getnodeDiskMap(newDisks)
		for node, disks := range nodeDisk {
			csp := nodeCspMap[node]
			if csp == nil {
				continue
			}
			pc.expandCsp(csp, disks)
		}
	}
}



func HandleDiskRemoval(pc *PoolConfig) (*apis.StoragePoolClaim, error) {
	removedDisks := pc.getRemovedDisks()
	for _, disk := range removedDisks {
		pc.removeDisk(disk)
	}
	for _, csp := range pc.cspList.Items {
		csp, err := pc.controller.updateCsp(&csp)
		if isTopVdevLost(csp) {
			err := pc.controller.deleteCsp(csp)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to delete csp %s for disk operations for spc %s", csp.Name, pc.spc.Name)
			}
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to update csp %s for disk operations for spc %s", csp.Name, pc.spc.Name)
		}
	}
	return pc.spc, nil
}

func HandleDiskAddition(pc *PoolConfig) (*apis.StoragePoolClaim, error) {
	for _, p := range diskAdditionOps {
		p()(pc)
		err := pc.updateCspList()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to update disk operations for spc %s", pc.spc.Name)
		}
	}
	return pc.spc, nil
}

func enqueueAddOperation(csp *apis.CStorPool, deviceIDs []string) *apis.CStorPool {
	newAdpcperation := &apis.CstorOperation{
		Action:   "Add",
		Status:   "Init",
		NewDisks: deviceIDs,
	}
	csp.Operations = append(csp.Operations, *newAdpcperation)
	return csp
}

func enqueueDeleteOperation(csp *apis.CStorPool) *apis.CStorPool {
	newAdpcperation := &apis.CstorOperation{
		Action: "Delete",
		Status: "Init",
	}
	csp.Operations = append(csp.Operations, *newAdpcperation)
	return csp
}

func (pc *PoolConfig) expandCsp(csp *apis.CStorPool, disks []DiskDetails) {
	var newGroup apis.DiskGroup
	var cspdisk []apis.CspDisk
	var deviceIDs []string
	defaultDiskCount := nodeselect.DefaultDiskCount[pc.spc.Spec.PoolSpec.PoolType]
	diskCount := 0
	if len(disks) >= defaultDiskCount {
		diskCount = (len(disks) / defaultDiskCount) * defaultDiskCount
		for i := 0; i < defaultDiskCount; i = i + defaultDiskCount {
			for j := 0; j < diskCount; j++ {
				var item apis.CspDisk
				item.Name = disks[j].DiskName
				item.DeviceID = disks[j].DeviceID
				item.InUseByPool = true
				deviceIDs = append(deviceIDs, disks[j].DeviceID)
				cspdisk = append(cspdisk, item)
			}
			newGroup.Item = cspdisk
			csp.Spec.Group = append(csp.Spec.Group, newGroup)
		}
		enqueueAddOperation(csp, deviceIDs)
		pc.updatePoolConfig(*csp)
	}
}

func (pc *PoolConfig) updatePoolConfig(csp apis.CStorPool) {
	for i, cspGot := range pc.cspList.Items {
		if cspGot.Name == csp.Name {
			pc.cspList.Items[i] = csp
		}
	}
}

func (pc *PoolConfig) getnodeDiskMap(disks []string) map[string][]DiskDetails {
	nodeDiskMap := make(map[string][]DiskDetails)
	for _, disk := range disks {
		gotDisk, err := pc.controller.clientset.OpenebsV1alpha1().Disks().Get(disk, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		if gotDisk == nil {
			return nil
		}
		devID := getDeviceId(gotDisk)
		disk := &DiskDetails{
			DiskName: gotDisk.Name,
			DeviceID: devID,
		}
		nodeDiskMap[gotDisk.Labels[string(apis.HostNameCPK)]] = append(nodeDiskMap[gotDisk.Labels[string(apis.HostNameCPK)]], *disk)
	}
	return nodeDiskMap
}

func (pc *PoolConfig) getCspNodeMap() map[string]*apis.CStorPool {
	cspNodeMap := make(map[string]*apis.CStorPool)
	for _, csp := range pc.cspList.Items {
		cspCopy := csp
		cspNodeMap[csp.Labels[string(apis.HostNameCPK)]] = &cspCopy
	}
	return cspNodeMap
}

func (pc *PoolConfig) reAttachDisk(diskName string) {
	for i, csp := range pc.cspList.Items {
		for j, group := range csp.Spec.Group {
			for k, disk := range group.Item {
				if disk.Name == diskName {
					pc.cspList.Items[i].Spec.Group[j].Item[k].InUseByPool = true
				}
			}
		}
	}
}

func (pc *PoolConfig) replaceDisk(diskName string) {
	for i, csp := range pc.cspList.Items {
		for j, group := range csp.Spec.Group {
			for k, disk := range group.Item {
				if disk.InUseByPool == false {
					pc.cspList.Items[i].Spec.Group[j].Item[k].Name = diskName
					pc.cspList.Items[i].Spec.Group[j].Item[k].InUseByPool = true
				}
			}
		}
	}
}

func (pc *PoolConfig) removeDisk(diskName string) {
	for i, csp := range pc.cspList.Items {
		for j, group := range csp.Spec.Group {
			for k, disk := range group.Item {
				if disk.Name == diskName {
					pc.cspList.Items[i].Spec.Group[j].Item[k].InUseByPool = false
				}
			}
		}
	}
}

// getSpcDisks returns map of spc disks present on SPC.
func (pc *PoolConfig) getSpcDisks() map[string]bool {
	// Make a map containing all the disks present in spc.
	spcDisks := make(map[string]bool)
	for _, disk := range pc.spc.Spec.Disks.DiskList {
		spcDisks[disk] = true
	}
	return spcDisks
}

func (pc *PoolConfig) getCspDisks() map[string]bool {
	// Make a map containing all the disks present in csp
	// Get all CSP corresponding to the SPC
	cspDisks := make(map[string]bool)
	for _, csp := range pc.cspList.Items {
		for _, group := range csp.Spec.Group {
			for _, disk := range group.Item {
				cspDisks[disk.Name] = true
			}
		}
	}
	return cspDisks
}

func (pc *PoolConfig) getDettachedCspDisks() map[string]bool {
	// Make a map containing all the disks present in csp whis in not present in SPC.
	spcDisks := pc.getSpcDisks()
	cspDisks := make(map[string]bool)
	for _, csp := range pc.cspList.Items {
		for _, group := range csp.Spec.Group {
			for _, disk := range group.Item {
				if spcDisks[disk.Name] == true && disk.InUseByPool == false {
					cspDisks[disk.Name] = true
				}
			}
		}
	}
	return cspDisks
}

// getRemovedDisks return a list of disks present on all CSPs for a given SPC, which has been removed from SPC
func (pc *PoolConfig) getRemovedDisks() []string {
	var removedDisk []string
	// Get the disks present on CSPs
	cspDisks := pc.getCspDisks()

	// get the disk present on SPC
	spcDisks := pc.getSpcDisks()
	for disk, _ := range cspDisks {
		if spcDisks[disk] == false {
			removedDisk = append(removedDisk, disk)
		}
	}
	return removedDisk
}

// getAddedDisks returns a list of disk that is added to SPC.
func (pc *PoolConfig) getAddedDisks() []string {
	var addedDisk []string
	// get the disks present on CSPs
	cspDisks := pc.getCspDisks()
	// get the disk present on SPC
	spcDisks := pc.getSpcDisks()
	for disk, _ := range spcDisks {
		if cspDisks[disk] == false {
			addedDisk = append(addedDisk, disk)
		}
	}
	return addedDisk
}


func (pc *PoolConfig) updateCspList() error {
	for i, csp := range pc.cspList.Items {
		csp, err := pc.controller.updateCsp(&csp)
		if err != nil {
			errors.Wrapf(err, "failed to update csp %s for disk replacement operations for spc %s", csp.Name, pc.spc.Name)
		}
		pc.cspList.Items[i] = *csp
	}
	return nil
}
// TODO: Patch using patch package.
func (c *Controller) patchSpcWithDiskHash(spc *apis.StoragePoolClaim) (*apis.StoragePoolClaim, error) {

	diskHash, _ := hash.Hash(spc.Spec.Disks)
	spcPatch := make([]Patch, 1)
	spcPatch[0].Op = PatchOperation
	// TODO: If there is no annotaion in SPC -- Create it
	if spc.Annotations == nil {
		return nil, errors.Errorf("No annotation found in spc %s", spc.Name)
	}
	if spc.Annotations[spcDiskHashKey] == "" {
		spcPatch[0].Op = PatchOperationAdd
	}
	spcPatch[0].Path = spcDiskHashKeyPath
	spcPatch[0].Value = diskHash
	spcPatchJSON, err := json.Marshal(spcPatch)
	spcGot, err := c.clientset.OpenebsV1alpha1().StoragePoolClaims().Patch(spc.Name, types.JSONPatchType, spcPatchJSON)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to patch spc %s with the new disk list hash", spc.Name)
	}

	return spcGot, nil
}


func (c *Controller) getCsp(spc *apis.StoragePoolClaim) (*apis.CStorPoolList, error) {
	cspList, err := c.clientset.OpenebsV1alpha1().CStorPools().List(metav1.ListOptions{LabelSelector: string(apis.StoragePoolClaimCPK) + "=" + spc.Name})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get csp objects for spc %s", spc.Name)
	}
	return cspList, nil

}

func (c *Controller) updateCsp(csp *apis.CStorPool) (*apis.CStorPool, error) {
	csp, err := c.clientset.OpenebsV1alpha1().CStorPools().Update(csp)
	return csp, err
}

func (c *Controller) deleteCsp(csp *apis.CStorPool) error {
	err := c.clientset.OpenebsV1alpha1().CStorPools().Delete(csp.Name, &metav1.DeleteOptions{})
	return err
}

func isTopVdevLost(csp *apis.CStorPool) bool {
	for _, group := range csp.Spec.Group {
		count := 0
		for _, disk := range group.Item {
			if disk.InUseByPool == false {
				count++
			}
		}
		if count >= 1 && csp.Spec.PoolSpec.PoolType == string(apis.PoolTypeStripedCPV) {
			return true
		}
		if count >= 2 && csp.Spec.PoolSpec.PoolType == string(apis.PoolTypeMirroredCPV) {
			return true
		}
	}
	return false
}

func getDeviceId(disk *apis.Disk) string {
	var DeviceID string
	if len(disk.Spec.DevLinks) != 0 && len(disk.Spec.DevLinks[0].Links) != 0 {
		DeviceID = disk.Spec.DevLinks[0].Links[0]
	} else {
		DeviceID = disk.Spec.Path
	}
	return DeviceID
}
