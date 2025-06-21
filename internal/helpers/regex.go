/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers

import (
	"regexp" 
	"fmt"
)

func ParseImage(image string) (hostname, repository, tag string, err error) {
	regexPattern := `^([a-zA-Z0-9.-]+(?::[0-9]+)?)/([^:]+)(?::(.+))?$`
	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(image)
	if len(matches) == 0 {
		err = fmt.Errorf("invalid image format")
		return
	}

	hostname = matches[1]
	repository = matches[2]
	tag = matches[3]
	if tag == "" {
		tag = "latest"
	}
	return
}