#!/bin/sh
set -eu

repo_url="${LCP_REPO_URL:-https://github.com/amirHdev/ebook-lcp-server.git}"
repo_dir="${LCP_QUICKSTART_DIR:-ebook-lcp-server}"

need() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1" >&2
    exit 1
  }
}

need docker
need curl
need jq

if [ ! -d "$repo_dir/.git" ]; then
  need git
  git clone "$repo_url" "$repo_dir"
fi

cd "$repo_dir"
docker compose up --build -d

echo "waiting for API readiness..."
i=0
until curl -fsS http://localhost:8080/readyz >/dev/null 2>&1; do
  i=$((i + 1))
  if [ "$i" -ge 90 ]; then
    echo "server did not become ready within 90 seconds" >&2
    exit 1
  fi
  sleep 1
done

sh scripts/demo-local.sh

cat <<'EOF'

Demo stack is ready:
  API:        http://localhost:8080
  Admin UI:   http://localhost:5173
  Swagger UI: http://localhost:8081

Admin login:
  username: admin
  password: admin
  2FA code: 123456
EOF
