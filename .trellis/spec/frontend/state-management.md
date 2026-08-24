# Frontend State Management

This project does not use a general-purpose global store. Choose the narrowest
existing owner for each kind of state.

| State kind | Owner and pattern |
| --- | --- |
| Backend resources | TanStack Query cache via feature/entity API modules |
| Mutations | `useMutation`, then `setQueryData` or targeted invalidation |
| Forms | React Hook Form with Zod resolver; DTO/form conversion stays explicit |
| Page interaction | Local `useState` for filters, selections, dialogs, cursors, and hover state |
| Authentication | `shared/auth/auth-context.tsx` and the token lifecycle in `shared/api/client.ts` |
| Route state | React Router route definitions/params in `app/router.tsx` |
| Deployment config | Immutable `shared/config/runtime-config.ts` loaded from `public/runtime-config.js` |
| Theme/i18n | Existing providers; never duplicate in feature state |

## Server and Form State

Server snapshots remain in Query cache. Do not copy a list response into a
global context. Use query keys that include every filter/cursor affecting the
response so cached results cannot cross scopes.

Forms may initialize/reset from a successful query. Keep conversion between
wire DTOs and editable values in a pure model module, as
`features/settings/settings-model.ts` does with `toSettingsForm` and
`toSettingsDTO`. Use the server revision for conflict-sensitive updates rather
than treating the form snapshot as authoritative.

## Local State

Keep transient UI state beside the page/component that owns it. When state
contains a `Set`, array, or object, replace it immutably. Derive rows/options
from query data instead of storing duplicate copies that can drift.

Do not store access tokens in feature state or persistent browser storage; the
API client keeps the access token in memory and refreshes through the secure
session flow. Do not mix service runtime state (readiness, jobs, request status)
with frontend control state: the backend response remains authoritative.

