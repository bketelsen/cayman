<script lang="ts">
  import { onMount } from "svelte";
  import type { MemoryInfo, Load } from "$lib/types";
  import { globalData } from "$lib/state.svelte";

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

  let load = $state<Load | null>(null);
  let cpu = $state<number | null>(null);
  let failedUnits = $state<number | null>(null);
  let activeUnits = $state<number | null>(null);
  let eventSource = $state<EventSource | undefined>(undefined);
  // break this out into its own state so we can use it in the header
  let hostname = $state<string | null>(null);
  let mem = $state<MemoryInfo | null>(null);
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
    if (mem) {
      return mem.total_bytes - mem.used_bytes;
    }
    return null;
  });
  let usedMemory = $derived(() => {
    if (mem) {
      return mem.total_bytes - mem.available_bytes;
    }
    return null;
  });
  let usedMemoryPercentage = $derived(() => {
    if (mem) {
      const used = usedMemory();
      if (used !== null) {
        return ((used / mem.total_bytes) * 100).toFixed(0);
      }
    }
    return null;
  });

  onMount(() => {
    // Initial data fetch
    fetch("/api/host/current")
      .then((response) => response.json())
      .then((data) => {
        load = data.load;
        cpu = data.cpu;
        failedUnits = data.unit_status.failed_count;
        activeUnits = data.unit_status.active_count;
        hostname = data.host_info.name;
        fqdn = data.fqdn;
        mem = data.memory_info;
        globalData.host = data.host_info;
      });

    // Set up SSE connection
    eventSource = new EventSource("/api/host/events");
    eventSource.addEventListener("load", (event) => {
      load = JSON.parse(event.data) as Load;
    });
    eventSource.addEventListener("cpu", (event) => {
      cpu = JSON.parse(event.data) as number;
    });
    eventSource.addEventListener("mem", (event) => {
      mem = JSON.parse(event.data) as MemoryInfo;
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
    <div class="stat-value">{cpu !== null ? `${cpu}%` : "Loading..."}</div>
    <div class="stat-desc">Current CPU usage percentage</div>
  </div>
  <div class="stat">
    <div class="stat-title">Load 1</div>
    <div class="stat-value">
      {load ? `${load.load1} ` : "Loading..."}
    </div>
    <div class="stat-desc">Current load average over 1 minute</div>
  </div>
  <div class="stat">
    <div class="stat-title">Load 5</div>
    <div class="stat-value">
      {load ? `${load.load5} ` : "Loading..."}
    </div>
    <div class="stat-desc">Current load average over 5 minutes</div>
  </div>
  <div class="stat">
    <div class="stat-title">Load 15</div>
    <div class="stat-value">
      {load ? `${load.load15} ` : "Loading..."}
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
        {globalData.host?.os !== null ? globalData.host.os.type : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Family
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os !== null
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
        {globalData.host?.os !== null
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
        {globalData.host?.os !== null ? globalData.host.os.name : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Version
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right"
      >
        {globalData.host?.os !== null
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
        {globalData.host?.os !== null
          ? globalData.host.os.major
          : "Loading..."}.{globalData.host?.os !== null
          ? globalData.host.os.minor
          : "Loading..."}.{globalData.host?.os !== null
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
        {globalData.host?.os !== null ? globalData.host.os.build : "Loading..."}
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
        {activeUnits !== null ? activeUnits : "Loading..."}
      </dd>
      <dt
        class="col-start-1 border-t border-primary/70 pt-3 text-neutral-content/60 first:border-none border-t border-primary/70 py-3"
      >
        Failed Units
      </dt>
      <dd
        class="pt-3 pb-3 col-start-2 col-span-2 border-t border-primary/70 border-t border-primary/70 py-3 nth-2:border-none text-right {failedUnits !==
          null && failedUnits > 0
          ? 'text-error'
          : ''}"
      >
        {failedUnits !== null ? failedUnits : "Loading..."}
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
              style="width: {cpu}%"
              class="h-2 rounded-full bg-primary"
            ></div>
          </div>
          <div class=" col-span-2 text-right">{cpu}% of X CPUs</div>
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
            {humanizeMemory(usedMemory() || 0)} of {humanizeMemory(
              mem?.total_bytes || 0,
            )}
            {usedMemoryPercentage()}%
          </div>
        </div>
      </dd>
    </dl>
    <div class="card-actions justify-end">
      <a href="/system" class="link">Details</a>
    </div>
  </div>
</div>
