from __future__ import annotations

from datetime import date
from typing import Any, Dict, List, Optional

from fastmcp import FastMCP

from .client import SkylightClient
from .config import get_settings, token_cache_path

mcp = FastMCP("skylight")


def _client() -> SkylightClient:
    settings = get_settings()
    return SkylightClient(
        email=settings.skylight_email,
        password=settings.skylight_password,
        token_cache_path=token_cache_path(),
        default_frame_id=settings.skylight_frame_id,
    )


@mcp.tool()
def get_frames() -> List[Dict[str, Any]]:
    return _client().get_frames()


@mcp.tool()
def get_lists(frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
    return _client().get_lists(frame_id)


@mcp.tool()
def get_list_items(frame_id: Optional[str] = None, list_id: Optional[str] = None, list_name: Optional[str] = None) -> List[Dict[str, Any]]:
    return _client().get_list_items(frame_id, list_id, list_name)


@mcp.tool()
def create_list_item(
    name: str,
    checked: bool = False,
    frame_id: Optional[str] = None,
    list_id: Optional[str] = None,
    list_name: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().create_list_item(frame_id, list_id, list_name, name, checked)


@mcp.tool()
def update_list_item(
    item_id: str,
    name: Optional[str] = None,
    checked: Optional[bool] = None,
    frame_id: Optional[str] = None,
    list_id: Optional[str] = None,
    list_name: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().update_list_item(frame_id, list_id, list_name, item_id, name, checked)


@mcp.tool()
def delete_list_items(
    ids: List[str],
    frame_id: Optional[str] = None,
    list_id: Optional[str] = None,
    list_name: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().delete_list_items(frame_id, list_id, list_name, ids)


@mcp.tool()
def get_meal_categories(frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
    return _client().get_meal_categories(frame_id)


@mcp.tool()
def get_meal_recipes(frame_id: Optional[str] = None) -> Dict[str, Any]:
    return _client().get_meal_recipes(frame_id)


@mcp.tool()
def get_meal_recipe(recipe_id: str, frame_id: Optional[str] = None) -> Dict[str, Any]:
    return _client().get_meal_recipe(frame_id, recipe_id)


@mcp.tool()
def create_meal_recipe(
    summary: str,
    description: Optional[str] = None,
    meal_category_id: Optional[str] = None,
    meal_type: Optional[str] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().create_meal_recipe(frame_id, summary, description, meal_category_id, meal_type)


@mcp.tool()
def update_meal_recipe(
    recipe_id: str,
    summary: str,
    description: Optional[str] = None,
    meal_category_id: Optional[str] = None,
    meal_type: Optional[str] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().update_meal_recipe(frame_id, recipe_id, summary, description, meal_category_id, meal_type)


@mcp.tool()
def add_recipe_to_grocery_list(recipe_id: str, frame_id: Optional[str] = None) -> Dict[str, Any]:
    return _client().add_recipe_to_grocery_list(frame_id, recipe_id)


@mcp.tool()
def get_meal_sittings(frame_id: Optional[str], date_min: date, date_max: date) -> Dict[str, Any]:
    return _client().get_meal_sittings(frame_id, date_min, date_max)


@mcp.tool()
def create_meal_sitting(
    date_value: date,
    meal_category_id: Optional[str] = None,
    meal_type: Optional[str] = None,
    meal_recipe_id: Optional[str] = None,
    summary: Optional[str] = None,
    note: Optional[str] = None,
    description: Optional[str] = None,
    add_to_grocery_list: bool = False,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().create_meal_sitting(
        frame_id,
        date_value,
        meal_category_id,
        meal_type,
        meal_recipe_id,
        summary,
        note,
        description,
        add_to_grocery_list,
    )


@mcp.tool()
def update_meal_sitting(
    sitting_id: str,
    date_value: date,
    meal_category_id: Optional[str] = None,
    meal_type: Optional[str] = None,
    meal_recipe_id: Optional[str] = None,
    summary: Optional[str] = None,
    note: Optional[str] = None,
    description: Optional[str] = None,
    add_to_grocery_list: bool = False,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().update_meal_sitting(
        frame_id,
        sitting_id,
        date_value,
        meal_category_id,
        meal_type,
        meal_recipe_id,
        summary,
        note,
        description,
        add_to_grocery_list,
    )


@mcp.tool()
def delete_meal_sitting(
    sitting_id: str,
    date_value: Optional[date] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    return _client().delete_meal_sitting(frame_id, sitting_id, date_value)


def main() -> None:
    mcp.run()


if __name__ == "__main__":
    main()
