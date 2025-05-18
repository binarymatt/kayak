import { useParams } from "@solidjs/router";
import { getClients } from "../client";
import { createResource, For } from "solid-js";
import { Show } from "solid-js";
import Loading from "../components/Loading";
export const Nodes = () => {
  const clients = getClients("http://localhost:8080");
  const [config] = createResource(async () => {
    return await clients.admin.stats({}).then((res) => res.stats);
  });
  return (
    <Show when={!config.loading} fallback={<Loading />}>
      <ul class="list bg-base-100 rounded-box shadow-md">
        <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">
          Raft Configuration
        </li>
        <For each={config().latest_configuration}>
          {(item, index) => <div>{item.Suffrage}</div>}
        </For>
      </ul>
      <pre>{config()}</pre>
      <pre>{config().latest_configuration}</pre>
      <ul class="list bg-base-100 rounded-box shadow-md">
        <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">
          Most played songs this week
        </li>

        <li class="list-row">
          <div>
            <img
              class="size-10 rounded-box"
              src="https://img.daisyui.com/images/profile/demo/1@94.webp"
            />
          </div>
          <div>
            <div>Dio Lupa</div>
            <div class="text-xs uppercase font-semibold opacity-60">
              Remaining Reason
            </div>
          </div>
          <p class="list-col-wrap text-xs">
            "Remaining Reason" became an instant hit, praised for its haunting
            sound and emotional depth. A viral performance brought it widespread
            recognition, making it one of Dio Lupa’s most iconic tracks.
          </p>
        </li>

        <li class="list-row">
          <div>
            <img
              class="size-10 rounded-box"
              src="https://img.daisyui.com/images/profile/demo/4@94.webp"
            />
          </div>
          <div>
            <div>Ellie Beilish</div>
            <div class="text-xs uppercase font-semibold opacity-60">
              Bears of a fever
            </div>
          </div>
          <p class="list-col-wrap text-xs">
            "Bears of a Fever" captivated audiences with its intense energy and
            mysterious lyrics. Its popularity skyrocketed after fans shared it
            widely online, earning Ellie critical acclaim.
          </p>
        </li>

        <li class="list-row">
          <div>
            <img
              class="size-10 rounded-box"
              src="https://img.daisyui.com/images/profile/demo/3@94.webp"
            />
          </div>
          <div>
            <div>Sabrino Gardener</div>
            <div class="text-xs uppercase font-semibold opacity-60">
              Cappuccino
            </div>
          </div>
          <p class="list-col-wrap text-xs">
            "Cappuccino" quickly gained attention for its smooth melody and
            relatable themes. The song’s success propelled Sabrino into the
            spotlight, solidifying their status as a rising star.
          </p>
        </li>
      </ul>
    </Show>
  );
};
/*
 * <For each={config().latest_configuration}>
          {(item) => (
            <li class="list-row">
              <div>{item.id}</div>
            </li>
          )}
        </For>
 */
