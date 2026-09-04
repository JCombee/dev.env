# dev.env Design System

A complete brand and design system for **dev.env** — a Go CLI tool that manages shared development services (databases, caches, search engines, queues) across all your projects via Docker. One global pool of containers, automatic per-project isolation, single-binary install.

This system was designed from scratch (no prior brand existed) for the **documentation site** as the primary surface, with foundations broad enough to apply to README graphics, social cards, slide decks, and future product UI.

## Product context

dev.env is a single Go binary that solves an annoying problem: every project ships its own `docker-compose.yml`, which means duplicate MySQL containers, duplicate Redis, wasted RAM, slow startup. dev.env runs **one shared pool** for all your projects and provisions isolated databases/users/indices per project automatically. Services that can't be safely shared can be marked `dedicated: true`. Any Docker image works; supported services (MySQL, Postgres, Mongo, Redis, Elasticsearch, Meilisearch, RabbitMQ, Kafka, MinIO, Soketi, Reverb, Mailpit, …) additionally get auto-provisioning and `.env` generation.

The CLI lives at `dev` and is the entire UX:

```
$ dev start
✓ mysql:8.0              already running (shared)
✓ redis:latest           started (shared)
✓ elasticsearch:8.11     started (dedicated)
✓ Database "my-project" created
✓ .env written
```

The audience is backend / full-stack developers — Laravel, Node, Go shops that juggle many services across many checkouts. The tone of the project is precise, lightly opinionated, and unfussy.

### Sources used

- **GitHub** — [JCombee/dev.env](https://github.com/JCombee/dev.env) — Go source, `README.md`, `CHANGELOG.md`, `CLAUDE.md`. Read these for a more accurate picture of the product than this design system can convey on its own.
  - [`README.md`](https://github.com/JCombee/dev.env/blob/main/README.md) — full feature & config reference
  - [`CLAUDE.md`](https://github.com/JCombee/dev.env/blob/main/CLAUDE.md) — internal architecture notes
  - [`CHANGELOG.md`](https://github.com/JCombee/dev.env/blob/main/CHANGELOG.md) — feature history
- No Figma, no existing brand assets, no slide template, no design system. Everything below is original work for this project.

## Index

| File / folder | What it contains |
|---|---|
| `README.md` | This file — brand context, content fundamentals, visual foundations, iconography |
| `SKILL.md` | Cross-compatible Agent Skill manifest for Claude / Claude Code |
| `colors_and_type.css` | Token CSS — colors, type scale, spacing, radii, shadows, semantic vars |
| `fonts/` | Webfont licence notes (IBM Plex is loaded from Google Fonts CDN) |
| `assets/logos/` | Wordmark + logomark SVGs in light/dark variants |
| `assets/icons/` | Custom service & status glyphs |
| `preview/` | Design-system tab cards (colors, type, components, etc.) |
| `ui_kits/docs/` | Documentation-site UI kit — homepage, docs page, CLI reference, changelog |

## Brand essentials

**Name.** Always lowercase: `dev.env`. The period is part of the wordmark — never replace it with a hyphen, slash, or space. When the name appears in body copy, treat it as a single token; do not split across lines.

**Tagline.** *Shared dev services. One binary. Zero config drift.* Use this on the homepage hero and social cards. Shorter alternatives: *One pool. Every project.* / *Stop duplicating containers.*

**What we are not.** dev.env is not a platform, not a SaaS, not a Kubernetes thing, not a replacement for Docker. Avoid corporate language ("enterprise", "platform", "solutions"). It's a small tool that does one thing.

---

## Content fundamentals

dev.env's voice is **precise, lightly opinionated, friendly without being chummy**. The product's source `README.md` is the canonical reference — match its rhythm.

### Voice

- **Direct.** Lead with the verb. *"Initialize the global directory."* not *"This command will initialize…"*.
- **Concrete.** Always pair a claim with the exact CLI output or YAML snippet that proves it. The product README never says *"easy to use"* — it shows a 5-line code block instead. Do the same.
- **Lightly opinionated.** When defaults matter, say so. *"Prefer a global override in `~/.dev.env/settings.yaml`."* — the README literally uses this phrasing and we should too.
- **No hype.** No "blazingly fast", "magical", "just works", "enterprise-grade". The tool is good because it's simple; the copy should be simple too.

### Person

- **Second person ("you")** for instructions: *"Place a `.dev.env.yaml` at the root of your project."*
- **First-person plural ("we")** is **never** used — there is no "team voice".
- **Imperative** for command descriptions: *"Start all services required by this project."*

### Casing

- Headings: **Sentence case.** *"Getting started"*, not *"Getting Started"*.
- The product name is always **`dev.env`** in body copy, **`DEV.ENV`** only when used as a heading-as-wordmark (rare; the source README uses this in its `# DEV.ENV` title).
- Command names are always lowercase backtick: `` `dev start` ``, `` `dev exec mysql` ``.
- File paths: backtick, with leading `~/` or `./` where relevant: `` `~/.dev.env/settings.yaml` ``.
- Env vars: backtick, ALL_CAPS: `` `DB_HOST` ``.

### Punctuation & formatting

- **Em-dashes** (—) with no surrounding spaces are fine in headers and short captions: *"One pool—every project"*. In body prose, prefer the spaced form *"one pool — every project"*.
- **Code blocks** lead the visual hierarchy on docs pages. Almost every section opens with a code block, then explains it. Never explain a CLI command without showing its output.
- **Tables** for field references — the product README uses them heavily. Mirror that.
- **Callouts** (`>` blockquote): used sparingly. The source README uses one for *"Not recommended."* warnings — that's the right density. About 1 per long page.

### Emoji

- **Never in body copy.** No 🚀 ✨ 🎉.
- The CLI uses ASCII glyphs (`✓`, `~`, `!`) — those *are* part of the brand. When rendered in HTML, use the same characters (not their emoji-presentation cousins). Treat them as type, not icons.
- Heading prefixes never get an emoji. If a section needs visual emphasis, use the monospaced section-label pattern (see Visual Foundations).

### Worked examples

> ✗ **"Welcome to dev.env! 🚀 Get started with our blazingly-fast service manager."**
> ✓ **"dev.env runs one shared pool of containers for all your projects."**

> ✗ **"Our platform makes managing dev environments effortless."**
> ✓ **"Every project needs MySQL, Redis, Elasticsearch. The typical solution is a docker-compose.yml per project. dev.env runs one shared pool instead."**

> ✗ **"To install, please run the following command:"**
> ✓ **"Install on macOS or Linux:"**

> ✗ **"Don't use per-project port overrides!"**
> ✓ **"Not recommended. A per-project port override forces the service to run as a dedicated container, bypassing the shared pool dev.env is optimised for."**

---

## Visual foundations

The brand is **typographic, monochromatic with one electric accent, and built on a tight grid**. Think technical-manual meets modern dev-tool docs (Bun, Astro, Linear changelog) but with sharper edges and more monospace.

### Colors

The palette is intentionally narrow: **ink, paper, one accent, semantic states**.

- **Lime `#C6F432`** — the only chromatic brand color. Use it for: the period in the wordmark, the active CLI cursor, the primary CTA, link underlines on hover, the `✓` success glyph in marketing surfaces. Never use it as a large flood fill — it's loud on purpose. About 5–8% coverage on any given page is the right amount.
- **Ink `#0E1116`** — near-black with a cool blue tint. Used for body type on light backgrounds and as the dark theme's main surface.
- **Paper `#F7F6F2`** — warm cream off-white. Default page surface in light mode. Slightly warm to offset the cool ink and balance the cool lime.
- **Stone scale** — 7-step neutral grays from `#E8E7E2` (subtle borders) to `#2A2D33` (high-contrast surface lines). Used for borders, dividers, secondary text, code-block backgrounds.
- **Semantic** — `#E5484D` (error / destructive), `#FFB224` (warning / dedicated badge), `#3E63DD` (info / link in light mode), `#46A758` (success accent for tables/charts; the CLI's `✓` itself stays lime in marketing, monochrome in real terminal output).

The lime is the only saturated color in the system. Reds, ambers, and blues are reserved for semantic meaning. Avoid introducing additional brand colors (no purples, no teals).

### Type

Single typeface family: **IBM Plex** (Sans + Mono, with Serif as an optional accent). Loaded from Google Fonts. Plex is humanist, technical, and reads well at both display and small mono sizes — appropriate for a tool that lives in terminals.

- **Display** — `IBM Plex Sans`, weight 600, negative tracking (-0.02em on h1, -0.015em on h2). Sizes step from 56 → 40 → 28 → 22 → 18.
- **Body** — `IBM Plex Sans`, weight 400, line-height 1.55, body size 16–17px.
- **Mono** — `IBM Plex Mono`, weight 400 for code, 500 for inline UI labels (button text, badges, section counters). Mono is heavily used — it's the brand's signal that this is a developer tool. Aim for 30–40% of all visible text on a marketing page to be monospaced.
- **Eyebrow labels** — Mono, 12–13px, uppercase, letter-spacing 0.06em, tinted with a mid-stone color. Used above section headings as `/ 01 — installation` or `/ services / databases`.

Optional serif (Plex Serif) is reserved for editorial pull-quotes in the blog/changelog only. Don't use it on docs pages.

### Spacing

8-pt base grid. Steps: 4, 8, 12, 16, 24, 32, 48, 64, 96, 128, 192. Most components snap to multiples of 8; tight typographic spacing uses 4.

Section padding: vertical 96–128px between major sections on marketing pages, 48–64px on docs pages.

Inline-text spacing inside a paragraph: never override the natural line-height. Spacing between paragraphs: 16–20px.

### Backgrounds

- **Default surface** is flat — paper in light mode, ink in dark mode. No gradients, no noise, no glassmorphism.
- **Hero / section dividers** can use a **light dot grid** (1px dots, 24px spacing, 6% opacity ink-on-paper or paper-on-ink). This is the only allowed decorative background and it's reserved for full-bleed hero sections and CTA strips. Never inside body content.
- **Code blocks** sit on a 1-step-darker stone background with a 1px solid stone-300 border. No tinted backgrounds.
- **Never use chromatic gradients.** No purple-to-blue, no neon glow. The lime never bleeds into a gradient.

### Borders & corners

- **Corner radius**: 4px (default), 6px (cards), 8px (large surfaces like the hero composer), 999px (capsule badges only). Buttons are 6px. **Never** use rounded-corner radii larger than 8px on rectangular surfaces — the brand reads as "engineering tool", not "consumer app".
- **Border width**: 1px solid. Use stone-200 for subtle dividers, stone-300 for visible card edges, ink at 60% for emphasized outlines.
- **Dashed borders** are allowed for "empty state" containers and for `dedicated: true` callout boxes (visually signals "private / set aside"). 1px dashed, 4px dash length.

### Shadows & elevation

The system is **mostly flat**. There are exactly three elevation tokens:

- **`elevation-0`** — no shadow. Default for cards and surfaces.
- **`elevation-1`** — `0 1px 0 0 rgb(14 17 22 / 0.04), 0 1px 2px 0 rgb(14 17 22 / 0.06)`. Used only for floating UI (dropdowns, popovers, command palette).
- **`elevation-2`** — `0 8px 24px -8px rgb(14 17 22 / 0.12), 0 2px 4px -2px rgb(14 17 22 / 0.06)`. Reserved for the homepage hero's terminal frame and modal dialogs.

No inner shadows. No layered/colored shadows. Depth is communicated by stacking borders and backgrounds, not by Z-axis blur.

### Hover & press states

- **Links (text)**: underline appears on hover, in `lime` on light bg, `lime` on dark bg. The text color itself does not change.
- **Primary button**: background shifts from ink → 90% ink + 10% lime mix on hover; on press, scales to 0.98 with a 80ms transition. The lime period in the wordmark stays lime.
- **Secondary button (outlined)**: background fills with stone-100 on hover; press uses stone-200 fill.
- **Cards / list rows**: background shifts +1 step on hover (paper → stone-100, or ink → stone-900). No transform.
- **Icon buttons**: opacity 0.7 → 1.0 on hover, ink fill 0.6 → 0.9 on press.
- **No glow effects.** No box-shadow on hover. No scale-up. Press states shrink (0.98), never grow.

### Animation

Animations are **functional and brief**. The brand's tempo is precise, not playful.

- **Easing**: a single custom cubic-bezier `cubic-bezier(0.32, 0.72, 0, 1)` (slight overshoot dampener) for transforms; standard `ease-out` for opacity-only fades.
- **Durations**: 80ms (press), 150ms (hover), 240ms (panel open / drawer), 400ms (route transition fades). Never longer than 400ms for UI; marketing scroll-reveals max 600ms.
- **Reveal pattern**: opacity 0 → 1 paired with `translateY(8px) → 0`. No scale-ups, no slide-from-left, no spring bounces.
- **Loading / pending**: the lime cursor block blinks at 1Hz exactly (60 frames at 60fps, 50% duty cycle). It's the brand's "thinking" indicator.

### Transparency & blur

Avoid both. No frosted-glass effects, no backdrop-blur. The one exception: the sticky docs header uses `background-color: rgb(247 246 242 / 0.85)` + `backdrop-filter: blur(8px)` so content slides under it cleanly. That's it.

### Cards

A card is a 1px stone-300 border on paper, 6px radius, no shadow. Padding 24px. The label "card" in this system means **a bordered container**, not a "floating tile" — there is no z-axis. Headers inside a card use a mono eyebrow then a sans heading.

### Imagery

There are no photographs in the brand. There are no illustrated mascots. **The only "imagery" is the terminal frame** — a stylized representation of the CLI output, framed in a window chrome with a dot-dot-dot top bar. Use it whenever the hero or a feature section needs a visual anchor.

Where a screenshot or diagram is truly needed (e.g. an architecture explainer), use **ASCII-art diagrams** rendered in a mono code block. They scale, they're searchable, they fit the brand.

### Layout rules

- **Max content width**: 1280px for marketing pages, 1080px for docs (incl. left nav).
- **Side gutters**: 32px on desktop, 16px on mobile.
- **Sticky top bar** on all surfaces: 56px tall, paper background with the translucent-blur exception above, 1px stone-200 bottom border.
- **Docs left nav**: 240px wide, sticky, scrolls independently. Right TOC: 200px wide, sticky.
- **No carousels**, no marquee tickers, no parallax. Content flows top-to-bottom.

---

## Iconography

dev.env's CLI uses **ASCII glyphs as its primary "icons"**: `✓` for success, `~` for unchanged-but-tracked, `!` for warning prompts. These are sacred to the brand — every marketing or docs surface that mimics CLI output must render them with the same characters. Don't substitute checkmark emoji.

**Outside the CLI**, the design system uses **[Lucide](https://lucide.dev)** as its icon set, loaded from the [unpkg CDN](https://unpkg.com/lucide-static/). Lucide is chosen because:

- 1.5px stroke weight matches the brand's overall line weight
- Geometric without being aggressive — no playful rounded ends
- MIT licensed, ~1500 icons, actively maintained

**Lucide substitution flag:** Lucide is the chosen icon set for this design system because no icon library shipped with the dev.env source. If the maintainers later commission custom service glyphs (one per supported Docker image), document the swap here and replace the `<lucide-*>` usages in the UI kit.

### Usage rules

- **Stroke icons only.** Never mix filled icons with stroke icons in the same surface.
- **Size scale**: 14px (inline-with-text), 16px (button), 20px (nav), 24px (feature card), 32px (hero). Stroke width stays 1.5px at every size — do not thicken at small sizes.
- **Color**: icons inherit `currentColor`. They default to `ink-600` (secondary text gray) and shift to `ink` on hover. They are **never** colored lime except inside the wordmark's period.
- **Spacing from text**: icon-and-label gap is always 8px.

### Service glyphs (custom)

The one place where pure Lucide isn't enough is representing **the supported Docker services** (MySQL, Redis, Elasticsearch, etc.). For these, the brand uses **a custom 24×24 monogram square** — a stone-200 rounded square with a 1.5px ink letter inside (e.g. "M" for MySQL). This is faster to ship than commissioning service-vendor logos, sidesteps trademark issues, and visually unifies the service table. The monogram square pattern is documented in `assets/icons/services/`.

### Emoji & Unicode

- **Emoji**: not used. Anywhere.
- **Unicode glyphs**: only the CLI set (`✓` `~` `!` `✗` `→` `…`). No `★`, no `❯`, no `▸`. The right-arrow `→` is allowed in CTAs ("Read the docs →") because it sets a softer tone than a Lucide chevron and matches the prose register.

---

## Where to go next

This file is the orientation document. For specifics:

- **Tokens** → `colors_and_type.css` (drop into any HTML for the full variable set)
- **Cards on the Design System tab** → `preview/` (one HTML file per card)
- **Production-quality UI components** → `ui_kits/docs/index.html` and the `*.jsx` files alongside it. That kit reproduces the documentation site's home, docs, CLI reference, and changelog views.
- **Skill for Claude Code** → `SKILL.md` ships this whole folder as a skill named `dev-env-design`.
