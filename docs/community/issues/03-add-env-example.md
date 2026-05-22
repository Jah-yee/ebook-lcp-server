# Add a root `.env.example` for local development

Labels: `good first issue`, `documentation`, `developer experience`, `help wanted`

## Summary

The README tells contributors to start from `.env.example`, but the repo does not currently include that file at the root. Add a safe local-development example file and make sure the docs point to the right source of truth.

## Why this matters

This is one of the first things a new contributor sees, so a missing starter file creates friction early.

## Suggested files

- `.env.example`
- `README.md`
- `CONTRIBUTING.md`

## Acceptance criteria

- A root `.env.example` exists with sensible local placeholders.
- Secrets are clearly placeholders, not real values.
- README and contributing docs reference the same setup flow.

## Nice to have

- Add short comments grouping variables by database, auth, storage, and Readium integration.
