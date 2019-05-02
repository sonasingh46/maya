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

package v1alpha2

import (
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
)

const (
	DiskActiveState   = "Active"
	DiskInactiveState = "Inactive"
)

// DefaultDiskCount is a map containing the default disk count of various raid types.
var DefaultDiskCount = map[string]int{
	string(apis.PoolTypeMirroredCPV): int(apis.MirroredDiskCountCPV),
	string(apis.PoolTypeStripedCPV):  int(apis.StripedDiskCountCPV),
	string(apis.PoolTypeRaidzCPV):    int(apis.RaidzDiskCountCPV),
	string(apis.PoolTypeRaidz2CPV):   int(apis.Raidz2DiskCountCPV),
}

// Disk encapsulates StoragePoolClaim api object.
type Disk struct {
	// actual spc object
	Object *apis.Disk
}

// DiskList holds the list of StoragePoolClaim api
type DiskList struct {
	// list of storagepoolclaims
	ObjectList *apis.DiskList
}

// Builder is the builder object for Disk.
type Builder struct {
	Disk *Disk
}

// ListBuilder is the builder object for DiskList.
type ListBuilder struct {
	DiskList *DiskList
}

// Predicate defines an abstraction to determine conditional checks against the provided spc instance.
type Predicate func(*Disk) bool

type predicateList []Predicate

// all returns true if all the predicates succeed against the provided disk instance.
func (l predicateList) all(c *Disk) bool {
	for _, pred := range l {
		if !pred(c) {
			return false
		}
	}
	return true
}

// HasAnnotation returns true if provided annotation key and value are present in the provided spc instance.
func HasAnnotation(key, value string) Predicate {
	return func(c *Disk) bool {
		val, ok := c.Object.GetAnnotations()[key]
		if ok {
			return val == value
		}
		return false
	}
}

// IsSparseType returns true if the spc is of sparse type.
func IsSparseType() Predicate {
	return func(d *Disk) bool {
		return d.Object.GetLabels()[string(apis.NdmDiskTypeCPK)] == string(apis.TypeSparseCPV)
	}
}

// IsSparseType returns true if the spc is of sparse type.
func IsDiskActive() Predicate {
	return func(d *Disk) bool {
		return d.Object.Status.State == DiskActiveState
	}
}

func IsUsableDisk(usedDiskMap map[string]int) Predicate {
	return func(d *Disk) bool {
		return usedDiskMap[d.Object.Name] == 0
	}
}

func IsUsableNodeDisk(usedNodeMap map[string]bool) Predicate {
	return func(d *Disk) bool {
		return !usedNodeMap[d.GetNodeName()]
	}
}

// IsDiskBelongToNode returns true if the spc is of sparse type.
func IsDiskBelongToNode(nodeName string) Predicate {
	return func(d *Disk) bool {
		return d.GetNodeName() == nodeName
	}
}

// IsType returns true if the disk is of type same as passed argument
func IsType(diskType string) Predicate {
	return func(d *Disk) bool {
		return d.Object.GetLabels()[string(apis.NdmDiskTypeCPK)] == diskType
	}
}

// IsValidPoolTopology returns true if the topology is valid.
func IsValidPoolTopology(poolType string, diskCount int) bool {
	return DefaultDiskCount[poolType]%diskCount == 0
}

// IsDiskType returns true if the spc is of disk type.
func IsDiskType() Predicate {
	return func(d *Disk) bool {
		return d.Object.GetLabels()[string(apis.NdmDiskTypeCPK)] == string(apis.TypeDiskCPV)
	}
}

// GetNodeName returns the node name to which the disk is attached
func (d *Disk) GetNodeName() string {
	return d.Object.GetLabels()[string(apis.HostNameCPK)]
}

// Filter will filter the disk instances if all the predicates succeed against that spc.
func (l *DiskList) Filter(p ...Predicate) *DiskList {
	var plist predicateList
	plist = append(plist, p...)
	if len(plist) == 0 {
		return l
	}

	filtered := NewListBuilder().List()
	for _, spcAPI := range l.ObjectList.Items {
		spcAPI := spcAPI // pin it
		Disk := BuilderForAPIObject(&spcAPI).Disk
		if plist.all(Disk) {
			filtered.ObjectList.Items = append(filtered.ObjectList.Items, *Disk.Object)
		}
	}
	return filtered
}

// NewBuilder returns an empty instance of the Builder object.
func NewBuilder() *Builder {
	return &Builder{
		Disk: &Disk{&apis.Disk{}},
	}
}

// BuilderForObject returns an instance of the Builder object based on spc object
func BuilderForObject(Disk *Disk) *Builder {
	return &Builder{
		Disk: Disk,
	}
}

// BuilderForAPIObject returns an instance of the Builder object based on spc api object.
func BuilderForAPIObject(spc *apis.Disk) *Builder {
	return &Builder{
		Disk: &Disk{spc},
	}
}

func (d *Disk) GetDeviceID() string {
	var deviceID string
	if len(d.Object.Spec.DevLinks) != 0 && len(d.Object.Spec.DevLinks[0].Links) != 0 {
		deviceID = d.Object.Spec.DevLinks[0].Links[0]
	} else {
		deviceID = d.Object.Spec.Path
	}
	return deviceID
}

// Build returns the Disk object built by this builder.
func (sb *Builder) Build() *Disk {
	return sb.Disk
}

// NewListBuilder returns a new instance of ListBuilder object.
func NewListBuilder() *ListBuilder {
	return &ListBuilder{DiskList: &DiskList{ObjectList: &apis.DiskList{}}}
}

// WithList builds the list based on the provided *DiskList instances.
func (b *ListBuilder) WithList(pools *DiskList) *ListBuilder {
	if pools == nil {
		return b
	}
	b.DiskList.ObjectList.Items = append(b.DiskList.ObjectList.Items, pools.ObjectList.Items...)
	return b
}

// WithAPIList builds the list based on the provided *apis.CStorPoolList.
func (b *ListBuilder) WithAPIList(pools *apis.DiskList) *ListBuilder {
	if pools == nil {
		return b
	}
	for _, pool := range pools.Items {
		pool := pool //pin it
		b.DiskList.ObjectList.Items = append(b.DiskList.ObjectList.Items, pool)
	}
	return b
}

// List returns the list of disk instances that were built by this builder.
func (b *ListBuilder) List() *DiskList {
	return b.DiskList
}

// ListBuilderForAPIObject returns a new instance of ListBuilderForApiList object based on csp api list.
func ListBuilderForAPIObject(diskAPIList *apis.DiskList) *ListBuilder {
	newLb := NewListBuilder()
	for _, obj := range diskAPIList.Items {
		// pin it
		obj := obj
		newLb.DiskList.ObjectList.Items = append(newLb.DiskList.ObjectList.Items, obj)
	}
	return newLb
}

// Len returns the length og DiskList.
func (l *DiskList) Len() int {
	return len(l.ObjectList.Items)
}
