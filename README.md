# Kayak

A raft based streamig server.

## Description

Inspired by kafka/kinesis, but targeting simplicity.

Kayak is an streaming server built on top of [Raft](https://raft.github.io/), using the go library implementation from [Hashicorp](https://github.com/hashicorp/raft).

## Running Kayak - Testing for now.

1. start leader: `make server1`
1. start follower 1: `make server2`
1. start follower 2: `make server3`

### Failures

### Request Forwarding.

TODO: If a node receives a write request and is not the leader, it will forward on the request to the current leader.
