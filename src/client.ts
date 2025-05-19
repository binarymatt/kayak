import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import { KayakService } from "./gen/kayak/v1/kayak_pb";
import { AdminService } from "./gen/kayak/v1/admin_pb";

export function getKayakClient(baseUrl: string) {
  const transport = createConnectTransport({
    baseUrl: baseUrl,
  });
  return createClient(KayakService, transport);
}

export function getClients(baseUrl: string) {
  const transport = createConnectTransport({
    baseUrl: baseUrl,
  });
  return {
    kayak: createClient(KayakService, transport),
    admin: createClient(AdminService, transport),
  };
}
