// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cache // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/cache"

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type cacheManager struct {
	cache      map[string]cachedResult
	cacheMutex sync.RWMutex
	logger     *zap.Logger
}

type cachedResult struct {
	Data     any
	CachedAt time.Time
	TTL      time.Duration
	IsValid  bool
}

func newCacheManager(logger *zap.Logger) cacheManager {
	return cacheManager{
		cache:  make(map[string]cachedResult),
		logger: logger,
	}
}

// UpdateCache stores collection results in cache
func (m *cacheManager) UpdateCache(collectionName string, data any) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache[collectionName] = cachedResult{
		Data:     data,
		CachedAt: time.Now(),
		TTL:      5 * time.Minute, // Configurable TTL
		IsValid:  true,
	}

	m.logger.Debug("Updated cache", zap.String("collection", collectionName))
}

// GetCachedResult retrieves cached collection result if valid
func (m *cacheManager) GetCachedResult(collectionName string) (any, bool) {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	cached, exists := m.cache[collectionName]
	if !exists || !cached.IsValid {
		return nil, false
	}

	// Check if cache is expired
	if time.Since(cached.CachedAt) > cached.TTL {
		return nil, false
	}

	return cached.Data, true
}

func (m *cacheManager) InvalidateCache(collectionName string) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	if cached, exists := m.cache[collectionName]; exists {
		cached.IsValid = false
		m.cache[collectionName] = cached
		m.logger.Debug("Invalidated cache", zap.String("collection", collectionName))
	}
}

func (m *cacheManager) GetCacheSize() int {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()
	return len(m.cache)
}
