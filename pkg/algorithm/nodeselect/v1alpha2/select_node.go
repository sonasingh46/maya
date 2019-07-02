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

package v1alpha2

import (
	"github.com/golang/glog"
	ndmapis "github.com/openebs/maya/pkg/apis/openebs.io/ndm/v1alpha1"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	bd "github.com/openebs/maya/pkg/blockdevice/v1alpha2"
	bdc "github.com/openebs/maya/pkg/blockdeviceclaim/v1alpha1"
	csp "github.com/openebs/maya/pkg/cstor/newpool/v1alpha3"
	nodeapis "github.com/openebs/maya/pkg/kubernetes/node/v1alpha1"
	"github.com/openebs/maya/pkg/volume"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//// GetCandidateNodeMap returns a map of all nodes where the pool needs to be created.
//func (ac *Config) GetCandidateNodeMap() (map[string]bool, error) {
//	candidateNodesMap := make(map[string]bool)
//	usedNodeMap, err := ac.GetUsedNodeMap()
//	if err != nil {
//		return nil, errors.Wrapf(err, "could not get candidate nodes for pool creation")
//	}
//	for _, pool := range ac.CSPC.Spec.Pools {
//		nodeName, err := ac.GetNodeFromLabelSelector(pool.NodeSelector)
//		if err != nil {
//			glog.Errorf("could not use node for selectors {%v}", pool.NodeSelector)
//			continue
//		}
//		if usedNodeMap[nodeName] == false {
//			candidateNodesMap[nodeName] = true
//		}
//	}
//	return candidateNodesMap, nil
//}

// SelectNode returns a node where pool should be created.
func (ac *Config) SelectNode() (*apis.PoolSpec, string, error) {
	usedNodes, err := ac.GetUsedNode()
	if err != nil {
		return nil, "", errors.Wrapf(err, "could not get used nodes list for pool creation")
	}
	for _, pool := range ac.CSPC.Spec.Pools {
		// pin it
		pool := pool
		nodeName, err := ac.GetNodeFromLabelSelector(pool.NodeSelector)
		if err != nil {
			glog.Errorf("could not use node for selectors {%v}", pool.NodeSelector)
			continue
		}
		if !usedNodes[nodeName] {
			return &pool, nodeName, nil
		}
	}
	return nil, "", errors.New("no node qualified for pool creation")
}

// GetNodeFromLabelSelector returns the node name selected by provided labels
func (ac *Config) GetNodeFromLabelSelector(labels map[string]string) (string, error) {
	nodeList, err := nodeapis.NewKubeClient().List(metav1.ListOptions{LabelSelector: getLabelSelectorString(labels)})
	if err != nil {
		return "", errors.Wrap(err, "failed to get node list from the node selector")
	}
	if len(nodeList.Items) != 1 {
		return "", errors.Errorf("could not get a unique node from the given node selectors")
	}
	if !nodeapis.NewBuilder().WithAPINode(&nodeList.Items[0]).Node.IsReady() {
		return "", errors.Errorf("node {%s} is not ready", nodeList.Items[0].Name)
	}
	return nodeList.Items[0].Name, nil
}

// getLabelSelectorString returns a string of label selector form label map to be used in
// list options.
func getLabelSelectorString(selector map[string]string) string {
	var selectorString string
	for key, value := range selector {
		selectorString = selectorString + key + "=" + value + ","
	}
	selectorString = selectorString[:len(selectorString)-len(",")]
	return selectorString
}

// GetUsedNode returns a map of node for which pool has already been created.
func (ac *Config) GetUsedNode() (map[string]bool, error) {
	usedNode := make(map[string]bool)
	cspList, err := csp.
		NewKubeClient().
		WithNamespace(ac.Namespace).
		List(
			metav1.
				ListOptions{LabelSelector: string(apis.CStorPoolClusterCPK) + "=" + ac.CSPC.Name},
		)
	if err != nil {
		return nil, errors.Wrap(err, "could not list already created csp(s)")
	}
	for _, cspObj := range cspList.Items {
		usedNode[cspObj.Labels[string(apis.HostNameCPK)]] = true
	}
	return usedNode, nil
}

// GetBDListForNode returns a list of BD from the pool spec.
// TODO : Move it to CStorPoolCluster packgage
func (ac *Config) GetBDListForNode(pool *apis.PoolSpec) []string {
	var BDList []string
	for _, group := range pool.RaidGroups {
		for _, bd := range group.BlockDevices {
			BDList = append(BDList, bd.BlockDeviceName)
		}
	}
	return BDList
}

// ClaimBDsForNode claims a given BlockDevice for node
// If the block device(s) is/are already claimed for any other CSPC it returns error.
// If the block device(s) is/are already calimed for the same CSPC -- it is left as it is and can be used for
// pool provisioning.
// If the block device(s) is/are unclaimed, then those are claimed.
func (ac *Config) ClaimBDsForNode(BD []string) error {
	for _, bdName := range BD {
		bdAPIObj, err := bd.NewKubeClient().WithNamespace(ac.Namespace).Get(bdName, metav1.GetOptions{})
		if err != nil {
			return errors.Wrapf(err, "error in getting details for BD {%s} whether it is claimed", bdName)
		}
		if bd.BuilderForAPIObject(bdAPIObj).BlockDevice.IsClaimed() {
			IsClaimedBDUsable, err := ac.IsClaimedBDUsable(bdAPIObj)
			if err != nil {
				return errors.Wrapf(err, "error in getting details for BD {%s} for usability", bdName)
			}
			if !IsClaimedBDUsable {
				return errors.Errorf("BD {%s} already in use", bdName)
			}
		}

		err = ac.ClaimBD(bdAPIObj)
		if err != nil {
			return errors.Wrapf(err, "Failed to claim BD {%s}", bdName)
		}
	}
	return nil
}

// ClaimBD claims a given BlockDevice
func (ac *Config) ClaimBD(bdObj *ndmapis.BlockDevice) error {
	newBDCObj, err := bdc.NewBuilder().
		WithName("bdc-" + string(bdObj.UID)).
		WithNamespace(ac.Namespace).
		WithLabels(map[string]string{string(apis.CStorPoolClusterCPK): ac.CSPC.Name}).
		WithBlockDeviceName(bdObj.Name).
		WithHostName(bdObj.Labels[string(apis.HostNameCPK)]).
		WithCapacity(volume.ByteCount(bdObj.Spec.Capacity.Storage)).
		WithCSPCOwnerReference(ac.CSPC).
		Build()

	if err != nil {
		return errors.Wrapf(err, "failed to build block device claim for bd {%s}", bdObj.Name)
	}

	_, err = bdc.NewKubeClient().WithNamespace(ac.Namespace).Create(newBDCObj.Object)
	if err != nil {
		return errors.Wrapf(err, "failed to create block device claim for bd {%s}", bdObj.Name)
	}
	return nil
}

// IsClaimedBDUsable returns true if the passed BD is already claimed and can be
// used for provisioning
func (ac *Config) IsClaimedBDUsable(bdAPIObj *ndmapis.BlockDevice) (bool, error) {
	bdObj := bd.BuilderForAPIObject(bdAPIObj)
	if bdObj.BlockDevice.IsClaimed() {
		bdcName := bdObj.BlockDevice.Object.Spec.ClaimRef.Name
		bdcAPIObject, err := bdc.NewKubeClient().WithNamespace(ac.Namespace).Get(bdcName, metav1.GetOptions{})
		if err != nil {
			return false, errors.Wrapf(err, "could not get block device claim for block device {%s}", bdAPIObj.Name)
		}
		bdcObj := bdc.BuilderForAPIObject(bdcAPIObject)
		if bdcObj.BDC.HasLabel(string(apis.CStorPoolClusterCPK), ac.CSPC.Name) {
			return true, nil
		}
	} else {
		return false, errors.Errorf("block device {%s} is not claimed", bdAPIObj.Name)
	}
	return false, nil
}

// ValidatePoolSpec validates the pool spec.
// TODO: Fix following function -- (Current is mock only )
func ValidatePoolSpec(pool *apis.PoolSpec) bool {
	return true
}
