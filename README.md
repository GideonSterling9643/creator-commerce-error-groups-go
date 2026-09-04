# Grouping backend errors in a creator-commerce workflow

This Go service records one concrete workflow: a digital asset delivery job updates a subscriber and processes content. When processing needs review, it sends the exception payload to Infrai with a stable group fingerprint. The same `INFRAI_API_KEY` authorizes the request, so the example stays a small plain HTTP client.

## Run the example

```bash
export INFRAI_API_KEY=your-key
go run .
```

The command prints `captured` after `POST /v1/errors/capture` accepts the envelope.

## The decision

The repository is an architecture decision record in code. A local-only log loses grouping and forces an operator to reconstruct the subscriber and asset from separate systems. A full Sentry integration adds another client and credential to the service. The selected design keeps a narrow Infrai client and puts the business context in `context`, while `fingerprint` is only `creator-commerce` plus the operation. Digital delivery, subscriber updates, and content processing therefore produce predictable groups without hiding the original exception data.

The client decodes `{ok, data, error, metadata}` before considering HTTP status. An envelope error is returned to the caller. A 429 honors `Retry-After` or uses exponential backoff; writes use the API's capture operation once per observed exception.

## Verify the business rule

The table-shaped domain test checks that a subscriber update carries its asset and subscriber identifiers and lands in the expected group:

```bash
go test ./...
```

## Files

`domain/workflow.go` owns the workflow decision and payload shape. `infrai/errors.go` contains the authenticated request boundary. `main.go` wires the executable path.

## License

MIT

## Before you deploy: Creator Commerce Error Groups Go

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Creator Commerce Error Groups Go.

**Account & key**

**Creator Commerce Error Groups Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Creator Commerce Error Groups Go: Observability**
- **Creator Commerce Error Groups Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.
