<script lang="ts">
  import { onMount } from "svelte";
  import type { Summary as ContainerSummary } from "$lib/dockercontainer";
  import type { Summary as ImageSummary } from "$lib/dockerimage";
  import { formatBytes, formatTimeAgo, formatContainerName } from "$lib/utils";

  import { dockerData } from "$lib/state.svelte";

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
    fetch("/api/virt/docker/current")
      .then((response) => response.json())
      .then((data) => {
        dockerData.containers = data.containers;
        dockerData.images = data.images;
      });

    // Set up SSE connection
    eventSource = new EventSource("/api/virt/docker/events");
    eventSource.addEventListener("containers", (event) => {
      dockerData.containers = JSON.parse(event.data) as ContainerSummary[];
    });
    eventSource.addEventListener("images", (event) => {
      dockerData.images = JSON.parse(event.data) as ImageSummary[];
    });

    return () => {
      if (eventSource) eventSource.close();
    };
  });
</script>

<Header title="Docker" />
<section
  class="stats stats-vertical xl:stats-horizontal bg-base-100 col-span-12 w-full shadow-xs"
>
  <div class="stat">
    <div class="stat-title">Running Containers</div>
    <div class="stat-value">
      {dockerData.containers !== null
        ? `${dockerData.containers.length}`
        : "Loading..."}
    </div>
    <div class="stat-desc">Number of currently running containers</div>
  </div>
  <div class="stat">
    <div class="stat-title">Local Images</div>
    <div class="stat-value">
      {dockerData.images !== null
        ? `${dockerData.images.length}`
        : "Loading..."}
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
        {#each dockerData.containers as container}
          <tr>
            <td>{formatContainerName(container.Names[0] || "")}</td>
            <td>{formatTimeAgo(container.Created)}</td>
            <td>
              {container.Status}
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
          <th>Repository</th>
          <th>Tag</th>
        </tr>
      </thead>
      <tbody>
        {#each dockerData.images as image}
          <tr>
            <td>{formatImageId(image.Id)}</td>
            <td>{formatTimeAgo(image.Created)}</td>
            <td>{formatBytes(image.Size)}</td>
            <td>{image.RepoTags[0]}</td>
            <td>{image.RepoTags[1]}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>
