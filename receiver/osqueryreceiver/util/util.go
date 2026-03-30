// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package util // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/util"

import "strconv"

func contains[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func getString(resultMap map[string]any, key string) string {
	if val, ok := resultMap[key].(string); ok {
		return val
	}
	return ""
}

func getInt(resultMap map[string]any, key string) int {
	switch val := resultMap[key].(type) {
	case float64:
		return int(val)
	case int:
		return val
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return 0
}

func getInt64(resultMap map[string]any, key string) int64 {
	switch val := resultMap[key].(type) {
	case float64:
		return int64(val)
	case int64:
		return val
	case int:
		return int64(val)
	case string:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return 0
}
