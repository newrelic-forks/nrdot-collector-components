// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package statemanager

import (
	"fmt"
	"testing"
)

func TestGetStateManager(t *testing.T) {
	tests := []struct {
		name         string
		managerType  string
		fileLocation string
		expectedType iStateManager
	}{
		{
			name:         "InMemory State Manager",
			managerType:  "inmemory",
			fileLocation: "",
			expectedType: &inMemoryStateManager{},
		},
		{
			name:         "Local State Manager",
			managerType:  "file",
			fileLocation: "/tmp/",
			expectedType: &fileStateManager{baseDir: "/tmp/"},
		},
		{
			name:         "Default to InMemory State Manager",
			managerType:  "unknown",
			fileLocation: "",
			expectedType: &inMemoryStateManager{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateManager := getStateManager(tt.managerType, nil, tt.fileLocation)
			if fmt.Sprintf("%T", stateManager) != fmt.Sprintf("%T", tt.expectedType) {
				t.Errorf("expected type %T, got %T", tt.expectedType, stateManager)
			}
		})
	}
}
