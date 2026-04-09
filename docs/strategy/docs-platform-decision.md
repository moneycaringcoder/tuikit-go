# Docs Platform Decision: MkDocs Material

**Decision date:** 2026-04-09
**Status:** Adopted

## Options Considered

| Option | Pros | Cons |
|--------|------|------|
| **MkDocs Material** | Zero JS build step, Markdown-native, beautiful dark/light themes, search out of the box, GitHub Pages deploy in one command | Python dependency |
| Docusaurus | React ecosystem, MDX support | Node.js build toolchain, heavier CI setup |
| Hugo | Fast, Go-native | Less polished default theme for API docs, more templating boilerplate |
| pkg.go.dev | Auto-generated godoc | No narrative docs, no examples gallery |

## Decision

**MkDocs Material** was chosen for the following reasons:

1. **Markdown-first** — all doc content is plain `.md` files; no JSX or template syntax. Contributors only need to know Markdown.
2. **GitHub Pages in one command** — `mkdocs gh-deploy` pushes to `gh-pages` branch automatically. The CI workflow is 10 lines.
3. **Navigation tabs** — `navigation.tabs` maps cleanly onto the Getting Started / Components / Reference structure tuikit-go needs.
4. **`content.code.copy`** — one-click copy on every code block is critical for a library reference site.
5. **Search** — `search.suggest` and `search.highlight` work without a paid plan or external index.
6. **Pinned dependency** — `mkdocs-material==9.5.18` in `requirements.txt` keeps CI reproducible.

## Site Structure

```
site/
  mkdocs.yml          # MkDocs config + nav
  requirements.txt    # pinned mkdocs-material
  docs/
    index.md          # landing page
    guides/           # Getting Started, Theming, Testing, CLI, Self-update, Layout
    components/       # one page per component with API tables + examples
```

## Deployment

The `docs.yml` workflow triggers on pushes to `main` that touch `site/**`. It installs Python, runs `pip install -r site/requirements.txt`, then calls `mkdocs gh-deploy --force` which commits the built HTML to the `gh-pages` branch.
