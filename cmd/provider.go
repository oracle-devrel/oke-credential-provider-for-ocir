/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package main

import (
	"flag"

	"github.com/devrocks/credential-provider-oke/internal/helpers"
	"github.com/devrocks/credential-provider-oke/internal/provider"
)

func main() {
	var configType string
	flag.StringVar(&configType, "config", "", "Type the absolute path to yaml file with provider configuration. If it's ommited, program reads from environment variables (REGISTRY_TOKEN_PATH, DEFAULT_USER, REGISTRY_PROTOCOL, OCIR_AUTH_METHOD).")
	flag.Parse()
	config := helpers.ReadConfig(configType)
	provider.GetCredentialProviderResponse(config)
}
