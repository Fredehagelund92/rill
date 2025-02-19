<script lang="ts">
  import { beforeNavigate } from "$app/navigation";
  import { page } from "$app/stores";
  import { initCloudMetrics } from "@rilldata/web-admin/features/telemetry/initCloudMetrics";
  import NotificationCenter from "@rilldata/web-common/components/notifications/NotificationCenter.svelte";
  import {
    featureFlags,
    retainFeaturesFlags,
  } from "@rilldata/web-common/features/feature-flags";
  import RillTheme from "@rilldata/web-common/layout/RillTheme.svelte";
  import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
  import { onMount } from "svelte";
  import ErrorBoundary from "../features/errors/ErrorBoundary.svelte";
  import {
    addJavascriptErrorListeners,
    createGlobalErrorCallback,
  } from "../features/errors/error-utils";
  import TopNavigationBar from "../features/navigation/TopNavigationBar.svelte";
  import { clearViewedAsUserAfterNavigate } from "../features/view-as-user/clearViewedAsUser";

  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnMount: false,
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
        retry: false,
      },
    },
  });
  // Motivation:
  // - https://tkdodo.eu/blog/breaking-react-querys-api-on-purpose#a-bad-api
  // - https://tkdodo.eu/blog/react-query-error-handling#the-global-callbacks
  queryClient.getQueryCache().config.onError =
    createGlobalErrorCallback(queryClient);

  featureFlags.set({
    // The admin server enables some dashboard features like scheduled reports and alerts
    adminServer: true,
    // Set read-only mode so that the user can't edit the dashboard
    readOnly: true,
  });

  beforeNavigate(retainFeaturesFlags);
  clearViewedAsUserAfterNavigate(queryClient);
  initCloudMetrics();

  onMount(() => addJavascriptErrorListeners());

  $: isEmbed = $page.url.pathname === "/-/embed";
</script>

<svelte:head>
  <meta content="Rill Cloud" name="description" />
</svelte:head>

<RillTheme>
  <QueryClientProvider client={queryClient}>
    <main class="flex flex-col h-screen">
      {#if !isEmbed}
        <TopNavigationBar />
      {/if}
      <div class="flex-grow overflow-hidden">
        <ErrorBoundary>
          <slot />
        </ErrorBoundary>
      </div>
    </main>
  </QueryClientProvider>

  <NotificationCenter />
</RillTheme>
