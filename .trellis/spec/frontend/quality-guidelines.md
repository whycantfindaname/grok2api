# Frontend Quality Guidelines

## Required Coverage

Add focused tests for pure behavior with meaningful branching, especially
formatting, protocol projection, state transitions, and missing/partial data.
Tests use Node's built-in `node:test` and `node:assert/strict`; see
`features/audits/audit-usage.test.ts` and `shared/lib/format.test.ts`.

For component/API changes, verify more than a happy-path render:

- response decoder rejects malformed external data;
- loading, empty, error, success, and mutation-pending states remain usable;
- query keys include relevant filters/cursors and cache updates target the
  correct resource;
- forms validate ranges and cross-field invariants before sending DTOs;
- keyboard/focus semantics and accessible names remain intact;
- narrow and desktop layouts preserve readable controls and tables.

## Anti-Patterns

- No explicit `any`, unchecked JSON assertions, or duplicate local API clients.
- No feature-local token refresh, Admin envelope parsing, or SSE framing.
- No inline user-visible strings when the surrounding feature uses i18next.
- No copied UI primitives or one-off design tokens when a shared primitive or
  semantic token already exists.
- No effect-driven fetch that duplicates TanStack Query.
- No generated `dist`, dependency directories, package caches, or build-info
  files in source changes.

`frontend/eslint.config.js` intentionally ignores `src/components/ui` because
those primitives follow their generated/Radix style; this is not permission to
ignore type errors or to place feature code there.

## Verification

```bash
cd frontend
pnpm test
pnpm lint
pnpm build
```

The CI currently runs frozen install, lint, and build in
`.github/workflows/ghcr-image.yml`; local `pnpm test` is also required for
changes covered by the repository's Node tests. A successful Vite build is
source/build evidence, not proof that the Go service or a downstream consumer
is deployed and live.
