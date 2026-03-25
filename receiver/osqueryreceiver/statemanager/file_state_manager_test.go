// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statemanager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestFileStateManagerSave(t *testing.T) {
	tests := map[string]struct {
		collectionName string
		data           any
		wantFileData   any
	}{
		"save map data": {
			collectionName: "system_info",
			data:           map[string]any{"hostname": "host1"},
			wantFileData:   map[string]any{"hostname": "host1"},
		},
		"save string data": {
			collectionName: "os_info",
			data:           "linux",
			wantFileData:   "linux",
		},
		"save slice data": {
			collectionName: "packages",
			data:           []any{map[string]any{"name": "pkg1"}, map[string]any{"name": "pkg2"}},
			wantFileData:   []any{map[string]any{"name": "pkg1"}, map[string]any{"name": "pkg2"}},
		},
		"overwrite existing entry": {
			collectionName: "system_info",
			data:           map[string]any{"hostname": "host2"},
			wantFileData:   map[string]any{"hostname": "host2"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			sm := NewFileStateManager(zap.NewNop(), dir+"/")

			sm.Save(tc.collectionName, tc.data)

			// Verify the file was written
			filePath := filepath.Join(dir, tc.collectionName+".json")
			fileData, err := os.ReadFile(filePath)
			require.NoError(t, err)

			var stored any
			require.NoError(t, json.Unmarshal(fileData, &stored))
			assert.Equal(t, tc.wantFileData, stored)
		})
	}
}

func TestFileStateManagerRetrieve(t *testing.T) {
	tests := map[string]struct {
		setup          func(dir string, sm IStateManager)
		collectionName string
		wantResult     any
	}{
		"retrieve existing saved data": {
			setup: func(dir string, sm IStateManager) {
				sm.Save("system_info", map[string]any{"hostname": "host1"})
			},
			collectionName: "system_info",
			wantResult:     map[string]any{"hostname": "host1"},
		},
		"retrieve non-existing file returns nil": {
			setup:          func(dir string, sm IStateManager) {},
			collectionName: "nonexistent",
			wantResult:     nil,
		},
		"retrieve after overwrite returns latest value": {
			setup: func(dir string, sm IStateManager) {
				sm.Save("system_info", map[string]any{"hostname": "host1"})
				sm.Save("system_info", map[string]any{"hostname": "host2"})
			},
			collectionName: "system_info",
			wantResult:     map[string]any{"hostname": "host2"},
		},
		"retrieve empty file returns nil": {
			setup: func(dir string, sm IStateManager) {
				filePath := filepath.Join(dir, "empty_collection.json")
				require.NoError(t, os.WriteFile(filePath, []byte{}, 0644))
			},
			collectionName: "empty_collection",
			wantResult:     nil,
		},
		"retrieve file with invalid JSON returns nil": {
			setup: func(dir string, sm IStateManager) {
				filePath := filepath.Join(dir, "bad_collection.json")
				require.NoError(t, os.WriteFile(filePath, []byte("not-valid-json{{{"), 0644))
			},
			collectionName: "bad_collection",
			wantResult:     nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			sm := NewFileStateManager(zap.NewNop(), dir+"/")
			tc.setup(dir, sm)

			result := sm.Retrieve(tc.collectionName)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}

func TestFileStateManagerComputeDiff(t *testing.T) {
	tests := map[string]struct {
		original  any
		modified  any
		checkDiff func(t *testing.T, diff []byte)
	}{
		"identical objects produce empty patch": {
			original: map[string]any{"a": 1.0},
			modified: map[string]any{"a": 1.0},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Empty(t, result)
			},
		},
		"changed value produces patch": {
			original: map[string]any{"hostname": "host1"},
			modified: map[string]any{"hostname": "host2"},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Equal(t, "host2", result["hostname"])
			},
		},
		"added key produces patch with new key": {
			original: map[string]any{"a": 1.0},
			modified: map[string]any{"a": 1.0, "b": 2.0},
			checkDiff: func(t *testing.T, diff []byte) {
				var result map[string]any
				require.NoError(t, json.Unmarshal(diff, &result))
				assert.Contains(t, result, "b")
				assert.NotContains(t, result, "a")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			sm := NewFileStateManager(zap.NewNop(), dir+"/")

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

func TestNewFileStateManagerCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "nested", "state") + "/"

	// Should create the directory without error
	sm := NewFileStateManager(zap.NewNop(), subDir)
	assert.NotNil(t, sm)

	_, err := os.Stat(filepath.Join(dir, "nested", "state"))
	assert.NoError(t, err)
}
