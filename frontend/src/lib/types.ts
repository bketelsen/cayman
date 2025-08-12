import type {HostInfo} from "./sysinfo.ts"


export interface Stat {
  title?: string;
  value?: number | string;
  description?: string;
}


export interface GlobalData {
  host: HostInfo;
}



