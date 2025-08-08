<script lang="ts">
	import "../app.css";
	import Aside from "$lib/components/aside.svelte";

	let { children } = $props();
</script>

<div class="drawer bg-base-200 lg:drawer-open min-h-screen">
	<input id="my-drawer" type="checkbox" class="drawer-toggle" />
	<main class="drawer-content">
		<div
			class="grid grid-cols-12 grid-rows-[min-content] gap-y-12 p-4 lg:gap-x-12 lg:p-10"
		>
			<!-- header -->

			<!-- /header -->
			{@render children()}
		</div>
	</main>
	<Aside />
</div>
