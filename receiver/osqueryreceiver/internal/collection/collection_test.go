// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection

import (
	"reflect"
	"testing"
)

func TestGetCollection(t *testing.T) {
	tests := []struct {
		name           string
		collectionName string
		factory        func() iCollection
		wantName       string
		wantErr        bool
	}{
		{
			name:           "system_info factory",
			collectionName: systemInfoCollectionName,
			factory:        newSystemInfoCollection,
			wantName:       systemInfoCollectionName,
		},
		{
			name:           "package_info factory",
			collectionName: packageInfoCollectionName,
			factory:        newPackageInfoCollection,
			wantName:       packageInfoCollectionName,
		},
		{
			name:           "os_info factory",
			collectionName: osInfoCollectionName,
			factory:        newOSInfoCollection,
			wantName:       osInfoCollectionName,
		},
		{
			name:           "secureboot_info factory",
			collectionName: secureBootCollectionName,
			factory:        newSecureBootCollection,
			wantName:       secureBootCollectionName,
		},
		{
			name:           "users_info factory",
			collectionName: userCollectionName,
			factory:        newUserCollection,
			wantName:       userCollectionName,
		},
		{
			name:           "invalid collection name",
			collectionName: "invalid_collection",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCollection(tt.collectionName)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			want := tt.factory()

			if reflect.TypeOf(got) != reflect.TypeOf(want) {
				t.Fatalf("expected type %T but got %T", want, got)
			}

			if got.GetName() != tt.wantName {
				t.Fatalf("expected name %q but got %q", tt.wantName, got.GetName())
			}

			// Check that the factory-provided collection exposes a query.
			if wantQuery := want.GetQuery(); wantQuery != "" && got.GetQuery() != wantQuery {
				t.Fatalf("expected query %q but got %q", wantQuery, got.GetQuery())
			}
		})
	}
}

func TestGetCustomCollection(t *testing.T) {
	const (
		name  = "custom_table"
		query = "select 1;"
	)

	got := getCustomCollection(name, query)

	if got.GetName() != name {
		t.Fatalf("expected name %q but got %q", name, got.GetName())
	}

	if got.GetQuery() != query {
		t.Fatalf("expected query %q but got %q", query, got.GetQuery())
	}

	// Custom collections should be pass-through for unmarshalling.
	inputs := []map[string]any{{"k": "v"}}
	if unmarshalled := got.Unmarshal(inputs); !reflect.DeepEqual(unmarshalled, inputs) {
		t.Fatalf("expected unmarshalled result %#v but got %#v", inputs, unmarshalled)
	}
}
