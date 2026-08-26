# kirolink Fix Plan — endpoint, cache, parser, logging

Scope: API_KEY auth only (no enterprise / external IdP / images). Text-to-text with
upstream-cacheable requests. Auth hardening deliberately OUT of scope (proxy runs
on localhost only; revisit if it ever binds a public interface).

Reference repo for patterns: `/home/ubuntu/9router` (open-sse/kiro* + executors/kiro.js).
All line numbers below are against current `kirolink.go` / `protocol/sse_parser.go`
(verified 2026-08-26).

---

## P0 — Endpoint routing for API_KEY

### Why
Kiro's legacy CodeWhisperer surface authenticates API keys but rejects the same
valid payload with `REQUEST_BODY_INVALID` ("Improperly formed request"). The
Amazon Q surface (`q.*`) is the correct host for API_KEY traffic. 9router routes
API_KEY to `q.us-east-1.amazonaws.com` first and documents the 400 behavior at
`open-sse/executors/kiro.js:265-272`.

kirolink hardcodes `codewhisperer.*` (`codeWhispererURL`, kirolink.go:1180) and
built an entire trim-and-retry loop (kirolink.go:1861-1878) around the resulting
"Improperly formed request" 400. Hypothesis: that retry loop is a symptom of the
wrong host, not oversized payloads. A real account with a working API key is
required to confirm; treat the q.* switch as the fix and the retry-loop removal
as conditional on it.

### Changes

1. **`codeWhispererURL(creds)`** (kirolink.go:1180):
   - `creds.AuthMethod == authMethodAPIKey` → `https://q.%s.amazonaws.com/generateAssistantResponse` (same path, different host).
   - Else keep `https://codewhisperer.%s.amazonaws.com/...`.
   - Region from `creds.Region` via existing `getKiroRegion`.

2. **`setCodeWhispererHeaders`** (kirolink.go:1184): drop `X-Amz-Target` when the
   host is NOT `codewhisperer`. 9router does exactly this
   (`executors/kiro.js:235-239`: set only when `url.includes("://codewhisperer.")`, delete otherwise).
   Signature will need the URL or a flag; simplest is to compute the same
   condition in both funcs (shared helper `isCodeWhispererHost(url string) bool`).

3. **Conditional removal of the 400 trim-retry** (kirolink.go:1861-1878): if q.*
   eliminates the 400 in live testing, delete the `keepMostRecentHistory` +
   tool-strip retry branch entirely and treat 400 as terminal (9router makes 400
   terminal on purpose — retrying an identical malformed body cannot repair it).

### Verify
- Run with a real API key against q.*; confirm no "Improperly formed request" 400.
- `go test ./...` still green.
- Keep the 403-token-sync branch (kirolink.go:1881) — it is correctly skipped
  for API_KEY already; no change needed there.

---

## P0 — Caching: make requests byte-stable + measurable

### Why
No explicit cache field exists for Kiro — caching is implicit via a byte-identical
request prefix (the `history` array + first user turn). kirolink currently
destroys cacheability two ways:

1. System prompt is wrapped into the **current message** every turn
   (`buildCurrentMessageContent`, kirolink.go:639), so `history[0]` on turn ≥2 is
   NOT byte-identical to what was sent as turn 1. Prefix changes every turn.
2. `ensurePayloadFits` trims history from the **front** (kirolink.go:570) and the
   streaming retry keeps only the last 2 turns (kirolink.go:1864) — both destroy
   the anchored prefix.
3. Even when caching does happen upstream, kirolink cannot see it: the parser
   never handles `metricsEvent`, which is where Kiro reports
   `cacheReadInputTokens` / `cacheCreationInputTokens`.

9router's `kiroSessionReplay.js` solves this with an in-memory Map that freezes
msg0. kirolink does NOT need that state: it never injects volatile content (no
timestamp), so a **positional, deterministic transform** is sufficient.

### Changes

1. **Add top-level `SystemPrompt string \`json:"systemPrompt,omitempty"\``** to
   `CodeWhispererRequest` (kirolink.go:253) and populate it in
   `buildCodeWhispererRequestWithCredentials` (kirolink.go:812) from
   `buildSystemContext(anthropicReq.System)`. 9router sends the system prompt
   both top-level AND as a content prefix because the CodeWhisperer surface does
   not always honor the top-level field (`claude-to-kiro.js:245-247`).

2. **Positional, deterministic system-prompt placement.** Replace
   `buildCurrentMessageContent` (kirolink.go:639) so the `<context>SYS</context>`
   wrapper is applied to the **first user turn regardless of position**:
   - Turn 1 (no history): system context goes in `currentMessage` content.
   - Turn ≥2: system context goes in `history[0]` (the first user turn) content;
     `currentMessage` carries only `<task>…</task>`.
   - All user turns use the same wrapper so the boundary is stable. No
     cross-request state, no timestamp → byte-identical prefix across turns.
   - Simplest correct implementation: strip system context out of
     `buildCurrentMessageContent`, and inject it into `HistoryUserMessage` at
     index 0 only when it is a user message. Skip when the first message is
     assistant/tool-shaped (fall back to top-level `systemPrompt` only).

3. **Pin `history[0]`.** In `ensurePayloadFits` (kirolink.go:546): do not trim the
   first history entry. Trim from index 1 onward (new `trimHistoryKeepFirst`),
   then tool-description truncation, then schema strip, then drop tools — the
   existing phase ladder but never touching `history[0]`. Update
   `trimOldestHistoryMessage` (kirolink.go:619) or add a sibling.

4. **Parse `metricsEvent` in the parser** (`protocol/sse_parser.go`, current
   metadata branch at lines ~88-105):
   - Add a `metricsEvent` struct: `inputTokens`, `outputTokens`,
     `cacheReadInputTokens` (also accept `cache_read_input_tokens`),
     `cacheCreationInputTokens` (also `cache_creation_input_tokens`) — mirror
     9router `executors/kiro.js:886-901`, which reads both spellings.
   - Emit a new `SSEEvent{Event: "metrics", Data: {...}}` so the caller can read
     usage.
   - In `assembleAnthropicResponse` / the handlers, replace the fake usage
     (`input_tokens` = byte length of content, kirolink.go:1995; hardcoded
     `output_tokens: 1`, kirolink.go:1673) with real values from `metricsEvent`.
     This makes cache hits visible: nonzero `cache_read_input_tokens` on turn ≥2
     = the strategy is working.

### Verify
- Turn 2+ request bodies: `history[0]` byte-identical to turn 1's first user
  message (unit-testable: build two requests from a two-turn Anthropic history
  and diff the `history[0]` content).
- After a live run, response `usage` includes `cache_read_input_tokens` > 0.
- `go test ./...` green; update `tool_mapping_test.go` / `response_translation_test.go`
  if payload shape changed.

---

## P0 — Parser hardening (panic / OOM on malformed frames)

### Why
`protocol/sse_parser.go:52-64` trusts frame sizes from the upstream stream:

- `headerLen` is an unbounded uint32 from the wire → `make([]byte, headerLen)`
  (line 58) can allocate up to ~4GB.
- If `headerLen > totalLen-12` while the stream still has bytes, `payloadLen`
  goes negative → `make([]byte, payloadLen)` (line 64) **panics**.
- `ParseEvents` runs on the streaming path (kirolink.go:1907) with no `recover`;
  the existing recover only wraps the non-streaming handler (kirolink.go:1396).

Upstream is AWS (trusted-ish), but any proxy/MITM in front (the repo historically
used a hardcoded local proxy) can inject a bad frame. Cheap to fix.

### Changes

1. Add a `maxHeaderBytes` const (e.g. 64 * 1024) matching 9router's
   `EVENTSTREAM_MAX_HEADERS_BYTES` (executors/kiro.js:934, 1207-1208).
2. Validate BEFORE allocating, right after reading `totalLen`/`headerLen`:
   - `headerLen > maxHeaderBytes` → log + break.
   - `headerLen > totalLen-12` → log + break (also prevents the negative payload).
3. Wrap `ParseEvents` call in the streaming handler (kirolink.go:1907) with a
   `defer recover` that logs and emits an SSE error event instead of dying.

### Verify
- Unit test: feed a hand-built frame where `headerLen > totalLen` (with enough
  trailing bytes to pass `io.ReadFull`) → parser returns cleanly, no panic.
- Feed `headerLen = 0xFFFFFFFF` → no 4GB allocation.

---

## P1 — Logging: stop leaking conversation + tokens

### Why
- kirolink.go:1810 (streaming request body), 1942 (non-stream request body),
  1978 (response body) dump FULL request/response bodies to stdout on every
  call. With docker `restart: unless-stopped`, conversation content accumulates
  in logs indefinitely. In API_KEY mode the request body also contains no key
  (it's a header), but conversation content is still private data.
- `kirolink read` prints the full access + refresh token (kirolink.go:961-962).

### Changes

1. Remove the three full-body `fmt.Printf`s, or gate behind an env var
   (e.g. `KIROLINK_DEBUG=1`). Default OFF.
2. Keep a one-line summary instead: method, status code, latency, payload byte
   size, `conversationId` prefix.
3. `kirolink read`: print token prefix only (`token[:12] + "…"`), or a
   `--full` flag for debugging.

### Verify
- No conversation text in stdout with default env.
- `go test ./...` green.

---

## Housekeeping (do alongside, low risk)

- **Delete dead constant** `profileArnIAM` (kirolink.go:300). It holds a
  profile ARN that belongs to 9router's shared social account and is never read;
  if someone wires it into the API_KEY path later it produces
  `403 bearer token invalid` (9router warns about this explicitly).
- `maxToolDescLen` (kirolink.go:309, = 200) is far below what upstream accepts
  (9router uses `KIRO_TOOL_DESCRIPTION_MAX_LENGTH = 10237`). Raise it after the
  endpoint fix is verified, since aggressive truncation also changes the cacheable
  prefix on tool-heavy conversations.
- Consider a single shared `*http.Client` (kirolink.go:1827, 1960 construct one
  per request) so connections pool; not correctness-critical.

---

## Execution order

1. Parser hardening (isolated, testable, zero behavior risk).
2. Endpoint routing (q.* for API_KEY + conditional `X-Amz-Target`). Live-test
   with a real key; if the 400 disappears, delete the trim-retry branch.
3. Caching (top-level `systemPrompt` + positional system context + pinned
   `history[0]` + `metricsEvent` parsing + real usage numbers). Unit-test byte
   stability; live-test `cache_read_input_tokens`.
4. Logging changes.

Build/run: `go build -o kirolink kirolink.go && go test ./...`; server:
`./kirolink server <port>`.
