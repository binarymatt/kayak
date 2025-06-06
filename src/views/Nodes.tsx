import { useParams } from "@solidjs/router";
import { getClients } from "../client";
import { createResource, For } from "solid-js";
import { Show } from "solid-js";
import Loading from "../components/Loading";
import { node } from "../stores";
export const Nodes = () => {
  const clients = getClients(node.address);
  const [config] = createResource(async () => {
    return await clients.admin.stats({});
  });
  return (
    <Show when={!config.loading} fallback={<Loading />}>
      <div>
        <div>Current Target</div>
        <div class="text-xs uppercase font-semibold opacity-60">
          name: {node.name}
        </div>
        <div class="text-xs uppercase font-semibold opacity-60">
          Address: {node.address}
        </div>
        <div class="text-xs uppercase font-semibold opacity-60">
          State: {config().state}
        </div>
      </div>
      <ul class="list bg-base-100 rounded-box shadow-md">
        <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">
          Raft Configuration
        </li>
        <For each={config().nodes}>
          {(item) => (
            <li class="list-row">
              <div>
                <div>{item.id}</div>
                <div class="text-xs uppercase font-semibold opacity-60">
                  {item.address}
                </div>
              </div>
              <div>{item.suffrage}</div>
              <div>
                <div>{item.isLeader ? "Leader" : ""}</div>
              </div>
            </li>
          )}
        </For>
      </ul>
    </Show>
  );
};
