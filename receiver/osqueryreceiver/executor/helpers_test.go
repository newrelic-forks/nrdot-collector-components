// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeChanges(t *testing.T) {
	tests := map[string]struct {
		previous   any
		current    any
		wantChange bool
		wantNil    bool // whether returned value should be nil
	}{
		// nil current cases
		"both nil - no change": {
			previous:   nil,
			current:    nil,
			wantChange: false,
			wantNil:    true,
		},
		"current nil, previous set - has change": {
			previous:   "some-state",
			current:    nil,
			wantChange: true,
			wantNil:    true,
		},

		// slice cases
		"slice: all new items (nil previous)": {
			previous:   nil,
			current:    []string{"a", "b"},
			wantChange: true,
			wantNil:    false,
		},
		"slice: no new items (same as previous)": {
			previous:   []string{"a", "b"},
			current:    []string{"a", "b"},
			wantChange: false,
			wantNil:    true,
		},
		"slice: one new item added": {
			previous:   []string{"a"},
			current:    []string{"a", "b"},
			wantChange: true,
			wantNil:    false,
		},
		"slice: item removed (not a new item)": {
			// removed items don't show as changes; only additions do
			previous:   []string{"a", "b"},
			current:    []string{"a"},
			wantChange: false,
			wantNil:    true,
		},
		"slice: completely different items": {
			previous:   []string{"a", "b"},
			current:    []string{"c", "d"},
			wantChange: true,
			wantNil:    false,
		},
		"slice: empty current, empty previous": {
			previous:   []string{},
			current:    []string{},
			wantChange: false,
			wantNil:    true,
		},
		"slice: empty current, non-empty previous": {
			previous:   []string{"a"},
			current:    []string{},
			wantChange: false,
			wantNil:    true,
		},
		"slice: maps - all new": {
			previous: nil,
			current: []map[string]any{
				{"name": "pkg1", "version": "1.0"},
			},
			wantChange: true,
			wantNil:    false,
		},
		"slice: maps - no change": {
			previous: []map[string]any{
				{"name": "pkg1", "version": "1.0"},
			},
			current: []map[string]any{
				{"name": "pkg1", "version": "1.0"},
			},
			wantChange: false,
			wantNil:    true,
		},
		"slice: maps - new item added": {
			previous: []map[string]any{
				{"name": "pkg1", "version": "1.0"},
			},
			current: []map[string]any{
				{"name": "pkg1", "version": "1.0"},
				{"name": "pkg2", "version": "2.0"},
			},
			wantChange: true,
			wantNil:    false,
		},

		// non-slice cases
		"non-slice: nil previous, has current - change": {
			previous:   nil,
			current:    map[string]any{"hostname": "host1"},
			wantChange: true,
			wantNil:    false,
		},
		"non-slice: equal structs - no change": {
			previous:   map[string]any{"hostname": "host1"},
			current:    map[string]any{"hostname": "host1"},
			wantChange: false,
			wantNil:    true,
		},
		"non-slice: different structs - has change": {
			previous:   map[string]any{"hostname": "host1"},
			current:    map[string]any{"hostname": "host2"},
			wantChange: true,
			wantNil:    false,
		},
		"non-slice: string equal - no change": {
			previous:   "same-state",
			current:    "same-state",
			wantChange: false,
			wantNil:    true,
		},
		"non-slice: string different - has change": {
			previous:   "old-state",
			current:    "new-state",
			wantChange: true,
			wantNil:    false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, hasChange := computeChanges(tc.previous, tc.current)
			assert.Equal(t, tc.wantChange, hasChange, "hasChange mismatch")
			if tc.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestComputeChangesSliceContent(t *testing.T) {
	// Verify that only new items are returned
	previous := []string{"a", "b"}
	current := []string{"a", "b", "c"}

	result, hasChange := computeChanges(previous, current)
	assert.True(t, hasChange)

	resultSlice, ok := result.([]string)
	assert.True(t, ok)
	assert.Equal(t, []string{"c"}, resultSlice)
}

func TestComparableValue(t *testing.T) {
	tests := map[string]struct {
		input    any
		wantJSON string
	}{
		"nil value": {
			input:    nil,
			wantJSON: "",
		},
		"string value": {
			input:    "hello",
			wantJSON: `"hello"`,
		},
		"integer value": {
			input:    42,
			wantJSON: "42",
		},
		"map value": {
			input:    map[string]any{"key": "value"},
			wantJSON: `{"key":"value"}`,
		},
		"empty map": {
			input:    map[string]any{},
			wantJSON: `{}`,
		},
		"boolean true": {
			input:    true,
			wantJSON: "true",
		},
		"boolean false": {
			input:    false,
			wantJSON: "false",
		},
		"float value": {
			input:    3.14,
			wantJSON: "3.14",
		},
		"empty string": {
			input:    "",
			wantJSON: `""`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := comparableValue(tc.input)
			assert.Equal(t, tc.wantJSON, result)
		})
	}
}

func TestCountRecords(t *testing.T) {
	tests := map[string]struct {
		input    any
		wantCount int
	}{
		"nil input": {
			input:     nil,
			wantCount: 0,
		},
		"empty slice": {
			input:     []string{},
			wantCount: 0,
		},
		"slice with one element": {
			input:     []string{"a"},
			wantCount: 1,
		},
		"slice with multiple elements": {
			input:     []string{"a", "b", "c"},
			wantCount: 3,
		},
		"slice of maps": {
			input: []map[string]any{
				{"name": "pkg1"},
				{"name": "pkg2"},
			},
			wantCount: 2,
		},
		"single non-slice value (string)": {
			input:     "single-item",
			wantCount: 1,
		},
		"single non-slice value (map)": {
			input:     map[string]any{"hostname": "host1"},
			wantCount: 1,
		},
		"single non-slice value (int)": {
			input:     42,
			wantCount: 1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := countRecords(tc.input)
			assert.Equal(t, tc.wantCount, result)
		})
	}
}
