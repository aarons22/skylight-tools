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
def get_categories(frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
    """List family member profiles (categories) for the frame."""
    return _client().get_categories(frame_id)


@mcp.tool()
def get_reward_points(frame_id: Optional[str] = None) -> Any:
    """Get reward point balances per family member profile."""
    return _client().get_reward_points(frame_id)


@mcp.tool()
def get_task_box_items(frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
    """List unscheduled tasks from the task bank."""
    return _client().get_task_box_items(frame_id)


@mcp.tool()
def create_task_box_item(
    summary: str,
    emoji_icon: Optional[str] = None,
    routine: bool = False,
    reward_points: Optional[int] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    """Create a new task in the task bank."""
    return _client().create_task_box_item(frame_id, summary, emoji_icon, routine, reward_points)


@mcp.tool()
def update_task_box_item(
    item_id: str,
    summary: Optional[str] = None,
    emoji_icon: Optional[str] = None,
    routine: Optional[bool] = None,
    reward_points: Optional[int] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    """Update a task in the task bank."""
    return _client().update_task_box_item(frame_id, item_id, summary, emoji_icon, routine, reward_points)


@mcp.tool()
def delete_task_box_item(item_id: str, frame_id: Optional[str] = None) -> Dict[str, Any]:
    """Delete a task from the task bank."""
    return _client().delete_task_box_item(frame_id, item_id)


@mcp.tool()
def get_chores(
    after: Optional[str] = None,
    before: Optional[str] = None,
    include_late: Optional[bool] = None,
    filter_profile: Optional[str] = None,
    frame_id: Optional[str] = None,
) -> Any:
    """Get chores for a date range. Optionally filter by profile name. Dates are YYYY-MM-DD."""
    return _client().get_chores(frame_id, after, before, include_late, filter_profile)


@mcp.tool()
def create_chore(
    summary: str,
    start: str,
    routine: bool = False,
    start_time: Optional[str] = None,
    recurrence_set: Optional[str] = None,
    emoji_icon: Optional[str] = None,
    recurring_until: Optional[str] = None,
    reward_points: Optional[int] = None,
    category_ids: Optional[List[str]] = None,
    category_id: Optional[str] = None,
    profile_name: Optional[str] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    """Create a scheduled chore or routine. Use profile_name to assign to a family member. start is YYYY-MM-DD. recurrence_set uses RRULE format."""
    return _client().create_chore(
        frame_id, summary, start, routine, start_time, recurrence_set,
        emoji_icon, recurring_until, reward_points, category_ids, category_id, profile_name,
    )


@mcp.tool()
def update_chore(
    chore_id: str,
    summary: Optional[str] = None,
    emoji_icon: Optional[str] = None,
    reward_points: Optional[int] = None,
    start: Optional[str] = None,
    start_time: Optional[str] = None,
    recurrence_set: Optional[str] = None,
    recurring_until: Optional[str] = None,
    routine: Optional[bool] = None,
    category_id: Optional[str] = None,
    profile_name: Optional[str] = None,
    up_for_grabs: Optional[bool] = None,
    frame_id: Optional[str] = None,
) -> Dict[str, Any]:
    """Update/reschedule a chore. Sends full object via PUT. For recurring chores, chore_id is composite: {series}-{date} or {series}-{date}-{time}."""
    return _client().update_chore(
        frame_id, chore_id, summary, emoji_icon, reward_points, start,
        start_time, recurrence_set, recurring_until, routine, category_id, profile_name, up_for_grabs,
    )


@mcp.tool()
def complete_chore(chore_id: str, frame_id: Optional[str] = None) -> Dict[str, Any]:
    """Mark a chore as complete."""
    return _client().complete_chore(frame_id, chore_id)


@mcp.tool()
def delete_chore(chore_id: str, apply_to: str = "one", frame_id: Optional[str] = None) -> Dict[str, Any]:
    """Delete a chore instance or entire series. apply_to: 'one' or 'all'."""
    return _client().delete_chore(frame_id, chore_id, apply_to)


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
    settings = get_settings()
    mcp.run(transport="http", host="127.0.0.1", port=settings.skylight_port)


if __name__ == "__main__":
    main()
