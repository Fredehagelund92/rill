<script lang="ts">
  import { goto } from "$app/navigation";
  import { notifications } from "@rilldata/web-common/components/notifications";
  import { renameFileArtifact } from "@rilldata/web-common/features/entity-management/actions";
  import { isDuplicateName } from "@rilldata/web-common/features/entity-management/name-utils";
  import { useAllNames } from "@rilldata/web-common/features/entity-management/resource-selectors";
  import { EntityType } from "@rilldata/web-common/features/entity-management/types";
  import { WorkspaceHeader } from "../../../layout/workspace";
  import { runtime } from "../../../runtime-client/runtime-store";
  import GoToDashboardButton from "./GoToDashboardButton.svelte";

  export let metricsDefName;

  $: runtimeInstanceId = $runtime.instanceId;
  $: allNamesQuery = useAllNames(runtimeInstanceId);

  const onChangeCallback = async (e) => {
    if (!e.target.value.match(/^[a-zA-Z_][a-zA-Z0-9_]*$/)) {
      notifications.send({
        message:
          "Dashboard name must start with a letter or underscore and contain only letters, numbers, and underscores",
      });
      e.target.value = metricsDefName; // resets the input
      return;
    }
    if (
      isDuplicateName(e.target.value, metricsDefName, $allNamesQuery.data ?? [])
    ) {
      notifications.send({
        message: `Name ${e.target.value} is already in use`,
      });
      e.target.value = metricsDefName; // resets the input
      return;
    }

    try {
      const toName = e.target.value;
      const type = EntityType.MetricsDefinition;
      await renameFileArtifact(runtimeInstanceId, metricsDefName, toName, type);
      goto(`/dashboard/${toName}/edit`, { replaceState: true });
    } catch (err) {
      console.error(err.response.data.message);
    }
  };

  $: titleInput = metricsDefName;
</script>

<WorkspaceHeader {...{ titleInput, onChangeCallback }}>
  <GoToDashboardButton {metricsDefName} slot="cta" />
</WorkspaceHeader>
