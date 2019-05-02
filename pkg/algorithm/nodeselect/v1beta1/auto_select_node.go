/*
Copyright 2018 The OpenEBS Authors

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

package v1beta1

import (
	apisv1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	apisv1beta1 "github.com/openebs/maya/pkg/apis/openebs.io/v1beta1"
	caspool "github.com/openebs/maya/pkg/caspool/v1alpha1"
	disk "github.com/openebs/maya/pkg/disk/v1alpha2"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

const (
	DefaultDiskCountMultiplier = 1
)

type nodeFilterPredicate func(nodeName string) bool
type nodeFilterPredicateList []nodeFilterPredicate

func (op *Operations) GetCasPoolForAutoProvisioning() (*apisv1alpha1.CasPool, error) {
	nodeName, diskList, err := op.SelectNode()
	diskGroups := op.GetDiskGroups(diskList)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get usable node for spc %s", op.SpcObject.Object.Name)
	}
	diskDeviceIDMap, err := op.GetDiskDeviceIDMapForDiskAPIList(diskList)
	if err != nil {
		return nil, errors.Wrapf(err, "could not form disk device ID map for %s", op.SpcObject.Object.Name)
	}

	spcObject := op.SpcObject.Object
	cp := caspool.NewBuilder().
		WithDiskType(spcObject.Spec.Type).
		WithPoolType(op.SpcObject.Object.Spec.PoolSpec.PoolType).
		WithAnnotations(op.SpcObject.GetAnnotations()).
		WithDiskGroup(diskGroups).
		WithCasTemplateName(op.SpcObject.GetCastName()).
		WithSpcName(op.SpcObject.Object.Name).
		WithNodeName(nodeName).
		WithDiskDeviceIDMap(diskDeviceIDMap).
		Build().Object
	return cp, nil
}

func (op *Operations) GetDisks() (*apisv1alpha1.DiskList, error) {
	diskAPIList, err := op.DiskClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "could not list disk")
	}

	usedDiskMap, err := op.GetUsedDiskMap()
	if err != nil {
		return nil, errors.Wrap(err, "could not get used disk map")
	}

	usedNodeMap, err := op.GetUsedNode()
	if err != nil {
		return nil, errors.Wrap(err, "could not get used node map")
	}

	newDiskAPIList := disk.ListBuilderForAPIObject(diskAPIList).DiskList.
		Filter(disk.IsDiskActive(), disk.IsType(op.SpcObject.Object.Spec.Type), disk.IsUsableDisk(usedDiskMap), disk.IsUsableNodeDisk(usedNodeMap)).
		ObjectList

	if len(newDiskAPIList.Items) == 0 {
		return nil, errors.New("no disks found")
	}

	return newDiskAPIList, nil
}

func (op *Operations) GetNodeDiskMap() (map[string]*apisv1alpha1.DiskList, error) {
	diskAPIList, err := op.GetDisks()

	if err != nil {
		return nil, errors.Wrap(err, "could not get usable disk")
	}
	nodeDiskMap := make(map[string]*apisv1alpha1.DiskList)
	for _, diskAPIObject := range diskAPIList.Items {
		// pin it
		diskAPIObject := diskAPIObject
		nodeName := disk.BuilderForAPIObject(&diskAPIObject).Disk.GetNodeName()
		if nodeDiskMap[nodeName] == nil {
			nodeDiskMap[nodeName] = &apisv1alpha1.DiskList{Items: []apisv1alpha1.Disk{diskAPIObject}}
		} else {
			nodeDiskMap[nodeName].Items = append(nodeDiskMap[nodeName].Items, diskAPIObject)
		}
	}
	filteredNodediskMap := FilterNodeDiskMap(nodeDiskMap, IsTopologyObeyed(nodeDiskMap, op.SpcObject.Object.Spec.PoolSpec.PoolType))

	if len(filteredNodediskMap) == 0 {
		return nil, errors.New("no node found with required number of disks")
	}

	return filteredNodediskMap, nil
}

func (op *Operations) SelectNode() (string, *apisv1alpha1.DiskList, error) {
	nodeDiskMap, err := op.GetNodeDiskMap()
	if err != nil {
		return "", nil, errors.Wrapf(err, "could not get node disk map")
	}
	for nodeName, diskList := range nodeDiskMap {
		return nodeName, diskList, nil
	}
	return "", nil, errors.Wrapf(err, "got empty node disk map")
}

// all returns true if all the predicates succeed against the provided disk instance.
func (l nodeFilterPredicateList) allFilterPredicates(nodeName string) bool {
	for _, pred := range l {
		if !pred(nodeName) {
			return false
		}
	}
	return true
}

// Filter will filter the disk instances if all the predicates succeed against that spc.
func FilterNodeDiskMap(nodeDiskMap map[string]*apisv1alpha1.DiskList, p ...nodeFilterPredicate) map[string]*apisv1alpha1.DiskList {
	var plist nodeFilterPredicateList
	plist = append(plist, p...)
	if len(plist) == 0 {
		return nodeDiskMap
	}

	filtered := make(map[string]*apisv1alpha1.DiskList)
	for node, diskAPIList := range nodeDiskMap {
		diskAPIList := diskAPIList // pin it
		if plist.allFilterPredicates(node) {
			filtered[node] = diskAPIList
		}
	}
	return filtered
}

func IsTopologyObeyed(nodeDiskMap map[string]*apisv1alpha1.DiskList, poolType string) nodeFilterPredicate {
	return func(nodeName string) bool {
		requiredDiskCount := GetRequiredDiskCount(poolType)
		if len(nodeDiskMap[nodeName].Items) >= requiredDiskCount {
			return true
		}
		return false
	}
}

func (op *Operations) GetDiskGroups(diskList *apisv1alpha1.DiskList) []apisv1beta1.StoragePoolClaimDiskGroups {
	var SPCDiskGroup []apisv1beta1.StoragePoolClaimDiskGroups
	var newDiskList []apisv1beta1.StoragePoolClaimDisk
	groupIndex := 0
	for i := 1; i <= GetRequiredDiskCount(op.SpcObject.Object.Spec.PoolSpec.PoolType); i++ {
		newSpcDisk := apisv1beta1.StoragePoolClaimDisk{
			Name: diskList.Items[i].Name,
			ID:   "",
		}
		newDiskList = append(newDiskList, newSpcDisk)
		if i%DefaultDiskCountMultiplier == 0 {
			newDiskGroup := apisv1beta1.StoragePoolClaimDiskGroups{
				Name:  "group" + strconv.Itoa(groupIndex),
				Disks: newDiskList,
			}
			SPCDiskGroup = append(SPCDiskGroup, newDiskGroup)
			groupIndex++
		}
	}
	return SPCDiskGroup
}

func GetRequiredDiskCount(poolType string) int {
	return disk.DefaultDiskCount[poolType] * DefaultDiskCountMultiplier
}
