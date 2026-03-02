<template>
  <div class="docs">
    <div class="docs-layout">
      <!-- Sidebar -->
      <aside class="sidebar">
        <nav>
          <div class="sidebar-group" v-for="g in toc" :key="g.title">
            <h4>{{ g.title }}</h4>
            <a
              v-for="item in g.items"
              :key="item.id"
              :href="'#' + item.id"
              :class="{ active: activeId === item.id }"
              @click.prevent="scrollTo(item.id)"
            >
              {{ item.label }}
            </a>
          </div>
        </nav>
      </aside>

      <!-- Content -->
      <div class="docs-content">
        <h1>Documentation</h1>
        <p class="lead">
          Everything you need to know about Pam's Database Drawer.
        </p>

        <!-- Installation -->
        <section id="installation">
          <h2>Installation</h2>

          <h3 id="go-install">Go Install</h3>
          <CodeBlock title="bash">{{ snippets.goInstall }}</CodeBlock>

          <h3 id="build-manually">Build Manually</h3>
          <CodeBlock title="bash">{{ snippets.buildManually }}</CodeBlock>

          <h3 id="nix">Nix / NixOS</h3>
          <CodeBlock title="bash">{{ snippets.nix }}</CodeBlock>

          <h3 id="releases">Binary Releases</h3>
          <p>
            Download pre-built binaries from the
            <a
              href="https://github.com/caiolandgraf/pam/releases"
              target="_blank"
              >releases page</a
            >. Available for Linux (amd64, arm64), macOS (amd64, arm64) and
            Windows.
          </p>
        </section>

        <!-- Connections -->
        <section id="connections">
          <h2>Connections</h2>

          <h3 id="init">Creating a connection</h3>
          <p>
            Use <code>pam init</code> to create and validate a new database
            connection. Run it with no arguments to launch the interactive TUI,
            which asks for each field individually and assembles the connection
            string live as you type. The database type is auto-inferred from the
            connection string when using CLI flags or positional arguments.
          </p>
          <CodeBlock title="bash">{{ snippets.init }}</CodeBlock>
          <CodeBlock title="Interactive TUI">{{ snippets.initTUI }}</CodeBlock>

          <h3 id="switch">Switching connections</h3>
          <CodeBlock title="bash">{{ snippets.switchConn }}</CodeBlock>

          <h3 id="status">Connection status</h3>
          <CodeBlock title="bash">{{ snippets.status }}</CodeBlock>

          <h3 id="list-connections">Listing connections</h3>
          <CodeBlock title="bash">{{ snippets.listConnections }}</CodeBlock>

          <h3 id="disconnect">Disconnecting</h3>
          <CodeBlock title="bash">{{ snippets.disconnect }}</CodeBlock>
        </section>

        <!-- Database Support -->
        <section id="databases">
          <h2>Database Support</h2>
          <div class="db-examples">
            <div
              v-for="db in dbExamples"
              :key="db.name"
              :id="'db-' + db.name.toLowerCase().replace(/[^a-z]/g, '')"
            >
              <h3>{{ db.name }}</h3>
              <CodeBlock title="bash">{{ db.code }}</CodeBlock>
            </div>
          </div>
        </section>

        <!-- Query Management -->
        <section id="queries">
          <h2>Query Management</h2>

          <h3 id="add-query">Adding queries</h3>
          <CodeBlock title="bash">{{ snippets.addQuery }}</CodeBlock>

          <h3 id="run-query">Running queries</h3>
          <CodeBlock title="bash">{{ snippets.runQuery }}</CodeBlock>

          <h3 id="query-table">Query a table</h3>
          <CodeBlock title="bash">{{ snippets.queryTable }}</CodeBlock>

          <h3 id="list-queries">Listing queries</h3>
          <CodeBlock title="bash">{{ snippets.listQueries }}</CodeBlock>

          <h3 id="remove-query">Removing queries</h3>
          <CodeBlock title="bash">{{ snippets.removeQuery }}</CodeBlock>

          <h3 id="remove-conn">Removing a connection</h3>
          <p>
            Use <code>pam remove --conn &lt;name&gt;</code> to permanently
            delete a saved connection. If it is the currently active connection,
            it will be cleared automatically.
          </p>
          <CodeBlock title="bash">{{ snippets.removeConn }}</CodeBlock>
        </section>

        <!-- Tables -->
        <section id="tables">
          <h2>Database Exploration</h2>

          <h3 id="tables-list">Tables</h3>
          <CodeBlock title="bash">{{ snippets.tablesList }}</CodeBlock>

          <h3 id="table-view">Table View</h3>
          <p>
            View and edit the structure (columns, types, constraints) of a
            table.
          </p>
          <CodeBlock title="bash">{{ snippets.tableView }}</CodeBlock>

          <h3 id="explore">Explore</h3>
          <CodeBlock title="bash">{{ snippets.explore }}</CodeBlock>

          <h3 id="explain">Explain</h3>
          <p>Visualize foreign key relationships between tables.</p>
          <CodeBlock title="bash">{{ snippets.explain }}</CodeBlock>

          <h3 id="info">Info</h3>
          <CodeBlock title="bash">{{ snippets.info }}</CodeBlock>
        </section>

        <!-- Import & Export -->
        <section id="export-import">
          <h2>Import &amp; Export</h2>

          <h3 id="export">Export</h3>
          <p>
            Export one or all tables as a SQL dump — generates
            <code>CREATE TABLE</code> and <code>INSERT</code> statements. SQL is
            written to <strong>stdout</strong> by default; use
            <code>--output</code> to write to a file. Status messages always go
            to <strong>stderr</strong>, so redirects like
            <code>pam export &gt; dump.sql</code> work cleanly.
          </p>
          <CodeBlock title="bash">{{ snippets.export }}</CodeBlock>

          <h3 id="import">Import</h3>
          <p>
            Execute a SQL dump against the active connection. Reads from a file
            or <strong>stdin</strong>. Statements are split on <code>;</code>,
            correctly handling single-quoted strings, <code>--</code> line
            comments and <code>/* */</code> block comments. By default the
            import stops on the first error; use
            <code>--continue-on-error</code> to collect all failures and report
            them at the end.
          </p>
          <CodeBlock title="bash">{{ snippets.import }}</CodeBlock>
        </section>

        <!-- Configuration -->
        <section id="config">
          <h2>Configuration</h2>
          <p>
            Pam stores configuration at <code>~/.config/pam/config.yaml</code>.
          </p>

          <h3 id="row-limit">Row Limit</h3>
          <p>
            All queries are limited by default. Set
            <code>default_row_limit</code> in config (default: 1000).
          </p>

          <h3 id="column-width">Column Width</h3>
          <p>
            Columns are dynamic and responsive — they adapt to content, terminal
            width, and header size. The fallback
            <code>default_column_width</code> (default: 15) is used for edge
            cases.
          </p>

          <h3 id="color-schemes">Color Schemes</h3>
          <p>
            Set <code>color_scheme</code> in your config. Available schemes:
          </p>
          <div class="scheme-chips">
            <span class="scheme-chip" v-for="s in colorSchemes" :key="s">{{
              s
            }}</span>
          </div>

          <h3 id="editor">Editor</h3>
          <p>
            Pam uses your <code>$EDITOR</code> environment variable for all
            editing operations (defaults to <code>vim</code>).
          </p>
          <CodeBlock title="bash">{{ snippets.editor }}</CodeBlock>

          <h3 id="edit-config">Editing config</h3>
          <CodeBlock title="bash">{{ snippets.editConfig }}</CodeBlock>
        </section>

        <!-- Shell Completion -->
        <section id="completion">
          <h2>Shell Completion</h2>
          <p>
            Pam supports tab-completion for bash, zsh and fish. Table names are
            fetched <strong>dynamically</strong> from the active database.
          </p>
          <CodeBlock title="bash">{{ snippets.completion }}</CodeBlock>
        </section>

        <!-- Commands Reference -->
        <section id="commands">
          <h2>Commands Reference</h2>
          <div class="commands-table-wrap">
            <table class="commands-table">
              <thead>
                <tr>
                  <th>Command</th>
                  <th>Alias</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in commands" :key="c.cmd">
                  <td>
                    <code>{{ c.cmd }}</code>
                  </td>
                  <td>
                    <code v-if="c.alias">{{ c.alias }}</code>
                  </td>
                  <td>{{ c.desc }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>

        <!-- TUI Navigation -->
        <section id="tui">
          <h2>TUI Table Navigation</h2>
          <div class="commands-table-wrap">
            <table class="commands-table">
              <thead>
                <tr>
                  <th>Key</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="k in tuiKeys" :key="k.key">
                  <td>
                    <kbd>{{ k.key }}</kbd>
                  </td>
                  <td>{{ k.action }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import CodeBlock from '../components/CodeBlock.vue'

const activeId = ref('')

const toc = [
  {
    title: 'Getting Started',
    items: [
      { id: 'installation', label: 'Installation' },
      { id: 'connections', label: 'Connections' },
      { id: 'databases', label: 'Database Support' }
    ]
  },
  {
    title: 'Usage',
    items: [
      { id: 'queries', label: 'Query Management' },
      { id: 'remove-conn', label: 'Removing a connection' },
      { id: 'tables', label: 'Database Exploration' },
      { id: 'export-import', label: 'Import & Export' }
    ]
  },
  {
    title: 'Configuration',
    items: [
      { id: 'config', label: 'Config' },
      { id: 'completion', label: 'Shell Completion' }
    ]
  },
  {
    title: 'Reference',
    items: [
      { id: 'commands', label: 'Commands' },
      { id: 'tui', label: 'TUI Keys' }
    ]
  }
]

const allIds = toc.flatMap(g => g.items.map(i => i.id))

function onScroll() {
  for (let i = allIds.length - 1; i >= 0; i--) {
    const el = document.getElementById(allIds[i])
    if (el && el.getBoundingClientRect().top <= 120) {
      activeId.value = allIds[i]
      return
    }
  }
  activeId.value = allIds[0]
}

function scrollTo(id) {
  const el = document.getElementById(id)
  if (el) {
    window.scrollTo({ top: el.offsetTop - 80, behavior: 'smooth' })
    activeId.value = id
  }
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})
onUnmounted(() => window.removeEventListener('scroll', onScroll))

// ---------------------------------------------------------------------------
// Code snippets — kept here as template literals so newlines are preserved
// ---------------------------------------------------------------------------
const snippets = {
  goInstall: `go install github.com/caiolandgraf/pam/cmd/pam@latest`,

  buildManually: `git clone https://github.com/caiolandgraf/pam
cd pam
go build -o pam ./cmd/pam`,

  nix: `# Run directly
nix run github:caiolandgraf/pam

# Install to profile
nix profile install github:caiolandgraf/pam`,

  init: `# Interactive TUI — recommended, guides you field by field
pam init

# Auto-inferred type from connection string
pam init mydb "postgres://user:pass@localhost:5432/mydb"

# Explicit type
pam init mydb postgres "postgres://user:pass@localhost:5432/mydb"

# With schema
pam init mydb oracle "user/pass@localhost:1521/XEPDB1" my_schema`,

  initTUI: `  Connection name  › mydb
  Database type    › postgres  ◀ ▶
  Host             › localhost
  Port             › 5432
  Username         › myuser
  Password         › ••••••
  Database         › mydb

  conn › postgres://myuser:secret@localhost:5432/mydb

  ↑/↓ Tab: navigate  Ctrl+P: show/hide password  Enter: confirm  Esc: cancel`,

  switchConn: `pam switch production
pam use dev          # alias`,

  status: `pam status`,

  listConnections: `pam list connections
pam ls               # alias`,

  disconnect: `pam disconnect`,

  addQuery: `# Inline
pam add list_users "SELECT * FROM users"

# With parameters and defaults
pam add emp_by_salary "SELECT * FROM employees WHERE salary > :min_sal|30000"

# Opens $EDITOR
pam add my_query`,

  runQuery: `# By name or ID
pam run list_users
pam run 2

# Inline SQL
pam run "SELECT * FROM products WHERE price > 100"

# Edit before running
pam run list_users --edit

# Re-run last query
pam run --last

# With named parameters
pam run emp_by_salary --min_sal 50000

# With positional parameters
pam run search_users Michael active`,

  queryTable: `# Default SELECT * FROM table
pam query --table=users

# Custom SQL with table context
pam query --table=users "SELECT * FROM users WHERE active = 1"

# Open editor
pam query --table=orders --edit`,

  listQueries: `pam list queries
pam list queries emp          # search
pam list queries --oneline    # compact`,

  removeQuery: `pam remove list_users
pam remove 3`,

  removeConn: `# Remove by name
pam remove --conn mydb
pam remove -c mydb       # short flag
pam delete --conn mydb   # alias`,

  tablesList: `# List all tables
pam tables

# Query a table directly
pam tables users

# One per line
pam tables --oneline`,

  tableView: `pam table-view users
pam tv users         # alias`,

  explore: `pam explore
pam explore employees --limit 100`,

  explain: `pam explain employees
pam explain employees --depth 2`,

  info: `pam info tables
pam info views`,

  editor: `export EDITOR=vim
export EDITOR=nano
export EDITOR=code`,

  editConfig: `pam edit config
pam edit queries`,

  completion: `# Persist completions (recommended)
pam completion bash --install
pam completion zsh --install
pam completion fish --install

# Or load for current session only
eval "$(pam completion bash)"
eval "$(pam completion zsh)"
pam completion fish | source`,

  export: `# Dump all tables to a file (schema + data)
pam export --output=backup.sql

# Dump a single table
pam export --table=users --output=users.sql
pam export users                              # shorthand

# Pipe to stdout (status messages go to stderr — clean redirect)
pam export --table=orders > orders.sql

# Schema only — CREATE TABLE statements, no INSERT
pam export --no-data --output=schema.sql

# Data only — INSERT statements, no CREATE TABLE
pam export --data-only > inserts.sql

# Add DROP TABLE IF EXISTS before each CREATE TABLE
pam export --drop --output=full.sql`,

  import: `# Import from a file
pam import dump.sql
pam import --file=dump.sql                    # same with explicit flag

# Read from stdin (pipe-friendly)
cat dump.sql | pam import

# Don't stop on the first error — collect and report all failures
pam import dump.sql --continue-on-error

# Dry run — parse and list every statement without executing
pam import dump.sql --dry-run

# Full migration pipeline between two connections
pam switch staging
pam export --table=users > users.sql
pam switch production
pam import users.sql`
}

// ---------------------------------------------------------------------------
// Data
// ---------------------------------------------------------------------------
const dbExamples = [
  {
    name: 'PostgreSQL',
    code: 'pam init pg-prod postgres "postgres://user:pass@localhost:5432/mydb?sslmode=disable"'
  },
  {
    name: 'MySQL / MariaDB',
    code: "pam init mysql-dev mysql 'user:pass@tcp(127.0.0.1:3306)/mydb'"
  },
  {
    name: 'SQL Server',
    code: 'pam init sqlserver-docker sqlserver "sqlserver://sa:MyStrongPass123@localhost:1433/master"'
  },
  {
    name: 'SQLite',
    code: 'pam init sqlite-local sqlite "file:///home/user/mydb.sqlite"'
  },
  {
    name: 'Oracle',
    code: 'pam init oracle-stg oracle "user/pass@localhost:1521/XEPDB1"'
  },
  {
    name: 'ClickHouse',
    code: 'pam init clickhouse-docker clickhouse "clickhouse://user:pass@localhost:9000/mydb"'
  },
  {
    name: 'Firebird',
    code: 'pam init firebird-docker firebird "user:masterkey@localhost:3050//var/lib/firebird/data/mydb"'
  }
]

const colorSchemes = [
  'default',
  'dracula',
  'gruvbox',
  'solarized',
  'nord',
  'monokai',
  'black-metal',
  'vesper',
  'catppuccin-mocha',
  'tokyo-night',
  'rose-pine',
  'terracotta'
]

const commands = [
  {
    cmd: 'init',
    alias: '',
    desc: 'Create and validate a new database connection'
  },
  { cmd: 'switch', alias: 'use', desc: 'Switch the active connection' },
  {
    cmd: 'disconnect',
    alias: 'clear, unset',
    desc: 'Disconnect from current database'
  },
  { cmd: 'add', alias: 'save', desc: 'Save a new named query' },
  {
    cmd: 'remove',
    alias: 'delete',
    desc: 'Remove a saved query by name or ID'
  },
  {
    cmd: 'remove --conn <name>',
    alias: 'delete --conn',
    desc: 'Remove a saved connection'
  },
  { cmd: 'run', alias: '', desc: 'Execute a saved query or inline SQL' },
  {
    cmd: 'query --table=<t>',
    alias: '',
    desc: 'Run SQL against a specific table'
  },
  { cmd: 'list', alias: '', desc: 'List connections or queries' },
  { cmd: 'ls', alias: '', desc: 'Shortcut for list connections' },
  {
    cmd: 'tables',
    alias: 't, explore',
    desc: 'List or query database tables'
  },
  { cmd: 'table-view', alias: 'tv', desc: 'View and edit table structure' },
  {
    cmd: 'info',
    alias: '',
    desc: 'Show tables or views in current connection'
  },
  { cmd: 'explain', alias: '', desc: 'Visualize foreign key relationships' },
  {
    cmd: 'export',
    alias: '',
    desc: 'Export tables as a SQL dump (CREATE TABLE + INSERTs)'
  },
  {
    cmd: 'import',
    alias: '',
    desc: 'Import a SQL dump file or stdin into the active connection'
  },
  { cmd: 'edit', alias: '', desc: 'Open config or queries in editor' },
  { cmd: 'status', alias: 'test', desc: 'Show current active connection' },
  {
    cmd: 'completion',
    alias: '',
    desc: 'Generate shell completion (bash, zsh, fish)'
  },
  { cmd: 'help', alias: '', desc: 'Show help for a command' }
]

const tuiKeys = [
  { key: 'h / ←', action: 'Move left' },
  { key: 'j / ↓', action: 'Move down' },
  { key: 'k / ↑', action: 'Move up' },
  { key: 'l / →', action: 'Move right' },
  { key: 'g', action: 'Jump to first row' },
  { key: 'G', action: 'Jump to last row' },
  { key: '0 / Home', action: 'Jump to first column' },
  { key: '$ / End', action: 'Jump to last column' },
  { key: 'Ctrl+u', action: 'Page up' },
  { key: 'Ctrl+d', action: 'Page down' },
  { key: 'y', action: 'Copy cell to clipboard' },
  { key: 'v', action: 'Enter visual selection mode' },
  { key: 'x', action: 'Export selected data' },
  { key: 'f', action: 'Toggle sort on column' },
  { key: 'u', action: 'Update current cell' },
  { key: 'D', action: 'Delete current row' },
  { key: 'e', action: 'Edit and re-run query' },
  { key: 's', action: 'Save current query' },
  { key: 'Enter', action: 'Detail view (JSON formatted)' },
  { key: 'q / Esc', action: 'Quit' }
]
</script>

<style scoped>
.docs {
  padding-top: 64px;
}
.docs-layout {
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  gap: 2rem;
}

/* ---------- SIDEBAR ---------- */
.sidebar {
  position: sticky;
  top: 80px;
  width: 220px;
  flex-shrink: 0;
  height: calc(100vh - 80px);
  overflow-y: auto;
  padding: 2rem 0 2rem 1.5rem;
}
.sidebar-group {
  margin-bottom: 1.5rem;
}
.sidebar-group h4 {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}
.sidebar-group a {
  display: block;
  padding: 4px 0 4px 12px;
  font-size: 0.85rem;
  color: var(--text-secondary);
  border-left: 2px solid transparent;
  transition: all 0.15s;
}
.sidebar-group a:hover {
  color: var(--text);
}
.sidebar-group a.active {
  color: var(--accent);
  border-left-color: var(--accent);
}

/* ---------- CONTENT ---------- */
.docs-content {
  flex: 1;
  min-width: 0;
  padding: 2rem 1.5rem 4rem;
}
.docs-content h1 {
  font-size: 2.2rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
}
.lead {
  color: var(--text-secondary);
  font-size: 1.1rem;
  margin-bottom: 3rem;
}
.docs-content section {
  margin-bottom: 3.5rem;
}
.docs-content h2 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
}
.docs-content h3 {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 1.5rem 0 0.5rem;
}
.docs-content p {
  color: var(--text-secondary);
  margin-bottom: 0.75rem;
}

/* schemes */
.scheme-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
}
.scheme-chip {
  background: var(--bg-code);
  border: 1px solid var(--border);
  padding: 3px 12px;
  border-radius: 100px;
  font-size: 0.8rem;
  font-family: var(--font-mono);
  color: var(--text-secondary);
}

/* commands table */
.commands-table-wrap {
  overflow-x: auto;
}
.commands-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}
.commands-table th {
  text-align: left;
  font-weight: 600;
  padding: 0.6rem 1rem;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
  color: var(--text-secondary);
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.commands-table td {
  padding: 0.5rem 1rem;
  border-bottom: 1px solid var(--border);
  color: var(--text-secondary);
}
.commands-table tbody tr:hover {
  background: var(--bg-card);
}
.commands-table kbd {
  font-family: var(--font-mono);
  font-size: 0.8rem;
  background: var(--bg-code);
  border: 1px solid var(--border);
  padding: 2px 8px;
  border-radius: 4px;
  color: var(--text);
}

@media (max-width: 900px) {
  .sidebar {
    display: none;
  }
  .docs-layout {
    display: block;
  }
}
</style>
