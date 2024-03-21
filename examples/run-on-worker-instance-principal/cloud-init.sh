#!/bin/bash
curl --fail -H "Authorization: Bearer Oracle" -L0 http://169.254.169.254/opc/v2/instance/metadata/oke_init_script | base64 --decode >/var/run/oke-init.sh

# download binaries on the worker node
wget https://frsxwtjslf35.objectstorage.eu-frankfurt-1.oci.customer-oci.com/p/vroemLQGZVbYNF8qgSKYW_mEyP9a5l-vxY_ASH20DuphPzbK-rkNlaSWlS0mFzpe/n/frsxwtjslf35/b/oke-credential-provider-for-ocir/o/oke-credential-provider-for-ocir-amd64 -P /usr/local/bin/
wget https://raw.githubusercontent.com/oracle-devrel/oke-credential-provider-for-ocir/main//examples/run-on-worker-instance-principal/credential-provider-config.yaml -P /etc/kubernetes/

# add permission to execute
sudo chmod 755 /usr/local/bin/oke-credential-provider-for-ocir-amd64

# configure kubelet with image credential provider
bash /var/run/oke-init.sh --kubelet-extra-args "--image-credential-provider-bin-dir=/usr/local/bin/ --image-credential-provider-config=/etc/kubernetes/credential-provider-config.yaml"