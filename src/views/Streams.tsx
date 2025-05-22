import { createResource, Resource, For, createSignal, batch } from "solid-js";
import { A } from "@solidjs/router";
import { getKayakClient } from "../client";
import { Stream } from "../gen/kayak/v1/model_pb";
import { addAlert, node } from "../stores";

const item = (req: Resource<Stream[]>) => {
  if (req.loading) {
    return false;
  }
  return req();
};
export const Streams = () => {
  const client = getKayakClient(node.address);
  const fetchStreams = async () => {
    return client.getStreams({}).then((res) => res.streams);
  };

  const [name, setName] = createSignal("");
  const [count, setCount] = createSignal(1);
  const [ttl, setTTL] = createSignal(300);
  const [streams, { refetch }] = createResource(fetchStreams);
  const deleteStream = async (name: string) => {
    batch(async () => {
      await client.deleteStream({ name });
      refetch();
    });
  };
  const saveContent = () => {
    batch(async () => {
      // TOOD: validation
      await client.createStream({
        name: name(),
        partitionCount: BigInt(count()),
        ttl: BigInt(ttl()),
      });
      setName("");
      setCount(1);
      setTTL(300);
      refetch();
    });
  };
  return (
    <>
      <div class="flex flex-row-reverse">
        <ul class="menu menu-vertical lg:menu-horizontal bg-base-200 rounded-box">
          <li>
            <a onclick="new_stream_modal.showModal()">Create Stream</a>
          </li>
        </ul>
      </div>
      <For each={item(streams)}>
        {(stream) => (
          <div class="card bg-base-100 w-150 shadow-sm p-4 pb-2">
            <div class="card-body">
              <div class="card-actions justify-end">
                <button
                  class="btn btn-square btn-sm"
                  onClick={() => deleteStream(stream.name)}
                >
                  <svg
                    fill="currentColor"
                    stroke-width="0"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 448 512"
                    height="1em"
                    width="1em"
                    style="overflow: visible; color: currentcolor;"
                  >
                    <path d="m170.5 51.6-19 28.4h145l-19-28.4c-1.5-2.2-4-3.6-6.7-3.6h-93.7c-2.7 0-5.2 1.3-6.7 3.6zm147-26.6 36.7 55H424c13.3 0 24 10.7 24 24s-10.7 24-24 24h-8v304c0 44.2-35.8 80-80 80H112c-44.2 0-80-35.8-80-80V128h-8c-13.3 0-24-10.7-24-24s10.7-24 24-24h69.8l36.7-55.1C140.9 9.4 158.4 0 177.1 0h93.7c18.7 0 36.2 9.4 46.6 24.9zM80 128v304c0 17.7 14.3 32 32 32h224c17.7 0 32-14.3 32-32V128H80zm80 64v208c0 8.8-7.2 16-16 16s-16-7.2-16-16V192c0-8.8 7.2-16 16-16s16 7.2 16 16zm80 0v208c0 8.8-7.2 16-16 16s-16-7.2-16-16V192c0-8.8 7.2-16 16-16s16 7.2 16 16zm80 0v208c0 8.8-7.2 16-16 16s-16-7.2-16-16V192c0-8.8 7.2-16 16-16s16 7.2 16 16z"></path>
                  </svg>
                </button>
              </div>
              <h2 class="card-title">
                <A href={"/stream/" + stream.name}>{stream.name}</A>
              </h2>
              <div class="stats ">
                <div class="stat place-items-center">
                  <div class="stat-title">TTL (time to live)</div>
                  <div class="stat-value">{stream.ttl + ""}</div>
                </div>
                <div class="stat place-items-center">
                  <div class="stat-title">Partition Count</div>
                  <div class="stat-value">{stream.partitionCount + ""}</div>
                </div>
              </div>
            </div>
          </div>
        )}
      </For>
      <dialog id="new_stream_modal" class="modal">
        <div class="modal-box">
          <h3 class="text-lg font-bold">Create Stream</h3>
          <p class="py-4">Press ESC key or click the button below to close</p>
          <fieldset class="fieldset">
            <legend class="fieldset-legend">Stream Name</legend>
            <input
              type="text"
              class="input"
              placeholder="Name here"
              value={name()}
              onInput={(e) => setName(e.currentTarget.value)}
            />
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-legend">Number of partitions</legend>
            <input
              type="number"
              class="input"
              placeholder="Partition count"
              value={count()}
              onInput={(e) => setCount(e.currentTarget.value)}
            />
            <p class="label">How many partitions to segment data by.</p>
          </fieldset>
          <fieldset class="fieldset">
            <legend class="fieldset-legend">Record TTL</legend>
            <input
              type="number"
              class="input"
              value={ttl()}
              onInput={(e) => setTTL(e.currentTarget.value)}
            />
            <p class="label">
              Time in seconds that a record is visible in stream.
            </p>
          </fieldset>
          <div class="modal-action">
            <form method="dialog">
              <button class="btn" onClick={saveContent}>
                Save
              </button>
            </form>
          </div>
        </div>
        <form method="dialog" class="modal-backdrop">
          <button>Submit</button>
        </form>
      </dialog>
    </>
  );
};
