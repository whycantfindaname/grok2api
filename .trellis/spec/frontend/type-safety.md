# Frontend Type Safety

## Compiler and Lint Contract

`frontend/tsconfig.app.json` enables strict mode, unused checks, isolated
modules, and case-sensitive paths. ESLint rejects explicit `any`. Preserve the
`@/` alias and use `import type` for type-only imports where appropriate.

Types normally live beside their owner:

- feature request/response DTOs in the feature `*-api.ts`;
- shared entity contracts in `entities/<entity>/types.ts`;
- inferred form types beside the Zod schema;
- ambient browser/runtime declarations only in `src/types/`.

## Runtime API Validation

Static types do not validate JSON. Every Admin API response passed to
`apiRequest` must provide an `ApiDecoder<T>`. Build decoders with
`createObjectDecoder`, `hasShape`, `isArrayOf`, `isOptional`, and `isOneOf`
from `shared/api/decoder.ts`. The single assertion inside
`createValidatedDecoder` is allowed because the shape was checked immediately
before it.

Keep external payloads as `unknown` until decoded. Do not write
`response.json() as SomeDTO`, add `any`, or scatter unchecked
`Record<string, unknown>` casts through components. If a protocol is dynamic,
centralize its guards/projections at the API boundary as
`shared/api/client.ts` does for error envelopes and SSE events.

## Forms and Constants

Infer form types with `z.infer<typeof schema>` and keep cross-field rules in
Zod refinements. Convert between form-friendly units and DTO units explicitly;
do not cast the form directly to a wire type. Use `as const` for stable query
keys and closed literal collections, as in `features/system/version-update.tsx`.

Verification is `pnpm lint`, `pnpm test`, and `pnpm build`; the build is the
authoritative full app type-check because it runs `tsc -b` before Vite.

