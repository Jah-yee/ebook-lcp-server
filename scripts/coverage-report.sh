#!/usr/bin/env sh

set -eu

out_dir="${1:-reports/coverage}"

mkdir -p "$out_dir"

profile="$out_dir/coverage.out"
filtered_profile="$out_dir/coverage.filtered.out"
html_report="$out_dir/index.html"
summary_txt="$out_dir/summary.txt"
summary_md="$out_dir/README.md"
summary_json="$out_dir/coverage-summary.json"
badge_svg="$out_dir/badge.svg"

go test ./... -coverprofile="$profile" -coverpkg=./... >"$out_dir/go-test.log"

exclude_pattern='/(cmd/lcpctl/main.go|cmd/server/main.go|internal/adapter/graphql/generated.go|internal/adapter/repository/lcp/postgres.go|internal/adapter/repository/lcp/license_repository.go|internal/adapter/repository/lcp/publication_repository.go|internal/adapter/rest/lcp.go|internal/lcp/encrypt/encrypt.go|internal/lcp/license/license.go|internal/storage/publication.go|internal/usecase/lcp/license/usecase.go):'
{
  sed -n '1p' "$profile"
  sed -n '2,$p' "$profile" | grep -Ev "$exclude_pattern" || true
} >"$filtered_profile"

go tool cover -html="$filtered_profile" -o "$html_report"
go tool cover -func="$filtered_profile" >"$summary_txt"

total="$(awk '/^total:/ {gsub("%", "", $3); print $3}' "$summary_txt")"

if [ -z "$total" ]; then
  echo "failed to calculate total coverage" >&2
  exit 1
fi

color="e05d44"
awk "BEGIN { exit !($total >= 40) }" || color="dfb317"
awk "BEGIN { exit !($total >= 60) }" || color="$color"
if awk "BEGIN { exit !($total >= 60) }"; then
  color="4c1"
fi

generated_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

cat >"$summary_json" <<EOF
{
  "schemaVersion": 1,
  "label": "coverage",
  "message": "${total}%",
  "color": "${color}",
  "total": ${total},
  "generatedAt": "${generated_at}"
}
EOF

cat >"$badge_svg" <<EOF
<svg xmlns="http://www.w3.org/2000/svg" width="104" height="20" role="img" aria-label="coverage: ${total}%">
  <title>coverage: ${total}%</title>
  <linearGradient id="s" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <clipPath id="r">
    <rect width="104" height="20" rx="3" fill="#fff"/>
  </clipPath>
  <g clip-path="url(#r)">
    <rect width="63" height="20" fill="#555"/>
    <rect x="63" width="41" height="20" fill="#${color}"/>
    <rect width="104" height="20" fill="url(#s)"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11">
    <text x="32.5" y="15" fill="#010101" fill-opacity=".3">coverage</text>
    <text x="32.5" y="14">coverage</text>
    <text x="82.5" y="15" fill="#010101" fill-opacity=".3">${total}%</text>
    <text x="82.5" y="14">${total}%</text>
  </g>
</svg>
EOF

{
  echo "# Coverage Report"
  echo
  echo "[![Coverage](badge.svg)](index.html)"
  echo
  echo "- Total coverage: \`${total}%\`"
  echo "- Generated at: \`${generated_at}\`"
  echo "- Raw profile: [coverage.out](coverage.out)"
  echo "- Filtered profile: [coverage.filtered.out](coverage.filtered.out)"
  echo "- HTML report: [index.html](index.html)"
  echo "- Exclusions: generated GraphQL code, main entrypoints, low-level storage/database adapters, the external encrypter wrapper, the LCP-core HTTP client, and the license orchestration use case"
  echo
  echo "## Package Summary"
  echo
  echo '```text'
  cat "$summary_txt"
  echo '```'
} >"$summary_md"
