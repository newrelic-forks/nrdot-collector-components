// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

import (
	"strconv"
	"strings"
)

// sanitizeRow trims string values and removes empty/zero values for the provided keys.
// It ensures that only meaningful data is kept so state comparisons ignore empty fields.
// It is not scalable, and needs to be changed, but leaving it for now as it is POC 😎
func sanitizeRow(resultMap map[string]any, stringKeys, intKeys, int64Keys []string) map[string]any {
	sanitized := make(map[string]any)

	for _, key := range stringKeys {
		value := strings.TrimSpace(getStringVal(resultMap, key))
		if value != "" {
			sanitized[key] = value
		}
	}

	for _, key := range intKeys {
		value := getIntVal(resultMap, key)
		if value > 0 {
			sanitized[key] = float64(value)
		}
	}

	for _, key := range int64Keys {
		value := getInt64Val(resultMap, key)
		if value > 0 {
			sanitized[key] = float64(value)
		}
	}

	return sanitized
}

func getStringVal(resultMap map[string]any, key string) string {
	if val, ok := resultMap[key].(string); ok {
		return val
	}
	return ""
}

func getIntVal(resultMap map[string]any, key string) int {
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

func getInt64Val(resultMap map[string]any, key string) int64 {
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
