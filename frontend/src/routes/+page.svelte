<script lang="ts">
  import { onMount } from "svelte";
  import type {  Load } from "$lib/cayman";
  import type { HostMemoryInfo } from "$lib/sysinfo";

  import { globalData, dashboardData } from "$lib/state.svelte";
 
  import Header from "$lib/components/header.svelte";

  // TODO: Move this to a utility file
  // This function converts bytes to a human-readable format
  function humanizeMemory(bytes: number): string {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const sizes = ["Bytes", "KiB", "MiB", "GiB", "TiB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(0)) + " " + sizes[i];
  }


   let eventSource = $state<EventSource | undefined>(undefined);
  // break this out into its own state so we can use it in the header
  let hostname = $state<string | null>(null);
  // let mem = $state<HostMemoryInfo | null>(null);
  let fqdn = $state<string | null>(null);
  /**
  Line 1 reflects physical memory, classified as:
    total          ( MemTotal )
    free           ( MemFree )
    used           ( MemTotal - MemAvailable )
    buff/cache     ( Buffers + Cached + SReclaimable )

   */
  // TODO: this is all broken/wrong... need to verify calculations
  // to match `top` output
  let freeMemory = $derived(() => {
    if (dashboardData.memory_info) {
      return dashboardData.memory_info.total_bytes - dashboardData.memory_info.used_bytes;
    }
    return null;
  });
  let usedMemory = $derived(() => {
    if (dashboardData.memory_info) {
      return dashboardData.memory_info.total_bytes - dashboardData.memory_info.available_bytes;
    }
    return null;
  });
  let usedMemoryPercentage = $derived(() => {
    if (dashboardData.memory_info) {
      const used = usedMemory();
      if (used !== null) {
        return ((used / dashboardData.memory_info.total_bytes) * 100).toFixed(0);
      }
    }
    return null;
  });

  onMount(() => {
    // Initial data fetch
    fetch("/api/dashboard/current")
      .then((response) => response.json())
      .then((data) => {
        dashboardData.cpu = data.cpu;
        dashboardData.cpu_count = data.cpu_count;
        dashboardData.load = data.load;
        dashboardData.unit_status = data.unit_status;
        dashboardData.physical_cores = data.physical_cores;
        dashboardData.logical_cores = data.logical_cores;
        dashboardData.host_info = data.host_info;
        dashboardData.memory_info = data.memory_info;
        hostname = data.host_info.name;
        fqdn = data.fqdn;
        globalData.host = data.host_info;
      });

    // Set up SSE connection
    eventSource = new EventSource("/api/dashboard/events");
    eventSource.addEventListener("load", (event) => {
      dashboardData.load = JSON.parse(event.data) as Load;
    });
    eventSource.addEventListener("cpu", (event) => {
      dashboardData.cpu = JSON.parse(event.data) as number;
    });
    eventSource.addEventListener("mem", (event) => {
      dashboardData.memory_info = JSON.parse(event.data) as HostMemoryInfo;
    });

    return () => {
      if (eventSource) eventSource.close();
    };
  });
</script>

<Header
  title="Dashboard"
  subtitle={hostname + " | " + fqdn || "Unknown Host"}
/>
<section
  class="stats stats-vertical xl:stats-horizontal bg-base-100 col-span-12 w-full shadow-xs"
>
  <div class="stat">
    <div class="stat-title">CPU Usage</div>
    <div class="stat-value">{dashboardData.cpu !== null ? `${dashboardData.cpu}%` : "Loading..."}</div>
    <div class="stat-desc">Current CPU usage percentage</div>
  </div>
  <div class="stat">
    <div class="stat-title">Load 1</div>
    <div class="stat-value">
      {dashboardData.load ? `${dashboardData.load.load1} ` : "Loading..."}
    </div>
    <div class="stat-desc">Current load average over 1 minute</div>
  </div>
  <div class="stat">
    <div class="stat-title">Load 5</div>
    <div class="stat-value">
      {dashboardData.load ? `${dashboardData.load.load5} ` : "Loading..."}
    </div>
    <div class="stat-desc">Current load average over 5 minutes</div>
  </div>
  <div class="stat">
    <div class="stat-title">Load 15</div>
    <div class="stat-value">
      {dashboardData.load ? `${dashboardData.load.load15} ` : "Loading..."}
    </div>
    <div class="stat-desc">Current load average over 15 minutes</div>
  </div>
</section>

<div class="card card-border bg-base-100 lg:col-span-4 col-span-12 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Operating System</h2>
    <dl class="grid grid-cols-3">
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        System
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os ? globalData.host.os.type : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Family
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os
          ? globalData.host.os.family
          : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Platform
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os
          ? globalData.host.os.platform
          : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Name
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os ? globalData.host.os.name : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Version
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os !== undefined && globalData.host?.os !== null
          ? globalData.host.os.version
          : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Version Details
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os
          ? globalData.host.os.major
          : "Loading..."}.{globalData.host?.os
          ? globalData.host.os.minor
          : "Loading..."}.{globalData.host?.os
          ? globalData.host.os.patch
          : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Build
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os ? globalData.host.os.build : "Loading..."}
      </dd>
    </dl>
  </div>
</div>
<div class="card card-border bg-base-100 lg:col-span-4 col-span-12 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">System Health</h2>
    <dl class="grid grid-cols-3">
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Active Units
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {dashboardData.unit_status.active_count !== null ? dashboardData.unit_status.active_count : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Failed Units
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right {dashboardData.unit_status.failed_count !==
          null && dashboardData.unit_status.failed_count > 0
          ? 'text-error'
          : ''}"
      >
        {dashboardData.unit_status.failed_count !== null ? dashboardData.unit_status.failed_count : "Loading..."}
      </dd>
    </dl>
    <div class="card-actions justify-end">
      <a href="/system" class="link">Details</a>
    </div>
  </div>
</div>
<div class="card card-border bg-base-100 lg:col-span-4 col-span-12 shadow-xl">
  <div class="card-body">
    <h2 class="card-title">Usage</h2>
    <dl class="grid grid-cols-3">
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        CPU
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none"
      >
        <div class="grid-cols-6 grid">
          <div
            class="col-span-4 my-2 h-2 overflow-hidden rounded-full bg-neutral"
          >
            <div
              style="width: {dashboardData.cpu}%"
              class="h-2 rounded-full bg-primary"
            ></div>
          </div>
          <div class=" col-span-2 text-right">{dashboardData.cpu}% of {dashboardData.cpu_count} CPUs</div>
        </div>
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Memory
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none"
      >
        <div class="grid-cols-6 grid">
          <div
            class="col-span-4 my-2 h-2 overflow-hidden rounded-full bg-neutral"
          >
            <div
              style="width: {usedMemoryPercentage()}%"
              class="h-2 rounded-full bg-primary"
            ></div>
          </div>
          <div class=" col-span-2 text-right">
            {humanizeMemory(usedMemory() || 0)} ({usedMemoryPercentage()}%) of {humanizeMemory(
              dashboardData.memory_info?.total_bytes || 0,
            )}
            
          </div>
        </div>
      </dd>
    </dl>
    <div class="card-actions justify-end">
      <a href="/system" class="link">Details</a>
    </div>
  </div>
</div>
