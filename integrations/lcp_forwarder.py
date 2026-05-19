#!/usr/bin/env python3
"""Forward EPUB/PDF files from self-hosted library tools into the LCP server."""

from __future__ import annotations

import argparse
import base64
import json
import os
import pathlib
import sys
import urllib.error
import urllib.request


def request_json(url: str, payload: dict, token: str | None = None) -> dict:
    body = json.dumps(payload).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/json")
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    try:
        with urllib.request.urlopen(req, timeout=30) as response:
            return json.loads(response.read().decode())
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode()
        raise SystemExit(f"{exc.code} {exc.reason}: {detail}") from exc


def login(base_url: str, username: str, password: str, two_factor: str) -> str:
    result = request_json(
        f"{base_url}/api/v1/auth/login",
        {"username": username, "password": password, "twoFactor": two_factor},
    )
    return result["token"]


def upload(base_url: str, token: str, path: pathlib.Path, title: str) -> dict:
    encoded = base64.b64encode(path.read_bytes()).decode()
    result = request_json(
        f"{base_url}/graphql",
        {
            "query": (
                "mutation UploadPublication($title: String!, $file: Upload!) "
                "{ uploadPublication(title: $title, file: $file) { id title } }"
            ),
            "variables": {"title": title, "file": encoded},
        },
        token,
    )
    if result.get("errors"):
        raise SystemExit(result["errors"][0]["message"])
    return result["data"]["uploadPublication"]


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("path", type=pathlib.Path, help="EPUB or PDF file to forward")
    parser.add_argument("--title", help="Catalog title override")
    parser.add_argument("--base-url", default=os.getenv("LCP_BASE_URL", "http://localhost:8080"))
    parser.add_argument("--username", default=os.getenv("LCP_USERNAME", "publisher"))
    parser.add_argument("--password", default=os.getenv("LCP_PASSWORD", "publisher"))
    parser.add_argument("--two-factor", default=os.getenv("LCP_2FA_CODE", ""))
    args = parser.parse_args()
    if not args.path.is_file():
        raise SystemExit(f"file not found: {args.path}")
    token = login(args.base_url.rstrip("/"), args.username, args.password, args.two_factor)
    print(json.dumps(upload(args.base_url.rstrip("/"), token, args.path, args.title or args.path.stem)))
    return 0


if __name__ == "__main__":
    sys.exit(main())
