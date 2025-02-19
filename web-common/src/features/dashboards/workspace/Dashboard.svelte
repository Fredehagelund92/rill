<script lang="ts">
  import { getEltSize } from "@rilldata/web-common/features/dashboards/get-element-size";
  import {
    useDashboard,
    useModelHasTimeSeries,
  } from "@rilldata/web-common/features/dashboards/selectors";
  import { featureFlags } from "@rilldata/web-common/features/feature-flags";
  import { createResizeListenerActionFactory } from "@rilldata/web-common/lib/actions/create-resize-listener-factory";
  import { getContext } from "svelte";
  import type { Tweened } from "svelte/motion";
  import { runtime } from "../../../runtime-client/runtime-store";
  import MeasuresContainer from "../big-number/MeasuresContainer.svelte";
  import { useDashboardStore } from "web-common/src/features/dashboards/stores/dashboard-stores";
  import DimensionDisplay from "../dimension-table/DimensionDisplay.svelte";
  import Filters from "../filters/Filters.svelte";
  import MockUserHasNoAccess from "../granular-access-policies/MockUserHasNoAccess.svelte";
  import { selectedMockUserStore } from "../granular-access-policies/stores";
  import LeaderboardDisplay from "../leaderboard/LeaderboardDisplay.svelte";
  import RowsViewerAccordion from "../rows-viewer/RowsViewerAccordion.svelte";
  import TimeControls from "../time-controls/TimeControls.svelte";
  import MetricsTimeSeriesCharts from "../time-series/MetricsTimeSeriesCharts.svelte";
  import DashboardCTAs from "./DashboardCTAs.svelte";
  import DashboardTitle from "./DashboardTitle.svelte";
  import TimeDimensionDisplay from "../time-dimension-details/TimeDimensionDisplay.svelte";

  export let metricViewName: string;
  export let leftMargin = undefined;

  let exploreContainerWidth;

  $: metricsExplorer = useDashboardStore(metricViewName);

  $: selectedDimensionName = $metricsExplorer?.selectedDimensionName;
  $: expandedMeasureName = $metricsExplorer?.expandedMeasureName;
  $: metricTimeSeries = useModelHasTimeSeries(
    $runtime.instanceId,
    metricViewName
  );
  $: hasTimeSeries = $metricTimeSeries.data;

  // flex-row flex-col
  $: dashboardAlignment = expandedMeasureName ? "col" : "row";

  // the navigationVisibilityTween is a tweened value that is used
  // to animate the extra padding that needs to be added to the
  // dashboard container when the navigation pane is collapsed
  const navigationVisibilityTween = getContext(
    "rill:app:navigation-visibility-tween"
  ) as Tweened<number>;

  const { observedNode, listenToNodeResize } =
    createResizeListenerActionFactory();

  $: exploreContainerWidth = getEltSize($observedNode, "x");

  $: leftSide = leftMargin
    ? leftMargin
    : `calc(${$navigationVisibilityTween * 24}px + 1.25rem)`;

  $: isRillDeveloper = $featureFlags.readOnly === false;

  // Check if the mock user (if selected) has access to the dashboard
  $: dashboard = useDashboard($runtime.instanceId, metricViewName);
  $: mockUserHasNoAccess =
    $selectedMockUserStore && $dashboard.error?.response?.status === 404;
</script>

<section
  class="flex flex-col gap-y-1 h-full overflow-x-auto overflow-y-hidden dashboard-theme-boundary"
  use:listenToNodeResize
>
  <div
    class="border-b mb-3 w-full flex flex-col"
    id="header"
    style:padding-left={leftSide}
  >
    {#if isRillDeveloper}
      <!-- FIXME: adding an -mb-3 fixes the spacing issue incurred by changes to the header 
        to accommodate the cloud dashboard. We should go back and reconcile these headers so we 
        don't need to do this. -->
      <div
        class="flex items-center justify-between -mb-3 w-full pl-1 pr-4"
        style:height="var(--header-height)"
      >
        <DashboardTitle {metricViewName} />
        <DashboardCTAs {metricViewName} />
      </div>
    {/if}

    {#if mockUserHasNoAccess}
      <div class="mb-3" />
    {:else}
      <div class="-ml-3 p-1 py-2 space-y-2">
        <TimeControls {metricViewName} />
        {#key metricViewName}
          <Filters />
        {/key}
      </div>
    {/if}
  </div>

  {#if mockUserHasNoAccess}
    <MockUserHasNoAccess />
  {:else}
    <div
      class="flex gap-x-1 h-full overflow-hidden flex-{dashboardAlignment}"
      style:padding-left={leftSide}
    >
      <div
        class:fixed-metric-height={expandedMeasureName}
        class="overflow-y-scroll pb-8 flex-none"
      >
        {#key metricViewName}
          {#if hasTimeSeries}
            <MetricsTimeSeriesCharts
              {metricViewName}
              workspaceWidth={exploreContainerWidth}
            />
          {:else}
            <MeasuresContainer {exploreContainerWidth} {metricViewName} />
          {/if}
        {/key}
      </div>

      <div class="overflow-y-hidden grow {expandedMeasureName ? '' : 'px-4'}">
        {#if expandedMeasureName}
          <TimeDimensionDisplay {metricViewName} />
        {:else if selectedDimensionName}
          <DimensionDisplay />
        {:else}
          <LeaderboardDisplay />
        {/if}
      </div>
    </div>

    {#if isRillDeveloper && !expandedMeasureName}
      <RowsViewerAccordion {metricViewName} />
    {/if}
  {/if}
</section>

<style>
  .fixed-metric-height {
    height: 280px;
  }
</style>
