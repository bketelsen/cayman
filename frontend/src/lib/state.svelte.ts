import type { GlobalData } from "./types";
import type {HostInfo} from "./sysinfo.ts"
import type {HostState} from "./cayman"
import type {HostMemoryInfo} from "./sysinfo"


// host is a field so it can be set in the store
// any additional fields should be done similarly
export const globalData: GlobalData = $state({
    host: {} as HostInfo
});

export const dashboardData: HostState = $state({
    hostname: "",
    fqdn: "",
    load: {
        load1: 0,
        load5: 0,
        load15: 0
    },
    cpu: 0,
    unit_status: {
        failed_count: 0,
        active_count: 0
    },
    physical_cores: 0,
    logical_cores: 0,
    host_info: {} as HostInfo,
    memory_info: {} as HostMemoryInfo
});