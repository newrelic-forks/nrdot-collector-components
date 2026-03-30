// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

type customCollection struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

func (c customCollection) GetName() string {
	return c.Name
}

func (c customCollection) GetQuery() string {
	return c.Query
}

func (customCollection) Unmarshal(results any) any {
	return results
}

func newCustomCollection(name, query string) iCollection {
	return customCollection{
		Name:  name,
		Query: query,
	}
}
