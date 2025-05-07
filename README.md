# LiteLLM-Go

A lightweight LLM gateway written in Go, inspired by [LiteLLM](https://github.com/BerriAI/litellm).

Features:

- OpenAI-compatible REST API (ChatCompletion) so existing OpenAI SDKs can point to this gateway.
- Route requests to multiple back-ends: OpenAI, Azure OpenAI, AWS Bedrock (more coming).
- Simple round-robin load balancing between multiple deployments of the same provider.
- Optional per-key rate limiting (token bucket, in-memory for now).
- YAML configuration file similar to `litellm_config.yaml`.

> ⚠️ Enterprise-only features of LiteLLM are **out-of-scope** for this project.

## Quick start

```bash
# install deps
go mod tidy

# run server with sample config
go run ./cmd/server --config examples/config.yaml
```

See `examples/config.yaml` for available options.

## TODO

- [x] Basic router
- [ ] Metrics / Prometheus
- [ ] Pluggable auth & logging callbacks 
