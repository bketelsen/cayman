package cayman

import "github.com/elastic/go-sysinfo/types"

// HostState is used on the dashboard
type HostState struct {
	Hostname   string     `json:"hostname"`
	FQDN       string     `json:"fqdn"` // Include FQDN
	Load       Load       `json:"load"`
	CPU        int        `json:"cpu"`         // Include CPU usage
	UnitStatus UnitStatus `json:"unit_status"` // Include unit status
	// CPUInfo       cpu.InfoStat             `json:"cpu_info"`       // Include CPU info
	PhysicalCores int                  `json:"physical_cores"` // Include physical cores
	LogicalCores  int                  `json:"logical_cores"`  // Include logical cores
	HostInfo      types.HostInfo       `json:"host_info"`      // Include host info
	MemoryInfo    types.HostMemoryInfo `json:"memory_info"`    // Include memory info
}

type UnitStatus struct {
	FailedCount int `json:"failed_count"`
	ActiveCount int `json:"active_count"` // Include active count
}

type Load struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}
