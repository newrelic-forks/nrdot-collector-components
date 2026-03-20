// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package models // import "github.com/newrelic/nrdot-collector-components/receiver/oracledbreceiver/models"

type SystemMetric struct {
	InstanceID string
	MetricName string
	Value      float64
}
