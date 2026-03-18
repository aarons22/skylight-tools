from __future__ import annotations

import os
from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    skylight_email: str
    skylight_password: str
    skylight_frame_id: str | None = None

    model_config = SettingsConfigDict(
        env_file=".env",
        env_prefix="",
        case_sensitive=False,
    )


def default_token_cache_path() -> Path:
    # Default to repo root .skylight_token.json
    cwd = Path.cwd()
    return cwd / ".skylight_token.json"


def get_settings() -> Settings:
    return Settings()


def token_cache_path() -> Path:
    return Path.cwd() / ".skylight_token.json"
