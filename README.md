# Skylight MCP Server

FastMCP server for interacting with the Skylight API.

## Quick Install (macOS)

```bash
curl -sSL https://raw.githubusercontent.com/aarons22/skylight-mcp/main/install.sh | bash
```

Then run:

```bash
skylight-mcp setup
skylight-mcp install
```

If `skylight-mcp` isn’t on your PATH, use:

```bash
$HOME/.local/bin/skylight-mcp --help
```

## CLI Commands

- `skylight-mcp setup`
- `skylight-mcp run`
- `skylight-mcp install`
- `skylight-mcp uninstall`
- `skylight-mcp status`
- `skylight-mcp logs`

## Config

- `~/Library/Application Support/skylight-mcp/config.toml`
- Token cache: `~/Library/Application Support/skylight-mcp/.skylight_token.json`

## Notes

- Auth uses email/password via `/sessions` and caches `user_id` + `user_token`.
- Requests use Token auth on `api.ourskylight.com` and retry once with Basic auth on `app.ourskylight.com` for 401/404.
- Deleting list items uses the bulk destroy endpoint only.
