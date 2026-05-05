# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.2.0] - 2026-05-05

### Added

#### From Upstream (squix merge)
- **Snowflake support** — connect to Snowflake with keypair authentication (`authenticator=snowflake_jwt`)
- **DuckDB improvements** — query CSV and JSON files directly as views; improved driver stability
- **Environment variable expansion** — use `${MY_VAR}` in any connection string; PAM expands at runtime
- **Interactive SQL REPL** — `pam shell` / `pam repl` with persistent connection, command history, multi-line input, and meta-commands (`\dt`, `\q`, `\h`, etc.)
- **`pam run --format`** — stream results to stdout as `csv`, `json`, `tsv`, `html`, `sql`, or `markdown`; fully pipe-friendly
- **`pam run --edit` / `-e`** — open query in `$EDITOR` before executing
- **`pam run --last` / `-l`** — re-run the last executed query
- **Visual line mode** — `V` key selects entire rows for copy/export (previously only cell-range `v` was available)
- **Export full table** — `X` key exports the entire result set to clipboard in the chosen format
- **`pam remove --connection <name>`** — remove a saved database connection from the config
- **`pam config`** — open the config file directly in `$EDITOR`
- **`pam import <file>`** — import a SQL dump file; reads from stdin when no file is given
- **`pam query --table=<name>`** — open a specific table directly in the TUI results view
- **Enhanced `pam explain`** — `--depth` / `-d` flag controls how many FK levels deep the tree renders
- **Shell completion `--install`** — writes completion scripts to the standard shell path automatically (fish/zsh/bash)
- **`pam tables` / `\dt` in shell** — list database tables from within the interactive REPL
- **Full-text search in table view** — `/` searches cell contents (`n`/`N` to cycle); `f` searches column headers (`;`/`,` to cycle)

#### CLI / TUI
- **Inline cell editor** — press `e` to edit a cell value directly inside the TUI using a `bubbles/textarea`, no external `$EDITOR` needed
- **Row marking** — press `m` to mark / unmark individual rows; marked rows are highlighted and queued for bulk operations
- **Multi-row delete** — `D` now deletes all marked rows (or the current visual selection range) in a single round-trip, with a confirmation prompt
- **In-app TableView DDL editor** — add, edit, rename, and drop columns using an inline textarea instead of spawning an external editor process

#### Docs / Playground
- **MemoPage component** — reusable layout with Dunder Mifflin letterhead, binder hole-punches, red ruled-paper margin line, and a form footer (`Page 1 of 1`)
- **Playground: pre-loaded queries** — `top_earners`, `by_dept`, `active_projects`, `recent_hires` are available immediately after load
- **Playground: `clear` command** — clears terminal output (alias for Ctrl+L)
- **Playground: `dm` easter egg** — prints a random Dunder Mifflin wisdom quote
- **Playground: Tab autocomplete** — pressing Tab completes the first token of any PAM command
- **Playground: categorized Quick Reference card** — chips organized into Connections / Schema / Saved Queries / Data sections
- **Playground: clickable Workflow demos** — three pre-built step-by-step flows (Explore Schema, Build Query Library, Cross-Table Analysis)
- **CSS utilities** — `laminated-card`, `annotation`, `hole-punch` / `hole-punch-row`, `dm-letterhead`, `memo-page-footer`, `red-margin`, `sticky-note-blue/green/pink` (with dark-mode variants)

### Changed

#### CLI / TUI
- Key rebinding: `e` → edit cell in-place (previously `u`); `E` → edit and rerun query (previously `e`)
- Delete (`D`) now dispatches to marked-rows path → visual-selection path → single-row path (in that priority order)

#### Docs
- **Navbar**: replaced fixed `height: 72px` with `min-height: 64px` + padding; added progressive breakpoints (1020 px hides nav-meta, 860 px hides brand subtitle, 768 px activates hamburger)
- **Playground**: sections reorganized — laminated reference card, sticky-note "About the Data" panel, workflow grid
- **Global theme**: distressed stamps (−8°, thicker border), tape strip on sticky notes, Dunder Mifflin letterhead on every page, page footers

### Fixed
- Navbar breaking / content overflowing at medium viewport widths (nav-meta `flex-wrap` inside fixed-height container)
- NavBar `nav-meta` no longer overlaps nav links when viewport is between 769 px and 1020 px

---

## [1.1.0] - 2025-04-XX

### Added
- SSL mode support in interactive `init` UI
- Interactive connection removal in the TUI
- Documentation for interactive init flow and connection management

### Changed
- Improved `pam init` help text and prompts

---

## [1.0.1] - 2025-XX-XX

### Fixed
- Minor bug fixes following the v1.0.0 release

---

## [1.0.0] - 2025-XX-XX

### Added
- Initial stable release
- Multi-database connection management (`init`, `switch`, `disconnect`)
- Named query library (`add`, `remove`, `run`, `list`)
- Interactive TUI table viewer (`tables`, `tv`, `query`)
- Visual mode, column sorting, detail view, cell copy
- Foreign-key relationship explorer (`explain`)
- SQL dump import / export
- Dynamic column editing (add, rename, drop via DDL)
- Shell completion for bash, zsh, and fish
- Docs site with interactive SQLite WASM playground
