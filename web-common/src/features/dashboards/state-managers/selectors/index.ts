import { sortingSelectors } from "./sorting";
import { derived, type Readable } from "svelte/store";
import type { MetricsExplorerEntity } from "../../stores/metrics-explorer-entity";
import type { ReadablesObj, SelectorFnsObj } from "./types";
import type { QueryObserverResult } from "@tanstack/svelte-query";
import type {
  RpcStatus,
  V1ColumnTimeRangeResponse,
  V1MetricsViewSpec,
} from "@rilldata/web-common/runtime-client";
import { formattingSelectors } from "./data-formatting";
import { contextColSelectors } from "./context-column";
import { activeMeasureSelectors } from "./active-measure";
import { dimensionSelectors } from "./dimensions";
import { dimensionFilterSelectors } from "./dimension-filters";
import { timeRangeSelectors } from "./time-range";
import { leaderboardQuerySelectors } from "./dashboard-queries";
import { comparisonSelectors } from "./comparisons";
import { dimensionTableSelectors } from "./dimension-table";
import { measureSelectors } from "./measures";

export type DashboardDataReadables = {
  dashboardStore: Readable<MetricsExplorerEntity>;
  metricsSpecQueryResultStore: Readable<
    QueryObserverResult<V1MetricsViewSpec, RpcStatus>
  >;
  timeRangeSummaryStore: Readable<
    QueryObserverResult<V1ColumnTimeRangeResponse, unknown>
  >;
};

export type StateManagerReadables = ReturnType<
  typeof createStateManagerReadables
>;

export const createStateManagerReadables = (
  dashboardDataReadables: DashboardDataReadables
) => {
  return {
    /**
     * Readables related to the sorting state of the dashboard.
     */
    sorting: createReadablesFromSelectors(
      sortingSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to number formatting for the dashboard.
     */
    numberFormat: createReadablesFromSelectors(
      formattingSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to the dashboard context column.
     */
    contextColumn: createReadablesFromSelectors(
      contextColSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to the primary active measure in the
     * leaderboard.
     */
    activeMeasure: createReadablesFromSelectors(
      activeMeasureSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to the dimensions available in the
     * leaderboard.
     */
    dimensions: createReadablesFromSelectors(
      dimensionSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to the dimension dimension.
     *
     * These are valid when the dimension table is visible, and
     * should only be used from within dimension table components.
     */
    dimensionTable: createReadablesFromSelectors(
      dimensionTableSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to selected (aka "filtered)
     * dimension values in the leaderboard, including
     * whether or not a dimension is in include or exclude mode.
     */
    dimensionFilters: createReadablesFromSelectors(
      dimensionFilterSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to the state of the time range selector
     * for the dashboard.
     */
    timeRangeSelectors: createReadablesFromSelectors(
      timeRangeSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to the dashboard comparison state
     */
    comparison: createReadablesFromSelectors(
      comparisonSelectors,
      dashboardDataReadables
    ),

    /**
     * Readables for query construction
     */
    dashboardQueries: createReadablesFromSelectors(
      leaderboardQuerySelectors,
      dashboardDataReadables
    ),

    /**
     * Readables related to dashboard measures
     */
    measures: createReadablesFromSelectors(
      measureSelectors,
      dashboardDataReadables
    ),
  };
};

function createReadablesFromSelectors<T extends SelectorFnsObj>(
  selectors: T,
  readables: DashboardDataReadables
): ReadablesObj<T> {
  return Object.fromEntries(
    Object.entries(selectors).map(([key, selectorFn]) => [
      key,
      derived(
        // Note: creating a svelte derived store from multiple stores
        // requires supplying a tuple of stores.
        // To simplify the selector function, we pack this into a single
        // selectorFnArgs object.
        [
          readables.dashboardStore,
          readables.metricsSpecQueryResultStore,
          readables.timeRangeSummaryStore,
        ],
        ([dashboard, metricsSpecQueryResult, timeRangeSummary]) =>
          selectorFn({
            dashboard,
            metricsSpecQueryResult,
            timeRangeSummary,
          })
      ),
    ])
  ) as ReadablesObj<T>;
}
