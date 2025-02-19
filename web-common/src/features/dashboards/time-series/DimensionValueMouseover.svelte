<script lang="ts">
  import WithGraphicContexts from "@rilldata/web-common/components/data-graphic/functional-components/WithGraphicContexts.svelte";
  import type { NumericPlotPoint } from "@rilldata/web-common/components/data-graphic/functional-components/types";
  import MultiMetricMouseoverLabel from "@rilldata/web-common/components/data-graphic/marks/MultiMetricMouseoverLabel.svelte";
  import { bisectData } from "@rilldata/web-common/components/data-graphic/utils";
  import type { DimensionDataItem } from "@rilldata/web-common/features/dashboards/time-series/multiple-dimension-queries";
  export let point: NumericPlotPoint;
  export let xAccessor: string;
  export let yAccessor: string;
  export let mouseoverFormat;
  export let dimensionData: DimensionDataItem[];
  export let dimensionValue: string | undefined;
  export let validPercTotal: number | null;
  export let hovered: boolean | undefined;

  $: x = point?.[xAccessor];

  function truncate(str) {
    const truncateLength = 34;

    if (str.length > truncateLength) {
      // Check if last character is space
      if (str[truncateLength - 1] === " ") {
        return str.slice(0, truncateLength - 1) + "...";
      }
      return str.slice(0, truncateLength) + "...";
    }
    return str;
  }

  let pointsData = dimensionData;
  $: if (dimensionValue !== undefined) {
    const higlighted = dimensionData.filter((d) => d.value === dimensionValue);

    if (higlighted.length) {
      pointsData = higlighted;
    }
  }
  $: yValues = pointsData.map((dimension) => {
    const y = bisectData(x, "center", xAccessor, dimension?.data)[yAccessor];
    return {
      y,
      fillClass: dimension?.fillClass,
      name: dimension?.value,
    };
  });

  $: points = yValues
    .map((dimension) => {
      const y = dimension.y;
      const currentPointIsNull = y === null;
      let value = mouseoverFormat(y);
      if (validPercTotal) {
        const percOfTotal = y / validPercTotal;
        value =
          mouseoverFormat(y) + ",  " + (percOfTotal * 100).toFixed(2) + "%";
      }
      return {
        x,
        y,
        value,
        yOverride: currentPointIsNull,
        yOverrideLabel: "no current data",
        yOverrideStyleClass: `fill-gray-600 italic`,
        key: dimension.name,
        label: hovered ? truncate(dimension.name) : "",
        pointColorClass: dimension.fillClass,
        valueStyleClass: "font-bold",
        valueColorClass: "fill-gray-600",
        labelColorClass: "fill-gray-600",
        labelStyleClass: "font-semibold",
      };
    })
    .filter((d) => !d.yOverride);

  /** get the final point set*/
  $: pointSet = points;
</script>

{#if pointSet.length}
  <WithGraphicContexts>
    <MultiMetricMouseoverLabel
      isDimension={true}
      attachPointToLabel
      direction="left"
      flipAtEdge="body"
      formatValue={mouseoverFormat}
      point={pointSet || []}
    />
  </WithGraphicContexts>
{/if}
