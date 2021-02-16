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

package assets

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"k8s.io/minikube/pkg/minikube/config"
	"k8s.io/minikube/pkg/minikube/constants"
	"k8s.io/minikube/pkg/minikube/out"
	"k8s.io/minikube/pkg/minikube/style"
	"k8s.io/minikube/pkg/minikube/vmpath"
	"k8s.io/minikube/pkg/version"
)

// Addon is a named list of assets, that can be enabled
type Addon struct {
	Assets    []*BinAsset
	enabled   bool
	addonName string
	Images    map[string]string

	// Registries currently only shows the default registry of images
	Registries map[string]string
}

// NewAddon creates a new Addon
func NewAddon(assets []*BinAsset, enabled bool, addonName string, images map[string]string, registries map[string]string) *Addon {
	a := &Addon{
		Assets:     assets,
		enabled:    enabled,
		addonName:  addonName,
		Images:     images,
		Registries: registries,
	}
	return a
}

// Name get the addon name
func (a *Addon) Name() string {
	return a.addonName
}

// IsEnabled checks if an Addon is enabled for the given profile
func (a *Addon) IsEnabled(cc *config.ClusterConfig) bool {
	status, ok := cc.Addons[a.Name()]
	if ok {
		return status
	}

	// Return the default unconfigured state of the addon
	return a.enabled
}

// Addons is the list of addons
// TODO: Make dynamically loadable: move this data to a .yaml file within each addon directory
var Addons = map[string]*Addon{
	"auto-pause": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/auto-pause/auto-pause.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"auto-pause.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/auto-pause/haproxy.cfg",
			"/var/lib/minikube/",
			"haproxy.cfg",
			"0640"),
		MustBinAsset(
			"deploy/addons/auto-pause/unpause.lua",
			"/var/lib/minikube/",
			"unpause.lua",
			"0640"),
		MustBinAsset(
			"deploy/addons/auto-pause/auto-pause.service",
			"/etc/systemd/system/",
			"auto-pause.service",
			"0640"),

		//GuestPersistentDir
	}, false, "auto-pause", map[string]string{
		"haproxy": "haproxy:2.3.5",
	}, map[string]string{
		"haproxy": "gcr.io",
	}),
	"dashboard": NewAddon([]*BinAsset{
		// We want to create the kubernetes-dashboard ns first so that every subsequent object can be created
		MustBinAsset("deploy/addons/dashboard/dashboard-ns.yaml", vmpath.GuestAddonsDir, "dashboard-ns.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-clusterrole.yaml", vmpath.GuestAddonsDir, "dashboard-clusterrole.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-clusterrolebinding.yaml", vmpath.GuestAddonsDir, "dashboard-clusterrolebinding.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-configmap.yaml", vmpath.GuestAddonsDir, "dashboard-configmap.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-dp.yaml.tmpl", vmpath.GuestAddonsDir, "dashboard-dp.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-role.yaml", vmpath.GuestAddonsDir, "dashboard-role.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-rolebinding.yaml", vmpath.GuestAddonsDir, "dashboard-rolebinding.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-sa.yaml", vmpath.GuestAddonsDir, "dashboard-sa.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-secret.yaml", vmpath.GuestAddonsDir, "dashboard-secret.yaml", "0640"),
		MustBinAsset("deploy/addons/dashboard/dashboard-svc.yaml", vmpath.GuestAddonsDir, "dashboard-svc.yaml", "0640"),
	}, false, "dashboard", map[string]string{
		"Dashboard":      "kubernetesui/dashboard:v2.1.0",
		"MetricsScraper": "kubernetesui/metrics-scraper:v1.0.4",
	}, nil),
	"default-storageclass": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/storageclass/storageclass.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"storageclass.yaml",
			"0640"),
	}, true, "default-storageclass", nil, nil),
	"pod-security-policy": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/pod-security-policy/pod-security-policy.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"pod-security-policy.yaml",
			"0640"),
	}, false, "pod-security-policy", nil, nil),
	"storage-provisioner": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/storage-provisioner/storage-provisioner.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"storage-provisioner.yaml",
			"0640"),
	}, true, "storage-provisioner", map[string]string{
		"StorageProvisioner": fmt.Sprintf("k8s-minikube/storage-provisioner:%s", version.GetStorageProvisionerVersion()),
	}, map[string]string{
		"StorageProvisioner": "gcr.io",
	}),
	"storage-provisioner-gluster": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/storage-provisioner-gluster/storage-gluster-ns.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"storage-gluster-ns.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/storage-provisioner-gluster/glusterfs-daemonset.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"glusterfs-daemonset.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/storage-provisioner-gluster/heketi-deployment.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"heketi-deployment.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/storage-provisioner-gluster/storage-provisioner-glusterfile.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"storage-privisioner-glusterfile.yaml",
			"0640"),
	}, false, "storage-provisioner-gluster", map[string]string{
		"Heketi":                 "heketi/heketi:latest",
		"GlusterfileProvisioner": "gluster/glusterfile-provisioner:latest",
		"GlusterfsServer":        "nixpanic/glusterfs-server:pr_fake-disk",
	}, map[string]string{
		"GlusterfsServer": "quay.io",
	}),
	"efk": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/efk/elasticsearch-rc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"elasticsearch-rc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/efk/elasticsearch-svc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"elasticsearch-svc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/efk/fluentd-es-rc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"fluentd-es-rc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/efk/fluentd-es-configmap.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"fluentd-es-configmap.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/efk/kibana-rc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"kibana-rc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/efk/kibana-svc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"kibana-svc.yaml",
			"0640"),
	}, false, "efk", map[string]string{
		"Elasticsearch":        "elasticsearch:v5.6.2",
		"FluentdElasticsearch": "fluentd-elasticsearch:v2.0.2",
		"Alpine":               "alpine:3.6",
		"Kibana":               "kibana/kibana:5.6.2",
	}, map[string]string{
		"Elasticsearch":        "k8s.gcr.io",
		"FluentdElasticsearch": "k8s.gcr.io",
		"Kibana":               "docker.elastic.co",
	}),
	"ingress": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/ingress/ingress-configmap.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ingress-configmap.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/ingress/ingress-rbac.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ingress-rbac.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/ingress/ingress-dp.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ingress-dp.yaml",
			"0640"),
	}, false, "ingress", map[string]string{
		"IngressController":        "k8s-artifacts-prod/ingress-nginx/controller:v0.40.2",
		"KubeWebhookCertgenCreate": "jettech/kube-webhook-certgen:v1.2.2",
		"KubeWebhookCertgenPatch":  "jettech/kube-webhook-certgen:v1.3.0",
	}, map[string]string{
		"IngressController": "us.gcr.io",
	}),
	"istio-provisioner": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/istio-provisioner/istio-operator.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"istio-operator.yaml",
			"0640"),
	}, false, "istio-provisioner", map[string]string{
		"IstioOperator": "istio/operator:1.5.0",
	}, nil),
	"istio": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/istio/istio-default-profile.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"istio-default-profile.yaml",
			"0640"),
	}, false, "istio", nil, nil),
	"kubevirt": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/kubevirt/pod.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"pod.yaml",
			"0640"),
	}, false, "kubevirt", map[string]string{
		"Kubectl": "bitnami/kubectl:1.17",
	}, nil),
	"metrics-server": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/metrics-server/metrics-apiservice.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"metrics-apiservice.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/metrics-server/metrics-server-deployment.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"metrics-server-deployment.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/metrics-server/metrics-server-service.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"metrics-server-service.yaml",
			"0640"),
	}, false, "metrics-server", map[string]string{
		"MetricsServer": fmt.Sprintf("metrics-server-%s:v0.2.1", runtime.GOARCH),
	}, map[string]string{
		"MetricsServer": "k8s.gcr.io",
	}),
	"olm": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/olm/crds.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"crds.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/olm/olm.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"olm.yaml",
			"0640"),
	}, false, "olm", map[string]string{
		"OLM":                        "operator-framework/olm:0.14.1",
		"UpstreamCommunityOperators": "operator-framework/upstream-community-operators:latest",
	}, map[string]string{
		"OLM":                        "quay.io",
		"UpstreamCommunityOperators": "quay.io",
	}),
	"registry": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/registry/registry-rc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"registry-rc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/registry/registry-svc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"registry-svc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/registry/registry-proxy.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"registry-proxy.yaml",
			"0640"),
	}, false, "registry", map[string]string{
		"Registry":          "registry:2.7.1",
		"KubeRegistryProxy": "google_containers/kube-registry-proxy:0.4",
	}, map[string]string{
		"KubeRegistryProxy": "gcr.io",
	}),
	"registry-creds": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/registry-creds/registry-creds-rc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"registry-creds-rc.yaml",
			"0640"),
	}, false, "registry-creds", map[string]string{
		"RegistryCreds": "upmcenterprises/registry-creds:1.10",
	}, nil),
	"registry-aliases": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/registry-aliases/registry-aliases-sa.tmpl",
			vmpath.GuestAddonsDir,
			"registry-aliases-sa.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/registry-aliases/registry-aliases-sa-crb.tmpl",
			vmpath.GuestAddonsDir,
			"registry-aliases-sa-crb.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/registry-aliases/registry-aliases-config.tmpl",
			vmpath.GuestAddonsDir,
			"registry-aliases-config.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/registry-aliases/node-etc-hosts-update.tmpl",
			vmpath.GuestAddonsDir,
			"node-etc-hosts-update.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/registry-aliases/patch-coredns-job.tmpl",
			vmpath.GuestAddonsDir,
			"patch-coredns-job.yaml",
			"0640"),
	}, false, "registry-aliases", map[string]string{
		"CoreDNSPatcher": "rhdevelopers/core-dns-patcher",
		"Alpine":         "alpine:3.11",
		"Pause":          "google_containers/pause-amd64:3.1",
	}, map[string]string{
		"CoreDNSPatcher": "quay.io",
		"Pause":          "gcr.io",
	}),
	"freshpod": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/freshpod/freshpod-rc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"freshpod-rc.yaml",
			"0640"),
	}, false, "freshpod", map[string]string{
		"FreshPod": "google-samples/freshpod:v0.0.1",
	}, map[string]string{
		"FreshPod": "gcr.io",
	}),
	"nvidia-driver-installer": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/gpu/nvidia-driver-installer.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"nvidia-driver-installer.yaml",
			"0640"),
	}, false, "nvidia-driver-installer", map[string]string{
		"NvidiaDriverInstaller": "minikube-nvidia-driver-installer:e2d9b43228decf5d6f7dce3f0a85d390f138fa01",
		"Pause":                 "pause:2.0",
	}, map[string]string{
		"NvidiaDriverInstaller": "k8s.gcr.io",
		"Pause":                 "k8s.gcr.io",
	}),
	"nvidia-gpu-device-plugin": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/gpu/nvidia-gpu-device-plugin.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"nvidia-gpu-device-plugin.yaml",
			"0640"),
	}, false, "nvidia-gpu-device-plugin", map[string]string{
		"NvidiaDevicePlugin": "nvidia/k8s-device-plugin:1.0.0-beta4",
	}, nil),
	"logviewer": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/logviewer/logviewer-dp-and-svc.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"logviewer-dp-and-svc.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/logviewer/logviewer-rbac.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"logviewer-rbac.yaml",
			"0640"),
	}, false, "logviewer", map[string]string{
		"LogViewer": "ivans3/minikube-log-viewer:latest",
	}, nil),
	"gvisor": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/gvisor/gvisor-pod.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"gvisor-pod.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/gvisor/gvisor-runtimeclass.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"gvisor-runtimeclass.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/gvisor/gvisor-config.toml",
			vmpath.GuestGvisorDir,
			constants.GvisorConfigTomlTargetName,
			"0640"),
	}, false, "gvisor", map[string]string{
		"GvisorAddon": "k8s-minikube/gvisor-addon:3",
	}, map[string]string{
		"GvisorAddon": "gcr.io",
	}),
	"helm-tiller": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/helm-tiller/helm-tiller-dp.tmpl",
			vmpath.GuestAddonsDir,
			"helm-tiller-dp.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/helm-tiller/helm-tiller-rbac.tmpl",
			vmpath.GuestAddonsDir,
			"helm-tiller-rbac.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/helm-tiller/helm-tiller-svc.tmpl",
			vmpath.GuestAddonsDir,
			"helm-tiller-svc.yaml",
			"0640"),
	}, false, "helm-tiller", map[string]string{
		"Tiller": "kubernetes-helm/tiller:v2.16.12",
	}, map[string]string{
		"Tiller": "gcr.io",
	}),
	"ingress-dns": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/ingress-dns/ingress-dns-pod.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ingress-dns-pod.yaml",
			"0640"),
	}, false, "ingress-dns", map[string]string{
		"IngressDNS": "cryptexlabs/minikube-ingress-dns:0.3.0",
	}, nil),
	"metallb": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/metallb/metallb.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"metallb.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/metallb/metallb-config.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"metallb-config.yaml",
			"0640"),
	}, false, "metallb", map[string]string{
		"Speaker":    "metallb/speaker:v0.8.2",
		"Controller": "metallb/controller:v0.8.2",
	}, nil),
	"ambassador": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/ambassador/ambassador-operator-crds.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ambassador-operator-crds.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/ambassador/ambassador-operator.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ambassador-operator.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/ambassador/ambassadorinstallation.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"ambassadorinstallation.yaml",
			"0640"),
	}, false, "ambassador", map[string]string{
		"AmbassadorOperator": "datawire/ambassador-operator:v1.2.3",
	}, map[string]string{
		"AmbassadorOperator": "quay.io",
	}),
	"gcp-auth": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/gcp-auth/gcp-auth-ns.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"gcp-auth-ns.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/gcp-auth/gcp-auth-service.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"gcp-auth-service.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/gcp-auth/gcp-auth-webhook.yaml.tmpl.tmpl",
			vmpath.GuestAddonsDir,
			"gcp-auth-webhook.yaml",
			"0640"),
	}, false, "gcp-auth", map[string]string{
		"KubeWebhookCertgen": "jettech/kube-webhook-certgen:v1.3.0",
		"GCPAuthWebhook":     "k8s-minikube/gcp-auth-webhook:v0.0.3",
	}, map[string]string{
		"GCPAuthWebhook": "gcr.io",
	}),
	"volumesnapshots": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/volumesnapshots/snapshot.storage.k8s.io_volumesnapshotclasses.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"snapshot.storage.k8s.io_volumesnapshotclasses.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/volumesnapshots/snapshot.storage.k8s.io_volumesnapshotcontents.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"snapshot.storage.k8s.io_volumesnapshotcontents.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/volumesnapshots/snapshot.storage.k8s.io_volumesnapshots.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"snapshot.storage.k8s.io_volumesnapshots.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/volumesnapshots/rbac-volume-snapshot-controller.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"rbac-volume-snapshot-controller.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/volumesnapshots/volume-snapshot-controller-deployment.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"volume-snapshot-controller-deployment.yaml",
			"0640"),
	}, false, "volumesnapshots", map[string]string{
		"SnapshotController": "k8s-staging-csi/snapshot-controller:v2.0.0-rc2",
	}, map[string]string{
		"SnapshotController": "gcr.io",
	}),
	"csi-hostpath-driver": NewAddon([]*BinAsset{
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/rbac/rbac-external-attacher.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"rbac-external-attacher.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/rbac/rbac-external-provisioner.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"rbac-external-provisioner.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/rbac/rbac-external-resizer.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"rbac-external-resizer.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/rbac/rbac-external-snapshotter.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"rbac-external-snapshotter.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-attacher.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-attacher.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-driverinfo.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-driverinfo.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-plugin.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-plugin.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-provisioner.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-provisioner.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-resizer.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-resizer.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-snapshotter.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-snapshotter.yaml",
			"0640"),
		MustBinAsset(
			"deploy/addons/csi-hostpath-driver/deploy/csi-hostpath-storageclass.yaml.tmpl",
			vmpath.GuestAddonsDir,
			"csi-hostpath-storageclass.yaml",
			"0640"),
	}, false, "csi-hostpath-driver", map[string]string{
		"Attacher":            "k8scsi/csi-attacher:v3.0.0-rc1",
		"NodeDriverRegistrar": "k8scsi/csi-node-driver-registrar:v1.3.0",
		"HostPathPlugin":      "k8scsi/hostpathplugin:v1.4.0-rc2",
		"LivenessProbe":       "k8scsi/livenessprobe:v1.1.0",
		"Resizer":             "k8scsi/csi-resizer:v0.6.0-rc1",
		"Snapshotter":         "k8scsi/csi-snapshotter:v2.1.0",
		"Provisioner":         "k8s-staging-sig-storage/csi-provisioner:v2.0.0-rc2",
	}, map[string]string{
		"Attacher":            "quay.io",
		"NodeDriverRegistrar": "quay.io",
		"HostPathPlugin":      "quay.io",
		"LivenessProbe":       "quay.io",
		"Resizer":             "quay.io",
		"Snapshotter":         "quay.io",
		"Provisioner":         "gcr.io",
	}),
}

// GenerateTemplateData generates template data for template assets
func GenerateTemplateData(addon *Addon, cfg config.KubernetesConfig) interface{} {

	a := runtime.GOARCH
	// Some legacy docker images still need the -arch suffix
	// for  less common architectures blank suffix for amd64
	ea := ""
	if runtime.GOARCH != "amd64" {
		ea = "-" + runtime.GOARCH
	}
	opts := struct {
		Arch                string
		ExoticArch          string
		ImageRepository     string
		LoadBalancerStartIP string
		LoadBalancerEndIP   string
		CustomIngressCert   string
		Images              map[string]string
		Registries          map[string]string
		CustomRegistries    map[string]string
	}{
		Arch:                a,
		ExoticArch:          ea,
		ImageRepository:     cfg.ImageRepository,
		LoadBalancerStartIP: cfg.LoadBalancerStartIP,
		LoadBalancerEndIP:   cfg.LoadBalancerEndIP,
		CustomIngressCert:   cfg.CustomIngressCert,
		Images:              addon.Images,
		Registries:          addon.Registries,
		CustomRegistries:    make(map[string]string),
	}
	if opts.ImageRepository != "" && !strings.HasSuffix(opts.ImageRepository, "/") {
		opts.ImageRepository += "/"
	}

	if opts.Images == nil {
		opts.Images = make(map[string]string) // Avoid nil access when rendering
	}

	images := viper.GetString(config.AddonImages)
	if images != "" {
		for _, image := range strings.Split(images, ",") {
			vals := strings.Split(image, "=")
			if len(vals) != 2 || vals[1] == "" {
				out.WarningT("Ignoring invalid custom image {{.conf}}", out.V{"conf": image})
				continue
			}
			if _, ok := opts.Images[vals[0]]; ok {
				opts.Images[vals[0]] = vals[1]
			} else {
				out.WarningT("Ignoring unknown custom image {{.name}}", out.V{"name": vals[0]})
			}
		}
	}

	if opts.Registries == nil {
		opts.Registries = make(map[string]string)
	}

	registries := viper.GetString(config.AddonRegistries)
	if registries != "" {
		for _, registry := range strings.Split(registries, ",") {
			vals := strings.Split(registry, "=")
			if len(vals) != 2 {
				out.WarningT("Ignoring invalid custom registry {{.conf}}", out.V{"conf": registry})
				continue
			}
			if _, ok := opts.Images[vals[0]]; ok { // check images map because registry map may omitted default registry
				opts.CustomRegistries[vals[0]] = vals[1]
			} else {
				out.WarningT("Ignoring unknown custom registry {{.name}}", out.V{"name": vals[0]})
			}
		}
	}

	// Append postfix "/" to registries
	for k, v := range opts.Registries {
		if v != "" && !strings.HasSuffix(v, "/") {
			opts.Registries[k] = v + "/"
		}
	}

	for k, v := range opts.CustomRegistries {
		if v != "" && !strings.HasSuffix(v, "/") {
			opts.CustomRegistries[k] = v + "/"
		}
	}

	for name, image := range opts.Images {
		if _, ok := opts.Registries[name]; !ok {
			opts.Registries[name] = "" // Avoid nil access when rendering
		}

		// Send messages to stderr due to some tests rely on stdout
		if override, ok := opts.CustomRegistries[name]; ok {
			out.ErrT(style.Option, "Using image {{.registry}}{{.image}}", out.V{
				"registry": override,
				"image":    image,
			})
		} else if opts.ImageRepository != "" {
			out.ErrT(style.Option, "Using image {{.registry}}{{.image}} (global image repository)", out.V{
				"registry": opts.ImageRepository,
				"image":    image,
			})
		} else {
			out.ErrT(style.Option, "Using image {{.registry}}{{.image}}", out.V{
				"registry": opts.Registries[name],
				"image":    image,
			})
		}
	}
	return opts
}
