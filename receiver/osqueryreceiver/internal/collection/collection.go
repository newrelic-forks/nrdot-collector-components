// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

import (
	"errors"
)

type iCollection interface {
	GetName() string
	GetQuery() string
	Unmarshal(any) any
}

func getCollection(name string) (iCollection, error) {
	switch name {
	case "system_info":
		return newSystemInfoCollection(), nil
	case "package_info":
		return newPackageInfoCollection(), nil
	case "os_info":
		return newOSInfoCollection(), nil
	case "secureboot_info":
		return newSecureBootCollection(), nil
	case "users_info":
		return newUserCollection(), nil
	default:
		return nil, errors.New("wrong collection name passed")
	}
}

func getCustomCollection(name, query string) iCollection {
	return newCustomCollection(name, query)
}

// New returns a collection by name. The concrete type implements GetName, GetQuery, and Unmarshal.
func New(name string) (any, error) {
	return getCollection(name)
}

// NewCustom returns a custom collection wrapping an arbitrary query.
func NewCustom(name, query string) any {
	return getCustomCollection(name, query)
}
