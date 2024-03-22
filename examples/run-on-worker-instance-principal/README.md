# Running the Provider on a Worker Node

To run the Provider on a worker node, follow the steps below. It's crucial to (1) [create](#1-create-dynamic-group-for-oke-worker-nodes) a dynamic group to represent worker nodes, (2) [create](#2-create-policy-to-pull-images) a Policy to authorize pulling from OCIR, and (3) [configure](#3-configure-cloud-init-for-oke-node-pool) a cloud-init script to do the heavy lifting.

## 1. Create Dynamic Group for OKE Worker Nodes
First, you must [create a dynamic group](https://docs.oracle.com/en-us/iaas/Content/Identity/dynamicgroups/To_create_a_dynamic_group.htm) representing worker nodes of a targeted OKE cluster. The example below adds all instances within the compartment to the dynamic group `oke-puller`.
```
All {instance.compartment.id = '<REPLACE_WITH_YOUR_COMPARTMENT_OCID>'}
```
Replace `<REPLACE_WITH_YOUR_COMPARTMENT_OCID>` with the actual compartment OCID where the worker nodes reside. Make sure that targeted worker nodes are a part of the compartment.

Additionally, you might expand the dynamic group expression to defined tags for better precision in selecting worker nodes.

## 2. Create Policy to Pull Images
[Create the policy](https://docs.oracle.com/en-us/iaas/Content/Identity/policymgmt/managingpolicies_topic-To_create_a_policy.htm) to authorize dynamic group `oke-puller` to pull images from OCIR. The example below allows the dynamic group `oke-puller` to pull images from any repository in a selected compartment.
```
Allow dynamic-group <REPLACE_WITH_YOUR_IAM_DOMAIN>/oke-puller to read repos in compartment <REPLACE_WITH_YOUR_COMPARTMENT_NAME>
```
Replace `<REPLACE_WITH_YOUR_IAM_DOMAIN>` with the IAM domain name. Replace `<REPLACE_WITH_YOUR_COMPARTMENT_NAME>` with the compartment name of the repositories with images.

## 3. Configure Cloud-Init for OKE Node Pool
The plugin is activated with the cloud-init executed at the boot time of every worker node. OKE uses cloud-init to set up the compute instances hosting managed nodes. Paste the [following](cloud-init.sh) cloud-init script to every node pool where you want to enable Image Credential Provider for OKE. Node pools must be a part of the dynamic group `oke-puller`.