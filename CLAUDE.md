# Skylight MCP Server

FastMCP server exposing the Skylight Calendar API as MCP tools.

## Repo Structure

- `skylight_mcp/client.py` — API client (HTTP calls, auth, resolver helpers)
- `skylight_mcp/server.py` — `@mcp.tool()` definitions (thin wrappers around client methods)
- `skylight_mcp/config.py` — Settings dataclass, TOML config, env var overrides
- `skylight_mcp/cli.py` — CLI: setup, run, install/uninstall LaunchAgent

## Architecture

Two-layer pattern:
1. **Client** (`SkylightClient`): handles API calls, auth, frame resolution, name-to-ID resolvers with caching
2. **Server** (`server.py`): thin `@mcp.tool()` wrappers that instantiate a fresh `_client()` per call

## Auth Flow

`email/password` → `POST /sessions` → `base64(user_id:user_token)` as `Token` header. Falls back to `Basic` auth on `app.ourskylight.com` if `api.ourskylight.com` fails.

## Key Conventions

- `frame_id` is always the last parameter and optional — auto-resolved via `resolve_frame_id()`
- Resolver pattern: `resolve_X_id(frame_id, id, name)` for name-to-ID lookups with caching
- Response parsing handles both `{data: [...]}` and raw list responses
- Client methods return raw API response dicts

## API Reference

See `API_REFERENCE.md` for endpoint details and request/response examples.

## Testing

See `AGENTS.md` for E2E validation sequence.

## Config

- Config file: `~/Library/Application Support/skylight-mcp/config.toml`
- Env var overrides: `SKYLIGHT_EMAIL`, `SKYLIGHT_PASSWORD`, `SKYLIGHT_FRAME_ID`, `SKYLIGHT_PORT`
- Token cache: `~/Library/Application Support/skylight-mcp/.skylight_token.json`
