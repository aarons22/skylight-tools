# Skylight

Tools for interacting with the [Skylight Calendar](https://www.ourskylight.com).

| | What it is |
|---|---|
| [`openapi.yaml`](./openapi.yaml) | OpenAPI 3.0.3 spec — machine-readable definition of the unofficial Skylight API |
| [`skylight/`](./skylight/) | Go CLI — resource-grouped commands with table + JSON output and shell completions |
| [`skylight_mcp/`](./skylight_mcp/) | Python FastMCP server — exposes Skylight data as MCP tools for AI agents |

---

## CLI

### Install

**Build from source** (requires Go 1.21+):

```bash
git clone https://github.com/aarons22/skylight-tools
cd skylight-tools/skylight
go build -o skylight .
mv skylight /usr/local/bin/
```

**go install** (once the repo is tagged):

```bash
go install github.com/aarons22/skylight-tools/skylight@latest
```

### Authenticate

```bash
skylight account login --email you@example.com --password yourpassword
```

The token is saved automatically to `~/.config/skylight/config.yaml` (mode 0600). All subsequent commands read it from there — no manual copy-paste required.

The token is also accepted via the `SKYLIGHT_TOKEN` environment variable. The frame ID can be set via `SKYLIGHT_FRAME_ID` or saved during login with `--save-frame-id <frameId>`.

### Commands

```
skylight account login                                             # authenticate and save token

skylight frames listFrames                                         # list all Skylight frames

skylight lists listLists                                           # all lists for a frame
skylight lists listItems --list-id <id>                            # items in a list
skylight lists createItem --list-id <id> --name "Milk"             # add a list item
skylight lists updateItem --list-id <id> --item-id <id>            # update a list item
skylight lists deleteItems --list-id <id> --ids "id1,id2"          # bulk delete items

skylight categories listCategories                                 # family member profiles

skylight reward-points get                                         # point balances per profile

skylight task-box listItems                                        # task bank
skylight task-box createItem --summary "Clean room"                # add a task
skylight task-box updateItem --item-id <id>                        # update a task
skylight task-box deleteItem --item-id <id>                        # delete a task

skylight chores listChores                                         # scheduled chores
skylight chores createChore --summary "..." --start 2026-03-29 --category-ids <id>
skylight chores updateChore --chore-id <id>                        # update a chore
skylight chores completeChore --chore-id <id>                      # mark complete
skylight chores deleteChore --chore-id <id>                        # delete a chore

skylight meals listCategories                                      # Breakfast/Lunch/Dinner/Snack
skylight meals listRecipes                                         # meal recipe bank
skylight meals getRecipe --recipe-id <id>                          # single recipe
skylight meals createRecipe --category-id <id> --summary "..."     # add a recipe
skylight meals listSittings --date-min 2026-03-01 --date-max 2026-03-31
skylight meals createSitting --category-id <id> --date 2026-03-01  # schedule a meal
skylight meals deleteSitting --sitting-id <id> --date 2026-03-01   # remove a sitting
```

### Output & global flags

```bash
# Raw JSON (pipe-friendly)
skylight frames listFrames --json | jq '.[].id'

# Specify frame ID (or set SKYLIGHT_FRAME_ID)
skylight lists listLists --frame-id <frameId>

# Disable colour
skylight frames listFrames --no-color

# Shell completions (bash, zsh, or fish)
skylight completion bash >> ~/.bashrc
```

---

## MCP Server

FastMCP server that exposes Skylight data as tools for AI agents (Claude, Cursor, etc.).

### Quick Install (macOS)

```bash
curl -sSL https://raw.githubusercontent.com/aarons22/skylight-tools/main/install.sh | bash
```

Then run:

```bash
skylight-mcp setup
skylight-mcp install
```

If `skylight-mcp` isn't on your PATH, use:

```bash
$HOME/.local/bin/skylight-mcp --help
```

### CLI Commands

| Command | Description |
|---------|-------------|
| `skylight-mcp setup` | Interactive credential and port setup |
| `skylight-mcp run` | Run the MCP server in the foreground |
| `skylight-mcp install` | Install as a macOS LaunchAgent (background service) |
| `skylight-mcp uninstall` | Remove the LaunchAgent |
| `skylight-mcp status` | Check LaunchAgent status |
| `skylight-mcp logs` | View server logs |

### Available Tools

| Tool | Description |
|------|-------------|
| `get_frames` | List all Skylight frames |
| `get_lists(frame_id?)` | List grocery/shopping lists |
| `get_list_items(list_id, frame_id?)` | List items in a list |
| `create_list_item(list_id, name, ...)` | Add an item to a list |
| `update_list_item(list_id, item_id, ...)` | Update a list item |
| `delete_list_items(list_id, ids, frame_id?)` | Bulk delete list items |
| `get_categories(frame_id?)` | List family member profile categories |
| `get_reward_points(frame_id?)` | Get reward point balances per profile |
| `get_task_box_items(frame_id?)` | List task bank items |
| `create_task_box_item(summary, ...)` | Add a task bank item |
| `update_task_box_item(item_id, ...)` | Update a task bank item |
| `delete_task_box_item(item_id, frame_id?)` | Delete a task bank item |
| `get_chores(frame_id?, after?, before?, ...)` | List scheduled chores |
| `create_chore(summary, start, category_ids, ...)` | Create chores |
| `update_chore(chore_id, ...)` | Update a chore |
| `complete_chore(chore_id, frame_id?)` | Mark a chore complete |
| `delete_chore(chore_id, frame_id?)` | Delete a chore |
| `get_meal_categories(frame_id?)` | List meal categories (Breakfast/Lunch/Dinner/Snack) |
| `get_meal_recipes(frame_id?)` | List meal recipes |
| `get_meal_recipe(recipe_id, frame_id?)` | Get a single meal recipe |
| `create_meal_recipe(meal_category_id, summary, ...)` | Add a meal recipe |
| `get_meal_sittings(date_min, date_max, frame_id?)` | List meal sittings in a date range |
| `create_meal_sitting(meal_category_id, date, ...)` | Schedule a meal |
| `delete_meal_sitting(sitting_id, date, frame_id?)` | Remove a meal sitting |

### Config & Logs

- Config: `~/Library/Application Support/skylight-mcp/config.toml`
- Token cache: `~/Library/Application Support/skylight-mcp/.skylight_token.json`
- Stdout log: `~/Library/Logs/skylight-mcp.out.log`
- Stderr log: `~/Library/Logs/skylight-mcp.err.log`

### Notes

- Frame ID is required for most operations; save a default with `skylight account login --save-frame-id <id>` or set `SKYLIGHT_FRAME_ID`
- The Skylight API is unofficial and undocumented; see [`API_REFERENCE.md`](./API_REFERENCE.md) for details
- Individual list item deletion is broken in the API; bulk destroy is used automatically
- Tokens are cached on disk and refreshed automatically on 401

---

## OpenAPI Spec

`openapi.yaml` is the authoritative machine-readable definition of the Skylight API. It covers authentication, frames, lists, categories, reward points, task box, chores, and meal planning.

Use it with any OpenAPI-compatible tooling — code generators, HTTP clients, or documentation renderers.

---

## Development

### Updating the OpenAPI spec

`openapi.yaml` is the source of truth. When it changes, regenerate the Go CLI using [commandspec](https://github.com/theaiteam-dev/commandspec), then reapply the hand-written customizations listed in `CLAUDE.md`. Do not edit `skylight/cmd/` files directly.

The Python MCP (`skylight_mcp/`) is edited directly — it is not generated.

### Validating changes

See `AGENTS.md` for the E2E validation sequence for each component.
