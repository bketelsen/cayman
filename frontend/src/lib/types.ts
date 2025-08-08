// place files you want to import through the `$lib` alias in this folder.
export interface Snapshot {
  cycles: number | undefined;
  ops: number | undefined;
}
export interface Stat {
  title?: string;
  value?: number | string;
  description?: string;
}

export interface Load {
  load1: number;
  load5: number;
  load15: number;
}
export interface HostState {
  load: Load | undefined;
  cpu: number | undefined;

}

/*
// OSInfo contains basic OS information
type OSInfo struct {
  Type     string `json:"type"`               // OS Type (one of linux, macos, unix, windows).
  Family   string `json:"family"`             // OS Family (e.g. redhat, debian, freebsd, windows).
  Platform string `json:"platform"`           // OS platform (e.g. centos, ubuntu, windows).
  Name     string `json:"name"`               // OS Name (e.g. Mac OS X, CentOS).
  Version  string `json:"version"`            // OS version (e.g. 10.12.6).
  Major    int    `json:"major"`              // Major release version.
  Minor    int    `json:"minor"`              // Minor release version.
  Patch    int    `json:"patch"`              // Patch release version.
  Build    string `json:"build,omitempty"`    // Build (e.g. 16G1114).
  Codename string `json:"codename,omitempty"` // OS codename (e.g. jessie).
}
  */
export interface OSInfo {
  type: string;
  family: string;
  platform: string;
  name: string;
  version: string;
  major: number;
  minor: number;
  patch: number;
  build?: string;
  codename?: string;
}

/*  Go Struct: HostInfo contains basic host information
type HostInfo struct {
  Architecture       string    `json:"architecture"`            // Process hardware architecture (e.g. x86_64, arm, ppc, mips).
  NativeArchitecture string    `json:"native_architecture"`     // Native OS hardware architecture (e.g. x86_64, arm, ppc, mips).
  BootTime           time.Time `json:"boot_time"`               // Host boot time.
  Containerized      *bool     `json:"containerized,omitempty"` // Is the process containerized.
  Hostname           string    `json:"name"`                    // Hostname.
  IPs                []string  `json:"ip,omitempty"`            // List of all IPs.
  KernelVersion      string    `json:"kernel_version"`          // Kernel version.
  MACs               []string  `json:"mac"`                     // List of MAC addresses.
  OS                 *OSInfo   `json:"os"`                      // OS information.
  Timezone           string    `json:"timezone"`                // System timezone.
  TimezoneOffsetSec  int       `json:"timezone_offset_sec"`     // Timezone offset (seconds from UTC).
  UniqueID           string    `json:"id,omitempty"`            // Unique ID of the host (optional).
}
  */
export interface HostInfo {
  architecture: string;
  nativeArchitecture: string;
  bootTime: string; // ISO 8601 format
  containerized?: boolean;
  hostname: string;
  ips?: string[];
  kernelVersion: string;
  macs?: string[];
  os: OSInfo;
  timezone: string;
  timezoneOffsetSec: number;
  uniqueID?: string; // Unique ID of the host (optional)
}

export interface GlobalData {
  host: HostInfo;
}

export interface MemoryInfo {
  total_bytes: number
  used_bytes: number
  available_bytes: number
  free_bytes: number
  virtual_total_bytes: number
  virtual_used_bytes: number
  virtual_free_bytes: number
  raw: Raw
}

export interface Raw {
  Active: number
  "Active(anon)": number
  "Active(file)": number
  AnonHugePages: number
  AnonPages: number
  Bounce: number
  Buffers: number
  Cached: number
  CmaFree: number
  CmaTotal: number
  CommitLimit: number
  Committed_AS: number
  DirectMap1G: number
  DirectMap2M: number
  DirectMap4k: number
  Dirty: number
  FileHugePages: number
  FilePmdMapped: number
  HardwareCorrupted: number
  HugePages_Free: number
  HugePages_Rsvd: number
  HugePages_Surp: number
  HugePages_Total: number
  Hugepagesize: number
  Hugetlb: number
  Inactive: number
  "Inactive(anon)": number
  "Inactive(file)": number
  KReclaimable: number
  KernelStack: number
  Mapped: number
  Mlocked: number
  NFS_Unstable: number
  PageTables: number
  Percpu: number
  SReclaimable: number
  SUnreclaim: number
  SecPageTables: number
  Shmem: number
  ShmemHugePages: number
  ShmemPmdMapped: number
  Slab: number
  SwapCached: number
  Unaccepted: number
  Unevictable: number
  VmallocChunk: number
  VmallocTotal: number
  VmallocUsed: number
  Writeback: number
  WritebackTmp: number
  Zswap: number
  Zswapped: number
}
