# AGENTS

## End-to-End Validation (Required)
When making changes that affect Skylight auth, network handling, or the MCP server runtime, run a local end-to-end validation before declaring success.

**Minimum E2E sequence:**
1. `python - <<'PY'`
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

2. `skylight-mcp run` and verify the server starts without errors in logs.

**Notes:**
- If DNS issues occur for `api.ourskylight.com`, ensure fallback to `app.ourskylight.com` is exercised.
- If auth fails, verify the login payloads (`user` wrapper and flat email/password) and JSON:API parsing.
