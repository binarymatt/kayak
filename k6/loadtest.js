import { Client, StatusOK } from "k6/net/grpc";
import { SharedArray } from "k6/data";
import { check } from "k6";
import exec from "k6/execution";
import http from "k6/http";
import { b64encode } from "k6/encoding";

export const options = {
  // A number specifying the number of VUs (virtual users) to run concurrently.
  // vus: 20,
  iterations: 1000,
  // A string specifying the total duration of the test run.
  // duration: "5m",
};
export default function() {
  const raw = { iteration: exec.vu.iterationInInstance, vu: exec.vu.idInTest };
  const bstring = b64encode(JSON.stringify(raw));
  const payload = JSON.stringify({
    stream_name: "test",
    records: [
      {
        payload: bstring,
      },
    ],
  });
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };
  const res = http.post(
    "http://localhost:8080/kayak.v1.KayakService/PutRecords",
    payload,
    params,
  );
  console.log(res);
  check(res, {
    "is status 200": (r) => r.status === 200,
  });
}
