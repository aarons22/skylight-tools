# Skylight MCP Server

FastMCP server and CLI for interacting with the Skylight Calendar API.

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

---

## Skylight CLI (Go)

A standalone CLI for the Skylight API, following the structure of [paprika-tools](https://github.com/aarons22/paprika-tools).

### Build

```bash
cd skylight
go build -o skylight .
```

### Login

```bash
./skylight account login --email you@example.com --password secret
# Optionally save your default frame ID:
./skylight account login --email you@example.com --password secret --save-frame-id <frameId>
```

Token is saved to `~/.config/skylight/config.yaml`. You can also set `SKYLIGHT_TOKEN` and `SKYLIGHT_FRAME_ID` env vars.

### Commands

| Group | Subcommands |
|---|---|
| `frames` | `listFrames` |
| `lists` | `listLists`, `listItems`, `createItem`, `updateItem`, `deleteItems` |
| `categories` | `listCategories` |
| `reward-points` | `get` |
| `task-box` | `listItems`, `createItem`, `updateItem`, `deleteItem` |
| `chores` | `listChores`, `createChore`, `updateChore`, `completeChore`, `deleteChore` |
| `meals` | `listCategories`, `listRecipes`, `getRecipe`, `createRecipe`, `listSittings`, `createSitting`, `deleteSitting` |
| `account` | `login` |
| `completion` | `bash`, `zsh`, `fish`, `powershell` |

### Global Flags

```
--json          Output raw JSON instead of a table
--no-color      Disable colored output
--frame-id      Override frame ID
--base-url      Override API base URL
--config        Override config file path
```

### Examples

```bash
# List frames
./skylight frames listFrames

# List grocery lists
./skylight lists listLists --frame-id <frameId>

# Add a grocery item
./skylight lists createItem --frame-id <frameId> --list-id <listId> --name "Milk"

# List this week's chores
./skylight chores listChores --frame-id <frameId> --after 2026-03-24 --before 2026-03-30

# Create a recurring chore
./skylight chores createChore --frame-id <frameId> \
  --summary "Clean room" --start 2026-03-29 \
  --category-ids "<catId>" --recurrence-set "RRULE:FREQ=WEEKLY;INTERVAL=1;BYDAY=SA"

# List meal plan for a date range
./skylight meals listSittings --frame-id <frameId> --date-min 2026-03-01 --date-max 2026-03-31

# Raw JSON output
./skylight --json frames listFrames
```

### API Spec

See [`openapi.yaml`](./openapi.yaml) for the full OpenAPI 3.0.3 specification.
