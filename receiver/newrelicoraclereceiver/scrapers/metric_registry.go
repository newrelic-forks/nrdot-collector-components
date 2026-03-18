// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package scrapers // import "github.com/newrelic/nrdot-collector-components/receiver/newrelicoraclereceiver/scrapers"

import (
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/newrelic/nrdot-collector-components/receiver/newrelicoraclereceiver/internal/metadata"
)

// MetricRecorderFunc defines a function type for recording metrics
type MetricRecorderFunc func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string)

// PdbMetricRecorderFunc defines a function type for recording PDB metrics
type PdbMetricRecorderFunc func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string)

// SystemMetricRegistry maps metric names to their recorder functions
type SystemMetricRegistry struct {
	recorders map[string]MetricRecorderFunc
}

// PdbMetricRegistry maps PDB metric names to their recorder functions
type PdbMetricRegistry struct {
	recorders map[string]PdbMetricRecorderFunc
}

// NewSystemMetricRegistry creates and initializes a system metric registry
func NewSystemMetricRegistry() *SystemMetricRegistry {
	registry := &SystemMetricRegistry{
		recorders: make(map[string]MetricRecorderFunc, 150),
	}
	registry.registerAll()
	return registry
}

// RecordMetric records a metric using the registered recorder function
func (r *SystemMetricRegistry) RecordMetric(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, metricName string, value float64, instanceID string) bool {
	recorder, exists := r.recorders[metricName]
	if !exists {
		return false
	}
	recorder(mb, ts, value, instanceID)
	return true
}

// NewPdbMetricRegistry creates and initializes a PDB metric registry
func NewPdbMetricRegistry() *PdbMetricRegistry {
	registry := &PdbMetricRegistry{
		recorders: make(map[string]PdbMetricRecorderFunc, 55),
	}
	registry.registerAll()
	return registry
}

// RecordMetric records a PDB metric using the registered recorder function
func (r *PdbMetricRegistry) RecordMetric(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, metricName string, value float64, instanceID, pdbName string) bool {
	recorder, exists := r.recorders[metricName]
	if !exists {
		return false
	}
	recorder(mb, ts, value, instanceID, pdbName)
	return true
}

// registerAll registers all system metric recorder functions
func (r *SystemMetricRegistry) registerAll() {
	// Cache and performance ratios
	r.recorders["Buffer Cache Hit Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemBufferCacheHitRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Memory Sorts Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemMemorySortsRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Redo Allocation Hit Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRedoAllocationHitRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Cursor Cache Hit Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCursorCacheHitRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Soft Parse Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemSoftParseRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["User Calls Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserCallsRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Row Cache Hit Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRowCacheHitRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Row Cache Miss Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRowCacheMissRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Library Cache Hit Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLibraryCacheHitRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Library Cache Miss Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLibraryCacheMissRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Database Wait Time Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDatabaseWaitTimeRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Database CPU Time Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDatabaseCPUTimeRatioDataPoint(ts, value, instanceID)
	}
	r.recorders["Execute Without Parse Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemExecuteWithoutParseRatioDataPoint(ts, value, instanceID)
	}

	// Transaction metrics
	r.recorders["User Transaction Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTransactionsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Reads Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Writes Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWritesPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Reads Direct Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadsDirectPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Writes Direct Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWritesDirectPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Reads Direct Lobs Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalLobsReadsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Writes Direct Lobs Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalLobsWritesPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Redo Generated Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRedoGeneratedBytesPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Logons Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLogonsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Open Cursors Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemOpenCursorsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["User Calls Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserCallsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Recursive Calls Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRecursiveCallsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Logical Reads Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLogicalReadsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Redo Writes Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRedoWritesPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Long Table Scans Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLongTableScansPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Table Scans Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalTableScansPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Full Index Scans Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemFullIndexScansPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Index Scans Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalIndexScansPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Parse Count Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalParseCountPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Hard Parse Count Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemHardParseCountPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Parse Failure Count Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemParseFailureCountPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Disk Sort Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDiskSortPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Timeouts Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueTimeoutsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Waits Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueWaitsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Deadlocks Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueDeadlocksPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Requests Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueRequestsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["DB Block Gets Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbBlockGetsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Consistent Read Gets Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemConsistentReadGetsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["DB Block Changes Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbBlockChangesPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Consistent Read Changes Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemConsistentReadChangesPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["CPU Usage Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCPUUsagePerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["CR Blocks Created Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCrBlocksCreatedPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["CR Undo Records Applied Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCrUndoRecordsAppliedPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["User Rollback Undo Records Applied Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserRollbackUndoRecordsAppliedPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Leaf Node Splits Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLeafNodeSplitsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Branch Node Splits Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemBranchNodeSplitsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["GC CR Block Received Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGcCrBlockReceivedPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["GC Current Block Received Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGcCurrentBlockReceivedPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Response Time Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemResponseTimePerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Executions Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemExecutionsPerTransactionDataPoint(ts, value, instanceID)
	}
	r.recorders["Txns Per Logon"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTransactionsPerLogonDataPoint(ts, value, instanceID)
	}

	// Per second metrics
	r.recorders["Physical Reads Direct Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadsDirectPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Writes Direct Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWritesDirectPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Reads Direct Lobs Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalLobsReadsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Writes Direct Lobs Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalLobsWritesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Redo Generated Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRedoGeneratedBytesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Open Cursors Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemOpenCursorsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["User Commits Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserCommitsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["User Commits Percentage"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserCommitsPercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["User Rollbacks Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserRollbacksPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["User Rollbacks Percentage"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserRollbacksPercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["User Calls Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserCallsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Recursive Calls Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRecursiveCallsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Logical Reads Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLogicalReadsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["DBWR Checkpoints Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbwrCheckpointsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Background Checkpoints Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemBackgroundCheckpointsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Redo Writes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRedoWritesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Long Table Scans Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLongTableScansPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Table Scans Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalTableScansPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Full Index Scans Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemFullIndexScansPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Index Scans Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalIndexScansPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Parse Count Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalParseCountPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Hard Parse Count Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemHardParseCountPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Parse Failure Count Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemParseFailureCountPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Disk Sort Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDiskSortPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Timeouts Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueTimeoutsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Waits Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueWaitsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Deadlocks Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueDeadlocksPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Enqueue Requests Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemEnqueueRequestsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["DB Block Gets Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbBlockGetsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Consistent Read Gets Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemConsistentReadGetsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["DB Block Changes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbBlockChangesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Consistent Read Changes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemConsistentReadChangesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["CPU Usage Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCPUUsagePerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["CR Blocks Created Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCrBlocksCreatedPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["CR Undo Records Applied Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCrUndoRecordsAppliedPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["User Rollback UndoRec Applied Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserRollbackUndoRecordsAppliedPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Leaf Node Splits Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLeafNodeSplitsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Branch Node Splits Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemBranchNodeSplitsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Read Total IO Requests Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadTotalIoRequestsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Read Total Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadTotalBytesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["GC CR Block Received Per Second"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGcCrBlockReceivedPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["GC Current Block Received Per Second"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGcCurrentBlockReceivedPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Write Total IO Requests Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWriteTotalIoRequestsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Write Total Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWriteTotalBytesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Write IO Requests Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWriteIoRequestsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Database Time Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDatabaseTimePerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Network Traffic Volume Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemNetworkTrafficVolumePerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Executions Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemExecutionsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Logons Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLogonsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Read Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadBytesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Read IO Requests Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadIoRequestsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Reads Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalReadsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Write Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWriteBytesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Physical Writes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPhysicalWritesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Background CPU Usage Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemBackgroundCPUUsagePerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Background Time Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemBackgroundTimePerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Host CPU Usage Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemHostCPUUsagePerSecondDataPoint(ts, value, instanceID)
	}

	// Single value metrics
	r.recorders["Rows Per Sort"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemRowsPerSortDataPoint(ts, value, instanceID)
	}
	r.recorders["Host CPU Utilization (%)"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemHostCPUUtilizationDataPoint(ts, value, instanceID)
	}
	r.recorders["Global Cache Average CR Get Time"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGlobalCacheAverageCrGetTimeDataPoint(ts, value, instanceID)
	}
	r.recorders["Global Cache Average Current Get Time"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGlobalCacheAverageCurrentGetTimeDataPoint(ts, value, instanceID)
	}
	r.recorders["Global Cache Blocks Corrupted"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGlobalCacheBlocksCorruptedDataPoint(ts, value, instanceID)
	}
	r.recorders["Global Cache Blocks Lost"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemGlobalCacheBlocksLostDataPoint(ts, value, instanceID)
	}
	r.recorders["Current Logons Count"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCurrentLogonsCountDataPoint(ts, value, instanceID)
	}
	r.recorders["Current Open Cursors Count"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCurrentOpenCursorsCountDataPoint(ts, value, instanceID)
	}
	r.recorders["User Limit %"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemUserLimitPercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["SQL Service Response Time"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemSQLServiceResponseTimeDataPoint(ts, value, instanceID)
	}
	r.recorders["Shared Pool Free %"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemSharedPoolFreePercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["PGA Cache Hit %"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemPgaCacheHitPercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["Process Limit %"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemProcessLimitPercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["Session Limit %"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemSessionLimitPercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["Temp Space Used"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTempSpaceUsedDataPoint(ts, value, instanceID)
	}
	r.recorders["Session Count"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemSessionCountDataPoint(ts, value, instanceID)
	}
	r.recorders["Captured user calls"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCapturedUserCallsDataPoint(ts, value, instanceID)
	}
	r.recorders["Current OS Load"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemCurrentOsLoadDataPoint(ts, value, instanceID)
	}
	r.recorders["Streams Pool Usage Percentage"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemStreamsPoolUsagePercentageDataPoint(ts, value, instanceID)
	}
	r.recorders["I/O Megabytes per Second"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemIoMegabytesPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["I/O Requests per Second"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemIoRequestsPerSecondDataPoint(ts, value, instanceID)
	}
	r.recorders["Average Active Sessions"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemAverageActiveSessionsDataPoint(ts, value, instanceID)
	}
	r.recorders["Active Serial Sessions"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemActiveSerialSessionsDataPoint(ts, value, instanceID)
	}
	r.recorders["Active Parallel Sessions"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemActiveParallelSessionsDataPoint(ts, value, instanceID)
	}

	// Per user call metrics
	r.recorders["DB Block Changes Per User Call"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbBlockChangesPerUserCallDataPoint(ts, value, instanceID)
	}
	r.recorders["DB Block Gets Per User Call"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemDbBlockGetsPerUserCallDataPoint(ts, value, instanceID)
	}
	r.recorders["Executions Per User Call"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemExecutionsPerUserCallDataPoint(ts, value, instanceID)
	}
	r.recorders["Logical Reads Per User Call"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemLogicalReadsPerUserCallDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Sorts Per User Call"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalSortsPerUserCallDataPoint(ts, value, instanceID)
	}
	r.recorders["Total Table Scans Per User Call"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID string) {
		mb.RecordOracledbSystemTotalTableScansPerUserCallDataPoint(ts, value, instanceID)
	}
}

// registerAll registers all PDB metric recorder functions
func (r *PdbMetricRegistry) registerAll() {
	r.recorders["Active Parallel Sessions"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbActiveParallelSessionsDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Active Serial Sessions"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbActiveSerialSessionsDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Average Active Sessions"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbAverageActiveSessionsDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Background CPU Usage Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbBackgroundCPUUsagePerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Background Time Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbBackgroundTimePerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["CPU Usage Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbCPUUsagePerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["CPU Usage Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbCPUUsagePerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Current Logons Count"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbCurrentLogonsDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Current Open Cursors Count"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbCurrentOpenCursorsDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Database CPU Time Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbCPUTimeRatioDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Database Wait Time Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbWaitTimeRatioDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["DB Block Changes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbBlockChangesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["DB Block Changes Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbBlockChangesPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Executions Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbExecutionsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Executions Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbExecutionsPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Hard Parse Count Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbHardParseCountPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Hard Parse Count Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbHardParseCountPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Logical Reads Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbLogicalReadsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Logical Reads Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbLogicalReadsPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Logons Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbLogonsPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Network Traffic Volume Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbNetworkTrafficBytePerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Open Cursors Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbOpenCursorsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Open Cursors Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbOpenCursorsPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Parse Failure Count Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbParseFailureCountPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Read Total Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbPhysicalReadBytesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Reads Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbPhysicalReadsPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Write Total Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbPhysicalWriteBytesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Writes Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbPhysicalWritesPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Redo Generated Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbRedoGeneratedBytesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Redo Generated Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbRedoGeneratedBytesPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Response Time Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbResponseTimePerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Session Count"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbSessionCountDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Soft Parse Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbSoftParseRatioDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["SQL Service Response Time"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbSQLServiceResponseTimeDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Total Parse Count Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbTotalParseCountPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Total Parse Count Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbTotalParseCountPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Calls Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbUserCallsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Calls Per Txn"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbUserCallsPerTransactionDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Commits Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbUserCommitsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Commits Percentage"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbUserCommitsPercentageDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Rollbacks Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbUserRollbacksPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Rollbacks Percentage"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbUserRollbacksPercentageDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["User Transaction Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbTransactionsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Execute Without Parse Ratio"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbExecuteWithoutParseRatioDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Logons Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbLogonsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Read Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbDbPhysicalReadBytesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Reads Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbDbPhysicalReadsPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Write Bytes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbDbPhysicalWriteBytesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
	r.recorders["Physical Writes Per Sec"] = func(mb *metadata.MetricsBuilder, ts pcommon.Timestamp, value float64, instanceID, pdbName string) {
		mb.RecordOracledbPdbDbPhysicalWritesPerSecondDataPoint(ts, value, instanceID, pdbName)
	}
}
