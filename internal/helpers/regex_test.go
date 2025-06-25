/*
 * Copyright (c) 2024 Oracle and/or its affiliates.
 * Use of this source code is governed by The Universal Permissive License (UPL), Version 1.0. that can be found in the LICENSE file.
 */

package helpers_test

import (
	"errors"
	"testing"

	"github.com/devrocks/credential-provider-oke/internal/helpers"
)

func Test_ParseImage(t *testing.T) {
	cases := map[string]struct {
		image      string
		hostname   string
		repository string
		tag        string
		err        error
	}{
		"CanonicalImageName": {
			image:      "fra.ocir.io/demo_namespace/demo_repo/demo_image:some_tag",
			hostname:   "fra.ocir.io",
			repository: "demo_namespace/demo_repo/demo_image",
			tag:        "some_tag",
			err:        nil,
		},
		"UntaggedCanonicalImageName": {
			image:      "fra.ocir.io/demo_namespace/demo_repo/demo_image",
			hostname:   "fra.ocir.io",
			repository: "demo_namespace/demo_repo/demo_image",
			tag:        "latest",
			err:        nil,
		},
		"WrongImageName": {
			image:      "some unexpected expression",
			hostname:   "",
			repository: "",
			tag:        "",
			err:        errors.New("invalid image format"),
		},
	}

	for name, test := range cases {
		t.Run(name, func(*testing.T) {
			hostname, repository, tag, err := helpers.ParseImage(test.image)
			if hostname != test.hostname {
				t.Errorf("Hostname (expected: '%s', actual: '%s')", test.hostname, hostname)
			}
			if repository != test.repository {
				t.Errorf("Repository (expected: '%s', actual: '%s')", test.repository, repository)
			}
			if tag != test.tag {
				t.Errorf("Tag (expected: '%s', actual: '%s')", test.tag, tag)
			}
			if err == nil && test.err != nil {
				t.Errorf("Expected err '%s' not produced", test.err.Error())
			}
			if err != nil && test.err == nil {
				t.Errorf("Unexpected Err: '%s')", err.Error())
			}
			if err != nil && test.err != nil && err.Error() != test.err.Error() {
				t.Errorf("Err (expected: '%s', actual: '%s')", test.err.Error(), err.Error())
			}
		})
	}
}
