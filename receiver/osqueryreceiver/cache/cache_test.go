// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewCacheManager(t *testing.T) {
	cm := newCacheManager(zap.NewNop())
	assert.Equal(t, 0, cm.GetCacheSize())
}

func TestUpdateCache(t *testing.T) {
	tests := map[string]struct {
		operations func(cm *cacheManager)
		wantSize   int
	}{
		"single entry": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("system_info", map[string]any{"hostname": "host1"})
			},
			wantSize: 1,
		},
		"multiple distinct entries": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("system_info", map[string]any{"hostname": "host1"})
				cm.UpdateCache("os_info", map[string]any{"name": "linux"})
			},
			wantSize: 2,
		},
		"overwrite same key": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("system_info", "first-value")
				cm.UpdateCache("system_info", "second-value")
			},
			wantSize: 1,
		},
		"nil data stored": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("system_info", nil)
			},
			wantSize: 1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cm := newCacheManager(zap.NewNop())
			tc.operations(&cm)
			assert.Equal(t, tc.wantSize, cm.GetCacheSize())
		})
	}
}

func TestUpdateCacheSetsValidAndTTL(t *testing.T) {
	cm := newCacheManager(zap.NewNop())
	before := time.Now()
	(&cm).UpdateCache("system_info", "data")
	after := time.Now()

	cm.cacheMutex.RLock()
	entry := cm.cache["system_info"]
	cm.cacheMutex.RUnlock()

	assert.True(t, entry.IsValid)
	assert.Equal(t, 5*time.Minute, entry.TTL)
	assert.False(t, entry.CachedAt.Before(before))
	assert.False(t, entry.CachedAt.After(after))
}

func TestGetCachedResult(t *testing.T) {
	tests := map[string]struct {
		setup    func(cm *cacheManager)
		key      string
		wantData any
		wantOk   bool
	}{
		"valid cache hit": {
			setup: func(cm *cacheManager) {
				cm.UpdateCache("system_info", "some-data")
			},
			key:      "system_info",
			wantData: "some-data",
			wantOk:   true,
		},
		"valid cache hit with map data": {
			setup: func(cm *cacheManager) {
				cm.UpdateCache("os_info", map[string]any{"name": "linux", "version": "22.04"})
			},
			key:      "os_info",
			wantData: map[string]any{"name": "linux", "version": "22.04"},
			wantOk:   true,
		},
		"cache miss - key does not exist": {
			setup:    func(_ *cacheManager) {},
			key:      "nonexistent",
			wantData: nil,
			wantOk:   false,
		},
		"expired cache entry": {
			setup: func(cm *cacheManager) {
				cm.cacheMutex.Lock()
				cm.cache["system_info"] = cachedResult{
					Data:     "old-data",
					CachedAt: time.Now().Add(-10 * time.Minute),
					TTL:      5 * time.Minute,
					IsValid:  true,
				}
				cm.cacheMutex.Unlock()
			},
			key:      "system_info",
			wantData: nil,
			wantOk:   false,
		},
		"invalidated cache entry": {
			setup: func(cm *cacheManager) {
				cm.UpdateCache("system_info", "some-data")
				cm.InvalidateCache("system_info")
			},
			key:      "system_info",
			wantData: nil,
			wantOk:   false,
		},
		"just-expired entry (exactly at TTL boundary)": {
			setup: func(cm *cacheManager) {
				cm.cacheMutex.Lock()
				cm.cache["system_info"] = cachedResult{
					Data:     "data",
					CachedAt: time.Now().Add(-(5*time.Minute + time.Millisecond)),
					TTL:      5 * time.Minute,
					IsValid:  true,
				}
				cm.cacheMutex.Unlock()
			},
			key:      "system_info",
			wantData: nil,
			wantOk:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cm := newCacheManager(zap.NewNop())
			tc.setup(&cm)
			data, ok := cm.GetCachedResult(tc.key)
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.wantData, data)
		})
	}
}

func TestInvalidateCache(t *testing.T) {
	tests := map[string]struct {
		setup           func(cm *cacheManager)
		key             string
		wantGetResultOk bool
		wantSizeAfter   int
	}{
		"invalidate existing valid entry": {
			setup: func(cm *cacheManager) {
				cm.UpdateCache("system_info", "data")
			},
			key:             "system_info",
			wantGetResultOk: false,
			wantSizeAfter:   1, // entry still in map, just marked invalid
		},
		"invalidate non-existent entry - no panic": {
			setup:           func(_ *cacheManager) {},
			key:             "nonexistent",
			wantGetResultOk: false,
			wantSizeAfter:   0,
		},
		"invalidate already-invalidated entry": {
			setup: func(cm *cacheManager) {
				cm.UpdateCache("system_info", "data")
				cm.InvalidateCache("system_info")
			},
			key:             "system_info",
			wantGetResultOk: false,
			wantSizeAfter:   1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cm := newCacheManager(zap.NewNop())
			tc.setup(&cm)
			cm.InvalidateCache(tc.key)
			_, ok := cm.GetCachedResult(tc.key)
			assert.Equal(t, tc.wantGetResultOk, ok)
			assert.Equal(t, tc.wantSizeAfter, cm.GetCacheSize())
		})
	}
}

func TestGetCacheSize(t *testing.T) {
	tests := map[string]struct {
		operations func(cm *cacheManager)
		wantSize   int
	}{
		"empty cache": {
			operations: func(_ *cacheManager) {},
			wantSize:   0,
		},
		"after one update": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("key1", "data")
			},
			wantSize: 1,
		},
		"after two updates with different keys": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("key1", "data1")
				cm.UpdateCache("key2", "data2")
			},
			wantSize: 2,
		},
		"after overwriting same key": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("key1", "data1")
				cm.UpdateCache("key1", "data2")
			},
			wantSize: 1,
		},
		"size unchanged after invalidation": {
			operations: func(cm *cacheManager) {
				cm.UpdateCache("key1", "data")
				cm.InvalidateCache("key1")
			},
			wantSize: 1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cm := newCacheManager(zap.NewNop())
			tc.operations(&cm)
			assert.Equal(t, tc.wantSize, cm.GetCacheSize())
		})
	}
}
