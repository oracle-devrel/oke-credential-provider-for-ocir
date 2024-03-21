# Image Credential Provider for OKE

<b>Image Credential Provider</b> (Provider) for  [Container Engine for Kubernetes (OKE)](https://www.oracle.com/cloud/cloud-native/container-engine-kubernetes/) is the implementation of [Kubelet CredentialProvider (v1) APIs](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/) for passwordless pulls from the [Container Registry (OCIR)](https://www.oracle.com/cloud/cloud-native/container-registry/) (OCIR). It's useful since OKE typically [requires](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengpullingimagesfromocir.htm) a stored Secret to pull private OCI images, referenced with `imagePullSecrets` in a manifest. With the provider in place, Kubelet will pull images using Instance Principal authentication, giving you a seamless image-pulling experience without hosting static docker credentials.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Contributing](#contributing)
- [License](#license)

## Prerequisites 
Your OKE Kubelet and API Server versions must be at least v1.26. To check the version, execute `kubectl version`.


## Installation
To install and run the Provider on a worker nodes, follow the steps described [here](/examples/run-on-worker-instance-principal/). 

It's crucial to (1) [create](#1-create-dynamic-group-for-oke-worker-nodes) a dynamic group to represent worker nodes, (2) [create](#2-create-policy-to-pull-images) a Policy to authorize pulling from OCIR, and (3) [configure](#3-configure-cloud-init-for-oke-node-pool) a cloud-init script to do the heavy lifting.

## How the Provider Works

The plugin implementation leverages the kubelet capability introduced in v1.26. The plugin is injected into kubelet via the extra `kubelet-extra-args`:
- `--image-credential-provider-config`set the path to the Image Credential Provider for OKE config file.
- `--image-credential-provider-bin-dir` sets the path to the directory where the Image Credential Provider for OKE binaries is located.

The cloud-init scripts act as glue, downloading the plugin and configuration and passing it through the kubelet.

## Contributing

If you find a bug or want to suggest an enhancement, please raise the Issue.

## License
Copyright (c) 2024 Oracle and/or its affiliates.

Licensed under the Universal Permissive License (UPL), Version 1.0.

See [LICENSE](LICENSE) for more details.