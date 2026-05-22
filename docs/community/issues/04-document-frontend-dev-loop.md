# Document the frontend development loop

Labels: `good first issue`, `documentation`, `frontend`, `help wanted`

## Summary

The repo includes a Vite frontend, but the contribution docs mostly focus on Docker and backend commands. Add a short frontend contributor section with install, dev, and build steps.

## Why this matters

New contributors interested in UI polish should not have to reverse-engineer the frontend workflow from `package.json`.

## Suggested files

- `CONTRIBUTING.md`
- `README.md`
- `frontend/package.json`

## Acceptance criteria

- Docs explain how to install frontend dependencies.
- Docs explain how to run the frontend in dev mode.
- Docs mention how the frontend connects to the local API stack.
- Docs still keep the quickstart concise.

## Nice to have

- Mention a common troubleshooting step for API URL or CORS confusion if needed.
