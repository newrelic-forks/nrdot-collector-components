// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statemanager // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/statemanager"

import (
	"sync"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"go.uber.org/zap"
)

type inMemoryStateManager struct {
	stateManager
	state sync.Map // map[collectionName]any
}

func newInMemoryStateManager(logger *zap.Logger) iStateManager {
	return &inMemoryStateManager{
		stateManager: stateManager{
			logger: logger,
		},
	}
}

func (ism *inMemoryStateManager) Save(collectionName string, data any) {
	ism.state.Store(collectionName, data)
}

func (ism *inMemoryStateManager) Retrieve(collectionName string) any {
	value, ok := ism.state.Load(collectionName)
	if !ok {
		return nil
	}
	return value
}

func (ism *inMemoryStateManager) ComputeDiff(originalState, newState []byte) []byte {
	patch, err := jsonpatch.CreateMergePatch(originalState, newState)
	if err != nil {
		ism.logger.Error("Failed to compute state diff", zap.Error(err))
		return nil
	}
	ism.logger.Debug("Computed state diff", zap.ByteString("diff", patch))
	return patch
}
