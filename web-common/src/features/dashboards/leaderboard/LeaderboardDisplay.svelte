<script lang="ts">
  import VirtualizedGrid from "@rilldata/web-common/components/VirtualizedGrid.svelte";
  import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
  import { onDestroy, onMount } from "svelte";
  import Leaderboard from "./Leaderboard.svelte";
  import LeaderboardControls from "./LeaderboardControls.svelte";

  const {
    selectors: {
      dimensions: { visibleDimensions },
    },
    metricsViewName,
  } = getStateManagers();

  /** Functionality for resizing the virtual leaderboard */
  let columns = 3;
  let availableWidth = 0;
  let leaderboardContainer: HTMLElement;
  let observer: ResizeObserver;

  function onResize() {
    if (!leaderboardContainer) return;
    availableWidth = leaderboardContainer.offsetWidth;
    columns = Math.max(1, Math.floor(availableWidth / (315 + 20)));
  }

  onMount(() => {
    onResize();
    const observer = new ResizeObserver(() => {
      onResize();
    });
    observer.observe(leaderboardContainer);
  });

  onDestroy(() => {
    observer?.disconnect();
  });
</script>

<svelte:window on:resize={onResize} />
<!-- container for the metrics leaderboard components and controls -->
<div
  bind:this={leaderboardContainer}
  class="flex flex-col overflow-hidden"
  style:height="calc(100vh - 130px - 4rem)"
  style:min-width="365px"
>
  <div
    class="grid grid-auto-cols justify-between grid-flow-col items-center pl-1 pb-3 flex-grow-0"
  >
    <LeaderboardControls metricViewName={$metricsViewName} />
  </div>
  <div class="grow overflow-hidden">
    {#if $visibleDimensions.length > 0}
      <VirtualizedGrid
        {columns}
        height="100%"
        items={$visibleDimensions}
        let:item
      >
        <!-- the single virtual element -->
        <Leaderboard dimensionName={item.name} />
      </VirtualizedGrid>
    {/if}
  </div>
</div>
