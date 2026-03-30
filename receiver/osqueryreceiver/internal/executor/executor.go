// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package executor // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/executor"

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/statemanager"
)

// collectionQuerier is a local interface for collection objects.
// Concrete types satisfying this interface are created by the collection package.
type collectionQuerier interface {
	GetName() string
	GetQuery() string
	Unmarshal(any) any
}

// stateStore is a local interface for state manager objects.
type stateStore interface {
	Save(collectionName string, data any)
	Retrieve(collectionName string) any
	ComputeDiff([]byte, []byte) []byte
}

// Runner is the exported interface for a collection executor used by the parent package.
type Runner interface {
	Register(any)
	CollectionCount() int
	ExecuteAll() map[string]CollectionResult
}

type collectionExecutor struct {
	logger       *zap.Logger
	collections  []collectionQuerier
	stateManager stateStore
}

// New creates a new collection executor.
func New(logger *zap.Logger, tmpDir string) Runner {
	return &collectionExecutor{
		logger:      logger,
		collections: []collectionQuerier{},
		// Using in-memory state manager for simplicity; will be replaced with persistent one later
		// stateManager: statemanager.New("inmemory", logger, "").(stateStore),
		stateManager: statemanager.New("file", logger, tmpDir).(stateStore),
	}
}

// Register adds a collection to the executor. The value must implement GetName, GetQuery, and Unmarshal.
func (e *collectionExecutor) Register(c any) {
	if q, ok := c.(collectionQuerier); ok {
		e.collections = append(e.collections, q)
	}
}

// CollectionCount returns the number of registered collections.
func (e *collectionExecutor) CollectionCount() int {
	return len(e.collections)
}

func (e *collectionExecutor) executeAll() map[string]queryExecution {
	results := make(map[string]queryExecution)
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		stateMu sync.Mutex
	)
	for _, coll := range e.collections {
		wg.Add(1)
		go func(coll collectionQuerier) {
			defer wg.Done()

			collectionName := coll.GetName()
			e.logger.Info("Executing collection", zap.String("collection", collectionName))

			query := coll.GetQuery()
			data, err := e.run(query)
			execResult := queryExecution{
				Query:      query,
				ExecutedAt: time.Now(),
			}

			if err != nil {
				e.logger.Error("Failed to execute query", zap.String("query", query), zap.Error(err))
				execResult.Error = err
				mu.Lock()
				results[collectionName] = execResult
				mu.Unlock()
				return
			}

			transformed := coll.Unmarshal(data)
			execResult.TransformInto = transformed
			execResult.State = transformed

			previousState := e.getCollectionState(collectionName)
			changedRows, hasChange := computeChanges(previousState, transformed)

			e.logger.Info("Collection execution completed", zap.String("collection", collectionName))

			if !hasChange {
				e.logger.Debug("No state change detected", zap.String("collection", collectionName))
				stateMu.Lock()
				e.updateCollectionState(collectionName, transformed)
				stateMu.Unlock()
				return
			}

			e.logStateChange(collectionName, previousState, transformed, changedRows)
			stateMu.Lock()
			e.updateCollectionState(collectionName, transformed)
			stateMu.Unlock()

			if changedRows == nil {
				return
			}

			execResult.TransformInto = changedRows
			execResult.ResultCount = countRecords(changedRows)

			mu.Lock()
			results[collectionName] = execResult
			mu.Unlock()
		}(coll)
	}
	wg.Wait()
	return results
}

// ExecuteAll runs all registered collections and returns the results.
func (e *collectionExecutor) ExecuteAll() map[string]CollectionResult {
	internal := e.executeAll()
	results := make(map[string]CollectionResult, len(internal))
	for k, v := range internal {
		results[k] = CollectionResult{
			Query:         v.Query,
			ExecutedAt:    v.ExecutedAt,
			ResultCount:   v.ResultCount,
			TransformInto: v.TransformInto,
			Error:         v.Error,
		}
	}
	return results
}

func (e *collectionExecutor) getCollectionState(collectionName string) any {
	return e.stateManager.Retrieve(collectionName)
}

func (e *collectionExecutor) updateCollectionState(collectionName string, latest any) {
	e.stateManager.Save(collectionName, latest)
}

func (e *collectionExecutor) logStateChange(collectionName string, previous, current, changed any) {
	payload := map[string]any{
		"collection":  collectionName,
		"changedRows": changed,
	}
	if previous != nil {
		payload["previousState"] = previous
	}
	if current != nil {
		payload["currentState"] = current
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		e.logger.Debug("State change detected", zap.String("collection", collectionName), zap.Any("changedRows", changed))
		return
	}

	e.logger.Debug("State change detected", zap.ByteString("state_change", encoded))
}

func computeChanges(previous, current any) (any, bool) {
	if current == nil {
		if previous != nil {
			return nil, true
		}
		return nil, false
	}

	currentValue := reflect.ValueOf(current)
	if currentValue.Kind() == reflect.Slice {
		changeSet := reflect.MakeSlice(currentValue.Type(), 0, currentValue.Len())
		previousKeys := make(map[string]struct{})

		if previous != nil {
			previousValue := reflect.ValueOf(previous)
			if previousValue.Kind() == reflect.Slice {
				for i := 0; i < previousValue.Len(); i++ {
					key := comparableValue(previousValue.Index(i).Interface())
					previousKeys[key] = struct{}{}
				}
			}
		}

		for i := 0; i < currentValue.Len(); i++ {
			elem := currentValue.Index(i)
			key := comparableValue(elem.Interface())
			if _, exists := previousKeys[key]; !exists {
				changeSet = reflect.Append(changeSet, elem)
			}
		}

		if changeSet.Len() == 0 {
			return nil, false
		}

		return changeSet.Interface(), true
	}

	if previous != nil && reflect.DeepEqual(previous, current) {
		return nil, false
	}

	return current, true
}

func comparableValue(value any) string {
	if value == nil {
		return ""
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}

	return string(encoded)
}

func countRecords(data any) int {
	if data == nil {
		return 0
	}

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Slice {
		return value.Len()
	}

	return 1
}

func (e *collectionExecutor) run(query string) (any, error) {
	e.logger.Debug("Executing osquery query", zap.String("query", query))
	cmd := exec.Command("osqueryi", "--json", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	e.logger.Debug("Osquery query executed successfully", zap.String("query", query))
	var outputData any
	if err := json.Unmarshal(output, &outputData); err != nil {
		e.logger.Error("Failed to unmarshal osquery output", zap.Error(err))
		return nil, err
	}
	e.logger.Debug("Unmarshalled osquery output", zap.Any("output", outputData))
	outputDataSlice, ok := outputData.([]any)
	if !ok {
		e.logger.Error("Failed to convert osquery output to slice", zap.String("query", query))
		return nil, nil
	}
	if len(outputDataSlice) == 0 {
		e.logger.Warn("No results returned from osquery", zap.String("query", query))
		return nil, nil
	}

	// Convert []any to []map[string]any for collections to unmarshal
	resultMaps := make([]map[string]any, 0, len(outputDataSlice))
	for _, item := range outputDataSlice {
		if itemMap, ok := item.(map[string]any); ok {
			resultMaps = append(resultMaps, itemMap)
		}
	}

	// If only one result, return as single map for single-row collections (like system_info)
	// Otherwise return as slice for multi-row collections (like package_info)
	if len(resultMaps) == 1 {
		return resultMaps[0], nil
	}
	return resultMaps, nil
}
