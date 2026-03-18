#!/usr/bin/env bash
set -euo pipefail

APP_NAME="skylight-mcp"
DEFAULT_REPO_URL="https://github.com/aarons22/skylight-mcp"
DEFAULT_INSTALL_DIR="$HOME/.local/skylight-mcp"

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

# Run guided setup if config missing, then install LaunchAgent
CONFIG_FILE="$HOME/Library/Application Support/skylight-mcp/config.toml"
if [[ ! -f "$CONFIG_FILE" ]]; then
  if [[ -n "$FRAME_ID" ]]; then
    "$INSTALL_DIR/.venv/bin/skylight-mcp" setup --frame-id "$FRAME_ID"
  else
    "$INSTALL_DIR/.venv/bin/skylight-mcp" setup
  fi
fi

"$INSTALL_DIR/.venv/bin/skylight-mcp" install

cat <<OUT
Installed and started $APP_NAME.
Logs:
  $HOME/Library/Logs/skylight-mcp.out.log
  $HOME/Library/Logs/skylight-mcp.err.log
OUT
