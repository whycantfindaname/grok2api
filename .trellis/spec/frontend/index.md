# Frontend Development Guidelines

The management UI is the React 19 + TypeScript application under `frontend/`.
It uses Vite, React Router, TanStack Query, React Hook Form + Zod, Radix-based
UI primitives, Tailwind CSS, and i18next. The package manager is pinned to
`pnpm@11.5.2` in `frontend/package.json`.

## Pre-Development Checklist

- Find the owning feature and inspect its `*-api.ts`, page/component, and any
  nearby test before adding a new abstraction.
- Check the backend Admin API envelope and add a runtime decoder for every new
  response shape.
- Classify state as server, form, local interaction, auth context, runtime
  config, or URL/route state before choosing where it lives.
- Reuse `components/ui` and `shared/components` before creating feature-local
  controls; add visible loading, empty, error, and disabled states.
- Put user-facing strings through i18next rather than inline literals.

## Guides

| Guide | Use it for |
| --- | --- |
| [Directory Structure](./directory-structure.md) | Feature, entity, shared, and app ownership |
| [Component Guidelines](./component-guidelines.md) | Composition, styling, feedback, and accessibility |
| [Hook Guidelines](./hook-guidelines.md) | Query/mutation and reusable hook patterns |
| [State Management](./state-management.md) | Local, server, auth, form, route, and runtime state |
| [Type Safety](./type-safety.md) | Strict TypeScript and API boundary validation |
| [Quality Guidelines](./quality-guidelines.md) | Lint, tests, build, and review expectations |

## Quality Check

Run from `frontend/`:

```bash
pnpm test
pnpm lint
pnpm build
```

`pnpm build` includes `tsc -b` before Vite. Do not install or regenerate the
lockfile merely to run checks; use the existing environment and report a
missing dependency/toolchain if verification cannot start.

