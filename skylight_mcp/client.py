from __future__ import annotations

import base64
import json
import os
from dataclasses import dataclass
from datetime import date
from pathlib import Path
from typing import Any, Dict, List, Optional

import httpx


@dataclass
class SkylightAuth:
    user_id: str
    user_token: str


class SkylightClient:
    DEFAULT_BASE_URL = "https://api.ourskylight.com/api"
    DEFAULT_FALLBACK_BASE_URL = "https://app.ourskylight.com/api"

    def __init__(
        self,
        email: str,
        password: str,
        token_cache_path: Path,
        default_frame_id: Optional[str] = None,
    ) -> None:
        self.email = email
        self.password = password
        self.base_url = self.DEFAULT_BASE_URL
        self.fallback_base_url = self.DEFAULT_FALLBACK_BASE_URL
        self.token_cache_path = token_cache_path
        self.default_frame_id = default_frame_id
        self._auth: Optional[SkylightAuth] = None
        self._client = httpx.Client(timeout=15.0)
        self._frames_cache: Optional[List[Dict[str, Any]]] = None
        self._lists_cache: Optional[List[Dict[str, Any]]] = None
        self._categories_cache: Optional[List[Dict[str, Any]]] = None

    def _load_cached_token(self) -> Optional[SkylightAuth]:
        if not self.token_cache_path.exists():
            return None
        try:
            data = json.loads(self.token_cache_path.read_text())
            if data.get("email") != self.email:
                return None
            user_id = str(data.get("user_id")) if data.get("user_id") else None
            user_token = data.get("user_token")
            if user_id and user_token:
                return SkylightAuth(user_id=user_id, user_token=user_token)
        except Exception:
            return None
        return None

    def _cache_token(self, auth: SkylightAuth) -> None:
        payload = {
            "email": self.email,
            "user_id": auth.user_id,
            "user_token": auth.user_token,
        }
        self.token_cache_path.write_text(json.dumps(payload))
        os.chmod(self.token_cache_path, 0o600)

    def _login(self, base_url: str) -> SkylightAuth:
        url = f"{base_url}/sessions"
        payloads = [
            {"user": {"email": self.email, "password": self.password}},
            {"email": self.email, "password": self.password},
        ]
        last_error: Optional[Exception] = None

        for payload in payloads:
            try:
                resp = self._client.post(url, json=payload)
                resp.raise_for_status()
                data = resp.json()

                # Legacy response
                user_id = data.get("user_id") or data.get("id")
                user_token = data.get("user_token") or data.get("token") or data.get("auth_token")

                # JSON:API response
                if not user_id or not user_token:
                    data_obj = data.get("data") if isinstance(data, dict) else None
                    attrs = data_obj.get("attributes", {}) if isinstance(data_obj, dict) else {}
                    user_id = user_id or data_obj.get("id") if isinstance(data_obj, dict) else None
                    user_token = user_token or attrs.get("user_token") or attrs.get("token") or attrs.get("auth_token")

                if user_id and user_token:
                    return SkylightAuth(user_id=str(user_id), user_token=str(user_token))
            except Exception as exc:
                last_error = exc

        if isinstance(last_error, httpx.RequestError):
            raise last_error
        if isinstance(last_error, httpx.HTTPStatusError):
            raise last_error
        raise ValueError("Login response missing user_id/user_token") from last_error

    def _ensure_auth(self) -> SkylightAuth:
        if self._auth:
            return self._auth
        cached = self._load_cached_token()
        if cached:
            self._auth = cached
            return cached
        try:
            auth = self._login(self.base_url)
        except httpx.RequestError:
            auth = self._login(self.fallback_base_url)
            self.base_url = self.fallback_base_url
        except httpx.HTTPStatusError as exc:
            status = exc.response.status_code if exc.response is not None else None
            if status in (404, 502, 503, 504):
                auth = self._login(self.fallback_base_url)
                self.base_url = self.fallback_base_url
            else:
                raise
        self._cache_token(auth)
        self._auth = auth
        return auth

    def _token_header(self, auth: SkylightAuth) -> str:
        token = base64.b64encode(f"{auth.user_id}:{auth.user_token}".encode()).decode()
        return f'Token token="{token}"'

    def _basic_header(self, auth: SkylightAuth) -> str:
        token = base64.b64encode(f"{auth.user_id}:{auth.user_token}".encode()).decode()
        return f"Basic {token}"

    def _request(
        self,
        method: str,
        endpoint: str,
        json_body: Optional[Dict[str, Any]] = None,
        params: Optional[Dict[str, Any]] = None,
    ) -> Any:
        auth = self._ensure_auth()
        url = f"{self.base_url}{endpoint}"
        use_basic = self.base_url == self.fallback_base_url
        headers = {
            "Accept": "application/json",
            "Content-Type": "application/json",
            "Authorization": self._basic_header(auth) if use_basic else self._token_header(auth),
            "User-Agent": "SkylightMCP",
        }

        try:
            resp = self._client.request(method, url, json=json_body, params=params, headers=headers)
        except httpx.RequestError:
            # Retry with fallback host + Basic auth
            fallback_url = f"{self.fallback_base_url}{endpoint}"
            headers["Authorization"] = self._basic_header(auth)
            resp = self._client.request(method, fallback_url, json=json_body, params=params, headers=headers)
            self.base_url = self.fallback_base_url

        if resp.status_code in (401, 404):
            # Retry once with fallback host and Basic auth
            fallback_url = f"{self.fallback_base_url}{endpoint}"
            headers["Authorization"] = self._basic_header(auth)
            resp = self._client.request(method, fallback_url, json=json_body, params=params, headers=headers)
            self.base_url = self.fallback_base_url

        if resp.status_code == 401:
            # token may be stale; clear and retry once
            self._auth = None
            if self.token_cache_path.exists():
                self.token_cache_path.unlink()
            auth = self._ensure_auth()
            headers["Authorization"] = (
                self._basic_header(auth) if self.base_url == self.fallback_base_url else self._token_header(auth)
            )
            url = f"{self.base_url}{endpoint}"
            resp = self._client.request(method, url, json=json_body, params=params, headers=headers)

        resp.raise_for_status()
        if not resp.text.strip():
            return {}
        return resp.json()

    # ---- Frames ----
    def get_frames(self) -> List[Dict[str, Any]]:
        # Some environments use /frames, others /frames/calendar
        try:
            result = self._request("GET", "/frames")
            if isinstance(result, dict):
                if "frames" in result:
                    return result["frames"]
                if "data" in result:
                    return result["data"]
            if isinstance(result, list):
                return result
        except httpx.HTTPStatusError:
            pass

        result = self._request("GET", "/frames/calendar")
        if isinstance(result, dict):
            return result.get("data", [])
        if isinstance(result, list):
            return result
        return []

    def resolve_frame_id(self, frame_id: Optional[str]) -> str:
        if frame_id:
            return frame_id
        if self.default_frame_id:
            return self.default_frame_id
        if not self._frames_cache:
            self._frames_cache = self.get_frames()
        if len(self._frames_cache) == 1:
            return str(self._frames_cache[0].get("id") or self._frames_cache[0].get("frame_id"))
        raise ValueError("Multiple frames found. Provide frame_id explicitly.")

    # ---- Lists ----
    def get_lists(self, frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
        frame_id = self.resolve_frame_id(frame_id)
        result = self._request("GET", f"/frames/{frame_id}/lists")
        if isinstance(result, dict):
            data = result.get("data")
            if isinstance(data, list):
                return data
            if "lists" in result:
                return result["lists"]
        if isinstance(result, list):
            return result
        return []

    def _cache_lists(self, frame_id: str) -> None:
        self._lists_cache = self.get_lists(frame_id)

    def resolve_list_id(self, frame_id: str, list_id: Optional[str], list_name: Optional[str]) -> str:
        if list_id:
            return list_id
        if not list_name:
            raise ValueError("list_id or list_name is required")
        if self._lists_cache is None:
            self._cache_lists(frame_id)
        for item in self._lists_cache or []:
            attrs = item.get("attributes", {})
            name = attrs.get("name") or attrs.get("label")
            if name == list_name:
                return str(item.get("id"))
        raise ValueError(f"List not found: {list_name}")

    def get_list_items(self, frame_id: Optional[str], list_id: Optional[str], list_name: Optional[str]) -> List[Dict[str, Any]]:
        frame_id = self.resolve_frame_id(frame_id)
        list_id = self.resolve_list_id(frame_id, list_id, list_name)

        # Prefer /items endpoint (documented). Fallback to /lists/{id} included.
        try:
            result = self._request("GET", f"/frames/{frame_id}/lists/{list_id}/items")
            if isinstance(result, dict) and isinstance(result.get("data"), list):
                return result["data"]
        except httpx.HTTPStatusError:
            pass

        result = self._request("GET", f"/frames/{frame_id}/lists/{list_id}")
        included = result.get("included", []) if isinstance(result, dict) else []
        items = [item for item in included if item.get("type") in ("list_item", "list_items")]
        return items

    def create_list_item(self, frame_id: Optional[str], list_id: Optional[str], list_name: Optional[str], name: str, checked: bool = False) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        list_id = self.resolve_list_id(frame_id, list_id, list_name)
        # Try documented JSON:API name/checked first
        payloads = [
            {
                "endpoint": f"/frames/{frame_id}/lists/{list_id}/items",
                "payload": {
                    "data": {
                        "type": "list_items",
                        "attributes": {"name": name, "checked": checked},
                    }
                },
            },
            {
                "endpoint": f"/frames/{frame_id}/lists/{list_id}/list_items",
                "payload": {
                    "data": {
                        "type": "list_item",
                        "attributes": {
                            "label": name,
                            "status": "completed" if checked else "pending",
                            "section": None,
                            "position": 1,
                        },
                    }
                },
            },
            {
                "endpoint": f"/frames/{frame_id}/lists/{list_id}/list_items",
                "payload": {
                    "list_item": {
                        "label": name,
                        "status": "completed" if checked else "pending",
                    }
                },
            },
        ]
        last_err: Optional[Exception] = None
        for attempt in payloads:
            try:
                return self._request("POST", attempt["endpoint"], attempt["payload"])
            except httpx.HTTPStatusError as exc:
                last_err = exc
                continue
        if last_err:
            raise last_err
        raise ValueError("Failed to create list item")

    def update_list_item(self, frame_id: Optional[str], list_id: Optional[str], list_name: Optional[str], item_id: str, name: Optional[str], checked: Optional[bool]) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        list_id = self.resolve_list_id(frame_id, list_id, list_name)
        attrs: Dict[str, Any] = {}
        if name is not None:
            attrs["name"] = name
        if checked is not None:
            attrs["checked"] = checked
        payload = {"data": {"type": "list_items", "id": item_id, "attributes": attrs}}
        try:
            return self._request("PATCH", f"/frames/{frame_id}/lists/{list_id}/items/{item_id}", payload)
        except httpx.HTTPStatusError:
            # Fallback to app-style payload
            status = None if checked is None else ("completed" if checked else "pending")
            body: Dict[str, Any] = {}
            if name is not None:
                body["label"] = name
            if status is not None:
                body["status"] = status
            return self._request("PUT", f"/frames/{frame_id}/lists/{list_id}/list_items/{item_id}", body)

    def delete_list_items(self, frame_id: Optional[str], list_id: Optional[str], list_name: Optional[str], ids: List[str]) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        list_id = self.resolve_list_id(frame_id, list_id, list_name)
        payload = {"ids": ids}
        return self._request("DELETE", f"/frames/{frame_id}/lists/{list_id}/list_items/bulk_destroy", payload)

    # ---- Categories (Family Member Profiles) ----
    def get_categories(self, frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
        frame_id = self.resolve_frame_id(frame_id)
        result = self._request("GET", f"/frames/{frame_id}/categories")
        if isinstance(result, dict):
            data = result.get("data")
            if isinstance(data, list):
                return data
        if isinstance(result, list):
            return result
        return []

    def resolve_category_id(
        self, frame_id: str, category_id: Optional[str], profile_name: Optional[str]
    ) -> Optional[str]:
        if category_id:
            return category_id
        if not profile_name:
            return None
        if self._categories_cache is None:
            self._categories_cache = self.get_categories(frame_id)
        target = profile_name.lower()
        for item in self._categories_cache:
            label = (item.get("attributes", {}).get("label") or "").lower()
            if label == target:
                return str(item.get("id"))
        raise ValueError(f"Profile not found: {profile_name}")

    # ---- Task Box Items ----
    def get_task_box_items(self, frame_id: Optional[str] = None) -> List[Dict[str, Any]]:
        frame_id = self.resolve_frame_id(frame_id)
        result = self._request("GET", f"/frames/{frame_id}/task_box/items")
        if isinstance(result, dict):
            data = result.get("data")
            if isinstance(data, list):
                return data
        if isinstance(result, list):
            return result
        return []

    def create_task_box_item(
        self,
        frame_id: Optional[str] = None,
        summary: str = "",
        emoji_icon: Optional[str] = None,
        routine: bool = False,
        reward_points: Optional[int] = None,
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        payload: Dict[str, Any] = {
            "summary": summary,
            "routine": routine,
        }
        if emoji_icon is not None:
            payload["emoji_icon"] = emoji_icon
        if reward_points is not None:
            payload["reward_points"] = reward_points
        return self._request("POST", f"/frames/{frame_id}/task_box/items", payload)

    def update_task_box_item(
        self,
        frame_id: Optional[str] = None,
        item_id: str = "",
        summary: Optional[str] = None,
        emoji_icon: Optional[str] = None,
        routine: Optional[bool] = None,
        reward_points: Optional[int] = None,
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        payload: Dict[str, Any] = {}
        if summary is not None:
            payload["summary"] = summary
        if emoji_icon is not None:
            payload["emoji_icon"] = emoji_icon
        if routine is not None:
            payload["routine"] = routine
        if reward_points is not None:
            payload["reward_points"] = reward_points
        return self._request("PATCH", f"/frames/{frame_id}/task_box/items/{item_id}", payload)

    def delete_task_box_item(self, frame_id: Optional[str] = None, item_id: str = "") -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        return self._request("DELETE", f"/frames/{frame_id}/task_box/items/{item_id}")

    # ---- Chores ----
    def get_chores(
        self,
        frame_id: Optional[str] = None,
        after: Optional[str] = None,
        before: Optional[str] = None,
        include_late: Optional[bool] = None,
        filter_profile: Optional[str] = None,
    ) -> Any:
        frame_id = self.resolve_frame_id(frame_id)
        params: Dict[str, Any] = {}
        if after is not None:
            params["after"] = after
        if before is not None:
            params["before"] = before
        if include_late is not None:
            params["include_late"] = str(include_late).lower()
        category_id = self.resolve_category_id(frame_id, None, filter_profile)
        if category_id:
            params["category_id"] = category_id
        return self._request("GET", f"/frames/{frame_id}/chores", params=params)

    def create_chore(
        self,
        frame_id: Optional[str] = None,
        summary: str = "",
        start: str = "",
        routine: bool = False,
        start_time: Optional[str] = None,
        recurrence_set: Optional[str] = None,
        emoji_icon: Optional[str] = None,
        recurring_until: Optional[str] = None,
        reward_points: Optional[int] = None,
        category_ids: Optional[List[str]] = None,
        category_id: Optional[str] = None,
        profile_name: Optional[str] = None,
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        resolved_cat = self.resolve_category_id(frame_id, category_id, profile_name)
        cat_ids = category_ids or ([resolved_cat] if resolved_cat else [])
        payload: Dict[str, Any] = {
            "summary": summary,
            "start": start,
            "routine": routine,
            "category_ids": cat_ids,
        }
        if start_time is not None:
            payload["start_time"] = start_time
        if recurrence_set is not None:
            payload["recurrence_set"] = recurrence_set
        if emoji_icon is not None:
            payload["emoji_icon"] = emoji_icon
        if recurring_until is not None:
            payload["recurring_until"] = recurring_until
        if reward_points is not None:
            payload["reward_points"] = reward_points
        return self._request("POST", f"/frames/{frame_id}/chores/create_multiple", payload)

    def update_chore(
        self,
        frame_id: Optional[str] = None,
        chore_id: str = "",
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
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        resolved_cat = self.resolve_category_id(frame_id, category_id, profile_name)
        payload: Dict[str, Any] = {}
        if summary is not None:
            payload["summary"] = summary
        if emoji_icon is not None:
            payload["emoji_icon"] = emoji_icon
        if reward_points is not None:
            payload["reward_points"] = reward_points
        if start is not None:
            payload["start"] = start
        if start_time is not None:
            payload["start_time"] = start_time
        if recurrence_set is not None:
            payload["recurrence_set"] = recurrence_set
        if recurring_until is not None:
            payload["recurring_until"] = recurring_until
        if routine is not None:
            payload["routine"] = routine
        if resolved_cat is not None:
            payload["category_id"] = resolved_cat
        if up_for_grabs is not None:
            payload["up_for_grabs"] = up_for_grabs
        return self._request("PUT", f"/frames/{frame_id}/chores/{chore_id}", payload)

    def complete_chore(self, frame_id: Optional[str] = None, chore_id: str = "") -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        return self._request("PUT", f"/frames/{frame_id}/chores/{chore_id}", {"status": "complete"})

    def delete_chore(
        self, frame_id: Optional[str] = None, chore_id: str = "", apply_to: str = "one"
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        return self._request("DELETE", f"/frames/{frame_id}/chores/{chore_id}", params={"apply_to": apply_to})

    # ---- Reward Points ----
    def get_reward_points(self, frame_id: Optional[str] = None) -> Any:
        frame_id = self.resolve_frame_id(frame_id)
        result = self._request("GET", f"/frames/{frame_id}/reward_points")
        if isinstance(result, dict):
            data = result.get("data")
            if isinstance(data, list):
                return data
        if isinstance(result, list):
            return result
        return result

    # ---- Meals ----
    def get_meal_categories(self, frame_id: Optional[str]) -> List[Dict[str, Any]]:
        frame_id = self.resolve_frame_id(frame_id)
        result = self._request("GET", f"/frames/{frame_id}/meals/categories")
        if isinstance(result, dict) and isinstance(result.get("data"), list):
            return result["data"]
        if isinstance(result, list):
            return result
        return []

    def resolve_meal_category_id(self, frame_id: str, meal_category_id: Optional[str], meal_type: Optional[str]) -> str:
        if meal_category_id:
            return meal_category_id
        if not meal_type:
            raise ValueError("meal_category_id or meal_type is required")
        categories = self.get_meal_categories(frame_id)
        target = meal_type.lower()
        for item in categories:
            if item.get("type") != "meal_category":
                continue
            label = (item.get("attributes", {}).get("label") or "").lower()
            if label == target:
                return str(item.get("id"))
        raise ValueError(f"Meal category not found for type: {meal_type}")

    def get_meal_recipes(self, frame_id: Optional[str]) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        return self._request("GET", f"/frames/{frame_id}/meals/recipes", params={"include": "meal_category"})

    def get_meal_recipe(self, frame_id: Optional[str], recipe_id: str) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        return self._request("GET", f"/frames/{frame_id}/meals/recipes/{recipe_id}", params={"include": "meal_category"})

    def create_meal_recipe(self, frame_id: Optional[str], summary: str, description: Optional[str], meal_category_id: Optional[str], meal_type: Optional[str]) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        category_id = self.resolve_meal_category_id(frame_id, meal_category_id, meal_type)
        payload = {
            "meal_category_id": category_id,
            "summary": summary,
            "description": description,
        }
        return self._request("POST", f"/frames/{frame_id}/meals/recipes", payload, params={"include": "meal_category"})

    def update_meal_recipe(self, frame_id: Optional[str], recipe_id: str, summary: str, description: Optional[str], meal_category_id: Optional[str], meal_type: Optional[str]) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        category_id = self.resolve_meal_category_id(frame_id, meal_category_id, meal_type)
        payload = {
            "meal_category_id": category_id,
            "summary": summary,
            "description": description,
        }
        return self._request("PATCH", f"/frames/{frame_id}/meals/recipes/{recipe_id}", payload, params={"include": "meal_category"})

    def add_recipe_to_grocery_list(self, frame_id: Optional[str], recipe_id: str) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        return self._request("POST", f"/frames/{frame_id}/meals/recipes/{recipe_id}/add_to_grocery_list")

    def get_meal_sittings(self, frame_id: Optional[str], date_min: date, date_max: date) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        params = {
            "date_min": date_min.isoformat(),
            "date_max": date_max.isoformat(),
            "include": "meal_category,meal_recipe",
        }
        return self._request("GET", f"/frames/{frame_id}/meals/sittings", params=params)

    def create_meal_sitting(
        self,
        frame_id: Optional[str],
        date_value: date,
        meal_category_id: Optional[str],
        meal_type: Optional[str],
        meal_recipe_id: Optional[str],
        summary: Optional[str],
        note: Optional[str],
        description: Optional[str],
        add_to_grocery_list: bool = False,
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        category_id = self.resolve_meal_category_id(frame_id, meal_category_id, meal_type)
        payload = {
            "meal_recipe_id": meal_recipe_id,
            "meal_category_id": category_id,
            "add_to_grocery_list": add_to_grocery_list,
            "date": date_value.isoformat(),
            "note": note,
            "rrule": None,
            "description": description,
        }
        # summary must be blank when meal_recipe_id is set
        if meal_recipe_id:
            payload["summary"] = None
        else:
            payload["summary"] = summary
        params = {
            "date_min": date_value.isoformat(),
            "date_max": date_value.isoformat(),
            "include": "meal_category,meal_recipe",
        }
        return self._request("POST", f"/frames/{frame_id}/meals/sittings", payload, params=params)

    def update_meal_sitting(
        self,
        frame_id: Optional[str],
        sitting_id: str,
        date_value: date,
        meal_category_id: Optional[str],
        meal_type: Optional[str],
        meal_recipe_id: Optional[str],
        summary: Optional[str],
        note: Optional[str],
        description: Optional[str],
        add_to_grocery_list: bool = False,
    ) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        category_id = self.resolve_meal_category_id(frame_id, meal_category_id, meal_type)
        payload = {
            "meal_recipe_id": meal_recipe_id,
            "meal_category_id": category_id,
            "add_to_grocery_list": add_to_grocery_list,
            "date": date_value.isoformat(),
            "note": note,
            "rrule": None,
            "description": description,
        }
        if meal_recipe_id:
            payload["summary"] = None
        else:
            payload["summary"] = summary

        # Try instance-based update first
        params = {
            "date_min": date_value.isoformat(),
            "date_max": date_value.isoformat(),
            "include": "meal_category,meal_recipe",
        }
        try:
            return self._request(
                "PATCH",
                f"/frames/{frame_id}/meals/sittings/{sitting_id}/instances/{date_value.isoformat()}",
                payload,
                params=params,
            )
        except httpx.HTTPStatusError:
            return self._request("PATCH", f"/frames/{frame_id}/meals/sittings/{sitting_id}", payload)

    def delete_meal_sitting(self, frame_id: Optional[str], sitting_id: str, date_value: Optional[date]) -> Dict[str, Any]:
        frame_id = self.resolve_frame_id(frame_id)
        if date_value:
            params = {
                "date_min": date_value.isoformat(),
                "date_max": date_value.isoformat(),
                "include": "meal_category,meal_recipe",
            }
            return self._request(
                "DELETE",
                f"/frames/{frame_id}/meals/sittings/{sitting_id}/instances/{date_value.isoformat()}",
                params=params,
            )
        return self._request("DELETE", f"/frames/{frame_id}/meals/sittings/{sitting_id}")
