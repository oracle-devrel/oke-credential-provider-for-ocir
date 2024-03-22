/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers

import "regexp"

func ExtractHostname(imageName string) string {
	regexPattern := `^(?:https?://)?(?:[^@/\n]+@)?(?:www\.)?([^:/\n]+)`
	regex := regexp.MustCompile(regexPattern)
	match := regex.FindStringSubmatch(imageName)
	var hostname string
	if len(match) > 1 {
		hostname = match[1]
	}
	return hostname
}
