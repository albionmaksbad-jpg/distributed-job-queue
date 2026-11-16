# distributed-job-queue

A small background job system in Go. The first version uses PostgreSQL as the durable queue and runs work in independent worker processes.

The design assumes a job can run more than once. Handlers must be idempotent, retries use backoff, and a worker crash must not leave a job stuck forever.

## Current state

Day one: repository bootstrap and a minimal API health endpoint. Job submission, persistence, and workers are not implemented yet.

## Run the API

```sh
go run ./cmd/api
```

The server listens on `:8080` by default. Set `api_addr` to change the address. `GET /healthz` returns `204`.
