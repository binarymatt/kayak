import { batch, Component, createSignal, For, JSX } from "solid-js";
import { clusters, setClusters, removeIndex, setNode, node } from "../stores";
import { VsTrash } from "solid-icons/vs";
import { getKayakClient } from "../client";

export const Dashboard: Component = () => {
  const [name, setName] = createSignal(node.name);
  const [address, setAddress] = createSignal(node.address);
  const addNode = () => {
    // setAddress(event.target.elements["address"].value);

    batch(() => {
      setNode({
        name: name(),
        address: address(),
      });
      //setName("");
      //setAddress("");
    });
  };
  return (
    <div class="h-screen flex items-center justify-center">
      <div class="card bg-base-100 w-96 shadow-sm">
        <div class="card-body">
          <h2 class="card-title">Current Node</h2>
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

          <p></p>
          <div class="card-actions justify-end">
            <button class="btn btn-primary" onClick={addNode}>
              Save
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
