# dev.env Documentation Site — UI kit

The documentation site is dev.env's primary brand surface. This kit recreates the homepage, docs reader, CLI reference, and changelog as a click-through prototype using React.

## Pages

| File | What it shows |
|---|---|
| `index.html` | Homepage — hero, install card, features, supported services |
| `docs.html` | Long-form docs page — left nav, body, right TOC |
| `reference.html` | CLI command reference — searchable command cards |
| `changelog.html` | Versioned changelog with conventional-commit groups |

Open `index.html` first — the header navigates between all four pages, and the theme toggle (top-right) flips ink/paper.

## Components

Each component lives in its own JSX file and is mounted to `window` so other files can use it.

- `Header.jsx` — sticky top bar with wordmark, page nav, theme toggle
- `Footer.jsx` — three-column footer with mono labels
- `Hero.jsx` — homepage hero with copy + live terminal
- `Terminal.jsx` — reusable CLI frame with traffic-light bar and the lime cursor
- `FeatureGrid.jsx` — bordered card grid; one card per pillar feature
- `InstallBlock.jsx` — install snippet with macOS/Linux/Windows tabs
- `ServiceTable.jsx` — supported-services table with monogram chips and badges
- `Sidebar.jsx` — left nav on the docs reader
- `DocsArticle.jsx` — rendered article with prose styles, callouts, tables
- `TOC.jsx` — right table-of-contents
- `CommandCard.jsx` — single CLI command spec card (used on the reference page)
- `ChangelogEntry.jsx` — one version's entries grouped by conventional-commit type
- `Badge.jsx` / `Button.jsx` / `Eyebrow.jsx` / `Kbd.jsx` — primitives

All components read from the design tokens in `../../colors_and_type.css`. Adding a new component? Stick to the variables there; if you need a new token, add it to that file first.
