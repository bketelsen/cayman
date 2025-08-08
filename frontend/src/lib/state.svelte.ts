import type { GlobalData } from "./types";
// host is a field so it can be set in the store
// any additional fields should be done similarly
export const globalData: GlobalData = $state(
    {
        host: {
            architecture: "",
            nativeArchitecture: "",
            bootTime: "",
            containerized: false,
            hostname: "",
            ips: [],
            kernelVersion: "",
            macs: [],
            os: {
                type: "",
                family: "",
                platform: "",
                name: "",
                version: "",
                major: 0,
                minor: 0,
                patch: 0,
                build: "",
                codename: ""
            },
            timezone: "",
            timezoneOffsetSec: 0,
            uniqueID: ""
        }
    });