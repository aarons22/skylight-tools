# Skylight MCP Server

FastMCP server for interacting with the Skylight API.

## Setup

1. Create a virtual environment and install dependencies:

```bash
python -m venv .venv
source .venv/bin/activate
pip install -e .
```

2. Configure environment variables (copy `.env.example` and edit):

```bash
cp .env.example .env
```

Required:
- `SKYLIGHT_EMAIL`
- `SKYLIGHT_PASSWORD`

Optional:
- `SKYLIGHT_FRAME_ID` (default: auto-pick if only one)

3. Run the server:

```bash
skylight-mcp
```

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

- Authentication uses email/password via `/sessions` and caches `user_id` + `user_token` in `.skylight_token.json` (repo root).
- Requests use Token auth on `api.ourskylight.com` and retry once with Basic auth on `app.ourskylight.com` for 401/404.
- Deleting list items uses the bulk destroy endpoint only.
