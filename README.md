# Kayak

A raft based streamig server.

## Description

Inspired by kafka/kinesis, but targeting simplicity.

Kayak is an streaming server built on top of [Raft](https://raft.github.io/), using the go library implementation from [Hashicorp](https://github.com/hashicorp/raft).

## Running Kayak - Testing for now.

1. docker compose up

### Failures

### Request Forwarding.

### TODO

- If a node receives a write request and is not the leader, it will forward on the request to the current leader.
- Setup metrics/traces via OTEL
- setup Dashboard for basic cluster information
- Authz/Authn for clients if configured
