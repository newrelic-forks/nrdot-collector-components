// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package executor // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/executor"

import (
	"time"
)

type queryExecution struct {
	Query         string
	ExecutedAt    time.Time
	ResultCount   int
	TransformInto any
	State         any
	Error         error
}

// CollectionResult is the externally-visible result of a collection execution.
// It contains the subset of queryExecution fields needed by callers outside this package.
type CollectionResult struct {
	Query         string
	ExecutedAt    time.Time
	ResultCount   int
	TransformInto any
	Error         error
}
