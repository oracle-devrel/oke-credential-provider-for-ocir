# Image Credential Provider for OKE

<b>Image Credential Provider</b> (Provider) for  [Container Engine for Kubernetes (OKE)](https://www.oracle.com/cloud/cloud-native/container-engine-kubernetes/) is the implementation of [Kubelet CredentialProvider (v1) APIs](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/) for passwordless pulls from the [Container Registry (OCIR)](https://www.oracle.com/cloud/cloud-native/container-registry/) (OCIR). It's useful since OKE typically [requires](https://docs.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengpullingimagesfromocir.htm) a stored Secret to pull private OCIR images, referenced with `imagePullSecrets` in a manifest. With the provider in place, Kubelet will pull images using instance principal authentication, giving you a seamless image-pulling experience without hosting static Docker credentials.

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
The plugin implementation leverages the Kubelet capability introduced in v1.26. Kubelet uses [CredentialProvider](https://kubernetes.io/docs/reference/config-api/kubelet-credentialprovider.v1/) APIs to fetch authentication credentials against Docker comaptible image registry and caches it on the worker node level. The plugin translates instance principal authentication into the JWT token that is used by Kubelet when pulling images from OCIR at runtime. In that case, you don't need to specify `imagePullSecrets` in a manifest, since Kubelet has JWT token based on instance principal auth locally.

The provider is injected into Kubelet via the extra `kubelet-extra-args`:
- `--image-credential-provider-config` sets the path to the Image Credential Provider for OKE config file.
- `--image-credential-provider-bin-dir` sets the path to the directory where the Image Credential Provider for OKE binary is located.

The cloud-init script act as glue, downloading the provider with the configuration file and passing it to the Kubelet. 

The current [cloud-init.sh](examples/run-on-worker-instance-principal/cloud-init.sh) example implementation uses the `wget` utility to download binaries on the worker nodes. Suppose you don't have access to the Internet (through NAT gateway) or your OS does not have a `wget`. In that case, you need to place binaries and configuration in the appropriate folders manually:
- The [provider binary (amd64)](https://github.com/oracle-devrel/oke-credential-provider-for-ocir/releases/latest/download/oke-credential-provider-for-ocir-linux-amd64) with the name `oke-credential-provider` must be in the following path: `/usr/local/bin`. Make sure the binary has permission mode to execute. You can enable it by executing `sudo chmod 755 /usr/local/bin/oke-credential-provider`.
- The kubelet configuration file [credential-provider-config.yaml](https://github.com/oracle-devrel/oke-credential-provider-for-ocir/releases/latest/download/credential-provider-config.yaml) must be placed in the path `/etc/kubernetes`.

Plugin binaries are avaialble both for OCI [arm64](https://github.com/oracle-devrel/oke-credential-provider-for-ocir/releases/latest/download/oke-credential-provider-for-ocir-linux-arm64) and [amd64](https://github.com/oracle-devrel/oke-credential-provider-for-ocir/releases/latest/download/oke-credential-provider-for-ocir-linux-amd64) architectures.

## Contributing

If you find a bug or want to suggest an enhancement, please raise the Issue.

## License
Copyright (c) 2024 Oracle and/or its affiliates.

Licensed under the Universal Permissive License (UPL), Version 1.0.

See [LICENSE](LICENSE) for more details.