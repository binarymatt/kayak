import { useParams } from "@solidjs/router";
import { getKayakClient } from "../client";
import { batch, createResource, For, Show } from "solid-js";
import Loading from "../components/Loading";
import { node } from "../stores";

const fetchStream = async (name: string) => {
  const client = getKayakClient(node.address);
  return await client
    .getStream({ name: name, includeStats: true })
    .then((res) => res.stream);
};
const clearRecords = async (name: string) => {
  const client = getKayakClient(node.address);
  return await client.deleteRecords({ name });
};
export const Stream = () => {
  const params = useParams();
  const name = params.name;
  const [stream, { refetch }] = createResource(name, fetchStream);
  const deleteStream = async (name: string) => {
    batch(async () => {
      await clearRecords(name);
      refetch();
    });
  };
  return (
    <Show when={!stream.loading} fallback={<Loading />}>
      <div class="breadcrumbs text-sm">
        <ul>
          <li>
            <a href="/streams">Streams</a>
          </li>
          <li>{stream().name}</li>
        </ul>
      </div>
      <div class="prose prose-sm">
        <h2>{stream().name} </h2>
        <h3>Stream Info</h3>
        <div class="stats shadow">
          <div class="stat">
            <div class="stat-title">Partitions</div>
            <div class="stat-value">{"" + stream().partitionCount}</div>
          </div>
          <div class="stat">
            <div class="stat-title">Record TTL</div>
            <div class="stat-value">{"" + stream().ttl}</div>
          </div>
          <div class="stat">
            <div class="stat-title">Visible Record Count</div>
            <div class="stat-value">{"" + stream().stats.recordCount}</div>
          </div>
        </div>
        <h3>Partition Information</h3>
        <For each={Object.entries(stream().stats.partitionCounts)}>
          {([key, value]) => (
            <p>
              {"" + key} - {"" + value}
            </p>
          )}
        </For>
        <h3>Group Information</h3>
        <button onClick={() => deleteStream(stream().name)} class="btn">
          Clear Records
        </button>
      </div>
    </Show>
  );
};
