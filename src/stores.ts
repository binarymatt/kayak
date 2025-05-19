import { createStore, type SetStoreFunction, type Store } from "solid-js/store";
import { createSignal, type Setter, Accessor } from "solid-js";
import { createEffect } from "solid-js";
export function createLocalStore<T extends object>(
  name: string,
  init: T,
): [Store<T>, SetStoreFunction<T>] {
  const localState = localStorage.getItem(name);
  const [state, setState] = createStore<T>(
    localState ? JSON.parse(localState) : init,
  );
  createEffect(() => localStorage.setItem(name, JSON.stringify(state)));
  return [state, setState];
}

function createSimpleStore(
  name: string,
  init: string,
): [Accessor<string>, Setter<string>] {
  const localState = localStorage.getItem(name);
  const [state, setState] = createSignal<string>(
    localState ? localState : init,
  );
  createEffect(() => localStorage.setItem(name, state()));
  return [state, setState];
}
type Cluster = { name: string; address: string };
export const [clusters, setClusters] = createLocalStore<Cluster[]>(
  "clusters",
  [],
);

export function removeIndex<T>(array: readonly T[], index: number): T[] {
  return [...array.slice(0, index), ...array.slice(index + 1)];
}
