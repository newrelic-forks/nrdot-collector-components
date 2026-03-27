// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statemanager

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInMemoryStateManagerSave(t *testing.T) {
	tests := map[string]struct {
		collectionName string
		data           any
	}{
		"save string data": {
			collectionName: "system_info",
			data:           "some-state",
		},
		"save map data": {
			collectionName: "os_info",
			data:           map[string]any{"name": "linux", "version": "22.04"},
		},
		"save nil data": {
			collectionName: "empty_collection",
			data:           nil,
		},
		"save slice data": {
			collectionName: "package_info",
			data:           []map[string]any{{"name": "pkg1"}, {"name": "pkg2"}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sm := newInMemoryStateManager(zap.NewNop())
			sm.Save(tc.collectionName, tc.data)
			result := sm.Retrieve(tc.collectionName)
			assert.Equal(t, tc.data, result)
		})
	}
}

func TestInMemoryStateManagerRetrieve(t *testing.T) {
	tests := map[string]struct {
		setup          func(sm iStateManager)
		collectionName string
		wantResult     any
	}{
		"retrieve existing key": {
			setup: func(sm iStateManager) {
				sm.Save("system_info", map[string]any{"hostname": "host1"})
			},
			collectionName: "system_info",
			wantResult:     map[string]any{"hostname": "host1"},
		},
		"retrieve non-existing key returns nil": {
			setup:          func(_ iStateManager) {},
			collectionName: "nonexistent",
			wantResult:     nil,
		},
		"retrieve after overwrite returns latest value": {
			setup: func(sm iStateManager) {
				sm.Save("system_info", "first-value")
				sm.Save("system_info", "second-value")
			},
			collectionName: "system_info",
			wantResult:     "second-value",
		},
		"retrieve one of multiple keys": {
			setup: func(sm iStateManager) {
				sm.Save("key1", "value1")
				sm.Save("key2", "value2")
			},
			collectionName: "key1",
			wantResult:     "value1",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sm := newInMemoryStateManager(zap.NewNop())
			tc.setup(sm)
			result := sm.Retrieve(tc.collectionName)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestInMemoryStateManagerComputeDiff(t *testing.T) {
	tests := map[string]struct {
		original  any
		modified  any
		checkDiff func(t *testing.T, diff []byte)
	}{
		"identical objects produce empty patch": {
			original: map[string]any{"a": 1},
			modified: map[string]any{"a": 1},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Empty(t, result)
			},
		},
		"changed value produces patch with new value": {
			original: map[string]any{"a": 1},
			modified: map[string]any{"a": 2},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Contains(t, result, "a")
			},
		},
		"added field produces patch with new field": {
			original: map[string]any{"a": 1},
			modified: map[string]any{"a": 1, "b": 2},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Contains(t, result, "b")
			},
		},
		"empty to non-empty produces full patch": {
			original: map[string]any{},
			modified: map[string]any{"a": 1},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Contains(t, result, "a")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sm := newInMemoryStateManager(zap.NewNop())

			origBytes, err := json.Marshal(tc.original)
			require.NoError(t, err)
			modBytes, err := json.Marshal(tc.modified)
			require.NoError(t, err)

			diff := sm.ComputeDiff(origBytes, modBytes)
			require.NotNil(t, diff)
			tc.checkDiff(t, diff)
		})
	}
}
