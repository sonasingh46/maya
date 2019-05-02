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
)

// GetCasPoolForManualProvisioning returns a CasPool object for manual provisioned pool.
func (op *Operations) GetCasPoolForManualProvisioning() (*apisv1alpha1.CasPool, error) {
	err := op.Validate()
	if err != nil {
		return nil, errors.Wrapf(err, "validation failed")
	}

	casPool, err := op.BuildCasPoolForManualProvisioning()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build cas pool object for spc %s", op.SpcObject.Object.Name)
	}
	return casPool, nil
}

// BuildCasPoolForManualProvisioning builds a CasPool object for pool creation.
func (op *Operations) BuildCasPoolForManualProvisioning() (*apisv1alpha1.CasPool, error) {
	nodeNames := op.SpcObject.GetNodeNames()
	usableNode, diskGroups, err := op.GetUsableNode(nodeNames)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get usable node for spc %s", op.SpcObject.Object.Name)
	}

	diskDeviceIDMap, err := op.GetDiskDeviceIDMap()
	if err != nil {
		return nil, errors.Wrapf(err, "could not form disk device ID map for %s", op.SpcObject.Object.Name)
	}

	spcObject := op.SpcObject.Object
	cp := caspool.NewBuilder().
		WithDiskType(spcObject.Spec.Type).
		WithPoolType(op.SpcObject.GetPoolType(usableNode)).
		WithAnnotations(op.SpcObject.GetAnnotations()).
		WithDiskGroup(diskGroups).
		WithCasTemplateName(op.SpcObject.GetCastName()).
		WithSpcName(op.SpcObject.Object.Name).
		WithNodeName(usableNode).
		WithDiskDeviceIDMap(diskDeviceIDMap).
		Build().Object
	return cp, nil
}

// GetUsableNode returns a node and disks attached to it where pool can be possibly provisioned.
func (op *Operations) GetUsableNode(nodes []string) (string, []apisv1beta1.StoragePoolClaimDiskGroups, error) {
	usedNodes, err := op.GetUsedNode()
	if err != nil {
		return "", []apisv1beta1.StoragePoolClaimDiskGroups{}, err
	}
	for _, node := range nodes {
		if !usedNodes[node] {
			return node, op.SpcObject.GetNodeDisk(node), nil
		}
	}
	return "", []apisv1beta1.StoragePoolClaimDiskGroups{}, errors.Errorf("no usable node found for spc %s", op.SpcObject.Object.Name)
}

// TODO: Find some better way to do validations
// Validate does validations for disk present in the spc.
func (op *Operations) Validate() error {
	for _, node := range op.SpcObject.Object.Spec.Nodes {
		for _, groups := range node.DiskGroups {

			if !disk.IsValidPoolTopology(node.PoolSpec.PoolType, len(groups.Disks)) {
				return errors.New("disk count is invalid for specified pool type")
			}

			for _, diskDetails := range groups.Disks {

				diskAPIObject, err := op.DiskClient.Get(diskDetails.Name, metav1.GetOptions{})
				if err != nil {
					return errors.Wrap(err, "could not get disk object")
				}

				diskObject := disk.BuilderForAPIObject(diskAPIObject).Disk

				if !disk.IsType(op.SpcObject.Object.Spec.Type)(diskObject) {
					return errors.New("disk type does not match type mentioned in spc")
				}

				if !disk.IsDiskActive()(diskObject) {
					return errors.New("disk is inactive")
				}

				if !disk.IsDiskBelongToNode(node.Name)(diskObject) {
					return errors.New("disk hostname is different then the one in which it is put in spc")
				}
			}
		}
	}
	return nil
}
