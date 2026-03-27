// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection

import (
	"reflect"
	"runtime"
	"testing"
)

func TestPackageInfoCollection_Metadata(t *testing.T) {
	col := packageInfoCollection{}

	if got := col.GetName(); got != packageInfoCollectionName {
		t.Fatalf("expected name %q but got %q", packageInfoCollectionName, got)
	}

	if got, want := col.GetQuery(), col.queryForOS(runtime.GOOS); got != want {
		t.Fatalf("expected query %q but got %q", want, got)
	}
}

func TestPackageInfoCollection_queryForOS(t *testing.T) {
	col := packageInfoCollection{}

	tests := []struct {
		os   string
		want string
	}{
		{os: "darwin", want: packageInfoCollectionQueryHomebrew},
		{os: "linux", want: packageInfoCollectionQueryDebian},
		{os: "windows", want: packageInfoCollectionQueryHomebrew},
	}

	for _, tt := range tests {
		if got := col.queryForOS(tt.os); got != tt.want {
			t.Fatalf("for os %q expected query %q but got %q", tt.os, tt.want, got)
		}
	}
}

func TestPackageInfoCollection_Unmarshal(t *testing.T) {
	col := packageInfoCollection{}

	tests := []struct {
		name string
		in   any
		want any
	}{
		{
			name: "mixed payload",
			in: []map[string]any{
				{
					"name":         "  package-one ",
					"version":      " 1.0.0 ",
					"type":         " formula ",
					"path":         " /usr/local ",
					"app_name":     " app ",
					"source":       " tap ",
					"vendor":       " vendor ",
					"arch":         " arm64 ",
					"install_time": " 1700000000 ",
					"size":         1234,
				},
				{
					"name": "", // should be dropped
				},
			},
			want: []map[string]any{
				{
					"name":         "package-one",
					"version":      "1.0.0",
					"type":         "formula",
					"path":         "/usr/local",
					"app_name":     "app",
					"source":       "tap",
					"vendor":       "vendor",
					"arch":         "arm64",
					"install_time": "1700000000",
					"size":         float64(1234),
				},
			},
		},
		{
			name: "non slice input",
			in:   map[string]any{"name": "pkg"},
			want: nil,
		},
		{
			name: "empty after sanitize",
			in: []map[string]any{
				{"name": "  "},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := col.Unmarshal(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected unmarshal result: got %#v want %#v", got, tt.want)
			}
		})
	}
}
