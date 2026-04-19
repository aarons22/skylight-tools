# Skylight Tools

Three components for interacting with the Skylight Calendar API.

| Component | What it is |
|---|---|
| `openapi.yaml` | OpenAPI 3.0.3 spec — source of truth for the unofficial Skylight API |
| `skylight/` | Go CLI — resource-grouped commands with table + JSON output |
| `skylight_mcp/` | Python FastMCP server — exposes Skylight data as MCP tools |

---

## Go CLI (`skylight/`)

### Structure

- `cmd/` — Cobra subcommands, one file per operation (e.g. `lists_listItems.go`)
- `internal/client/client.go` — HTTP client, auth, config/YAML storage, token and frame ID resolution
- `internal/output/output.go` — table and JSON output formatting, JSON:API flattening
- `main.go` — entry point

### Auth & Config

Token saved to `~/.config/skylight/config.yaml` after `skylight account login`. All commands resolve token via `client.ResolveToken()` (env `SKYLIGHT_TOKEN` → config file). Frame ID resolves via `client.ResolveFrameID()` (flag → env `SKYLIGHT_FRAME_ID` → config file).

---

## OpenAPI → Go CLI Workflow

`openapi.yaml` is the source of truth. The `cmd/*.go` files are generated from it using [commandspec](https://github.com/theaiteam-dev/commandspec). **Do not edit `cmd/` files directly** — changes will be overwritten on the next regeneration.

When `openapi.yaml` changes:
1. Run commandspec to regenerate `skylight/cmd/`
2. Reapply the customizations below

### Customizations to Preserve

These files are hand-written and must be restored after any commandspec regeneration:

| File | What to preserve |
|---|---|
| `internal/client/client.go` | Entire file — Config/YAML storage, `Login`, `ResolveToken`, `ResolveFrameID` |
| `internal/output/output.go` | Entire file — table/JSON output, JSON:API flattening |
| `cmd/root.go` | Persistent flags: `--json`, `--no-color`, `--base-url`, `--frame-id`, `--config` |
| `cmd/account_login.go` | YAML token + frame-id save logic |

---

## Python MCP (`skylight_mcp/`)

### Structure

- `client.py` — `SkylightClient`: API calls, auth, frame resolution, name-to-ID resolvers with caching
- `server.py` — `@mcp.tool()` definitions: thin wrappers, one per operation
- `config.py` — `Settings` dataclass, TOML config, env var overrides
- `cli.py` — `skylight-mcp` CLI: setup, run, install/uninstall LaunchAgent

### Architecture

Two-layer pattern:
1. **Client** (`SkylightClient`): handles API calls, auth, frame resolution, name-to-ID resolvers with caching
2. **Server** (`server.py`): thin `@mcp.tool()` wrappers that instantiate a fresh `_client()` per call

### Auth Flow

`email/password` → `POST /sessions` → `base64(user_id:user_token)` as `Token` header. Falls back to `Basic` auth on `app.ourskylight.com` if `api.ourskylight.com` fails.

### Key Conventions

- `frame_id` is always the last parameter and optional — auto-resolved via `resolve_frame_id()`
- Resolver pattern: `resolve_X_id(frame_id, id, name)` for name-to-ID lookups with caching
- Response parsing handles both `{data: [...]}` and raw list responses
- Client methods return raw API response dicts

### Config

- Config file: `~/Library/Application Support/skylight-mcp/config.toml`
- Env var overrides: `SKYLIGHT_EMAIL`, `SKYLIGHT_PASSWORD`, `SKYLIGHT_FRAME_ID`, `SKYLIGHT_PORT`
- Token cache: `~/Library/Application Support/skylight-mcp/.skylight_token.json`

---

## API Reference

See `API_REFERENCE.md` for endpoint details and request/response examples.
