<script lang="ts">
    import type { Stat } from "$lib/types";
    export let stats: Stat[] = [];
</script>

<section
    class="stats stats-vertical xl:stats-horizontal bg-base-100 col-span-12 w-full shadow-xs"
>
    {#each stats as stat}
        <div class="stat">
            <div class="stat-title">{stat.title}</div>
            <div class="stat-value">{stat.value}</div>
            <div class="stat-desc">{stat.description}</div>
        </div>
    {/each}
</section>
