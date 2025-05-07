import { createResource, Resource, For } from "solid-js";
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
  const [streams] = createResource(fetchStreams);
  return (
    <>
      <For each={item(streams)}>
        {(stream) => (
          <div class="card bg-base-100 w-150 shadow-sm p-4 pb-2">
            <div class="card-body">
              <h2 class="card-title"><A href={"/stream/" + stream.name}>{stream.name}</A></h2>
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
    </>
  );
};
