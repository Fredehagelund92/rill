<script lang="ts">
  import Base from "./Base.svelte";
  import { PERC_DIFF, isPercDiff } from "./type-utils";
  import type { NumberParts } from "@rilldata/web-common/lib/number-formatting/humanizer-types";
  export let isNull = false;
  export let inTable = false;
  export let dark = false;
  export let customStyle = "";
  export let value: number | undefined | null | NumberParts | PERC_DIFF;
  export let tabularNumber = true;

  let diffIsNegative = false;
  let intValue: string;
  let negSign = "";
  let approxSign = "";
  let suffix = "";

  $: isNoData = isPercDiff(value) || value === null || value === undefined;

  $: if (
    !isNoData &&
    // expanding this out in full provides type narrowing
    !isPercDiff(value) &&
    value !== null &&
    value !== undefined &&
    typeof value !== "number"
  ) {
    // in this case, we have a NumberParts object.
    // We have a couple cases to consider:
    // * If the NumberParts object has approxZero===true,
    // we want to show e.g. "~0%" WITHOUT a negative sign
    // * However, in this case we show the number in red to indicate a
    // small negative change.
    //
    // Otherwise, we format the number as usual.
    intValue = value.int;
    diffIsNegative = value?.neg === "-";
    negSign = diffIsNegative && !value?.approxZero ? "-" : "";
    approxSign = value?.approxZero ? "~" : "";
    suffix = value?.suffix ?? "";
  } else if (typeof value === "number") {
    // FIXME: this seems to only come up in the tool tip,
    // for percentages in the dimension table,
    // but this whole thing is a mess and needs to be cleaned up.

    intValue = Math.round(100 * value).toString();
    approxSign = Math.abs(value) < 0.005 ? "~" : "";
    negSign = "";
    suffix = "";
  }
</script>

<Base
  {isNull}
  classes="{tabularNumber
    ? 'ui-copy-number'
    : ''} font-normal {customStyle} {inTable && 'block text-right'}"
  {dark}
>
  <slot name="value">
    {#if isNoData}
      <span class="opacity-50 italic" style:font-size=".925em">no data</span>
    {:else if value !== null}
      <span class:text-red-500={diffIsNegative}>
        {approxSign}{negSign}{intValue}{suffix}<span class="opacity-50">%</span>
      </span>
    {/if}
  </slot>
</Base>
