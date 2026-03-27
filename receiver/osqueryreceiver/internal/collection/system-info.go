// Copyright New Relic, Inc. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package collection // import "github.com/newrelic/nrdot-collector-components/receiver/osqueryreceiver/internal/collection"

// systemInfoCollection represents the system_info collection
// https://github.com/osquery/osquery/blob/master/specs/system_info.table
type systemInfoCollection struct {
	Hostname         string `json:"hostname"`
	UUID             string `json:"uuid"`
	CPUType          string `json:"cpu_type"`
	CPUSubtype       string `json:"cpu_subtype"`
	CPUBrand         string `json:"cpu_brand"`
	CPUPhysicalCores int    `json:"cpu_physical_cores"`
	CPULogicalCores  int    `json:"cpu_logical_cores"`
	PhysicalMemory   string `json:"physical_memory"`
	HardwareVendor   string `json:"hardware_vendor"`
	HardwareModel    string `json:"hardware_model"`
	ComputerName     string `json:"computer_name,omitempty"`
	EmulatedCPUType  string `json:"emulated_cpu_type,omitempty"`
}

func (systemInfoCollection) GetName() string {
	return systemInfoCollectionName
}

func (systemInfoCollection) GetQuery() string {
	return systemInfoCollectionQuery
}

func (systemInfoCollection) Unmarshal(result any) any {
	resultMap, ok := result.(map[string]any)
	if !ok {
		return nil
	}

	sanitized := sanitizeRow(
		resultMap,
		[]string{
			"hostname",
			"uuid",
			"cpu_type",
			"cpu_subtype",
			"cpu_brand",
			"physical_memory",
			"hardware_vendor",
			"hardware_model",
			"computer_name",
			"emulated_cpu_type",
		},
		[]string{
			"cpu_physical_cores",
			"cpu_logical_cores",
		},
		nil,
	)

	if len(sanitized) == 0 {
		return nil
	}

	return sanitized
}

func newSystemInfoCollection() iCollection {
	return systemInfoCollection{}
}
