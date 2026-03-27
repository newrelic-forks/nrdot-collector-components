// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statemanager // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/statemanager"

import (
	"encoding/json"
	"os"
	"sync"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"go.uber.org/zap"
)

type fileStateManager struct {
	stateManager
	baseDir string
	locks   sync.Map // map[collectionName]*sync.RWMutex
}

func newFileStateManager(logger *zap.Logger, baseDir string) iStateManager {
	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		logger.Error("Failed to create state directory", zap.String("baseDir", baseDir), zap.Error(err))
	}
	return &fileStateManager{
		stateManager: stateManager{
			logger: logger,
		},
		baseDir: baseDir,
	}
}

func (fsm *fileStateManager) getLock(collectionName string) *sync.RWMutex {
	lock, _ := fsm.locks.LoadOrStore(collectionName, &sync.RWMutex{})
	return lock.(*sync.RWMutex)
}

func (fsm *fileStateManager) getFilePath(collectionName string) string {
	return fsm.baseDir + collectionName + ".json"
}

func (fsm *fileStateManager) Save(collectionName string, data any) {
	lock := fsm.getLock(collectionName)
	lock.Lock()
	defer lock.Unlock()

	filePath := fsm.getFilePath(collectionName)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		fsm.logger.Error("Failed to marshal state to JSON",
			zap.String("collection", collectionName),
			zap.String("filePath", filePath),
			zap.Error(err))
		return
	}
	// TODO: reduce security permissions when moving to beta.
	err = os.WriteFile(filePath, dataBytes, 0o644) //nolint:gosec
	if err != nil {
		fsm.logger.Error("Failed to save state to file",
			zap.String("collection", collectionName),
			zap.String("filePath", filePath),
			zap.Error(err))
	}
}

func (fsm *fileStateManager) Retrieve(collectionName string) any {
	lock := fsm.getLock(collectionName)
	lock.RLock()
	defer lock.RUnlock()

	filePath := fsm.getFilePath(collectionName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fsm.logger.Error("Failed to read state from file",
				zap.String("collection", collectionName),
				zap.String("filePath", filePath),
				zap.Error(err))
		}
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	var state any
	err = json.Unmarshal(data, &state)
	if err != nil {
		fsm.logger.Error("Failed to unmarshal state from file",
			zap.String("collection", collectionName),
			zap.String("filePath", filePath),
			zap.Error(err))
		return nil
	}
	return state
}

func (fsm *fileStateManager) ComputeDiff(originalState, newState []byte) []byte {
	patch, err := jsonpatch.CreateMergePatch(originalState, newState)
	if err != nil {
		fsm.logger.Error("Failed to compute state diff", zap.Error(err))
		return nil
	}
	fsm.logger.Debug("Computed state diff", zap.ByteString("diff", patch))
	return patch
}
