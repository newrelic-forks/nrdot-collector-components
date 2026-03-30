// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

// https://github.com/osquery/osquery/blob/master/specs/secureboot.table
type secureBootCollection struct {
	SecureBoot       int    `json:"secure_boot"`
	SetupMode        int    `json:"setup_mode,omitempty"`
	SecureMode       int    `json:"secure_mode,omitempty"`
	Description      string `json:"description,omitempty"`
	KernelExtensions int    `json:"kernel_extensions,omitempty"`
	MDMOperations    int    `json:"mdm_operations,omitempty"`
}

func (secureBootCollection) GetName() string {
	return secureBootCollectionName
}

func (secureBootCollection) GetQuery() string {
	return secureBootCollectionQuery
}

func (secureBootCollection) Unmarshal(result any) any {
	resultMap, ok := result.(map[string]any)
	if !ok {
		return nil
	}

	sanitized := sanitizeRow(
		resultMap,
		[]string{
			"description",
		},
		[]string{
			"secure_boot",
			"setup_mode",
			"secure_mode",
			"kernel_extensions",
			"mdm_operations",
		},
		nil,
	)

	if len(sanitized) == 0 {
		return nil
	}

	return sanitized
}

func newSecureBootCollection() iCollection {
	return &secureBootCollection{}
}
