# GofkGo

A lightweight Kafka-compatible message broker written in Go, focused on adhering closely to the [Kafka wire protocol](https://kafka.apache.org/protocol.html) for interoperability and learning purposes.

## Goals

- 🧪 Minimal and educational — designed for developers who want to understand how Kafka works under the hood
- ⚡ Protocol-faithful — parses and responds using Kafka's actual binary protocol
- 🚀 Written in idiomatic Go — simple concurrency, channels, and TCP-based transport

## Project Structure
```bash
gofkgo/
├── README.md
├── go.mod
├── cmd/
│   └── gofkgo/          # main entrypoint (e.g., broker executable)
│       └── main.go
├── pkg/
│   ├── protocol/        # Kafka protocol encoding/decoding
│   ├── broker/          # core broker logic (topics, partitions, etc.)
│   └── transport/       # TCP connection handling, framing, etc.
├── internal/
│   └── testutil/        # internal helpers for testing
└── docs/
    └── protocol.md      # reference to Kafka wire protocol details
```

## Getting Started

```bash
git clone https://github.com/Jack4Code/gofkgo.git
cd gofkgo
go run ./cmd/gofkgo
```

## Status
Currently in development — expect bugs, wire format weirdness, and imcomplete feature support. PRs welcome.

## License
MIT