<script lang="ts">
  import { page } from "$app/stores";
  import ContentContainer from "@rilldata/web-admin/components/layout/ContentContainer.svelte";
  import VerticalScrollContainer from "@rilldata/web-common/layout/VerticalScrollContainer.svelte";
  import { createAdminServiceGetProject } from "../../../client";
  import DashboardsTable from "../../../features/dashboards/listing/DashboardsTable.svelte";
  import RedeployProjectCta from "../../../features/projects/RedeployProjectCTA.svelte";

  $: organization = $page.params.organization;
  $: project = $page.params.project;

  $: proj = createAdminServiceGetProject(organization, project);
  $: isProjectDeployed = $proj?.data && $proj.data.prodDeployment;
  $: isProjectHibernating = $proj?.data && !$proj.data.prodDeployment;
</script>

<svelte:head>
  <title>{project} overview - Rill</title>
</svelte:head>
<VerticalScrollContainer>
  {#if isProjectHibernating}
    <RedeployProjectCta {organization} {project} />
  {:else if isProjectDeployed}
    <ContentContainer>
      <div class="flex flex-col items-center gap-y-4">
        <DashboardsTable {organization} {project} />
      </div>
    </ContentContainer>
  {/if}
</VerticalScrollContainer>
