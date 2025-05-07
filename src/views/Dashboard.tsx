import { batch, Component, createSignal, For, JSX } from "solid-js";
import { clusters, setClusters, removeIndex } from "../stores";
import { VsTrash } from "solid-icons/vs";

export const Dashboard: Component = () => {
  const [name, setName] = createSignal("");
  const [address, setAddress] = createSignal("");
  const addCluster = (event: SubmitEvent) => {
    event.preventDefault();
    // setAddress(event.target.elements["address"].value);
    batch(() => {
      setClusters(clusters.length, {
        name: name(),
        address: address(),
      });
      setName("");
      setAddress("");
    });
  };
  return (
    <div class="h-screen flex items-center justify-center">
      <ul class="list bg-base-100 rounded-box shadow-md">
        <li class="p-4 pb-2 text-xs opacity-60 tracking-wide">Clusters</li>
        <For each={clusters}>
          {(cluster, i) => (
            <li class="list-row">
              <div class="text-4xl font-thin opacity-30 tabular-nums">
                {i()}
              </div>
              <div class="list-col-grow">
                <div>{cluster.name}</div>
                <div class="text-xs uppercase font-semibold opacity-60">
                  {cluster.address}
                </div>
              </div>
              <button
                class="btn btn-square btn-ghost"
                onClick={() => setClusters((t) => removeIndex(t, i()))}
              >
                <svg
                  fill="currentColor"
                  stroke-width="0"
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 16 16"
                  height="1em"
                  width="1em"
                  style="overflow: visible; color: currentcolor;"
                >
                  <path
                    fill-rule="evenodd"
                    d="M10 3h3v1h-1v9l-1 1H4l-1-1V4H2V3h3V2a1 1 0 0 1 1-1h3a1 1 0 0 1 1 1v1zM9 2H6v1h3V2zM4 13h7V4H4v9zm2-8H5v7h1V5zm1 0h1v7H7V5zm2 0h1v7H9V5z"
                    clip-rule="evenodd"
                  ></path>
                </svg>
              </button>
            </li>
          )}
        </For>
        <li>
          <form onSubmit={addCluster}>
            <div class="form-control">
              <input
                required
                type="text"
                placeholder="Cluster Name"
                class="input input-bordered"
                value={name()}
                onInput={(e) => setName(e.currentTarget.value)}
              />
            </div>
            <div class="form-control">
              <input
                required
                type="text"
                placeholder="Cluster Address"
                class="input input-bordered"
                value={address()}
                onInput={(e) => setAddress(e.currentTarget.value)}
              />
            </div>
            <button type="submit" class="btn btn-neutral">
              Save
            </button>
          </form>
        </li>
      </ul>
    </div>
  );
};

