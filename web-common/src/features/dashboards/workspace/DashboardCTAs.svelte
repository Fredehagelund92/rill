<script lang="ts">
  import { goto } from "$app/navigation";
  import MetricsIcon from "@rilldata/web-common/components/icons/Metrics.svelte";
  import PanelCTA from "@rilldata/web-common/components/panel/PanelCTA.svelte";
  import { useDashboard } from "@rilldata/web-common/features/dashboards/selectors";
  import { V1ReconcileStatus } from "@rilldata/web-common/runtime-client";
  import { Button } from "../../../components/button";
  import Tooltip from "../../../components/tooltip/Tooltip.svelte";
  import TooltipContent from "../../../components/tooltip/TooltipContent.svelte";
  import { behaviourEvent } from "../../../metrics/initMetrics";
  import { BehaviourEventMedium } from "../../../metrics/service/BehaviourEventTypes";
  import {
    MetricsEventScreenName,
    MetricsEventSpace,
  } from "../../../metrics/service/MetricsTypes";
  import { runtime } from "../../../runtime-client/runtime-store";
  import { featureFlags } from "../../feature-flags";
  import { useDashboardPolicyCheck } from "../granular-access-policies/useDashboardPolicyCheck";
  import ViewAsButton from "../granular-access-policies/ViewAsButton.svelte";
  import DeployDashboardCta from "./DeployDashboardCTA.svelte";

  export let metricViewName: string;

  $: dashboardPolicyCheck = useDashboardPolicyCheck(
    $runtime.instanceId,
    metricViewName
  );

  $: isEditableDashboard = $featureFlags.readOnly === false;

  $: dashboardQuery = useDashboard($runtime.instanceId, metricViewName);
  $: dashboardIsIdle =
    $dashboardQuery.data?.meta?.reconcileStatus ===
    V1ReconcileStatus.RECONCILE_STATUS_IDLE;

  function viewMetrics(metricViewName: string) {
    goto(`/dashboard/${metricViewName}/edit`);

    behaviourEvent.fireNavigationEvent(
      metricViewName,
      BehaviourEventMedium.Button,
      MetricsEventSpace.Workspace,
      MetricsEventScreenName.Dashboard,
      MetricsEventScreenName.MetricsDefinition
    );
  }

  let showDeployDashboardModal = false;
</script>

<PanelCTA side="right">
  {#if $dashboardPolicyCheck.data}
    <ViewAsButton />
  {/if}
  {#if isEditableDashboard}
    <Tooltip distance={8}>
      <Button
        disabled={!dashboardIsIdle}
        on:click={() => viewMetrics(metricViewName)}
        type="secondary"
      >
        Edit Metrics <MetricsIcon size="16px" />
      </Button>
      <TooltipContent slot="tooltip-content">
        {#if !dashboardIsIdle}
          Dependencies are being ingested
        {:else}
          Edit this dashboard's metrics & settings
        {/if}
      </TooltipContent>
    </Tooltip>
    <Tooltip distance={8}>
      <Button on:click={() => (showDeployDashboardModal = true)} type="primary">
        Deploy
      </Button>
      <TooltipContent slot="tooltip-content">
        Deploy this dashboard to Rill Cloud
      </TooltipContent>
    </Tooltip>
  {/if}
</PanelCTA>

<DeployDashboardCta
  on:close={() => (showDeployDashboardModal = false)}
  open={showDeployDashboardModal}
/>
