// +build integration

/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package integration

import (
	"strings"
	"testing"

	"github.com/docker/machine/libmachine/state"
	"k8s.io/minikube/test/integration/util"
)

func TestFunctional(t *testing.T) {
	p := "minikube"
	mk := NewMinikubeRunner(t, p)
	mk.EnsureRunning()
	// This one is not parallel, and ensures the cluster comes up
	// before we run any other tests.
	t.Run("Status", testClusterStatus)
	t.Run("ProfileList", testProfileList)
	t.Run("DNS", testClusterDNS)
	t.Run("Logs", testClusterLogs)
	t.Run("Addons", testAddons)
	t.Run("Registry", testRegistry)
	t.Run("Dashboard", testDashboard)
	t.Run("ServicesList", testServicesList)
	t.Run("Provisioning", testProvisioning)
	t.Run("Tunnel", testTunnel)

	if !usingNoneDriver(mk) {
		t.Run("EnvVars", testClusterEnv)
		t.Run("SSH", testClusterSSH)
		t.Run("IngressController", testIngressController)
		t.Run("Mounting", testMounting)
	}
}

func TestFunctionalContainerd(t *testing.T) {
	p := "minikube"
	mk := NewMinikubeRunner(t, p)

	if usingNoneDriver(mk) {
		t.Skip("Can't run containerd backend with none driver")
	}

	if mk.GetStatus() != state.None.String() {
		mk.RunCommand("delete", true)
	}
	mk.Start("--container-runtime=containerd", "--docker-opt containerd=/var/run/containerd/containerd.sock")
	t.Run("Gvisor", testGvisor)
	t.Run("GvisorRestart", testGvisorRestart)
	mk.RunCommand("delete", true)
}

// usingNoneDriver returns true if using the none driver
func usingNoneDriver(r util.MinikubeRunner) bool {
	return strings.Contains(r.StartArgs, "--vm-driver=none")
}
