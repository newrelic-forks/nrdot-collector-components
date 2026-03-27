// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package osqueryreceiver // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type osqueryReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Logs
	config       *Config
}

func (o *osqueryReceiver) Start(ctx context.Context, host component.Host) error {
	o.host = host
	ctx, o.cancel = context.WithCancel(ctx)

	interval, _ := time.ParseDuration(o.config.CollectionInterval)

	osQueryManager, err := newOSQueryManager(o.config, o.logger)
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				o.logger.Info("Starting collection")
				if err := osQueryManager.collect(o.nextConsumer); err != nil {
					o.logger.Error("Collection failed", zap.Error(err))
				}
			case <-ctx.Done():
				o.logger.Info("Shutting down osquery receiver collection")
				return
			}
		}
	}()

	return nil
}

func (o *osqueryReceiver) Shutdown(_ context.Context) error {
	if o.cancel != nil {
		o.cancel()
	}
	return nil
}
