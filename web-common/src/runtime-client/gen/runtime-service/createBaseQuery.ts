import {
  notifyManager,
  type QueryKey,
  type QueryObserver,
} from "@tanstack/query-core";
import type {
  CreateBaseQueryOptions as _CreateBaseQueryOptions,
  CreateBaseQueryResult,
  QueryClient,
} from "@tanstack/svelte-query";
import { useQueryClient } from "@tanstack/svelte-query";
import { derived, readable } from "svelte/store";

type CreateBaseQueryOptions<
  T,
  R,
  S,
  U,
  W extends QueryKey
> = _CreateBaseQueryOptions<T, R, S, U, W> & {
  queryClient?: QueryClient;
};

export function createBaseQuery<
  TQueryFnData,
  TError,
  TData,
  TQueryData,
  TQueryKey extends QueryKey
>(
  options: CreateBaseQueryOptions<
    TQueryFnData,
    TError,
    TData,
    TQueryData,
    TQueryKey
  >,
  Observer: typeof QueryObserver
): CreateBaseQueryResult<TData, TError> {
  const queryClient = options.queryClient ?? useQueryClient();
  const defaultedOptions = queryClient.defaultQueryOptions(options);
  defaultedOptions._optimisticResults = "optimistic";

  let observer = new Observer<
    TQueryFnData,
    TError,
    TData,
    TQueryData,
    TQueryKey
  >(queryClient, defaultedOptions);

  // Include callbacks in batch renders
  if (defaultedOptions.onError) {
    defaultedOptions.onError = notifyManager.batchCalls(
      defaultedOptions.onError
    );
  }

  if (defaultedOptions.onSuccess) {
    defaultedOptions.onSuccess = notifyManager.batchCalls(
      defaultedOptions.onSuccess
    );
  }

  if (defaultedOptions.onSettled) {
    defaultedOptions.onSettled = notifyManager.batchCalls(
      defaultedOptions.onSettled
    );
  }

  readable(observer).subscribe(($observer) => {
    observer = $observer;
    // Do not notify on updates because of changes in the options because
    // these changes should already be reflected in the optimistic result.
    observer.setOptions(defaultedOptions, { listeners: false });
  });

  const result = readable(observer.getCurrentResult(), (set) => {
    return observer.subscribe(notifyManager.batchCalls(set));
  });

  const { subscribe } = derived(result, ($result) => {
    $result = observer.getOptimisticResult(defaultedOptions);
    return !defaultedOptions.notifyOnChangeProps
      ? observer.trackResult($result)
      : $result;
  });

  return { subscribe };
}
