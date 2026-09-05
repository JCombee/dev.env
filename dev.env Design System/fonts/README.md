# Fonts

The dev.env design system uses **IBM Plex** in three flavors:

| Family | Used for | Weights |
|---|---|---|
| IBM Plex Sans | All UI text — headings, body, captions | 400, 500, 600, 700 |
| IBM Plex Mono | Code, CLI output, eyebrow labels, badges | 400, 500, 600 |
| IBM Plex Serif | Editorial pull-quotes (rare) | 400, 500, 400 italic |

IBM Plex is licensed under the **SIL Open Font License 1.1** — free for commercial and open-source use including embedding in apps.

## Loading

`colors_and_type.css` imports IBM Plex from Google Fonts via `@import url(...)`. No local font files are checked into this repo to keep it lean. For offline-first or production environments, self-host by downloading from:

- [IBM Plex on GitHub](https://github.com/IBM/plex) (official, all weights, woff2 included)
- [Google Fonts → IBM Plex Sans](https://fonts.google.com/specimen/IBM+Plex+Sans)
- [Google Fonts → IBM Plex Mono](https://fonts.google.com/specimen/IBM+Plex+Mono)
- [Google Fonts → IBM Plex Serif](https://fonts.google.com/specimen/IBM+Plex+Serif)

## ⚠ Substitution flag

dev.env had **no specified typeface** when this design system was created — IBM Plex was chosen by the design agent as a close match for the brand's "technical, humanist, monospace-friendly" needs. If the dev.env maintainers prefer a different family (Geist, Söhne, Berkeley Mono, GT America, etc.), swap the `@import` in `colors_and_type.css` and update the `--font-*` variables. Everything downstream will inherit.
