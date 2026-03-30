// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statemanager // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/statemanager"

import (
	"go.uber.org/zap"
)

type iStateManager interface {
	Save(collectionName string, data any)
	Retrieve(collectionName string) any
	ComputeDiff([]byte, []byte) []byte
}

type stateManager struct {
	logger *zap.Logger
}

func getStateManager(managerType string, logger *zap.Logger, fileLocation string) iStateManager {
	switch managerType {
	case "inmemory":
		return newInMemoryStateManager(logger)
	case "file":
		return newFileStateManager(logger, fileLocation)
	default:
		return newInMemoryStateManager(logger)
	}
}

// New returns a state manager of the given type. The concrete type implements Save, Retrieve, and ComputeDiff.
func New(managerType string, logger *zap.Logger, fileLocation string) any {
	return getStateManager(managerType, logger, fileLocation)
}
