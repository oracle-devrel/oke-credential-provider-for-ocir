/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers

import (
	"fmt"
	"time"
)

func FormatTimeDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds = seconds % 60
	return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
}
