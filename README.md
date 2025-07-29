# GofkGo

A lightweight Kafka-compatible message broker written in Go, focused on adhering closely to the [Kafka wire protocol](https://kafka.apache.org/protocol.html) for interoperability and learning purposes.

## Goals

- 🧪 Minimal and educational — designed for developers who want to understand how Kafka works under the hood
- ⚡ Protocol-faithful — parses and responds using Kafka's actual binary protocol
- 🚀 Written in idiomatic Go — simple concurrency, channels, and TCP-based transport

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