<script lang="ts" context="module">
  // only one at a time
  const globalActiveMenu = writable<string | undefined>(undefined);
</script>

<script lang="ts">
  import {
    createEventDispatcher,
    getContext,
    onMount,
    setContext,
  } from "svelte";
  import { Writable, writable } from "svelte/store";
  import { fade } from "svelte/transition";
  import { clickOutside } from "../../../lib/actions/click-outside";
  import { guidGenerator } from "../../../lib/guid";

  export let dark: boolean | undefined = undefined;
  export let maxWidth: string | undefined = undefined;
  export let minWidth = "300px";
  export let minHeight: string | undefined = undefined;
  export let maxHeight: string | undefined = undefined;
  export let paddingTop = 2;
  export let paddingBottom = 2;
  export let rounded = true;
  export let focusOnMount = true;
  export let role = "menu";
  /** used for selector-style menus */
  export let multiselectable = false;
  export let label: string | undefined = undefined;

  if (dark) {
    setContext("rill:menu:dark", dark);
  }
  const dispatch = createEventDispatcher();

  const menuID = guidGenerator();

  let key;

  function getPreviousAvailableItem(availableItems, currentItemId) {
    // get next largest thing for currentItemId;
    let nextIndex = currentItemId;
    const reversedItems = [...availableItems];
    if (currentItemId === undefined)
      return reversedItems.filter((item) => !item.disabled)?.[0]?.id;

    reversedItems.reverse();
    for (let item of reversedItems) {
      if (item.id < currentItemId && !item.disabled) {
        nextIndex = item.id;
        break;
      }
    }
    return nextIndex;
  }

  function getNextAvailableItem(availableItems, currentItemId) {
    if (currentItemId === undefined)
      return availableItems.filter((item) => !item.disabled)?.[0]?.id;
    let nextIndex = currentItemId;
    for (let item of availableItems) {
      if (item.id > currentItemId && !item.disabled) {
        nextIndex = item.id;
        break;
      }
    }
    return nextIndex;
  }

  function handleKeydown(event) {
    key = event.key;

    if (key === "Escape") {
      dispatch("escape");
    }

    if (key === "ArrowDown") {
      $currentItem = getNextAvailableItem($menuItems, $currentItem);
    }
    if (key === "ArrowUp") {
      $currentItem = getPreviousAvailableItem($menuItems, $currentItem);
    }
  }

  function onSelect() {
    dispatch("item-select");
  }

  const menuItems = writable<{ id: number; disabled: boolean }[]>([]);
  const currentItem = writable<number | undefined>(undefined);

  setContext("rill:menu:onSelect", onSelect);
  setContext("rill:menu:menuItems", menuItems);
  setContext("rill:menu:currentItem", currentItem);

  const menuTrigger: Writable<HTMLElement> =
    getContext("rill:menu:menuTrigger") || writable(undefined);

  let mounted = false;
  onMount(() => {
    $globalActiveMenu = menuID;
    mounted = true;
  });

  // once open, we should select the first menu item.
  $: if (focusOnMount && mounted) {
    $currentItem = $menuItems.find((item) => !item.disabled)?.id;
  }

  // This will effectively close any additional menus that might be open.
  $: if ($globalActiveMenu !== menuID) {
    dispatch("escape");
  }

  /** Accessibility properties */
  let ariaProperties = {};
  $: if (role === "menu") {
    ariaProperties = { role };
  } else if (role === "listbox") {
    ariaProperties = { role, ["aria-multiselectable"]: multiselectable };
  }

  // hints for tailwind re: spacing
  // pt-1, pt-2, pt-3, pb-1, pb-2, pb-3
</script>

<svelte:window on:keydown={handleKeydown} />

<div
  role="menu"
  tabindex="0"
  style:max-width={maxWidth}
  style:min-height={minHeight}
  style:max-height={maxHeight}
  style:min-width={minWidth}
  transition:fade={{ duration: 50 }}
  on:mouseleave={() => {
    $currentItem = undefined;
  }}
  use:clickOutside={[
    [$menuTrigger],
    () => {
      dispatch("click-outside");
    },
  ]}
  class:rounded
  class="
        pt-{paddingTop} 
        pb-{paddingBottom}
        w-max
        flex
        flex-col
        outline-none
        overflow-y-auto
        {dark
    ? 'bg-gray-800 dark:bg-gray-700 border-none dark:border-none shadow'
    : 'bg-white dark:bg-gray-700 border-[.5px] border-gray-300 dark:border-none shadow-lg dark:shadow-xl'}
        "
  style:outline="none"
  aria-label={label}
  {...ariaProperties}
>
  <slot />
</div>
