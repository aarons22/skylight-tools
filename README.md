# Skylight MCP Server

FastMCP server for interacting with the Skylight API.

## Setup

1. Create a virtual environment and install dependencies:

```bash
python -m venv .venv
source .venv/bin/activate
pip install -e .
```

2. Run interactive setup to create your config:

```bash
skylight-mcp setup
```

This writes `~/Library/Application Support/skylight-mcp/config.toml`.
You can edit this file later if needed.

Environment overrides (optional):
```bash
SKYLIGHT_EMAIL=you@example.com SKYLIGHT_PASSWORD=your_password skylight-mcp run
```

3. Run the server:

```bash
skylight-mcp run
```

## CLI Commands

- `skylight-mcp setup`
- `skylight-mcp run`
- `skylight-mcp install`
- `skylight-mcp uninstall`
- `skylight-mcp status`
- `skylight-mcp logs`

## One-Command Install (macOS user login)

This installs a LaunchAgent that runs on user login.

```bash
curl -sSL https://raw.githubusercontent.com/aarons22/skylight-mcp/main/install.sh | bash
```

Optional:
```bash
SKYLIGHT_FRAME_ID=your_frame_id curl -sSL https://raw.githubusercontent.com/aarons22/skylight-mcp/main/install.sh | bash
```

It installs to `$HOME/.local/skylight-mcp` and runs as the current user.

## Tools

- `get_frames`
- `get_lists`
- `get_list_items`
- `create_list_item`
- `update_list_item`
- `delete_list_items`
- `get_meal_categories`
- `get_meal_recipes`
- `get_meal_recipe`
- `create_meal_recipe`
- `update_meal_recipe`
- `add_recipe_to_grocery_list`
- `get_meal_sittings`
- `create_meal_sitting`
- `update_meal_sitting`
- `delete_meal_sitting`

## Notes

- Authentication uses email/password via `/sessions` and caches `user_id` + `user_token` in `~/Library/Application Support/skylight-mcp/.skylight_token.json`.
- Requests use Token auth on `api.ourskylight.com` and retry once with Basic auth on `app.ourskylight.com` for 401/404.
- Deleting list items uses the bulk destroy endpoint only.
