/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func init() {
	file, err := os.OpenFile("credential-provider-oke.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	Logger = log.New(file, "", log.LstdFlags)
}

func Log(message string) {
	Logger.Println(message)
}
