#!/usr/bin/env bash
set -euxo pipefail

# This script sets up an environment for running integration tests with a local
# microk8s cluster. Assumes that make, gcc/clang, and docker are already installed

sudo snap install microk8s --classic
sudo snap install go --channel 1.19/stable --classic
sudo snap install kubectl --classic
sudo microk8s enable registry
sudo usermod -aG microk8s ubuntu
