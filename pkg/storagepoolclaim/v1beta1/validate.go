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

package v1beta1

import (
	apisv1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/pkg/errors"
)

var (
	// supportedPool is a map holding the supported raid configurations.
	supportedPool = map[string]bool{
		string(apisv1alpha1.PoolTypeStripedCPV):  true,
		string(apisv1alpha1.PoolTypeMirroredCPV): true,
		string(apisv1alpha1.PoolTypeRaidzCPV):    true,
		string(apisv1alpha1.PoolTypeRaidz2CPV):   true,
	}
)

// ValidateFunc is typed function for spc validation functions.
type ValidateFunc func(*SPC) error

// ValidateFuncList holds a list of validate functions for spc
var ValidateFuncList = []ValidateFunc{
	ValidatePoolType,
	ValidateDiskType,
	ValidateAutoSpcMaxPool,
}

func Validate(spc *SPC) error {
	for _, v := range ValidateFuncList {
		err := v(spc)
		//TODO: Add context to error
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateDiskType validates the disk types in spc.
func ValidateDiskType(spc *SPC) error {
	diskType := spc.Object.Spec.Type
	if !(diskType == "sparse" || diskType == "disk") {
		return errors.Errorf("specified type on spc %s is %s which is invalid", spc.Object.Name, diskType)
	}
	return nil
}

// ValidateAutoSpcMaxPool validates the max pool count in auto spc
func ValidateAutoSpcMaxPool(spc *SPC) error {
	spcName := spc.Object.Name
	if IsAutoProvisioning()(spc) {
		maxPools := spc.Object.Spec.MaxPools
		if maxPools == nil {
			return errors.Errorf("maxpool value is nil for spc %s which is invalid", spcName)
		}
		if *maxPools < 0 {
			return errors.Errorf("maxpool value is %v for spc %s which is invalid", maxPools, spcName)
		}
	}
	return nil
}

func ValidatePoolType(spc *SPC) error {
	for _, node := range spc.Object.Spec.Nodes {
		if !supportedPool[node.PoolSpec.PoolType] {
			return errors.Errorf("pool type is %s for node %s in spc %s", node.PoolSpec.PoolType, node.Name, spc.Object.Name)
		}
	}
	return nil
}
