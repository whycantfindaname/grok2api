# Frontend Hook Guidelines

## Server Data Hooks

Use TanStack Query for Admin API server state. Queries need a stable array key
and a query function from the owning feature/entity API module. Mutations
update the exact cache entry or invalidate the affected key after success.
`features/system/version-update.tsx` demonstrates shared query keys and cache
replacement; `features/settings/use-settings.ts` demonstrates mutation,
invalidating related system info, resetting a form, and toast feedback.

Keep global defaults in `app/providers.tsx`: queries retry once, mutations do
not retry automatically, stale time is 15 seconds, and focus does not refetch.
Override these only for a concrete endpoint behavior.

## Custom Hooks

- Name hooks `useX` and keep them in the owning feature unless they are truly
  cross-feature.
- Return the underlying query/mutation objects when callers need status or
  errors; do not hide pending/error states behind booleans that discard detail.
- Effects must clean up timers, subscriptions, streams, and abort controllers.
  `shared/hooks/use-debounced-value.ts` is the timer cleanup reference.
- Derive display data with plain functions or `useMemo`; do not mirror query
  data into local state unless the user is editing a snapshot (for example, a
  form reset from settings data).

## Authentication and Streams

Use `shared/api/client.ts` for authenticated requests and Admin SSE. It owns
access-token refresh serialization, browser locking, session invalidation,
envelope decoding, SSE chunk boundaries, buffer limits, and inactivity
timeouts. Feature hooks must not duplicate token refresh or parse SSE chunks
with ad hoc `split` logic.

Avoid effects that perform the same fetch as a query, unstable object query
keys, and mutation retries for non-idempotent operations. When a request may be
canceled during search/filter changes, propagate an `AbortSignal` and treat an
abort separately from a user-visible failure, following the account page's
abort handling.

