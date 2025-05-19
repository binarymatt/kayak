import { useParams } from "@solidjs/router";
import { getKayakClient } from "../client";
import { createResource, Show } from "solid-js";
import Loading from "../components/Loading";

const fetchStream = async (name: string) => {
  const client = getKayakClient("http://localhost:8080");
  return await client.getStream({ name }).then((res) => res.stream);
};
export const Stream = () => {
  const params = useParams();
  const name = params.name;
  const [stream] = createResource(name, fetchStream);
  return (
    <Show when={!stream.loading} fallback={<Loading />}>
      {stream().name} stream detail coming soon
    </Show>
  );
};
