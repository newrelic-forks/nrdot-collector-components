// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

// Package oracledbreceiver implements a receiver that can fetch Oracle database metrics
// and emit them in a format compatible with New Relic's monitoring infrastructure.
// This receiver provides a simplified interface to collect key Oracle metrics
// for New Relic observability platform integration.
package oracledbreceiver // import "github.com/newrelic/nrdot-collector-components/receiver/oracledbreceiver"
