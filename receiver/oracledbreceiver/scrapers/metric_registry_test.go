// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package scrapers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/newrelic/nrdot-collector-components/receiver/oracledbreceiver/internal/metadata"
)

func TestNewSystemMetricRegistry(t *testing.T) {
	registry := NewSystemMetricRegistry()

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.recorders)
	assert.NotEmpty(t, registry.recorders)
}

func TestSystemMetricRegistry_RecordMetric_Success(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemBufferCacheHitRatio.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))
	success := registry.RecordMetric(mb, ts, "Buffer Cache Hit Ratio", 95.5, "instance1")

	assert.True(t, success)
}

func TestSystemMetricRegistry_RecordMetric_UnknownMetric(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))
	success := registry.RecordMetric(mb, ts, "Unknown Metric Name", 100.0, "instance1")

	assert.False(t, success)
}

func TestSystemMetricRegistry_RecordMultipleMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemBufferCacheHitRatio.Enabled = true
	config.Metrics.OracledbSystemMemorySortsRatio.Enabled = true
	config.Metrics.OracledbSystemRedoAllocationHitRatio.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	metrics := []struct {
		name  string
		value float64
	}{
		{"Buffer Cache Hit Ratio", 95.5},
		{"Memory Sorts Ratio", 98.2},
		{"Redo Allocation Hit Ratio", 99.9},
	}

	for _, metric := range metrics {
		success := registry.RecordMetric(mb, ts, metric.name, metric.value, "instance1")
		assert.True(t, success, "Failed to record metric: %s", metric.name)
	}
}

func TestSystemMetricRegistry_CacheRatioMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemCursorCacheHitRatio.Enabled = true
	config.Metrics.OracledbSystemSoftParseRatio.Enabled = true
	config.Metrics.OracledbSystemRowCacheHitRatio.Enabled = true
	config.Metrics.OracledbSystemLibraryCacheHitRatio.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Cursor Cache Hit Ratio",
		"Soft Parse Ratio",
		"Row Cache Hit Ratio",
		"Library Cache Hit Ratio",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 90.0, "instance1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestSystemMetricRegistry_TransactionMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemTransactionsPerSecond.Enabled = true
	config.Metrics.OracledbSystemPhysicalReadsPerTransaction.Enabled = true
	config.Metrics.OracledbSystemPhysicalWritesPerTransaction.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"User Transaction Per Sec",
		"Physical Reads Per Txn",
		"Physical Writes Per Txn",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 50.0, "instance1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestSystemMetricRegistry_PerSecondMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemPhysicalReadsDirectPerSecond.Enabled = true
	config.Metrics.OracledbSystemOpenCursorsPerSecond.Enabled = true
	config.Metrics.OracledbSystemUserCommitsPerSecond.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Physical Reads Direct Per Sec",
		"Open Cursors Per Sec",
		"User Commits Per Sec",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 100.0, "instance1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestSystemMetricRegistry_SingleValueMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemRowsPerSort.Enabled = true
	config.Metrics.OracledbSystemHostCPUUtilization.Enabled = true
	config.Metrics.OracledbSystemCurrentLogonsCount.Enabled = true
	config.Metrics.OracledbSystemSessionCount.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Rows Per Sort",
		"Host CPU Utilization (%)",
		"Current Logons Count",
		"Session Count",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 50.0, "instance1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestNewPdbMetricRegistry(t *testing.T) {
	registry := NewPdbMetricRegistry()

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.recorders)
	assert.NotEmpty(t, registry.recorders)
}

func TestPdbMetricRegistry_RecordMetric_Success(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbActiveParallelSessions.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))
	success := registry.RecordMetric(mb, ts, "Active Parallel Sessions", 10.0, "instance1", "pdb1")

	assert.True(t, success)
}

func TestPdbMetricRegistry_RecordMetric_UnknownMetric(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))
	success := registry.RecordMetric(mb, ts, "Unknown PDB Metric", 100.0, "instance1", "pdb1")

	assert.False(t, success)
}

func TestPdbMetricRegistry_RecordMultipleMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbActiveParallelSessions.Enabled = true
	config.Metrics.OracledbPdbActiveSerialSessions.Enabled = true
	config.Metrics.OracledbPdbAverageActiveSessions.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	metrics := []struct {
		name  string
		value float64
	}{
		{"Active Parallel Sessions", 5.0},
		{"Active Serial Sessions", 15.0},
		{"Average Active Sessions", 20.0},
	}

	for _, metric := range metrics {
		success := registry.RecordMetric(mb, ts, metric.name, metric.value, "instance1", "pdb1")
		assert.True(t, success, "Failed to record metric: %s", metric.name)
	}
}

func TestPdbMetricRegistry_CPUMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbBackgroundCPUUsagePerSecond.Enabled = true
	config.Metrics.OracledbPdbCPUUsagePerSecond.Enabled = true
	config.Metrics.OracledbPdbCPUUsagePerTransaction.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Background CPU Usage Per Sec",
		"CPU Usage Per Sec",
		"CPU Usage Per Txn",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 75.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_DatabaseMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbCPUTimeRatio.Enabled = true
	config.Metrics.OracledbPdbWaitTimeRatio.Enabled = true
	config.Metrics.OracledbPdbBlockChangesPerSecond.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Database CPU Time Ratio",
		"Database Wait Time Ratio",
		"DB Block Changes Per Sec",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 80.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_TransactionMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbExecutionsPerSecond.Enabled = true
	config.Metrics.OracledbPdbExecutionsPerTransaction.Enabled = true
	config.Metrics.OracledbPdbTransactionsPerSecond.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Executions Per Sec",
		"Executions Per Txn",
		"User Transaction Per Sec",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 100.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_ParseMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbHardParseCountPerSecond.Enabled = true
	config.Metrics.OracledbPdbHardParseCountPerTransaction.Enabled = true
	config.Metrics.OracledbPdbSoftParseRatio.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Hard Parse Count Per Sec",
		"Hard Parse Count Per Txn",
		"Soft Parse Ratio",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 25.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_IOMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbPhysicalReadBytesPerSecond.Enabled = true
	config.Metrics.OracledbPdbPhysicalReadsPerTransaction.Enabled = true
	config.Metrics.OracledbPdbPhysicalWriteBytesPerSecond.Enabled = true
	config.Metrics.OracledbPdbPhysicalWritesPerTransaction.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Physical Read Total Bytes Per Sec",
		"Physical Reads Per Txn",
		"Physical Write Total Bytes Per Sec",
		"Physical Writes Per Txn",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 1024.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_UserMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbUserCommitsPerSecond.Enabled = true
	config.Metrics.OracledbPdbUserCommitsPercentage.Enabled = true
	config.Metrics.OracledbPdbUserRollbacksPerSecond.Enabled = true
	config.Metrics.OracledbPdbUserRollbacksPercentage.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"User Commits Per Sec",
		"User Commits Percentage",
		"User Rollbacks Per Sec",
		"User Rollbacks Percentage",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 50.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestSystemMetricRegistry_AllCacheMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()

	cacheMetrics := []string{
		"Buffer Cache Hit Ratio",
		"Cursor Cache Hit Ratio",
		"Row Cache Hit Ratio",
		"Row Cache Miss Ratio",
		"Library Cache Hit Ratio",
		"Library Cache Miss Ratio",
	}

	for _, metricName := range cacheMetrics {
		_, exists := registry.recorders[metricName]
		assert.True(t, exists, "Metric not registered: %s", metricName)
	}
}

func TestSystemMetricRegistry_AllEnqueueMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()

	enqueueMetrics := []string{
		"Enqueue Timeouts Per Txn",
		"Enqueue Waits Per Txn",
		"Enqueue Deadlocks Per Txn",
		"Enqueue Requests Per Txn",
		"Enqueue Timeouts Per Sec",
		"Enqueue Waits Per Sec",
		"Enqueue Deadlocks Per Sec",
		"Enqueue Requests Per Sec",
	}

	for _, metricName := range enqueueMetrics {
		_, exists := registry.recorders[metricName]
		assert.True(t, exists, "Metric not registered: %s", metricName)
	}
}

func TestPdbMetricRegistry_AllRegisteredMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()

	// Verify we have the expected number of PDB metrics registered
	assert.Greater(t, len(registry.recorders), 40, "Expected more than 40 PDB metrics registered")
}

func TestSystemMetricRegistry_AllRegisteredMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()

	// Verify we have the expected number of system metrics registered
	assert.Greater(t, len(registry.recorders), 100, "Expected more than 100 system metrics registered")
}

func TestSystemMetricRegistry_PerUserCallMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemDbBlockChangesPerUserCall.Enabled = true
	config.Metrics.OracledbSystemDbBlockGetsPerUserCall.Enabled = true
	config.Metrics.OracledbSystemExecutionsPerUserCall.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"DB Block Changes Per User Call",
		"DB Block Gets Per User Call",
		"Executions Per User Call",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 10.0, "instance1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_SessionAndCursorMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbCurrentLogons.Enabled = true
	config.Metrics.OracledbPdbCurrentOpenCursors.Enabled = true
	config.Metrics.OracledbPdbSessionCount.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Current Logons Count",
		"Current Open Cursors Count",
		"Session Count",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 50.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestSystemMetricRegistry_GlobalCacheMetrics(t *testing.T) {
	registry := NewSystemMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbSystemGlobalCacheAverageCrGetTime.Enabled = true
	config.Metrics.OracledbSystemGlobalCacheAverageCurrentGetTime.Enabled = true
	config.Metrics.OracledbSystemGlobalCacheBlocksCorrupted.Enabled = true
	config.Metrics.OracledbSystemGlobalCacheBlocksLost.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Global Cache Average CR Get Time",
		"Global Cache Average Current Get Time",
		"Global Cache Blocks Corrupted",
		"Global Cache Blocks Lost",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 5.0, "instance1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}

func TestPdbMetricRegistry_RedoAndLogMetrics(t *testing.T) {
	registry := NewPdbMetricRegistry()
	config := metadata.DefaultMetricsBuilderConfig()
	config.Metrics.OracledbPdbRedoGeneratedBytesPerSecond.Enabled = true
	config.Metrics.OracledbPdbRedoGeneratedBytesPerTransaction.Enabled = true
	config.Metrics.OracledbPdbLogonsPerSecond.Enabled = true
	config.Metrics.OracledbPdbLogonsPerTransaction.Enabled = true
	settings := receivertest.NewNopSettings(metadata.Type)
	mb := metadata.NewMetricsBuilder(config, settings)

	ts := pcommon.NewTimestampFromTime(time.Unix(1234567890, 0))

	tests := []string{
		"Redo Generated Per Sec",
		"Redo Generated Per Txn",
		"Logons Per Sec",
		"Logons Per Txn",
	}

	for _, metricName := range tests {
		success := registry.RecordMetric(mb, ts, metricName, 2048.0, "instance1", "pdb1")
		assert.True(t, success, "Failed to record: %s", metricName)
	}
}
