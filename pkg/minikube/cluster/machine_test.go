/*
Copyright 2019 The Kubernetes Authors All rights reserved.

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

package cluster

import (
	"io/ioutil"
	"k8s.io/minikube/pkg/minikube/constants"
	"os"
	"path/filepath"
	"testing"
)

func TestListMachines(t *testing.T) {
	const (
		numberOfValidMachines   = 2
		numberOfInValidMachines = 3
		totalNumberOfMachines   = numberOfValidMachines + numberOfInValidMachines
	)

	testMinikubeDir := "./testdata/machine"
	miniDir, err := filepath.Abs(testMinikubeDir)

	if err != nil {
		t.Errorf("error getting dir path for %s : %v", testMinikubeDir, err)
	}

	err = os.Setenv(constants.MinikubeHome, miniDir)
	if err != nil {
		t.Errorf("error setting up test environment. could not set %s", constants.MinikubeHome)
	}

	files, _ := ioutil.ReadDir(filepath.Join(constants.GetMinipath(), "machines"))
	numberOfMachineDirs := len(files)

	validMachines, inValidMachines, err := ListMachines()

	if err != nil {
		t.Error(err)
	}

	if numberOfValidMachines != len(validMachines) {
		t.Errorf("expected %d valid machines, got %d", numberOfValidMachines, len(validMachines))
	}

	if numberOfInValidMachines != len(inValidMachines) {
		t.Errorf("expected %d invalid machines, got %d", numberOfInValidMachines, len(inValidMachines))
	}

	if totalNumberOfMachines != len(validMachines)+len(inValidMachines) {
		t.Errorf("expected %d total machines, got %d", totalNumberOfMachines, len(validMachines)+len(inValidMachines))
	}

	if numberOfMachineDirs != len(validMachines)+len(inValidMachines) {
		t.Error("expected number of machine directories to be equal to the number of total machines")
	}
}
