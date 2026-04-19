# AGENTS

End-to-end validation procedures. Run the relevant section after making changes.

---

## Go CLI

After changes to `skylight/`, verify the binary builds and authenticates:

```bash
cd skylight && go build -o skylight . && echo "build ok"
```

Smoke test (requires `~/.config/skylight/config.yaml` with a saved token):

```bash
./skylight frames listFrames
./skylight lists listLists
```

If no config exists, authenticate first:

```bash
./skylight account login --email "$SKYLIGHT_EMAIL" --password "$SKYLIGHT_PASSWORD"
```

---

## OpenAPI → commandspec Regeneration

After editing `openapi.yaml` and running commandspec to regenerate `skylight/cmd/`:

1. Verify the following files still contain their hand-written customizations (see CLAUDE.md for details):
   - `skylight/internal/client/client.go`
   - `skylight/internal/output/output.go`
   - `skylight/cmd/root.go` (persistent flags)
   - `skylight/cmd/account_login.go` (YAML token save)

2. Rebuild and run the CLI smoke test above.

---

## Python MCP Server

After changes to `skylight_mcp/`, run a client smoke test:

```bash
python - <<'PY'
from skylight_mcp.client import SkylightClient
from skylight_mcp.config import get_settings, token_cache_path

settings = get_settings()
client = SkylightClient(
    email=settings.skylight_email,
    password=settings.skylight_password,
    token_cache_path=token_cache_path(),
    default_frame_id=settings.skylight_frame_id,
)
frames = client.get_frames()
print("frames:", frames)
print("count:", len(frames) if isinstance(frames, list) else "n/a")
PY
```

Then verify the server starts:

```bash
skylight-mcp run
```

**Notes:**
- `SKYLIGHT_EMAIL` and `SKYLIGHT_PASSWORD` may be available as env vars or in a local `.env` file.
- If DNS resolution fails for `api.ourskylight.com`, the client should fall back to `app.ourskylight.com`.
- If auth fails, verify the login payload (flat `email`/`password`) and JSON:API response parsing.
