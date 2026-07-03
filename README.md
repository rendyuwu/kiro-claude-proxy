<p align="left">
<pre style="line-height: 0.9;">
 ████  ███  ██                   ████      ████          ████    
▒▒██  ██▒  ▒▒                   ▒▒██      ▒▒██          ▒▒██     
 ▒██ ██   ███  ██████  █████     ▒██       ▒██  ██████   ▒██ ████
 ▒█████  ▒▒██ ▒▒██▒▒█ ██▒▒██     ▒██       ▒██ ▒▒██▒▒█   ▒██▒▒██ 
 ▒██▒▒██  ▒██  ▒██ ▒ ▒██ ▒██     ▒██       ▒██  ▒██ ▒▒   ▒█████  
 ▒██ ▒▒██ ▒██  ▒██   ▒██ ▒██     ▒██     █ ▒██  ▒██ ██   ▒██▒▒██ 
 ███  ▒▒█████ ████   ▒▒█████     ███████▒█ ████ ███▒███  ███ ████
▒▒▒    ▒▒▒▒▒ ▒▒▒▒     ▒▒▒▒▒     ▒▒▒▒▒▒▒ ▒ ▒▒▒▒ ▒▒▒ ▒▒▒  ▒▒▒ ▒▒▒▒                            
</pre>
</p>

# KiroLink ![Stars](https://img.shields.io/github/stars/alexandephilia/kiro-claude-proxy)

[![Claude Code](https://img.shields.io/badge/Claude_Code-Override-964B00?logo=anthropic&logoColor=white)](https://claude.ai/code)
[![Kiro](https://img.shields.io/badge/Kiro-Authentication-7B2D8B?logo=data:image/svg%2Bxml;base64,PHN2ZyB3aWR0aD0iMTIwMCIgaGVpZ2h0PSIxMjAwIiB2aWV3Qm94PSIwIDAgMTIwMCAxMjAwIiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8cmVjdCB3aWR0aD0iMTIwMCIgaGVpZ2h0PSIxMjAwIiByeD0iMjYwIiBmaWxsPSIjOTA0NkZGIi8+CjxtYXNrIGlkPSJtYXNrMF8xMTA2XzQ4NTYiIHN0eWxlPSJtYXNrLXR5cGU6bHVtaW5hbmNlIiBtYXNrVW5pdHM9InVzZXJTcGFjZU9uVXNlIiB4PSIyNzIiIHk9IjIwMiIgd2lkdGg9IjY1NSIgaGVpZ2h0PSI3OTYiPgo8cGF0aCBkPSJNOTI2LjU3OCAyMDIuNzkzSDI3Mi42MzdWOTk3Ljg1N0g5MjYuNTc4VjIwMi43OTNaIiBmaWxsPSJ3aGl0ZSIvPgo8L21hc2s+CjxnIG1hc2s9InVybCgjbWFzazBfMTEwNl80ODU2KSI+CjxwYXRoIGQ9Ik0zOTguNTU0IDgxOC45MTRDMzE2LjMxNSAxMDAxLjAzIDQ5MS40NzcgMTA0Ni43NCA2MjAuNjcyIDk0MC4xNTZDNjU4LjY4NyAxMDU5LjY2IDgwMS4wNTIgOTcwLjQ3MyA4NTIuMjM0IDg3Ny43OTVDOTY0Ljc4NyA2NzMuNTY3IDkxOS4zMTggNDY1LjM1NyA5MDcuNjQgNDIyLjM3NEM4MjcuNjM3IDEyOS40NDMgNDI3LjYyMyAxMjguOTQ2IDM1OC44IDQyMy44NjVDMzQyLjY1MSA0NzUuNTQ0IDM0Mi40MDIgNTM0LjE4IDMzMy40NTggNTk1LjA1MUMzMjguOTg2IDYyNS44NiAzMjUuNTA3IDY0NS40ODggMzEzLjgzIDY3Ny43ODVDMzA2Ljg3MyA2OTYuNDI0IDI5Ny42OCA3MTIuODE5IDI4Mi43NzMgNzQwLjY0NUMyNTkuOTE1IDc4My44ODEgMjY5LjYwNCA4NjcuMTEzIDM4Ny44NyA4MjMuODgzTDM5OS4wNTEgODE4LjkxNEgzOTguNTU0WiIgZmlsbD0id2hpdGUiLz4KPHBhdGggZD0iTTYzNi4xMjMgNTQ5LjM1M0M2MDMuMzI4IDU0OS4zNTMgNTk4LjM1OSA1MTAuMDk3IDU5OC4zNTkgNDg2Ljc0MkM1OTguMzU5IDQ2NS42MjMgNjAyLjA4NiA0NDguOTc3IDYwOS4yOTMgNDM4LjI5M0M2MTUuNTA0IDQyOC44NTIgNjI0LjY5NyA0MjQuMTMxIDYzNi4xMjMgNDI0LjEzMUM2NDcuNTU1IDQyNC4xMzEgNjU3LjQ5MiA0MjguODUyIDY2NC40NDcgNDM4LjU0MUM2NzIuMzk4IDQ0OS40NzQgNjc2LjYyMyA0NjYuMTIgNjc2LjYyMyA0ODYuNzQyQzY3Ni42MjMgNTI1Ljk5OCA2NjEuNDcxIDU0OS4zNTMgNjM2LjM3NSA1NDkuMzUzSDYzNi4xMjNaIiBmaWxsPSJibGFjayIvPgo8cGF0aCBkPSJNNzcxLjI0IDU0OS4zNTNDNzM4LjQ0NSA1NDkuMzUzIDczMy40NzcgNTEwLjA5NyA3MzMuNDc3IDQ4Ni43NDJDNzMzLjQ3NyA0NjUuNjIzIDczNy4yMDMgNDQ4Ljk3NyA3NDQuNDEgNDM4LjI5M0M3NTAuNjIxIDQyOC44NTIgNzU5LjgxNCA0MjQuMTMxIDc3MS4yNCA0MjQuMTMxQzc4Mi42NzIgNDI0LjEzMSA3OTIuNjA5IDQyOC44NTIgNzk5LjU2NCA0MzguNTQxQzgwNy41MTYgNDQ5LjQ3NCA4MTEuNzQgNDY2LjEyIDgxMS43NCA0ODYuNzQyQzgxMS43NCA1MjUuOTk4IDc5Ni41ODggNTQ5LjM1MyA3NzEuNDkyIDU0OS4zNTNINzcxLjI0WiIgZmlsbD0iYmxhY2siLz4KPC9nPgo8L3N2Zz4K&logoColor=white)](https://kiro.dev/)
![OpenClaw](https://img.shields.io/badge/OpenClaw-Compatible-FF4B4B?logo=data:image/svg%2Bxml;base64,PHN2ZyB2aWV3Qm94PSIwIDAgMTIwIDEyMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8cGF0aCBkPSJNNjAgMTAgQzMwIDEwIDE1IDM1IDE1IDU1IEMxNSA3NSAzMCA5NSA0NSAxMDAgTDQ1IDExMCBMNTUgMTEwIEw1NSAxMDAgQzU1IDEwMCA2MCAxMDIgNjUgMTAwIEw2NSAxMTAgTDc1IDExMCBMNzUgMTAwIEM5MCA5NSAxMDUgNzUgMTA1IDU1IEMxMDUgMzUgOTAgMTAgNjAgMTBaIiBmaWxsPSJ3aGl0ZSI+PC9wYXRoPgogIDxwYXRoIGQ9Ik0yMCA0NSBDNSA0MCAwIDUwIDUgNjAgQzEwIDcwIDIwIDY1IDI1IDU1IEMyOCA0OCAyNSA0NSAyMCA0NVoiIGZpbGw9IndoaXRlIj48L3BhdGg+CiAgPHBhdGggZD0iTTEwMCA0NSBDMTE1IDQwIDEyMCA1MCAxMTUgNjAgQzExMCA3MCAxMDAgNjUgOTUgNTUgQzkyIDQ4IDk1IDQ1IDEwMCA0NVoiIGZpbGw9IndoaXRlIj48L3BhdGg+CiAgPHBhdGggZD0iTTQ1IDE1IFEzNSA1IDMwIDgiIHN0cm9rZT0id2hpdGUiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIj48L3BhdGg+CiAgPHBhdGggZD0iTTc1IDE1IFE4NSA1IDkwIDgiIHN0cm9rZT0id2hpdGUiIHN0cm9rZS13aWR0aD0iMiIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIj48L3BhdGg+CiAgPGNpcmNsZSBjeD0iNDUiIGN5PSIzNSIgcj0iNiIgZmlsbD0iI0ZGNEI0QiI+PC9jaXJjbGU+CiAgPGNpcmNsZSBjeD0iNzUiIGN5PSIzNSIgcj0iNiIgZmlsbD0iI0ZGNEI0QiI+PC9jaXJjbGU+Cjwvc3ZnPgo=&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.23.3-00ADD8?logo=go&logoColor=white)
</br>
![API](https://img.shields.io/badge/API-Compatible-964B00?logo=anthropic&logoColor=white)
![AWS CodeWhisperer](https://img.shields.io/badge/AWS-CodeWhisperer-232F3E?logo=amazon-aws&logoColor=white)

`kirolink` is a small Go CLI that exposes an Anthropic-shaped local API for Kiro. It can use a Kiro/CodeWhisperer API key from `.env` or environment variables, with the old local SSO cache path kept as a fallback.

In practice: your client sends `POST /v1/messages` to this proxy, the proxy translates the request to Kiro's CodeWhisperer runtime, sends it with Kiro API-key compatible headers, and translates the response back on the way out.

## What this thing actually does

- Reads Kiro API-key credentials from `.env`, process env, or request headers
- Falls back to `~/.aws/sso/cache/kiro-auth-token.json` when no API key is configured
- Prints shell-ready `ANTHROPIC_*` environment variable setup
- Starts a local server on port `8080` by default
- Exposes these endpoints:
  - `POST /v1/messages`
  - `GET /v1/models`
  - `GET /health`
- Includes a `claude` helper command that edits `~/.claude.json`

<details>
<summary><b>🔍 Troubleshoot: Token Not Found</b></summary>

If you get a "file not found" error, ensure you've run the Kiro login first.
Default path: `~/.aws/sso/cache/kiro-auth-token.json`

</details>

## Quick start

### 1. Build it

```bash
go build -o kirolink kirolink.go
```

### 2. Configure a Kiro API key

Create `.env` in the repo root:

```env
KIRO_API_KEY=your_kiro_api_key
KIRO_REGION=us-east-1
KIRO_PROFILE_ARN=
```

`KIRO_PROFILE_ARN` is optional. If you set it, make sure it belongs to the API key's account.

You can also pass the API key per request with `Authorization: Bearer ...`, `x-api-key`, or `X-Kiro-API-Key`.

If no API key is configured, `kirolink` falls back to the old local SSO cache path:

```text
~/.aws/sso/cache/kiro-auth-token.json
```

### 3. Export the Anthropic env vars

On macOS/Linux, you can eval the output directly:

```bash
eval "$(./kirolink export)"
```

On Windows, the command prints both CMD and PowerShell variants for you to copy:

```bash
./kirolink export
```

By default this sets:

- `ANTHROPIC_BASE_URL=http://localhost:8080`
- `ANTHROPIC_API_KEY=<current access token>`

With `KIRO_API_KEY` in `.env`, the server ignores the local SSO cache for chat requests. `ANTHROPIC_API_KEY` can be any client-side placeholder unless you want to pass the Kiro key through request headers.

### 4. Start the proxy

```bash
./kirolink server
```

Custom port:

```bash
./kirolink server 9000
```

<p align="left">
  <img src="Claude-Code.jpg" alt="Claude Code using kirolink" width="600" style="border-radius: 12px;">
</p>

> **Note:** Make sure your server proxy is running before using Claude Code with kirolink.

If you use a custom port, set `ANTHROPIC_BASE_URL` manually — the `export` command always prints `http://localhost:8080`.

### 5. Point your client at it

Claude Code and other Anthropic-compatible clients can use the exported env vars. You can also hit the proxy directly:

```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"model":"claude-haiku-4-5","messages":[{"role":"user","content":"Hello"}],"max_tokens":256}'
```

## Commands

| Command                    | What it does                                                                                 |
| -------------------------- | -------------------------------------------------------------------------------------------- |
| `./kirolink read`          | Reads and prints the fallback cached token data.                                             |
| `./kirolink refresh`       | Syncs the fallback token from Kiro CLI sqlite and writes the token cache.                    |
| `./kirolink export`        | Prints environment variable commands for the current OS/shell style.                         |
| `./kirolink claude`        | Updates `~/.claude.json` and sets `hasCompletedOnboarding=true` plus `kirolink=true`.        |
| `./kirolink server [port]` | Starts the local Anthropic-compatible proxy server.                                          |

## HTTP surface

When the server is running, these routes are available:

- `POST /v1/messages` — main Anthropic-compatible message endpoint
- `GET /v1/models` — returns the currently exposed model aliases
- `GET /health` — returns `OK`

Example:

```bash
curl http://localhost:8080/v1/models
```

## Model aliases

The proxy accepts multiple Anthropic-style aliases, including:

- `default`
- `claude-sonnet-4-6`
- `claude-sonnet-4-5`
- `claude-opus-4-8`
- `claude-opus-4-7`
- `claude-opus-4-6`
- `claude-haiku-4-5`
- `claude-haiku-4-5-20251001`
- `claude-4-sonnet`
- `claude-4-opus`

In API-key mode these map to Kiro upstream IDs like `claude-opus-4.6`, `claude-sonnet-4.6`, and `claude-haiku-4.5`.

## How it works

1. Load `.env`, then resolve credentials from `KIRO_API_KEY`, request headers, or the fallback SSO cache.
2. Accept Anthropic-style requests over HTTP.
3. Translate them into the existing CodeWhisperer request payload.
4. Send them to `https://codewhisperer.<region>.amazonaws.com/generateAssistantResponse`.
5. In API-key mode, send `Authorization: Bearer <key>` plus `tokentype: API_KEY`.
6. Translate the response back into Anthropic-style JSON or SSE.

## Using with OpenClaw

OpenClaw is an open-source AI assistant framework that supports multiple LLM providers. You can configure it to use this proxy as a custom Claude provider.

### Configuration

Add a custom provider to your `openclaw.json` or environment config:

```json
{
  "providers": {
    "kiro-claude": {
      "api": "anthropic-messages",
      "baseURL": "http://localhost:8080",
      "apiKey": "<your-kiro-token>"
    }
  }
}
```

## Development

Build:

```bash
go build -o kirolink kirolink.go
```

Run tests:

```bash
go test ./...
```

Run protocol tests only:

```bash
go test ./protocol -v
```

## Rough edges you should know about

- API-key mode expects `KIRO_API_KEY` in `.env`, process env, or request headers.
- Fallback SSO mode still depends on a local Kiro token file.
- `refresh` writes back to `~/.aws/sso/cache/kiro-auth-token.json` for fallback SSO mode.
- `claude` modifies `~/.claude.json`; that's convenient, but it's still changing your config, so don't run it blindly.
- The documented export path is hardcoded to `http://localhost:8080`.
- The upstream CodeWhisperer endpoint is region-aware through `KIRO_REGION`, defaulting to `us-east-1`.

| Feature            | Supported natively | Handled by `kirolink` | Notes                    |
| :----------------- | :----------------: | :-------------------: | :----------------------- |
| Standard Messaging |         ❌         |          ✅           | Translated cleanly       |
| Streaming (SSE)    |         ❌         |          ✅           | Handled dynamically      |
| Local Auth         |         ❌         |          ✅           | API key or SSO fallback  |
| Tool Use           |         ❌         |          ✅           | Tool-use/Mcp            |

## Credits

<table>
  <tr>
    <td align="center" valign="top">
      <a href="https://github.com/alexandephilia">
        <img src="https://avatars.githubusercontent.com/u/43126944?v=4" width="100px;" height="100px;" style="border-radius: 30%;" alt=""/><br />
        <sub><b>Alexandephilia</b></sub>
      </a><br />
      <sub>Vibing</sub>
    </td>
    <td align="center" valign="top">
      <img src="https://matthiasroder.com/content/images/2026/01/Claude.png" width="100px;" height="100px;" style="border-radius: 30%;" alt=""/><br />
      <sub><b>Claude</b></sub><br />
      <sub>Implementation</sub>
    </td>
    <td align="center" valign="top">
      <img src="https://i.logos-download.com/114346/31977-s2560-1379b284e07d56ea2516b6dedb07d436.png/Jules_Logo_2025_favicon-s2560.png?dl" width="100px;" height="100px;" style="border-radius: 30%;" alt=""/><br />
      <sub><b>Jules</b></sub><br />
      <sub>Planning</sub>
    </td>
  </tr>
</table>

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=alexandephilia/kiro-claude-proxy&type=Date)](https://star-history.com/#alexandephilia/kiro-claude-proxy)
