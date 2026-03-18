from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, Optional


try:
    import tomllib as toml
except ModuleNotFoundError:  # Python < 3.11
    import tomli as toml


CONFIG_DIR = Path.home() / "Library" / "Application Support" / "skylight-mcp"
CONFIG_FILE = CONFIG_DIR / "config.toml"
TOKEN_CACHE_FILE = CONFIG_DIR / ".skylight_token.json"


@dataclass
class Settings:
    skylight_email: str
    skylight_password: str
    skylight_frame_id: Optional[str] = None


def ensure_config_dir() -> None:
    CONFIG_DIR.mkdir(parents=True, exist_ok=True)


def load_config_file() -> Dict[str, Any]:
    if not CONFIG_FILE.exists():
        return {}
    data = toml.loads(CONFIG_FILE.read_text())
    return data if isinstance(data, dict) else {}


def save_config(email: str, password: str, frame_id: Optional[str]) -> None:
    ensure_config_dir()
    lines = ["[skylight]", f"email = \"{email}\"", f"password = \"{password}\""]
    if frame_id:
        lines.append(f"frame_id = \"{frame_id}\"")
    CONFIG_FILE.write_text("\n".join(lines) + "\n")
    CONFIG_FILE.chmod(0o600)


def get_settings(overrides: Optional[Dict[str, str]] = None) -> Settings:
    data = load_config_file()
    skylight = data.get("skylight", {}) if isinstance(data.get("skylight"), dict) else {}

    email = skylight.get("email")
    password = skylight.get("password")
    frame_id = skylight.get("frame_id")

    # Env overrides
    import os

    email = os.getenv("SKYLIGHT_EMAIL", email)
    password = os.getenv("SKYLIGHT_PASSWORD", password)
    frame_id = os.getenv("SKYLIGHT_FRAME_ID", frame_id)

    # Explicit overrides (CLI args)
    if overrides:
        email = overrides.get("email", email)
        password = overrides.get("password", password)
        frame_id = overrides.get("frame_id", frame_id)

    if not email or not password:
        raise ValueError("Missing Skylight credentials. Run 'skylight-mcp setup'.")

    return Settings(
        skylight_email=email,
        skylight_password=password,
        skylight_frame_id=frame_id,
    )


def token_cache_path() -> Path:
    ensure_config_dir()
    return TOKEN_CACHE_FILE
