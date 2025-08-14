<script lang="ts">
  import { onMount } from "svelte";
  import type { InstanceFull, Image } from "$lib/incus";
  import { formatBytes, formatTimeAgo, formatContainerName } from "$lib/utils";

  import { incusData } from "$lib/state.svelte";

  import Header from "$lib/components/header.svelte";

  // Function to format Docker image ID by removing "sha256:" prefix and returning first 12 characters
  function formatImageId(id: string): string {
    if (!id) return "";

    // Remove "sha256:" prefix if present
    const cleanId = id.startsWith("sha256:") ? id.substring(7) : id;

    // Return first 12 characters
    return cleanId.substring(0, 12);
  }

  let eventSource = $state<EventSource | undefined>(undefined);
  // break this out into its own state so we can use it in the header

  onMount(() => {
    // Initial data fetch
    fetch("/api/virt/incus/current")
      .then((response) => response.json())
      .then((data) => {
        incusData.instances = data.instances;
        incusData.images = data.images;
      });

    // Set up SSE connection
    eventSource = new EventSource("/api/virt/incus/events");
    eventSource.addEventListener("instances", (event) => {
      incusData.instances = JSON.parse(event.data) as InstanceFull[];
    });
    eventSource.addEventListener("images", (event) => {
      incusData.images = JSON.parse(event.data) as Image[];
    });

    return () => {
      if (eventSource) eventSource.close();
    };
  });
</script>

<Header title="Incus" />
<section
  class="stats stats-vertical xl:stats-horizontal bg-base-100 col-span-12 w-full shadow-xs"
>
  <div class="stat">
    <div class="stat-title">Running Instances</div>
    <div class="stat-value">
      {incusData.instances !== null
        ? `${incusData.instances.length}`
        : "Loading..."}
    </div>
    <div class="stat-desc">Number of currently running instances</div>
  </div>
  <div class="stat">
    <div class="stat-title">Local Images</div>
    <div class="stat-value">
      {incusData.images !== null ? `${incusData.images.length}` : "Loading..."}
    </div>
    <div class="stat-desc">Number of locally available images</div>
  </div>
</section>

<section class="card bg-base-100 col-span-12 overflow-hidden shadow-xs">
  <div class="card-body grow-0">
    <div class="flex justify-between gap-2">
      <h2 class="card-title">
        <a class="link-hover link">Containers</a>
      </h2>

      <select class="select">
        <option disabled selected>Filter by Status</option>
        <option>Running</option>
        <option>Stopped</option>
        <option>Exited</option>
      </select>
    </div>
  </div>
  <div class="overflow-x-auto">
    <table class="table-zebra table">
      <thead>
        <tr>
          <th>Container Name</th>
          <th>Created</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {#each incusData.instances as instance}
          <tr>
            <td>{formatContainerName(instance.name || "")}</td>
            <td>{instance.created_at}</td>
            <td>
              {instance.status}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>
<section class="card bg-base-100 col-span-12 overflow-hidden shadow-xs">
  <div class="card-body grow-0">
    <h2 class="card-title">
      <a class="link-hover link">Images</a>
    </h2>
  </div>
  <div class="overflow-x-auto">
    <table class="table-zebra table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Created</th>
          <th>Size</th>
        </tr>
      </thead>
      <tbody>
        {#each incusData.images as image}
          <tr>
            <td>{formatImageId(image.fingerprint)}</td>
            <td>{image.created_at}</td>
            <td>{formatBytes(image.size)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>
