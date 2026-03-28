#!/usr/bin/env bash
set -euo pipefail

APP_NAME="skylight-mcp"
DEFAULT_REPO_URL="https://github.com/aarons22/skylight-tools"
DEFAULT_INSTALL_DIR="$HOME/.local/skylight-tools"

REPO_URL="${REPO_URL:-$DEFAULT_REPO_URL}"
INSTALL_DIR="${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}"
FRAME_ID="${SKYLIGHT_FRAME_ID:-}"

if [[ $(uname -s) != "Darwin" ]]; then
  echo "This installer currently supports macOS only." >&2
  exit 1
fi

if [[ $EUID -eq 0 ]]; then
  echo "Please run as your normal user (no sudo) so the agent installs in your user profile." >&2
  exit 1
fi

if ! command -v git >/dev/null 2>&1; then
  echo "git is required but not installed." >&2
  exit 1
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "python3 is required but not installed." >&2
  exit 1
fi

# Install repo
mkdir -p "$(dirname "$INSTALL_DIR")"
if [[ ! -d "$INSTALL_DIR/.git" ]]; then
  git clone "$REPO_URL" "$INSTALL_DIR"
else
  git -C "$INSTALL_DIR" pull --ff-only
fi

# Setup venv + install
python3 -m venv "$INSTALL_DIR/.venv"
"$INSTALL_DIR/.venv/bin/pip" install -e "$INSTALL_DIR"

# Ensure CLI is on PATH via ~/.local/bin
BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"
ln -sf "$INSTALL_DIR/.venv/bin/skylight-mcp" "$BIN_DIR/skylight-mcp"

CLI_BIN="$INSTALL_DIR/.venv/bin/skylight-mcp"
if [[ ! -x "$CLI_BIN" ]]; then
  echo "skylight-mcp entrypoint not found at $CLI_BIN" >&2
  echo "Try re-running install or ensure dependencies installed." >&2
  exit 1
fi

# Install LaunchAgent if config exists; otherwise leave setup to user
"$CLI_BIN" install || true

cat <<OUT
Installed and started $APP_NAME.
CLI:
  Ensure $BIN_DIR is on your PATH, then run: skylight-mcp --help
Next:
  Run 'skylight-mcp setup' to configure credentials and port.
Logs:
  $HOME/Library/Logs/skylight-mcp.out.log
  $HOME/Library/Logs/skylight-mcp.err.log
OUT
