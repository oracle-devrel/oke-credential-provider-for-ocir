# Running the Provider Locally

You can use the following steps to run the provider locally and test the functionality. The script assumes you have OCI CLI installed and configured.

Execute the following command:
```
echo '{
  "apiVersion": "credentialprovider.kubelet.k8s.io/v1",
  "kind": "CredentialProviderRequest",
  "image": "fra.ocir.io/demo_namespace/demo_repo/demo_image"
}' | go run cmd/provider.go -config ./examples/run-locally-yaml/config.yaml
```

If you have OCI CLI installed and configured properly, you should see the output below:
```
{"apiVersion":"credentialprovider.kubelet.k8s.io/v1","kind":"CredentialProviderResponse","cacheKeyType":"Registry","cacheDuration":"0h59m0s","auth":{"fra.ocir.io":{"username":"","password":""}}}
```

That's it!

