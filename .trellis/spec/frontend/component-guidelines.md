# Frontend Component Guidelines

## Composition

Use named function components for feature/product components and typed props
near the component. `features/dashboard/dashboard-panel.tsx` is the compact
reference: a semantic `<section>`, an `aria-labelledby` relationship, optional
actions, and `cn` for controlled class extension.

Low-level primitives under `components/ui` follow the existing Radix/shadcn
shape: `React.forwardRef`, native element props, `class-variance-authority`
variants, `asChild` where composition requires it, and `cn` for class merging.
Do not copy a primitive into a feature to change one class; extend its supported
variant or pass `className` when that remains coherent.

## Styling and Responsive Behavior

- Use Tailwind utilities and existing semantic tokens (`bg-card`,
  `text-muted-foreground`, `text-destructive`) rather than inline colors.
- Preserve dark-mode variants and responsive prefixes already used by the
  surrounding page.
- Use `cn` from `shared/lib/cn.ts` for conditional classes.
- Keep desktop and narrow layouts usable; large tables should preserve their
  existing scroll/virtualization behavior rather than compressing every cell.

## Interaction Feedback

Every asynchronous interaction must expose its state. Reuse
`shared/components/data-state.tsx` for loading, empty, and retryable error
views; disable mutation controls while pending and show `Spinner` where the
action originated. Mutation success/error feedback uses `sonner`, as in
`features/settings/use-settings.ts`.

Do not optimistically claim success unless rollback is implemented. Existing
mutations normally update or invalidate the TanStack Query cache only after a
successful response.

## Accessibility and Content

- Prefer semantic buttons, links, headings, labels, sections, and tables.
- Give icon-only controls an accessible name; keep focus-visible behavior from
  the shared primitives.
- Connect dialogs/forms to labels and descriptions through the Radix/form
  APIs; do not replace buttons with clickable `<div>` elements.
- Use i18next for visible copy. Preserve `title` text for truncated values when
  the full value otherwise becomes inaccessible.
- External links opened in a new tab use `rel="noreferrer"`, as in
  `features/system/version-update.tsx`.

Avoid monolithic page additions when a panel, pure formatter, or API contract
has independent ownership. Split on real reuse or testability, not arbitrary
line counts.

