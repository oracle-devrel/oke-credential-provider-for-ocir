/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// ConfigMaps always use strings, even for YAML; so everything MUST be string
	RegistryTokenPath string `yaml:"registryTokenPath"`
	DefaultUser       string `yaml:"defaultUser"`
	RegistryProtocol  string `yaml:"registryProtocol"`
	OCIRAuthMethod    string `yaml:"ocirAuthMethod"`
	TokenValidation   string `yaml:"tokenValidation"`
}

func (config Config) IsTokenValidationEnabled() bool {
	return bytes.EqualFold([]byte("enabled"), []byte(config.TokenValidation))
}

func ReadConfig(configPath string) Config {
	var config Config
	if len(configPath) > 0 {
		Log("Loading config from YAML file: " + configPath)
		config = readConfigFromYaml(configPath)
	} else {
		Log("Loading config from environment variables.")
		config = readConfigFromEnv()
	}

	// Setting the default values for critical configs if unset via configuration
	if config.RegistryProtocol == "" {
		config.RegistryProtocol = "https"
	}
	if config.OCIRAuthMethod == "" {
		config.OCIRAuthMethod = "INSTANCE_PRINCIPAL"
	}
	Log(fmt.Sprintf("Configuration Loaded: %+v", config))

	return config
}

func readConfigFromEnv() Config {
	return Config{
		RegistryTokenPath: os.Getenv("REGISTRY_TOKEN_PATH"),
		DefaultUser:       os.Getenv("DEFAULT_USER"),
		RegistryProtocol:  os.Getenv("REGISTRY_PROTOCOL"),
		OCIRAuthMethod:    os.Getenv("OCIR_AUTH_METHOD"),
		TokenValidation:   os.Getenv("TOKEN_VALIDATION"),
	}
}

func readConfigFromYaml(configPath string) Config {
	configYaml, err := os.ReadFile(configPath)
	FatalIfErrorDescription(err, "Program couldn't locate config.yaml file. Make sure to place it in the folder with the binary")

	var config Config
	err = yaml.Unmarshal(configYaml, &config)
	FatalIfErrorDescription(err, "Program couldn't load config.yaml file. There is a problem with a yaml structure")

	return config
}
