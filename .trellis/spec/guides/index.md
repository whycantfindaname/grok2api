# Cross-Cutting Guides

These guides cover decisions that span the Go backend, React frontend, provider
adapters, generated API documentation, and runtime boundaries.

| Guide | Read when |
| --- | --- |
| [Code Reuse and Ownership](./code-reuse-thinking-guide.md) | A new helper, decoder, provider behavior, or UI primitive resembles existing code |
| [Cross-Layer Contract Changes](./cross-layer-thinking-guide.md) | A route, DTO, persisted field, provider capability, config field, or runtime status crosses ownership boundaries |

Before editing a value or contract, use `rg` to find all readers, writers,
tests, docs, and generated representations. Reuse is correct only when the
candidate has the same owner and semantics; similar Build, Web, and Console
payloads are not automatically one abstraction.

## Quality Check

- Verify source, generated artifacts, and tests independently.
- Keep runtime secrets/data outside Git.
- Do not describe build, health, readiness, or a configured consumer as live
  request evidence; `STRUCTURE.md` defines these proof boundaries.

