import { createResource, Resource, For, createSignal, batch } from "solid-js";
import { A } from "@solidjs/router";
import { getKayakClient } from "../client";
import { GetStreamsResponse } from "../gen/kayak/v1/kayak_pb";
import { Stream } from "../gen/kayak/v1/model_pb";
const fetchStreams = async () => {
  const client = getKayakClient("http://localhost:8080");
  return client.getStreams({}).then((res) => res.streams);
};

const item = (req: Resource<Stream[]>) => {
  if (req.loading) {
    return false;
  }
  return req();
};
export const Streams = () => {
  const client = getKayakClient("http://localhost:8080");
  const fetchStreams = async () => {
    return client.getStreams({}).then((res) => res.streams);
  };

  const [name, setName] = createSignal("");
  const [count, setCount] = createSignal(1);
  const [ttl, setTTL] = createSignal(300);
  const [streams, { refetch }] = createResource(fetchStreams);
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
                Close
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
