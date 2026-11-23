# GofkGo

> **Work in progress.** Not usable as a broker yet. No Kafka client can talk to it.

A lightweight Kafka-compatible message broker written in Go, focused on adhering closely to the [Kafka wire protocol](https://kafka.apache.org/protocol.html) for interoperability and learning purposes.

## Goals

- Minimal and educational -- designed for developers who want to understand how Kafka works under the hood
- Protocol-faithful -- parses and responds using Kafka's actual binary protocol
- Written in Go -- simple concurrency, channels, and TCP-based transport

## Project Structure
```bash
gofkgo/
├── cmd/gofkgo/          # entrypoint
├── internal/
│   ├── app/             # wires dependencies together
│   ├── models/          # request types
│   ├── ports/           # interfaces
│   ├── protocol/        # framing and header parsing
│   └── server/          # TCP listener and connection handling
└── compose.yaml         # reference Kafka for comparing responses
```

## Getting Started

```bash
git clone https://github.com/Jack4Code/gofkgo.git
cd gofkgo
go run ./cmd/gofkgo
```

## Reference Kafka

`compose.yaml` runs a real single-node Kafka (KRaft) on `localhost:9092` to compare responses against. gofkgo stays on `:9090`.

```bash
docker compose up -d --wait
docker compose down
```

Send a raw ApiVersions v0 request and dump the response:

```bash
printf '\x00\x00\x00\x0e\x00\x12\x00\x00\x00\x00\x00\x01\x00\x04test' | timeout 3 nc localhost 9092 | xxd
```

Point the same command at `9090` to compare with gofkgo.

## Status
Early development. Right now it accepts TCP connections and parses request headers, but doesn't send real Kafka responses yet. Next up is ApiVersions, then Metadata, Produce and Fetch. PRs welcome.

## License
MIT
