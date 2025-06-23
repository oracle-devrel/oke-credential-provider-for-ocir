/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/devrocks/credential-provider-oke/internal/helpers"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
)

type OcirDockerToken struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type CredentialProviderRequest struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Image      string `json:"image"`
}

type CredentialProviderResponse struct {
	APIVersion    string                `json:"apiVersion"`
	Kind          string                `json:"kind"`
	CacheKeyType  string                `json:"cacheKeyType"`
	CacheDuration string                `json:"cacheDuration"`
	Auth          map[string]AuthConfig `json:"auth"`
}

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func signUserPrincipalRequest(tokenRequest *http.Request) {
	provider := common.DefaultConfigProvider()
	signer := common.DefaultRequestSigner(provider)
	signer.Sign(tokenRequest)
}

func signInstancePrincipalRequest(tokenRequest *http.Request) {
	provider, err := auth.InstancePrincipalConfigurationProvider()
	helpers.FatalIfError(err)
	signer := common.DefaultRequestSigner(provider)
	signer.Sign(tokenRequest)
}

func getDockerToken(urlTokenIssuer string, config helpers.Config) OcirDockerToken {
	tokenRequest, err := http.NewRequest("GET", urlTokenIssuer, nil)
	helpers.FatalIfError(err)
	tokenRequest.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	if config.OCIRAuthMethod == "INSTANCE_PRINCIPAL" {
		signInstancePrincipalRequest(tokenRequest)
	} else if config.OCIRAuthMethod == "USER_PRINCIPAL" {
		signUserPrincipalRequest(tokenRequest)
	} else {
		helpers.FatalIfDescription("OCIR Authentication method is not properly set. It should be one of INSTANCE_PRINCIPAL or USER_PRINCIPAL")
	}
	provider := common.DefaultConfigProvider()
	signer := common.DefaultRequestSigner(provider)
	signer.Sign(tokenRequest)

	client := http.Client{}
	tokenResponse, err := client.Do(tokenRequest)
	helpers.FatalIfError(err)
	defer tokenResponse.Body.Close()

	body, err := io.ReadAll(tokenResponse.Body)
	helpers.FatalIfError(err)

	// helpers.Log(fmt.Sprintf("Token retrieved from %s:\n%s", urlTokenIssuer, string(body)))
	helpers.Log(fmt.Sprintf("Token retrieved from %s", urlTokenIssuer))

	var dockerToken OcirDockerToken
	json.Unmarshal(body, &dockerToken)

	// Reduce lifespan 1m for kubelet cache reasons (e.g. Token with validity 1h0m0s will be shorten to 0h59m0s)
	if dockerToken.ExpiresIn > 60 {
		dockerToken.ExpiresIn -= 60
	}

	return dockerToken
}

func newCredentialProviderResponse(dockerToken OcirDockerToken, image string, cacheDuration string) CredentialProviderResponse {
	return CredentialProviderResponse{
		Kind:          "CredentialProviderResponse",
		APIVersion:    "credentialprovider.kubelet.k8s.io/v1",
		CacheKeyType:  "Registry",
		CacheDuration: cacheDuration,
		Auth: map[string]AuthConfig{
			image: {
				Username: "BEARER_TOKEN",
				Password: dockerToken.Token,
			},
		},
	}
}

func readCredentialProviderRequestFromStdin() CredentialProviderRequest {
	var credentialProviderRequest CredentialProviderRequest
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		input, err := io.ReadAll(os.Stdin)
		helpers.FatalIfError(err)
		err = json.Unmarshal([]byte(input), &credentialProviderRequest)
		helpers.FatalIfErrorDescription(err, "Couldn't read the input CredentialProviderRequest")
	} else {
		helpers.FatalIfDescription("Stdin data is missing. Please supply the program with the proper CredentialProviderRequest json")
	}

	return credentialProviderRequest
}

func GetCredentialProviderResponse(config helpers.Config) {
	registryTokenPath := config.RegistryTokenPath

	credentialProviderRequest := readCredentialProviderRequestFromStdin()
	helpers.Log(fmt.Sprintf("credentialProviderRequest: %s", credentialProviderRequest))
	registryHostname, image, tag, err := helpers.ParseImage(credentialProviderRequest.Image)

	repositoryUrl := fmt.Sprintf("%s://%s", config.RegistryProtocol, registryHostname)
	repositoryEndpoint := fmt.Sprintf("%s%s", repositoryUrl, registryTokenPath)

	issuedToken := getDockerToken(repositoryEndpoint, config)

	duration := issuedToken.ExpiresIn
	if config.IsTokenValidationEnabled() {
		duration = 120 // Initiate the provider with a short lived cache when validating
	}
	cacheDuration := helpers.FormatTimeDuration(duration)

	valid := false
	var status string
	if config.IsTokenValidationEnabled() && !strings.Contains(image, "oke-public") {
		imageManifest := fmt.Sprintf("%s/v2/%s/manifests/%s", repositoryUrl, image, tag)
		valid, status = tokenValidate(imageManifest, AuthConfig{
			Username: "BEARER_TOKEN",
			Password: issuedToken.Token,
		})
		helpers.Log(fmt.Sprintf("Token validation for %s results. Valid: %t; Status: %v", imageManifest, valid, status))
		if valid {
			// Token is working for private regs; increase the cache
			cacheDuration = helpers.FormatTimeDuration(issuedToken.ExpiresIn)
		}
	}
	helpers.Log(fmt.Sprintf("Provider cacheDuration: %s", cacheDuration))
	credentialProviderResponse := newCredentialProviderResponse(issuedToken, registryHostname, cacheDuration)

	result, err := json.Marshal(credentialProviderResponse)
	helpers.FatalIfError(err)

	// helpers.Log(fmt.Sprintf("credentialProviderResponse: %s", credentialProviderResponse))
	fmt.Print(string(result))
}

// Verifies that the token can be used to authenticate to the remote registry.
func tokenValidate(manifestUrl string, auth AuthConfig) (valid bool, status string) {
	req, err := http.NewRequest("GET", manifestUrl, nil)
	if err != nil {
		return false, fmt.Sprintf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+auth.Password)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Sprintf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return true, fmt.Sprintf("successful response: %s", resp.Status)
	case http.StatusForbidden:
		return false, fmt.Sprintf("failed response: %s", resp.Status)
	default:
		return false, fmt.Sprintf("failed response: %s", resp.Status)
	}
}
