/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers

import (
	"fmt"
)

func FatalIfDescription(description string) {
	FatalIfErrorDescription(fmt.Errorf("business error"), description)
}

func FatalIfErrorDescription(err error, description string) {
	if err != nil {
		var message string
		if len(description) > 0 {
			message = fmt.Sprintf("%s:\n%s", description, err.Error())
		} else {
			message = err.Error()
		}
		panic(message)
	}
}

func FatalIfError(err error) {
	FatalIfErrorDescription(err, "")
}
