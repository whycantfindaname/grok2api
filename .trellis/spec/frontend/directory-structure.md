# Frontend Directory Structure

## Ownership Map

| Path | Ownership |
| --- | --- |
| `frontend/src/app/` | Root providers, auth boundaries, route table, lazy pages, and app shell |
| `frontend/src/features/<feature>/` | Page behavior, feature API module, forms, panels, and feature-local types |
| `frontend/src/entities/<entity>/` | API/types reused by more than one feature |
| `frontend/src/components/ui/` | Low-level Radix/shadcn-style primitives and variants |
| `frontend/src/shared/components/` | Product-level reusable components such as data states and pagination |
| `frontend/src/shared/api/` | Authenticated client, SSE handling, errors, and runtime decoders |
| `frontend/src/shared/auth/` | Admin session context and access-token lifecycle |
| `frontend/src/shared/hooks/` | Cross-feature hooks with no product feature ownership |
| `frontend/src/shared/lib/` | Pure formatting, period, sort, and class-name helpers |
| `frontend/src/shared/config/` | Deployment-injected runtime config reader |
| `frontend/src/types/` | Ambient/global TypeScript declarations only |

`frontend/src/main.tsx` mounts the app. `app/providers.tsx` owns the Query,
theme, auth, tooltip, and toast providers. `app/router.tsx` owns routes and auth
boundaries; lazy imports stay centralized in `app/deferred-pages.tsx`.

## Placement Rules

- Keep a feature's wire calls and decoders in its `*-api.ts`, as in
  `features/dashboard/dashboard-api.ts` and `features/settings/settings-api.ts`.
- Promote a type/API to `entities/` only when multiple features consume the
  same entity contract. Do not create a parallel global `types.ts` dump.
- A reusable visual primitive belongs in `components/ui`; a product-aware
  reusable composition belongs in `shared/components`.
- Keep pure transforms outside a large page component so they can be tested.
  `features/audits/audit-usage.ts` plus `audit-usage.test.ts` is the reference.

Files use kebab-case; exported React components use PascalCase; hooks use
`useX`; the `@/` alias points to `frontend/src/`. Avoid deep relative imports
that cross feature boundaries.

Generated `frontend/dist`, `node_modules`, `.pnpm-store`, and TypeScript/Vite
caches are not source and must not be edited or committed.

