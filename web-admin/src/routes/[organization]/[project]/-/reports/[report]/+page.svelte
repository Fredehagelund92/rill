<script lang="ts">
  import { page } from "$app/stores";
  import ContentContainer from "@rilldata/web-admin/components/layout/ContentContainer.svelte";
  import ReportHistoryTable from "@rilldata/web-admin/features/scheduled-reports/history/ReportHistoryTable.svelte";
  import ReportMetadata from "@rilldata/web-admin/features/scheduled-reports/metadata/ReportMetadata.svelte";
  import VerticalScrollContainer from "@rilldata/web-common/layout/VerticalScrollContainer.svelte";

  $: organization = $page.params.organization;
  $: project = $page.params.project;
  $: report = $page.params.report;
</script>

<VerticalScrollContainer>
  <ContentContainer>
    <div class="flex justify-center">
      <div class="w-[960px] flex flex-col items-start gap-y-9">
        <ReportMetadata {organization} {project} {report} />
        <ReportHistoryTable {report} />
      </div>
    </div>
  </ContentContainer>
</VerticalScrollContainer>
