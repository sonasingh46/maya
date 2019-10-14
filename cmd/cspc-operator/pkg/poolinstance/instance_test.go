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

package poolinstance

import (
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

type injection struct {
	InjectPoolInstanceSpecGettingFailure bool
	InjectPoolInstanceCreationFailure    bool
}

var inject = injection{
	InjectPoolInstanceSpecGettingFailure: false,
	InjectPoolInstanceCreationFailure:    false,
}

type fakePoolInstanceController struct {
	fakeCSPC *apis.CStorPoolCluster
}
type fakeAlgorithmConfig struct {
}

func NewFakePoolInstanceController(CSPC *apis.CStorPoolCluster) *fakePoolInstanceController {
	return &fakePoolInstanceController{fakeCSPC: CSPC}
}
func (c *fakePoolInstanceController) Get(name string, opts metav1.GetOptions) (*apis.CStorPoolInstance, error) {
	return nil, nil
}

func (c *fakePoolInstanceController) List(opts metav1.ListOptions) (*apis.CStorPoolInstanceList, error) {
	return nil, nil
}

func (c *fakePoolInstanceController) Create(CSPI *apis.CStorPoolInstance) (*apis.CStorPoolInstance, error) {
	if inject.InjectPoolInstanceCreationFailure == true {
		return nil, errors.New("failed to get create CSPI")
	} else {
		return &apis.CStorPoolInstance{}, nil
	}
}

func (c *fakePoolInstanceController) Delete(name string, opts *metav1.DeleteOptions) error {
	return nil
}
func (fac *fakeAlgorithmConfig) GetPendingPoolCount() (int, error) {
	return 0, nil
}
func (c *fakeAlgorithmConfig) GetPoolInstanceSpec() (*apis.CStorPoolInstance, error) {
	if inject.InjectPoolInstanceSpecGettingFailure == true {
		return nil, errors.New("failed to get pool instance spec")
	} else {
		return &apis.CStorPoolInstance{}, nil
	}
}
func TestConfig_Provision(t *testing.T) {
	newFC := &Config{
		PoolInstance:    NewFakePoolInstanceController(&apis.CStorPoolCluster{}),
		AlgorithmConfig: &fakeAlgorithmConfig{},
	}
	tests := map[string]struct {
		inject        injection
		expectedCSPI  *apis.CStorPoolInstance
		expectedError error
		errorOccurs   bool
	}{
		"#1:Success": {
			expectedCSPI:  &apis.CStorPoolInstance{},
			expectedError: nil,
			errorOccurs:   false,
		},
		"#2:Failure in getting CSPI Spec": {
			inject: injection{
				InjectPoolInstanceSpecGettingFailure: true,
				InjectPoolInstanceCreationFailure:    false,
			},
			expectedCSPI:  nil,
			expectedError: errors.New("failed to get pool instance spec"),
			errorOccurs:   true,
		},
		"#3:Failure in creating CSPI Spec": {
			inject: injection{
				InjectPoolInstanceSpecGettingFailure: false,
				InjectPoolInstanceCreationFailure:    true,
			},
			expectedCSPI:  nil,
			expectedError: errors.New("failed to get create CSPI"),
			errorOccurs:   true,
		},
		"#4:Failure in both": {
			inject: injection{
				InjectPoolInstanceSpecGettingFailure: true,
				InjectPoolInstanceCreationFailure:    true,
			},
			expectedCSPI:  nil,
			expectedError: errors.New("failed to get pool instance spec"),
			errorOccurs:   true,
		},
	}

	for name, test := range tests {
		test := test //pin it
		t.Run(name, func(t *testing.T) {
			inject = test.inject
			gotCSPI, gotErr := newFC.Provision()
			if !reflect.DeepEqual(gotCSPI, test.expectedCSPI) {
				t.Errorf("Expected CSPI : {%#v } but got: {%#v}", test.expectedCSPI, gotCSPI)
			}

			if test.errorOccurs {
				if !(test.expectedError.Error() == gotErr.Error()) {
					t.Errorf("Expected error: {%s} but got error: {%s}", test.expectedError, gotErr)
				}
			} else {
				if gotErr != nil {
					t.Errorf("Expected not error but got error: {%s}", gotErr.Error())
				}
			}
		})
	}
}
