#!/bin/bash

# Copyright 2016 The Kubernetes Authors All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# This script runs the integration tests on a Linux machine for the none Driver

# The script expects the following env variables:
# MINIKUBE_LOCATION: GIT_COMMIT from upstream build.
# COMMIT: Actual commit ID from upstream build
# EXTRA_BUILD_ARGS (optional): Extra args to be passed into the minikube integrations tests
# access_token: The Github API access token. Injected by the Jenkins credential provider. 


set -e

OS_ARCH="linux-amd64"
VM_DRIVER="none"
JOB_NAME="none_Linux"
EXTRA_ARGS="--bootstrapper=kubeadm"
EXPECTED_DEFAULT_DRIVER="kvm2"

SUDO_PREFIX="sudo -E "
export KUBECONFIG="/root/.kube/config"

# "none" driver specific cleanup from previous runs.
sudo kubeadm reset -f || true
# kubeadm reset may not stop pods immediately
docker rm -f $(docker ps -aq) >/dev/null 2>&1 || true

# Cleanup data directory
sudo rm -rf /data/*
# Cleanup old Kubernetes configs
sudo rm -rf /etc/kubernetes/*
# Cleanup old minikube files
sudo rm -rf /var/lib/minikube/*
# Stop any leftover kubelets
systemctl is-active --quiet kubelet \
  && echo "stopping kubelet" \
  && sudo systemctl stop kubelet

mkdir -p cron && gsutil -m rsync "gs://minikube-builds/${MINIKUBE_LOCATION}/cron" cron || echo "FAILED TO GET CRON FILES"
sudo install cron/cleanup_and_reboot_Linux.sh /etc/cron.hourly/cleanup_and_reboot || echo "FAILED TO INSTALL CLEANUP"

## download and install gopogh
wget -O gopogh https://github.com/medyagh/gopogh/releases/download/v0.0.14/gopogh-linux-amd64 || echo "Failed to download gopogh"
sudo install gopogh /usr/local/bin/ || echo "Failed to install Gopogh"


# Download files and set permissions
source ./common.sh
