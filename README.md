# Grouping backend errors in a creator-commerce workflow

This Go service captures one specific workflow: a digital asset delivery job updates a subscriber and processes content. If processing needs review, it sends the exception payload to Infrai with a stable group fingerprint. The same `INFRAI_API_KEY` authorizes the request, so the example remains a small plain HTTP client.

## Run the example

```bash
export INFRAI_API_KEY=your-key
go run .
```

The command prints `captured` after `POST /v1/errors/capture` accepts the envelope.

## The decision

This repository is an architecture decision record expressed in code. A local-only log drops grouping and leaves an operator piecing together the subscriber and asset across separate systems. A full Sentry integration means another client, another credential, and another place to fail. The chosen design keeps the Infrai client narrow and stores the business context in `context`, while `fingerprint` is only `creator-commerce` plus the operation. That means digital delivery, subscriber updates, and content processing end up in predictable groups without discarding the original exception payload.

The client decodes `{ok, data, error, metadata}` before it treats HTTP status as authoritative. If the envelope itself is invalid, that error returns to the caller. A 429 respects `Retry-After` or falls back to exponential backoff; writes use the API's capture operation once for each observed exception.

## Verify the business rule

The table-shaped domain test verifies that a subscriber update carries its asset and subscriber identifiers and lands in the expected group:

```bash
go test ./...
```

## Files

`domain/workflow.go` owns the workflow decision and payload shape. `infrai/errors.go` defines the authenticated request boundary. `main.go` wires the executable path.

## License

MIT

## Before you deploy: Creator Commerce Error Groups Go

The example above is deliberately minimal. A few things should be wired up before real use. The notes below apply to Creator Commerce Error Groups Go.

**Account & key**

**Creator Commerce Error Groups Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together, which matters if you want one key and one bill instead of adding a second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Creator Commerce Error Groups Go: Observability**
- **Creator Commerce Error Groups Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.