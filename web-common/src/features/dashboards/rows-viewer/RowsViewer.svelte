<script lang="ts">
  import { getStateManagers } from "@rilldata/web-common/features/dashboards/state-managers/state-managers";
  import { useTimeControlStore } from "@rilldata/web-common/features/dashboards/time-controls/time-control-store";
  import { runtime } from "@rilldata/web-common/runtime-client/runtime-store";
  import { useMetaQuery } from "../selectors";
  import {
    createQueryServiceMetricsViewRows,
    createQueryServiceTableColumns,
  } from "@rilldata/web-common/runtime-client";
  import { useDashboardStore } from "web-common/src/features/dashboards/stores/dashboard-stores";
  import PreviewTable from "@rilldata/web-common/components/preview-table/PreviewTable.svelte";
  import type { VirtualizedTableColumns } from "@rilldata/web-local/lib/types";
  import { writable } from "svelte/store";

  export let metricViewName = "";

  const SAMPLE_SIZE = 10000;
  const FALLBACK_SAMPLE_SIZE = 1000;

  $: dashboardStore = useDashboardStore(metricViewName);
  const timeControlsStore = useTimeControlStore(getStateManagers());

  $: modelName = useMetaQuery<string>(
    $runtime.instanceId,
    metricViewName,
    (data) => data.table ?? ""
  );

  $: name = $modelName?.data ?? "";

  let limit = writable(SAMPLE_SIZE);

  $: tableQuery = createQueryServiceMetricsViewRows(
    $runtime?.instanceId,
    metricViewName,
    {
      limit: $limit,
      filter: $dashboardStore.filters,
      timeStart: $timeControlsStore.timeStart,
      timeEnd: $timeControlsStore.timeEnd,
    },
    {
      query: {
        enabled: $timeControlsStore.ready && !!$dashboardStore?.filters,
      },
    }
  );

  // If too much date is requested, limit the query to 1000 rows
  $: if (
    // @ts-ignore
    $tableQuery?.error?.response?.data?.code === 8 &&
    $limit > FALLBACK_SAMPLE_SIZE
  ) {
    // SK: Have to set the limit on the next tick or the tableQuery does not update. Not sure why, seems like a svelte-query issue.
    setTimeout(() => {
      limit.set(FALLBACK_SAMPLE_SIZE);
    });
  }

  let rows;
  $: {
    if ($tableQuery.isSuccess) {
      rows = $tableQuery.data.data;
    }
  }

  $: profileColumnsQuery = createQueryServiceTableColumns(
    $runtime?.instanceId,
    name,
    {}
  );
  $: profileColumns = $profileColumnsQuery?.data
    ?.profileColumns as VirtualizedTableColumns[];

  let rowOverscanAmount = 0;
  let columnOverscanAmount = 0;

  const configOverride = {
    indexWidth: 72,
    rowHeight: 32,
  };
</script>

<div class="h-72 overflow-y-auto bg-gray-100 border-t border-gray-200">
  {#if rows}
    <PreviewTable
      {rows}
      columnNames={profileColumns}
      {rowOverscanAmount}
      {columnOverscanAmount}
      {configOverride}
    />
  {/if}
</div>
