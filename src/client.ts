import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import { KayakService } from "./gen/kayak/v1/kayak_pb";

export function getKayakClient(baseUrl: string) {
  const transport = createConnectTransport({
    baseUrl: baseUrl,
  });
  return createClient(KayakService, transport);
}
