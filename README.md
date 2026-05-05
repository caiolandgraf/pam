<div align="center">

# 🗂️ PAM
### Pam's Database Drawer — SQL Query Management for the Scranton Branch

[![MIT License](https://img.shields.io/badge/license-MIT-white.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![GitHub Release](https://img.shields.io/github/v/release/caiolandgraf/pam)](https://github.com/caiolandgraf/pam/releases)

**A minimal CLI tool for managing and executing SQL queries across multiple databases.**
**Written in Go, made beautiful with BubbleTea. Filed, labeled, and ready when you need it.**

[Quick Start](#-quick-start) • [Configuration](docs/configuration.md) • [Commands](docs/commands.md) • [Keybindings](docs/keybindings.md) • [Features](docs/features.md) • [Completion](docs/completion.md) • [Databases](docs/databases.md) • [Roadmap](#%EF%B8%8F-roadmap) • [Contributing](CONTRIBUTING.md)

> 🟨 **Sticky Note:** This project is currently in beta — please report unexpected behavior through the Issues tab.

</div>

<table>
  <tr>
    <td width="60%">
      <strong>📎 INTEROFFICE MEMO</strong><br/>
      <strong>To:</strong> Everyone running queries<br/>
      <strong>From:</strong> Pam Beesly, Receptionist & DBA<br/>
      <strong>Subject:</strong> Please use the drawer. I organized it myself.<br/>
      <strong>Status:</strong> 🟡 Beta — file unexpected behavior in Issues
    </td>
    <td width="40%">
      <strong>🗂 OFFICE DIRECTORY</strong>
      <ul>
        <li><a href="#-quick-start">Quick Start</a></li>
        <li><a href="#-installation">Installation</a></li>
        <li><a href="#%EF%B8%8F-database-support">Database Support</a></li>
        <li><a href="#-features">Features</a></li>
        <li><a href="#-all-commands">All Commands</a></li>
        <li><a href="#%EF%B8%8F-tui-keybindings">TUI Keybindings</a></li>
        <li><a href="#%EF%B8%8F-configuration">Configuration</a></li>
        <li><a href="#%EF%B8%8F-roadmap">Roadmap</a></li>
        <li><a href="#-contributing">Contributing</a></li>
      </ul>
    </td>
  </tr>
</table>

<pre>
┌──────────────────────────────┐
│  DUNDER MIFFLIN • SCRANTON   │
│   FILED ✔   AUTHORIZED ✔     │
└──────────────────────────────┘
</pre>

---

> *"I have a system. Every query has a name, every name has a drawer, every drawer has a label. It's not complicated, it's just how I do things."* — Pam Beesly, Scranton Branch Receptionist & Unofficial DBA

PAM (Pam's Database Drawer) is a keyboard-first CLI tool for saving, organizing, and running SQL queries across multiple databases. Think of it as a filing cabinet for your SQL — instead of hunting through terminal history or a dozen `.sql` files, you `pam add`, `pam list`, and `pam run`. The interactive TUI table viewer lets you explore results, edit cells in-place, export data, and visualize schema relationships, all without leaving the terminal.

---

## 🎬 Demo

![pamdemo](https://github.com/user-attachments/assets/b62bec1d-2255-4d02-9b7f-1c99afbeb664)

---

## ✨ Features

- **Query Library** — Save, label, and organize queries like a well-maintained filing cabinet; search by name or content
- **Multi-Database** — PostgreSQL, MySQL/MariaDB, SQLite, Oracle, SQL Server, ClickHouse, Firebird, DuckDB, and **Snowflake**
- **Interactive TUI** — Vim-style keyboard navigation in a beautiful BubbleTea table viewer
- **In-Place Editing** — Update cells, delete rows, and edit SQL directly from the results table
- **Interactive Shell** — `pam shell` / `pam repl` for a persistent SQL REPL with history, multi-line input, and meta-commands
- **Flexible Export** — `pam run --format <csv|json|tsv|html|sql|markdown>` streams results to stdout, pipe-friendly
- **Edit Before Run** — `pam run --edit` / `-e` opens the query in `$EDITOR` before executing
- **Repeat Last Query** — `pam run --last` / `-l` re-runs the last executed query without retyping
- **Visual Line Mode** — `V` selects entire rows; `v` selects cell ranges — both copyable with `y`
- **Export Full Table** — `X` exports the entire result set to clipboard in your chosen format
- **Full-Text Search** — `/` searches cell contents; `f` searches column headers; `n`/`N` cycles matches
- **Row Marking** — `m` marks rows for bulk operations; `D` deletes all marked rows in one round-trip
- **Inline Cell Edit** — `e` edits a cell value in-place; `E` edits the query and reruns it
- **Connection Management** — `pam remove --connection <name>` removes a saved connection
- **Config Editor** — `pam config` opens the config file in `$EDITOR`
- **SQL Import** — `pam import <file>` imports SQL dumps; `pam export` creates them
- **Table Query Shortcut** — `pam query --table=<name>` for quick table access
- **Enhanced Explain** — `pam explain --depth <n>` visualizes FK relationships up to N levels deep
- **Shell Completion** — `pam completion --install` writes completion scripts to the standard path automatically
- **`pam tables` / `\dt`** — list tables directly from the interactive shell
- **Environment Variable Expansion** — use `${MY_VAR}` in connection strings; PAM expands them at runtime
- **Database Exploration** — browse schema, visualize foreign key relationships with `pam explore` and `pam explain`
- **Parameterized Queries** — `:param|default` syntax; pass values with `--param` flags or positional args

See [Features](docs/features.md) for details and examples

---

## 🗄️ Database Support

| Database | Type String | Notes |
|----------|------------|-------|
| PostgreSQL | `postgres` | Schema selection supported |
| MySQL / MariaDB | `mysql` / `mariadb` | |
| SQLite | `sqlite` | Local file-based |
| Oracle | `oracle` | Schema selection supported |
| SQL Server | `sqlserver` | |
| ClickHouse | `clickhouse` | |
| Firebird | `firebird` | |
| DuckDB | `duckdb` | CSV/JSON file queries; requires CGO |
| Snowflake | `snowflake` | Keypair authentication supported |

See connection init examples in [Database Support](docs/databases.md)

---

## 🚀 Quick Start

*Set up your desk at reception in under 2 minutes.*

```bash
# Create your first connection (PostgreSQL example)
pam init mydb postgres "postgresql://user:pass@localhost:5432/mydb"

# Use environment variables in connection strings — PAM expands them at runtime
pam init mydb postgres "postgresql://${DB_USER}:${DB_PASS}@localhost:5432/mydb"

# Add a saved query
pam add list_users "SELECT * FROM users"

# List your saved queries
pam list queries

# Run it — opens the interactive table viewer
pam run list_users

# Or run inline SQL
pam run "SELECT * FROM products WHERE price > 100"

# Export results without opening the TUI
pam run list_users --format csv > users.csv
pam run list_users --format json

# Start an interactive SQL shell
pam shell
```

---

## 📦 Installation

Go to [the releases page](https://github.com/caiolandgraf/pam/releases) and download the binary for your system. Make sure it is executable and in a directory on your `$PATH`.

<details>
<summary>Go install</summary>

```bash
go install github.com/caiolandgraf/pam/cmd/pam@latest
```

This puts the `pam` binary in your `$GOBIN` path (usually `~/go/bin`).
</details>

<details>
<summary>Build Manually</summary>

```bash
git clone https://github.com/caiolandgraf/pam
go build -o pam ./cmd/pam
```

The `pam` binary will be available in the project root directory.

DuckDB requires CGO and is included in the default build. To build without DuckDB:

```bash
CGO_ENABLED=0 go build -o pam ./cmd/pam
```
</details>

<details>
<summary>Nix / NixOS (Flake)</summary>

PAM is available as a Nix flake for easy installation on NixOS and systems with Nix.

#### Run directly without installing
```bash
nix run github:caiolandgraf/pam
```

#### Install to user profile
```bash
nix profile install github:caiolandgraf/pam
```

#### Enter development shell
```bash
nix develop github:caiolandgraf/pam
```

#### NixOS System-wide

Add to your flake-based `configuration.nix` or `flake.nix`:

```nix
{
  description = "My NixOS config";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    pam.url = "github:caiolandgraf/pam";
  };

  outputs = { self, nixpkgs, pam, ... }: {
    nixosConfigurations.myHostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        {
          environment.systemPackages = [
            pam.packages.x86_64-linux.default
          ];
        }
      ];
    };
  };
}
```

Then rebuild: `sudo nixos-rebuild switch`

#### Home Manager

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nix-unstable";
    pam.url = "github:caiolandgraf/pam";
  };

  outputs = { self, nixpkgs, pam, ... }: {
    homeConfigurations."username" = {
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      modules = [
        {
          home.packages = [
            pam.packages.x86_64-linux.default
          ];
        }
      ];
    };
  };
}
```

Then apply: `home-manager switch`
</details>

<details>
<summary>Arch (AUR) - Unofficial</summary>

An unofficial AUR package is available at: [pam-bin](https://aur.archlinux.org/packages/pam-bin)
</details>

---

## 📋 All Commands

See [Commands](docs/commands.md) for the full command reference and database init examples.

### Connection Management

| Command | Description | Example |
|---------|-------------|---------|
| `init <name> <type> <conn>` | Create a new database connection | `pam init mydb postgres "postgresql://..."` |
| `use` / `switch <name>` | Switch active connection | `pam use production` |
| `status` | Show current connection info | `pam status` |
| `list connections` | List all configured connections | `pam list connections` |
| `remove --connection <name>` | Remove a saved connection | `pam remove --connection dev4` |

### Query Operations

| Command | Description | Example |
|---------|-------------|---------|
| `add <name> [sql]` | Save a new query | `pam add users "SELECT * FROM users"` |
| `remove <name\|id>` | Remove a saved query | `pam remove users` |
| `list queries` | List all saved queries | `pam list queries` |
| `list queries --oneline` | One query per line | `pam list -o` |
| `list queries <term>` | Search queries by name or SQL | `pam list employees` |
| `run <name\|id\|sql>` | Execute a query | `pam run users` or `pam run 2` |
| `run` | Create and run a new query | `pam run` |
| `run --edit` / `-e` | Edit query before running | `pam run users --edit` |
| `run --last` / `-l` | Re-run last executed query | `pam run --last` |
| `run --format <fmt>` | Output as csv/json/tsv/html/sql/markdown | `pam run users --format json` |
| `run --param` | Run with named parameters | `pam run emp --name Michael` |
| `shell` / `repl` | Interactive SQL REPL with history | `pam shell` |

### Database Exploration

| Command | Description | Example |
|---------|-------------|---------|
| `explore` | List all tables and views | `pam explore` |
| `explore <table> [-l N]` | Query a table with optional row limit | `pam explore employees --limit 100` |
| `explain <table>` | Visualize foreign key relationships | `pam explain employees` |
| `explain <table> -d N` | FK relationships up to depth N | `pam explain employees --depth 2` |
| `tables` | Open tables in the TUI results view | `pam tables` |
| `query --table=<name>` | Quick table query in TUI | `pam query --table=employees` |

### Configuration & Utilities

| Command | Description | Example |
|---------|-------------|---------|
| `config` | Edit config file in `$EDITOR` | `pam config` |
| `edit` | Edit all queries for current connection | `pam edit` |
| `edit <name\|id>` | Edit a single named query | `pam edit 3` |
| `import <file>` | Import a SQL dump from a file | `pam import dump.sql` |
| `export` | Dump all tables to stdout | `pam export > backup.sql` |
| `export --table=<t>` | Dump a single table | `pam export --table=users` |
| `export --output=<f>` | Write dump to a file | `pam export --output=dump.sql` |
| `export --no-data` | Schema only (no INSERT statements) | `pam export --no-data` |
| `export --data-only` | Data only (no CREATE TABLE) | `pam export --data-only > inserts.sql` |
| `export --drop` | Prepend DROP TABLE IF EXISTS | `pam export --drop --output=full.sql` |
| `completion --install` | Install shell completion scripts | `pam completion --install` |
| `help [command]` | Show help information | `pam help run` |

### Command Aliases

| Alias | Full Command | Description |
|-------|--------------|-------------|
| `use` | `switch` | Switch active connection |
| `save` | `add` | Save a new query |
| `delete` | `remove` | Remove a saved query or connection |
| `ls` | `list connections` | List all connections |
| `t`, `explore` | `tables` | List or query tables |
| `tv` | `table-view` | Inspect and edit table structure |
| `test` | `status` | Show current connection |
| `clear`, `unset` | `disconnect` | Disconnect from database |
| `repl` | `shell` | Interactive SQL REPL |

---

## ⌨️ TUI Keybindings

See [Keybindings](docs/keybindings.md) for the full reference.

Once your query results appear you can navigate, edit, and export without leaving the TUI:

```
# Navigation (Vim-style)
j / k      Move down / up
h / l      Move left / right
g / G      Jump to first / last row
0 / $      Jump to first / last column
Ctrl+u/d   Page up / down

# Copy & export
y          Yank (copy) current cell
v          Visual selection mode (cell range)
V          Visual line mode (full rows)
x          Export selection (csv/tsv/json/sql/markdown/html)
X          Export the entire table

# Edit data
e          Edit cell value in-place
E          Edit query and re-run
m          Mark / unmark row for bulk ops
D          Delete all marked rows (or current row)

# Search
/          Search cell contents  (n / N to cycle)
f          Search column headers (; / , to cycle)

# Other
Enter      Detail view (JSON-formatted)
s          Save current query
?          Toggle keybindings help
q / Esc    Quit table view
```

---

## ⚙️ Shell Completion

PAM provides dynamic tab completion for bash, zsh, and fish — automatically including your saved queries and connections.

```bash
# Install to the standard path for your shell (recommended)
pam completion --install

# Or manually:
echo 'eval "$(pam completion bash)"' >> ~/.bashrc        # Bash
echo 'eval "$(pam completion zsh)"' >> ~/.zshrc          # Zsh
pam completion fish > ~/.config/fish/completions/pam.fish # Fish
```

See [Shell Completion](docs/completion.md)

---

## ⚙️ Configuration

Row limits, column widths, color schemes, and UI visibility options are configured at `~/.config/pam/config.yaml`.

```yaml
default_row_limit: 1000
default_column_width: 15
color_scheme: "dracula"   # dracula, gruvbox, catppuccin-mocha, tokyo-night, nord, rose-pine, ...

ui_visibility:
  query_name: true          # Show query name header
  query_sql: true           # Show SQL query display
  type_display: true        # Show column type indicators
  key_icons: true           # Show primary key (⚿) and foreign key (⚭) icons
  footer_cell_content: true # Show current cell preview in footer
  footer_stats: true        # Show row/col count and position in footer
  footer_keymaps: true      # Show keybindings help in footer
```

Open and edit the config directly with:

```bash
pam config
```

See [Configuration](docs/configuration.md)

---

## 🗺️ Roadmap

> *"I have a lot of questions. Number one: how dare you ship without tests."* — Dwight Schrute (probably)

### v1.2.0 — The Merge
- [x] Snowflake database support with keypair authentication
- [x] DuckDB improvements (CSV/JSON file queries, improved driver)
- [x] Environment variable expansion in connection strings (`${MY_VAR}`)
- [x] Interactive SQL REPL (`pam shell` / `pam repl`) with history, multi-line, meta-commands
- [x] `pam run --format <csv|json|tsv|html|sql|markdown>` — pipe-friendly export
- [x] `pam run --edit` / `-e` — open query in editor before running
- [x] `pam run --last` / `-l` — repeat last executed query
- [x] Visual line mode (`V` key)
- [x] Export full table (`X` key)
- [x] `pam remove --connection <name>` — remove saved connections
- [x] `pam config` — edit config file in `$EDITOR`
- [x] `pam import <file>` — import SQL dumps
- [x] `pam query --table=<name>` — quick table query
- [x] Enhanced `pam explain --depth <n>` — N-level FK visualization
- [x] Shell completion with `--install` flag (writes to standard paths)
- [x] `pam tables` / `\dt` in interactive shell
- [x] Full-text search in table view (`/`) and column header search (`f`)
- [x] Row marking (`m`) and multi-row delete (`D`)
- [x] Inline cell edit (`e`) and edit+rerun (`E`)

### v1.3.0 — Schrute's Farm
- [ ] Configurable keybinds
- [ ] Migrate to Bubble Tea v2
- [ ] Return more info on exec statements (INSERT, UPDATE, DELETE row counts)
- [ ] Homebrew custom tap and nixpkgs entry
- [ ] More options to encrypt data in the config file

---

## 🤖 For Robots

PAM ships a `SKILL.md` file in the repo root — a simple reference for AI coding agents (Claude Code, Copilot, etc.) to use PAM non-interactively. It covers safe commands, format flags, parameterized queries, and which commands to avoid (TUI/editor). Point your agent at it if you want it to run SQL queries as part of an automated workflow.

---

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed instructions.

*"We keep it like a Dunder Mifflin memo: concise, kind, and well-labeled."*

Thanks to all contributors:

<a href="https://github.com/DeprecatedLuar"><img src="https://github.com/DeprecatedLuar.png" width="40" /></a>
<a href="https://github.com/caiolandgraf"><img src="https://github.com/caiolandgraf.png" width="40" /></a>
<a href="https://github.com/g4brielklein"><img src="https://github.com/g4brielklein.png" width="40" /></a>
<a href="https://github.com/eduardofuncao"><img src="https://github.com/eduardofuncao.png" width="40" /></a>
<a href="https://github.com/udirona"><img src="https://github.com/udirona.png" width="40" /></a>
<a href="https://github.com/Leosallin"><img src="https://github.com/Leosallin.png" width="40" /></a>

PAM builds on these fantastic projects:

- **[naggie/dstask](https://github.com/naggie/dstask)** — elegant CLI design patterns and file-based data storage
- **[DeprecatedLuar/better-curl-saul](https://github.com/DeprecatedLuar/better-curl-saul)** — a simple and genius approach to CLI tooling
- **[eduardofuncao/squix](https://github.com/eduardofuncao/squix)** — upstream project whose features were merged into PAM
- **[dbeaver](https://github.com/dbeaver/dbeaver)** — the OG database management tool

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the Go standard library.

<div align="center">

**Made with 🗂️ for the Scranton Branch**

*"I am ready to face any challenges that might be foolish enough to face me."* — Dwight Schrute

</div>
