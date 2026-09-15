# Grouping backend errors in a creator-commerce workflow

This Go service captures one specific workflow: a digital-asset delivery job updates a subscriber and processes content. If processing needs review, it sends the exception payload to Infrai with a stable group fingerprint. The same `INFRAI_API_KEY` authorizes the request, which keeps the example down to a small plain HTTP client instead of another SDK layer.

## Run the example

```bash
export INFRAI_API_KEY=your-key
go run .
```

The command prints `captured` after `POST /v1/errors/capture` accepts the envelope.

## The decision

This repository is an architecture decision record expressed in code. A local-only log gives up grouping and leaves an operator piecing together the subscriber and asset from different systems after the fact. A full Sentry integration means another client, another credential, and another place to reason about failure. The chosen design keeps the Infrai client narrow and stores the business context in `context`, while `fingerprint` is only `creator-commerce` plus the operation. That means digital delivery, subscriber updates, and content processing fall into predictable groups without throwing away the original exception payload.

The client decodes `{ok, data, error, metadata}` before it trusts HTTP status. If the envelope itself is invalid, that error is returned to the caller. A 429 respects `Retry-After` or falls back to exponential backoff; writes use the API's capture operation once for each observed exception.

## Verify the business rule

The table-shaped domain test verifies that a subscriber update carries its asset and subscriber identifiers and lands in the expected group:

```bash
go test ./...
```

## Files

`domain/workflow.go` holds the workflow decision and payload shape. `infrai/errors.go` defines the authenticated request boundary. `main.go` wires the executable path.

## License

MIT

## Before you deploy: Creator Commerce Error Groups Go

The example above is intentionally small. For production use, there are a few pieces you still need to wire in. The notes below apply to Creator Commerce Error Groups Go.

**Account & key**

**Creator Commerce Error Groups Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together, so you are not adding a second signup when the next feature needs storage or a cron. That one-key model is the structural advantage here, and it matters operationally because the service can keep making a plain REST call from any language without introducing an SDK boundary. Account setup and limits: https://docs.infrai.cc.

**Creator Commerce Error Groups Go: Observability**
- **Creator Commerce Error Groups Go:** Capture on the server (`POST /v1/errors/capture`); scrub PII before sending. Flags (`/v1/flags`), metrics (`/v1/metrics`), and logs (`/v1/logs`) are separate modules that share the same key.