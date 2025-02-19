<script lang="ts">
  import { page } from "$app/stores";
  import { IconButton } from "@rilldata/web-common/components/button";
  import HideRightSidebar from "@rilldata/web-common/components/icons/HideRightSidebar.svelte";
  import SlidingWords from "@rilldata/web-common/components/tooltip/SlidingWords.svelte";
  import Tooltip from "@rilldata/web-common/components/tooltip/Tooltip.svelte";
  import TooltipContent from "@rilldata/web-common/components/tooltip/TooltipContent.svelte";
  import { createResizeListenerActionFactory } from "@rilldata/web-common/lib/actions/create-resize-listener-factory";
  import { dynamicTextInputWidth } from "@rilldata/web-common/lib/actions/dynamic-text-input-width";
  import { getContext } from "svelte";
  import type { Tweened } from "svelte/motion";
  import type { Writable } from "svelte/store";
  import SourceUnsavedIndicator from "../../features/sources/editor/SourceUnsavedIndicator.svelte";
  import type { LayoutElement } from "./types";

  export let onChangeCallback;
  export let titleInput;
  export let editable = true;
  export let showInspectorToggle = true;

  let titleInputElement;
  let editingTitle = false;

  let tooltipActive;

  const { listenToNodeResize, observedNode } =
    createResizeListenerActionFactory();

  const inspectorLayout = getContext(
    "rill:app:inspector-layout"
  ) as Writable<LayoutElement>;

  const navigationVisibilityTween = getContext(
    "rill:app:navigation-visibility-tween"
  ) as Tweened<number>;

  function onKeydown(event) {
    if (editingTitle && event.key === "Enter") {
      titleInputElement.blur();
    }
  }

  $: width = $observedNode?.getBoundingClientRect()?.width;

  function onInput() {
    if (editable) {
      editingTitle = true;
    }
  }
</script>

<svelte:window on:keydown={onKeydown} />
<header
  class="grid items-center content-stretch justify-between pl-4 border-b border-gray-300"
  style:grid-template-columns="[title] minmax(0, 1fr) [controls] auto"
  style:height="var(--header-height)"
  use:listenToNodeResize
>
  <div style:padding-left="{$navigationVisibilityTween * 24}px">
    {#if titleInput !== undefined && titleInput !== null}
      <h1
        style:font-size="16px"
        class="grid grid-flow-col justify-start items-center gap-x-1 overflow-hidden"
      >
        <Tooltip
          distance={8}
          alignment="start"
          bind:active={tooltipActive}
          suppress={editingTitle || !editable}
        >
          <input
            autocomplete="off"
            disabled={!editable}
            id="model-title-input"
            bind:this={titleInputElement}
            on:focus={() => {
              editingTitle = true;
            }}
            on:input={onInput}
            class="bg-transparent border border-transparent border-2 {editable
              ? 'hover:border-gray-400 cursor-pointer'
              : ''} rounded pl-2 pr-2"
            class:font-bold={editingTitle === false}
            on:blur={() => {
              editingTitle = false;
            }}
            value={titleInput}
            use:dynamicTextInputWidth
            on:change={onChangeCallback}
          />
          <TooltipContent slot="tooltip-content">
            <div class="flex items-center gap-x-2">Edit</div>
          </TooltipContent>
        </Tooltip>

        {#if $page.url.pathname.startsWith("/source")}
          <SourceUnsavedIndicator sourceName={titleInput} />
        {/if}
      </h1>
    {/if}
  </div>

  <div class="flex items-center mr-4">
    <slot name="workspace-controls" {width} />
    {#if showInspectorToggle}
      <IconButton
        on:click={() => {
          inspectorLayout.update((state) => {
            state.visible = !state.visible;
            return state;
          });
        }}
      >
        <span class="text-gray-500">
          <HideRightSidebar size="18px" />
        </span>
        <svelte:fragment slot="tooltip-content">
          <SlidingWords
            active={$inspectorLayout?.visible}
            direction="horizontal"
            reverse>inspector</SlidingWords
          >
        </svelte:fragment>
      </IconButton>
    {/if}

    <div class="pl-4">
      <slot name="cta" {width} />
    </div>
  </div>
</header>
