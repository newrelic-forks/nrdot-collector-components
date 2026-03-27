// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

import (
	"runtime"
)

var userCollectionQueryMap = map[string]string{
	"linux":   userCollectionQueryLinux,
	"darwin":  userCollectionQueryDarwin,
	"windows": userCollectionQueryWindows,
}

type userCollection struct {
	Username string `json:"username"`
	Groups   string `json:"groups,omitempty"`
}

func (userCollection) GetName() string {
	return userCollectionName
}

func (userCollection) GetQuery() string {
	return userCollectionQueryMap[runtime.GOOS]
}

func (userCollection) Unmarshal(result any) any {
	resultSlice, ok := result.([]map[string]any)
	if !ok {
		return nil
	}

	usersList := make([]map[string]any, 0, len(resultSlice))
	for _, resultMap := range resultSlice {
		sanitized := sanitizeRow(
			resultMap,
			[]string{
				"username",
				"groups",
			},
			nil,
			nil,
		)
		if len(sanitized) == 0 {
			continue
		}
		usersList = append(usersList, sanitized)
	}

	if len(usersList) == 0 {
		return nil
	}

	return usersList
}

func newUserCollection() iCollection {
	return &userCollection{}
}
