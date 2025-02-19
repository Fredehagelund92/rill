<script lang="ts">
  /**
   * DimensionDisplay.svelte
   * -------------------------
   * Create a table with the selected dimension and measures
   * to be displayed in explore
   */
  import { cancelDashboardQueries } from "@rilldata/web-common/features/dashboards/dashboard-queries";

  import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
  import { useTimeControlStore } from "@rilldata/web-common/features/dashboards/time-controls/time-control-store";
  import {
    createQueryServiceMetricsViewComparison,
    createQueryServiceMetricsViewTotals,
  } from "@rilldata/web-common/runtime-client";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { getDimensionFilterWithSearch } from "./dimension-table-utils";
  import DimensionHeader from "./DimensionHeader.svelte";
  import DimensionTable from "./DimensionTable.svelte";
  import { notifications } from "@rilldata/web-common/components/notifications";
  import { metricsExplorerStore } from "../stores/dashboard-stores";

  const stateManagers = getStateManagers();
  const {
    dashboardStore,
    selectors: {
      dashboardQueries: {
        dimensionTableSortedQueryBody,
        dimensionTableTotalQueryBody,
      },
      comparison: { isBeingCompared },
      dimensions: { dimensionTableDimName, dimensionTableColumnName },
      dimensionTable: {
        virtualizedTableColumns,
        selectedDimensionValueNames,
        prepareDimTableRows,
      },
      activeMeasure: { activeMeasureName },
    },
    metricsViewName,
    runtime,
  } = stateManagers;

  // cast is safe because dimensionTableDimName must be defined
  // for the dimension table to be open
  $: dimensionName = $dimensionTableDimName as string;
  $: dimensionColumnName = $dimensionTableColumnName(dimensionName) as string;

  let searchText = "";

  const queryClient = useQueryClient();

  $: instanceId = $runtime.instanceId;

  const timeControlsStore = useTimeControlStore(stateManagers);

  $: filterSet = getDimensionFilterWithSearch(
    $dashboardStore?.filters,
    searchText,
    dimensionName
  );

  $: totalsQuery = createQueryServiceMetricsViewTotals(
    instanceId,
    $metricsViewName,
    $dimensionTableTotalQueryBody,
    {
      query: {
        enabled: $timeControlsStore.ready,
      },
    }
  );

  $: unfilteredTotal = $totalsQuery?.data?.data?.[$activeMeasureName] ?? 0;

  $: columns = $virtualizedTableColumns($totalsQuery);

  $: sortedQuery = createQueryServiceMetricsViewComparison(
    $runtime.instanceId,
    $metricsViewName,
    $dimensionTableSortedQueryBody,
    {
      query: {
        enabled: $timeControlsStore.ready && !!filterSet,
      },
    }
  );

  $: tableRows = $prepareDimTableRows($sortedQuery, unfilteredTotal);

  $: areAllTableRowsSelected = tableRows.every((row) =>
    $selectedDimensionValueNames.includes(row[dimensionColumnName] as string)
  );

  function onSelectItem(event) {
    const label = tableRows[event.detail][dimensionColumnName] as string;
    cancelDashboardQueries(queryClient, $metricsViewName);
    metricsExplorerStore.toggleFilter($metricsViewName, dimensionName, label);
  }

  function toggleComparisonDimension(dimensionName, isBeingCompared) {
    metricsExplorerStore.setComparisonDimension(
      $metricsViewName,
      isBeingCompared ? undefined : dimensionName
    );
  }

  function toggleAllSearchItems() {
    const labels = tableRows.map((row) => row[dimensionColumnName] as string);
    cancelDashboardQueries(queryClient, $metricsViewName);

    if (areAllTableRowsSelected) {
      metricsExplorerStore.deselectItemsInFilter(
        $metricsViewName,
        dimensionName,
        labels
      );

      notifications.send({
        message: `Removed ${labels.length} items from filter`,
      });
      return;
    } else {
      const newValuesSelected = metricsExplorerStore.selectItemsInFilter(
        $metricsViewName,
        dimensionName,
        labels
      );
      notifications.send({
        message: `Added ${newValuesSelected} items to filter`,
      });
    }
  }

  function handleKeyDown(e) {
    // Select all items on Meta+A
    if ((e.ctrlKey || e.metaKey) && e.key === "a") {
      if (e.target.tagName === "INPUT") return;
      e.preventDefault();
      if (areAllTableRowsSelected) return;
      toggleAllSearchItems();
    }
  }
</script>

{#if sortedQuery}
  <div class="h-full flex flex-col" style:min-width="365px">
    <div class="flex-none" style:height="50px">
      <DimensionHeader
        {dimensionName}
        {areAllTableRowsSelected}
        isRowsEmpty={!tableRows.length}
        isFetching={$sortedQuery?.isFetching}
        on:search={(event) => {
          searchText = event.detail;
        }}
        on:toggle-all-search-items={() => toggleAllSearchItems()}
      />
    </div>

    {#if tableRows && columns.length && dimensionName}
      <div class="grow" style="overflow-y: hidden;">
        <DimensionTable
          on:select-item={(event) => onSelectItem(event)}
          on:toggle-dimension-comparison={() =>
            toggleComparisonDimension(dimensionName, isBeingCompared)}
          isFetching={$sortedQuery?.isFetching}
          {dimensionName}
          {columns}
          selectedValues={$selectedDimensionValueNames}
          rows={tableRows}
        />
      </div>
    {/if}
  </div>
{/if}

<svelte:window on:keydown={handleKeyDown} />
