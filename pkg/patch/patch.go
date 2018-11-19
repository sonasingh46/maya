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

package patch

import (
	"encoding/json"
	"fmt"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/maya/pkg/client/k8s"
	"k8s.io/apimachinery/pkg/types"
)

// PatchPayloadCSP struct is ussed to patch CSP object.
// Similarly, for other objects (if required to patch) we can have structs for them
// to have a implementation if patch function.
type PatchPayloadCSP struct {
	// 'Object' is the object which needs to be patched.
	Object *apis.CStorPool
	// PatchPayloadCSP is the payload to patch CSP.
	PatchPayloadCSP []Patch
}

// Patch is the struct based on standards of JSON patch.
type Patch struct {
	// Op defines the operation
	Op string `json:"op"`
	// Path defines the key path
	// eg. for
	// {
	//  	"Name": "openebs"
	//	    Category: {
	//		  "Inclusive": "v1",
	//		  "Rank": "A"
	//	     }
	// }
	// The path of 'Inclusive' would be
	// "/Name/Category/Inclusive"
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// Patcher interface has Patch functions which can be implemented for several objects that needs to be patched.
type Patcher interface {
	Patch(string, types.PatchType) (interface{}, error)
}

// NewPatchPayload constructs the patch payload fo any type of object.
func NewPatchPayload(operation string, path string, value interface{}) (payload []Patch) {
	PatchPayload := make([]Patch, 1)
	PatchPayload[0].Op = operation
	PatchPayload[0].Path = path
	PatchPayload[0].Value = value
	return PatchPayload
}

// Patch is the specific implementation if Patch() interface for patching CSP objects.
// Similarly, we can have for other objects, if required.
func (payload *PatchPayloadCSP) Patch(namesapce string, patchType types.PatchType) (interface{}, error) {
	newK8sClient, err := k8s.NewK8sClient(namesapce)
	if err != nil {
		return nil, fmt.Errorf("Unable to get clientset for patch operation:%v", err)
	}
	PatchJSON, err := json.Marshal(payload.PatchPayloadCSP)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal patch payload for csp :%v", err)
	}
	cspObject, err := newK8sClient.PatchOEV1alpha1CSPAsRaw(payload.Object.Name, patchType, PatchJSON)
	return cspObject, err
}
