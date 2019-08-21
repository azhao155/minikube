// +build integration

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

// a_download_only_test.go filename starts with a, for the purpose that it runs before all parallel tests and downloads the images and caches them.
package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"k8s.io/minikube/pkg/minikube/constants"
	pkgutil "k8s.io/minikube/pkg/util"
	"k8s.io/minikube/pkg/util/retry"
)

// Note this test runs before all because filename is alphabetically first
// is used to cache images and binaries used by other parallel tests to avoid redownloading.
// TestDownloadOnly tests the --download-only option
func TestDownloadOnly(t *testing.T) {
	p := profileName(t)
	mk := NewMinikubeRunner(t, p)
	if !isTestNoneDriver(t) { // none driver doesnt need to be deleted
		defer mk.TearDown(t)
	}
	minHome := constants.GetMinipath()

	t.Run("Oldest", func(t *testing.T) {
		mk.StartWithFail("--download-only", fmt.Sprintf("--kubernetes-version=%s", constants.OldestKubernetesVersion))
	})

	t.Run("Newest", func(t *testing.T) {
		v := constants.NewestKubernetesVersion
		mk.StartWithFail("--download-only", fmt.Sprintf("--kubernetes-version=%s", v))

		// checking binaries downloaded
		_, imgs := constants.GetKubeadmCachedImages("", v)
		for _, img := range imgs {
			_, err := os.Stat(filepath.Join(minHome, fmt.Sprintf("images/%s", img)))
			if err != nil {
				t.Errorf("error expected download-only to cachne image %q but got error %v", img, err)
			}
		}

		// checking binaries downloaded (kubelet,kubeadm)
		for _, bin := range constants.GetKubeadmCachedBinaries() {
			_, err := os.Stat(filepath.Join(minHome, fmt.Sprintf("cache/%s/%s", v, bin)))
			if err != nil {
				t.Errorf("error expected download-only to cachne binary %q but got error %v", bin, err)
			}
		}

		// checking binaries downloaded
		for _, bin := range []string{"kublet,kbueadm"} {
			_, err := os.Stat(filepath.Join(minHome, fmt.Sprintf("cache/%s/%s", v, bin)))
			if err != nil {
				t.Errorf("error expected download-only to cachne binary %q but got error %v", bin, err)
			}
		}

	})

	// this downloads the latest published binary from where we publish the minikube binary
	t.Run("DownloadLatestRelease", func(t *testing.T) {
		dest := filepath.Join(*testdataDir, fmt.Sprintf("minikube-%s-%s-latest-stable", runtime.GOOS, runtime.GOARCH))
		err := downloadMinikubeBinary(t, dest, "latest")
		if err != nil {
			t.Errorf("erorr downloading the latest minikube release %v", err)
		}
	})
}

// downloadMinikubeBinary downloads the minikube binary from github used by TestVersionUpgrade
// acts as a test setup for TestVersionUpgrade
func downloadMinikubeBinary(t *testing.T, dest string, version string) error {
	t.Helper()
	// Grab latest release binary
	url := pkgutil.GetBinaryDownloadURL(version, runtime.GOOS)
	download := func() error {
		return getter.GetFile(dest, url)
	}

	if err := retry.Expo(download, 3*time.Second, 3*time.Minute); err != nil {
		return errors.Wrap(err, "Failed to get latest release binary")
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(dest, 0700); err != nil {
			return err
		}
	}
	return nil
}
