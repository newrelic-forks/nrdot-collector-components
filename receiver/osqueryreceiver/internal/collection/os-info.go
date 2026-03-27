// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

// https://github.com/osquery/osquery/blob/master/specs/os_version.table
type osInfoCollection struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Build        string `json:"build"`
	Platform     string `json:"platform"`
	PlatformLike string `json:"platform_like"`
	Codename     string `json:"codename,omitempty"`
	Arch         string `json:"arch,omitempty"`
}

func (osInfoCollection) GetName() string {
	return osInfoCollectionName
}

func (osInfoCollection) GetQuery() string {
	return osInfoCollectionQuery
}

func (osInfoCollection) Unmarshal(result any) any {
	resultMap, ok := result.(map[string]any)
	if !ok {
		return nil
	}

	sanitized := sanitizeRow(
		resultMap,
		[]string{
			"name",
			"platform",
			"platform_like",
			"build",
			"version",
			"codename",
			"arch",
		},
		nil,
		nil,
	)

	if len(sanitized) == 0 {
		return nil
	}

	return sanitized
}

func newOSInfoCollection() iCollection {
	return &osInfoCollection{}
}
