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

package v1alpha1

import (
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	apisbeta "github.com/openebs/maya/pkg/apis/openebs.io/v1beta1"
)

// CasPool encapsulates CasPool object.
type CasPool struct {
	// actual CasPool object
	Object *apis.CasPool
}

// Builder is the builder object for CasPool.
type Builder struct {
	CasPoolObject *CasPool
}

// NewBuilder returns an empty instance of the Builder object.
func NewBuilder() *Builder {
	return &Builder{
		CasPoolObject: &CasPool{&apis.CasPool{}},
	}
}

// Build returns the CasPool object built by this builder.
func (cp *Builder) Build() *CasPool {
	return cp.CasPoolObject
}

func (cb *Builder) WithCasTemplateName(casTemplateName string) *Builder {
	//casTemplateName := spc.Annotations[string(v1alpha1.CreatePoolCASTemplateKey)]
	cb.CasPoolObject.Object.CasCreateTemplate = casTemplateName
	return cb
}

func (cb *Builder) WithSpcName(name string) *Builder {
	cb.CasPoolObject.Object.StoragePoolClaim = name
	return cb
}

func (cb *Builder) WithNodeName(nodeName string) *Builder {
	cb.CasPoolObject.Object.NodeName = nodeName
	return cb
}

func (cb *Builder) WithPoolType(poolType string) *Builder {
	cb.CasPoolObject.Object.PoolType = poolType
	return cb
}

func (cb *Builder) WithMaxPool(spc *apisbeta.StoragePoolClaim) *Builder {
	cb.CasPoolObject.Object.MaxPools = *spc.Spec.MaxPools
	return cb
}

func (cb *Builder) WithDiskType(diskType string) *Builder {
	cb.CasPoolObject.Object.Type = diskType
	return cb
}

func (cb *Builder) WithAnnotations(annotations map[string]string) *Builder {
	cb.CasPoolObject.Object.Annotations = annotations
	return cb
}

func (cb *Builder) WithDiskGroup(diskGroup []apisbeta.StoragePoolClaimDiskGroups) *Builder {
	cb.CasPoolObject.Object.DiskGroups = diskGroup
	return cb
}

func (cb *Builder) WithDiskDeviceIDMap(diskDeviceIDMap map[string]string) *Builder {
	cb.CasPoolObject.Object.DiskDeviceIDMap = diskDeviceIDMap
	return cb
}
