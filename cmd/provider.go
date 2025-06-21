/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package main

import (
	"flag"
	"os"
	"time"

	"github.com/gofrs/flock"
	"github.com/devrocks/credential-provider-oke/internal/helpers"
	"github.com/devrocks/credential-provider-oke/internal/provider"
)

func main() {
	var configType string
	flag.StringVar(&configType, "config", "", "Path to config YAML or environment config")
	flag.Parse()

	lockFilePath := "/var/lib/kubelet/credential-provider.lock"
	fileLock := flock.New(lockFilePath)

	locked, err := fileLock.TryLock()
	if err != nil {
		os.Exit(1)
	}

	if !locked {
		// Lock is held by another process, sleep then exit
		// that process will then be able to use the cache and reduce requests
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}
	defer fileLock.Unlock()

	config := helpers.ReadConfig(configType)
	provider.GetCredentialProviderResponse(config)
}
