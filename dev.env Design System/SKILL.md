---
name: dev-env-design
description: Use this skill to generate well-branded interfaces and assets for dev.env — a Go CLI that manages shared Docker dev services. Use this for production work, mocks, slide decks, README graphics, social cards, or throwaway prototypes. Contains brand guidelines, colors, type, fonts, logos, service iconography, and a documentation-site UI kit.
user-invocable: true
---

# dev.env design skill

Read `README.md` in this directory first — it covers brand context, content fundamentals, visual foundations, and iconography. Then explore:

- `colors_and_type.css` — token CSS; drop into any HTML and you have the full system (colors, type, spacing, radii, shadows, motion).
- `fonts/README.md` — typeface info (IBM Plex via Google Fonts) + a flagged substitution note.
- `assets/logos/` — wordmark, logomark, full lockup. Light + dark variants.
- `preview/` — small Design-System-tab cards that double as living-styleguide references.
- `ui_kits/docs/` — production-quality documentation-site UI kit. `index.html`, `docs.html`, `reference.html`, `changelog.html` plus the JSX components. Mirror this structure when building new pages.
- `SKILL.md` — this file.

## When to use

- **Visual artifacts** (slides, README diagrams, social cards, mocks, throwaway prototypes): copy assets out, write static HTML files that load `colors_and_type.css`, and let the user view them.
- **Production code**: read the rules in `README.md`, lift tokens from `colors_and_type.css` into the target project's variable system, and reuse JSX components from `ui_kits/docs/` as a reference (they're cosmetic React, not production-grade — copy the markup and styles, not the file structure).

## Working principles

- The brand is **typographic + monochromatic + one electric accent** (lime `#C6F432`). Resist adding more colors.
- **30–40% of visible text on marketing surfaces should be monospace** (IBM Plex Mono). That's the brand's signal.
- Use the CLI's ASCII glyphs (`✓` `~` `!` `✗`) — never substitute emoji.
- The lime is loud on purpose. ~5–8% coverage per surface. Never as a flood fill.
- No gradients. No glassmorphism. No purple. No emoji.

## If the user invokes this skill bare

Ask what they want to build — a docs page, a release-notes graphic, a slide deck, a marketing surface. Ask which target (HTML file for review, or production code in another codebase). Then act as an expert designer.
