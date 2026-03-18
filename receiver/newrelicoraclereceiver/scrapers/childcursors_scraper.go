// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package scrapers // import "github.com/newrelic/nrdot-collector-components/receiver/newrelicoraclereceiver/scrapers"

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/zap"

	"github.com/newrelic/nrdot-collector-components/receiver/newrelicoraclereceiver/client"
	"github.com/newrelic/nrdot-collector-components/receiver/newrelicoraclereceiver/internal/metadata"
	"github.com/newrelic/nrdot-collector-components/receiver/newrelicoraclereceiver/models"
)

// ChildCursorsScraper handles scraping of child cursor metrics from V$SQL
type ChildCursorsScraper struct {
	client                client.OracleClient
	mb                    *metadata.MetricsBuilder
	logger                *zap.Logger
	metricsBuilderConfig  metadata.MetricsBuilderConfig
	enableAdvancedMetrics bool
}

// NewChildCursorsScraper creates a new child cursors scraper
func NewChildCursorsScraper(oracleClient client.OracleClient, mb *metadata.MetricsBuilder, logger *zap.Logger, metricsBuilderConfig metadata.MetricsBuilderConfig, enableAdvancedMetrics bool) *ChildCursorsScraper {
	return &ChildCursorsScraper{
		client:                oracleClient,
		mb:                    mb,
		logger:                logger,
		metricsBuilderConfig:  metricsBuilderConfig,
		enableAdvancedMetrics: enableAdvancedMetrics,
	}
}

func (s *ChildCursorsScraper) ScrapeChildCursorsForIdentifiers(ctx context.Context, identifiers []models.SQLIdentifier) ([]models.SQLIdentifier, []error) {
	var errs []error
	s.logger.Debug("Starting child cursors scrape")
	now := pcommon.NewTimestampFromTime(time.Now())
	metricsEmitted := 0

	for i := range identifiers {
		cursor, err := s.client.QuerySpecificChildCursor(ctx, identifiers[i].SQLID, identifiers[i].ChildNumber)
		if err != nil {
			s.logger.Warn("Failed to fetch specific child cursor from V$SQL",
				zap.String("sql_id", identifiers[i].SQLID),
				zap.Int64("child_number", identifiers[i].ChildNumber),
				zap.Error(err))
			errs = append(errs, err)
			continue
		}

		if cursor == nil {
			s.logger.Debug("No child cursor found in V$SQL",
				zap.String("sql_id", identifiers[i].SQLID),
				zap.Int64("child_number", identifiers[i].ChildNumber))
			continue
		}

		if !cursor.HasValidIdentifier() {
			s.logger.Debug("Child cursor has invalid identifier, skipping",
				zap.String("sql_id", identifiers[i].SQLID),
				zap.Int64("child_number", identifiers[i].ChildNumber))
			continue
		}

		identifiers[i].PlanHash = fmt.Sprintf("%d", cursor.GetPlanHashValue())
		s.recordChildCursorMetrics(now, cursor)
		metricsEmitted++
	}

	s.logger.Info("Child cursors scrape completed", zap.Int("metrics_emitted", metricsEmitted))

	return identifiers, errs
}

// recordChildCursorMetrics records UI-critical metrics for a single child cursor
func (s *ChildCursorsScraper) recordChildCursorMetrics(now pcommon.Timestamp, cursor *models.ChildCursor) {
	collectionTimestamp := cursor.GetCollectionTimestamp()
	sqlID := cursor.GetSQLID()
	childNumber := cursor.GetChildNumber()
	planHashValue := fmt.Sprintf("%d", cursor.GetPlanHashValue())
	databaseName := cursor.GetDatabaseName()

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsElapsedTime.Enabled {
		s.mb.RecordOracledbChildCursorsElapsedTimeDataPoint(
			now,
			cursor.GetElapsedTime(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsDetails.Enabled {
		s.mb.RecordOracledbChildCursorsDetailsDataPoint(
			now,
			1,
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
			cursor.GetFirstLoadTime(),
			cursor.GetLastLoadTime(),
		)
	}

	if s.enableAdvancedMetrics {
		s.recordAdvancedChildCursorMetrics(now, cursor, collectionTimestamp, databaseName, sqlID, childNumber, planHashValue)
	}
}

func (s *ChildCursorsScraper) recordAdvancedChildCursorMetrics(now pcommon.Timestamp, cursor *models.ChildCursor, collectionTimestamp, databaseName, sqlID string, childNumber int64, planHashValue string) {
	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsCPUTime.Enabled {
		s.mb.RecordOracledbChildCursorsCPUTimeDataPoint(
			now,
			cursor.GetCPUTime(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsUserIoWaitTime.Enabled {
		s.mb.RecordOracledbChildCursorsUserIoWaitTimeDataPoint(
			now,
			cursor.GetUserIOWaitTime(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsExecutions.Enabled {
		s.mb.RecordOracledbChildCursorsExecutionsDataPoint(
			now,
			cursor.GetExecutions(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsDiskReads.Enabled {
		s.mb.RecordOracledbChildCursorsDiskReadsDataPoint(
			now,
			cursor.GetDiskReads(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsBufferGets.Enabled {
		s.mb.RecordOracledbChildCursorsBufferGetsDataPoint(
			now,
			cursor.GetBufferGets(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}

	if s.metricsBuilderConfig.Metrics.OracledbChildCursorsInvalidations.Enabled {
		s.mb.RecordOracledbChildCursorsInvalidationsDataPoint(
			now,
			cursor.GetInvalidations(),
			collectionTimestamp,
			databaseName,
			sqlID,
			childNumber,
			planHashValue,
		)
	}
}
