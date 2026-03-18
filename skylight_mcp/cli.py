from __future__ import annotations

import argparse
import getpass
import os
import shutil
import subprocess
import sys
from pathlib import Path

from .client import SkylightClient
from .config import CONFIG_FILE, get_settings, save_config, token_cache_path
from .server import main as server_main


PLIST_PATH = Path.home() / "Library" / "LaunchAgents" / "com.skylight.mcp.plist"
LOG_DIR = Path.home() / "Library" / "Logs"
OUT_LOG = LOG_DIR / "skylight-mcp.out.log"
ERR_LOG = LOG_DIR / "skylight-mcp.err.log"


def _print(msg: str) -> None:
    print(msg)


def cmd_setup(args: argparse.Namespace) -> int:
    email = args.email or input("Skylight email: ").strip()
    password = args.password or getpass.getpass("Skylight password: ").strip()
    frame_id = args.frame_id

    if not email or not password:
        _print("❌ Email and password are required.")
        return 1

    save_config(email=email, password=password, frame_id=frame_id)

    # Verify auth
    try:
        client = SkylightClient(email=email, password=password, token_cache_path=token_cache_path())
        client.get_frames()
        _print("✅ Credentials verified and config saved.")
        _print(f"Config: {CONFIG_FILE}")
        return 0
    except Exception as exc:
        _print(f"⚠️ Saved config, but auth failed: {exc}")
        _print("Check credentials and try again.")
        return 1


def _resolve_entrypoint() -> str:
    path = shutil.which("skylight-mcp")
    if not path:
        raise RuntimeError("Could not find 'skylight-mcp' in PATH.")
    return path


def _write_plist(entrypoint: str) -> None:
    LOG_DIR.mkdir(parents=True, exist_ok=True)
    PLIST_PATH.parent.mkdir(parents=True, exist_ok=True)
    content = f"""<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">
<plist version=\"1.0\">
  <dict>
    <key>Label</key>
    <string>com.skylight.mcp</string>

    <key>ProgramArguments</key>
    <array>
      <string>{entrypoint}</string>
      <string>run</string>
    </array>

    <key>WorkingDirectory</key>
    <string>{Path.home()}</string>

    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>

    <key>StandardOutPath</key>
    <string>{OUT_LOG}</string>
    <key>StandardErrorPath</key>
    <string>{ERR_LOG}</string>
  </dict>
</plist>
"""
    PLIST_PATH.write_text(content)
    PLIST_PATH.chmod(0o644)


def _launchctl(*args: str) -> int:
    result = subprocess.run(["launchctl", *args], check=False)
    return result.returncode


def cmd_install(_: argparse.Namespace) -> int:
    # Ensure config exists
    try:
        get_settings()
    except Exception as exc:
        _print(f"❌ {exc}")
        _print("Run 'skylight-mcp setup' first.")
        return 1

    try:
        entrypoint = _resolve_entrypoint()
    except Exception as exc:
        _print(f"❌ {exc}")
        return 1

    _write_plist(entrypoint)

    domain = f"gui/{os.getuid()}"
    _launchctl("bootout", domain, str(PLIST_PATH))
    _launchctl("bootstrap", domain, str(PLIST_PATH))
    _launchctl("enable", f"{domain}/com.skylight.mcp")
    _launchctl("kickstart", "-k", f"{domain}/com.skylight.mcp")

    _print("✅ LaunchAgent installed and started.")
    _print(f"Plist: {PLIST_PATH}")
    return 0


def cmd_uninstall(_: argparse.Namespace) -> int:
    domain = f"gui/{os.getuid()}"
    _launchctl("bootout", domain, str(PLIST_PATH))
    if PLIST_PATH.exists():
        PLIST_PATH.unlink()
    _print("✅ LaunchAgent removed.")
    return 0


def cmd_status(_: argparse.Namespace) -> int:
    domain = f"gui/{os.getuid()}"
    code = _launchctl("print", f"{domain}/com.skylight.mcp")
    if code != 0:
        _print("❌ LaunchAgent not running.")
        return 1
    _print("✅ LaunchAgent appears to be running.")
    return 0


def cmd_logs(args: argparse.Namespace) -> int:
    if args.follow:
        return subprocess.run(["tail", "-f", str(OUT_LOG), str(ERR_LOG)], check=False).returncode
    subprocess.run(["tail", "-n", "200", str(OUT_LOG)], check=False)
    subprocess.run(["tail", "-n", "200", str(ERR_LOG)], check=False)
    return 0


def cmd_run(_: argparse.Namespace) -> int:
    server_main()
    return 0


def create_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="skylight-mcp",
        description="Skylight MCP server CLI",
    )

    subparsers = parser.add_subparsers(dest="command", help="Commands")

    setup_parser = subparsers.add_parser("setup", help="Interactive setup")
    setup_parser.add_argument("--email", help="Skylight email")
    setup_parser.add_argument("--password", help="Skylight password")
    setup_parser.add_argument("--frame-id", help="Default frame id")

    subparsers.add_parser("run", help="Run MCP server in foreground")
    subparsers.add_parser("install", help="Install LaunchAgent")
    subparsers.add_parser("uninstall", help="Uninstall LaunchAgent")

    subparsers.add_parser("status", help="Show LaunchAgent status")

    logs_parser = subparsers.add_parser("logs", help="Show logs")
    logs_parser.add_argument("--follow", action="store_true", help="Follow logs")

    return parser


def main() -> int:
    parser = create_parser()
    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return 0

    handlers = {
        "setup": cmd_setup,
        "run": cmd_run,
        "install": cmd_install,
        "uninstall": cmd_uninstall,
        "status": cmd_status,
        "logs": cmd_logs,
    }

    handler = handlers.get(args.command)
    if not handler:
        parser.print_help()
        return 1

    return handler(args)


if __name__ == "__main__":
    sys.exit(main())
