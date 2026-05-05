<template>
  <div class="playground-page">
    <MemoPage
      title="Playground"
      subtitle="Try PAM commands in your browser — powered by SQLite (WASM). No server needed. Everything runs locally."
      to="Training Room"
      from="Scranton Branch"
      subject="Playground drills"
      status="Hands-on"
      stamp="TRAINING"
      badge="Form TR-04"
    >
      <template #tags>
        <span class="label-chip">SQLite WASM</span>
        <span class="label-chip">Local Only</span>
        <span class="label-chip">Safe Sandbox</span>
      </template>
      <template #meta>
        <div class="form-row">
          <div class="form-field">
            <span class="form-label">Session</span> Local
          </div>
          <div class="form-field">
            <span class="form-label">Engine</span> SQLite WASM
          </div>
          <div class="form-field">
            <span class="form-label">Mode</span> Sandbox
          </div>
        </div>
      </template>
      <template #note>
        Pre-loaded: <strong>run top_earners</strong> ·
        <strong>run by_dept</strong> · <strong>run active_projects</strong> — or
        type <strong>help</strong>
      </template>

      <!-- Terminal-style frame -->
      <div class="terminal-frame">
        <div class="terminal-bar">
          <span class="dot red"></span>
          <span class="dot yellow"></span>
          <span class="dot green"></span>
          <span class="terminal-title"
            >pam (playground) ▸ scranton branch terminal</span
          >
        </div>

        <div class="terminal-body" ref="terminalBody">
          <!-- Output history -->
          <div class="terminal-output">
            <div
              v-for="(entry, i) in outputHistory"
              :key="i"
              class="output-entry"
            >
              <div v-if="entry.type === 'command'" class="output-command">
                <span class="prompt">{{ entry.prompt }}</span>
                <span class="command-text">{{ entry.text }}</span>
              </div>
              <div v-else-if="entry.type === 'info'" class="output-info">
                {{ entry.text }}
              </div>
              <div v-else-if="entry.type === 'success'" class="output-success">
                {{ entry.text }}
              </div>
              <div v-else-if="entry.type === 'error'" class="output-error">
                {{ entry.text }}
              </div>
              <div v-else-if="entry.type === 'table'" class="output-table-wrap">
                <table class="output-table">
                  <thead>
                    <tr>
                      <th class="row-num">#</th>
                      <th v-for="col in entry.columns" :key="col">{{ col }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(row, ri) in entry.rows" :key="ri">
                      <td class="row-num">{{ ri + 1 }}</td>
                      <td
                        v-for="col in entry.columns"
                        :key="col"
                        :class="{ 'null-val': row[col] === null }"
                      >
                        {{ row[col] === null ? 'NULL' : row[col] }}
                      </td>
                    </tr>
                  </tbody>
                </table>
                <div class="table-meta">
                  {{ entry.rows.length }} row{{
                    entry.rows.length !== 1 ? 's' : ''
                  }}
                  · {{ entry.columns.length }} column{{
                    entry.columns.length !== 1 ? 's' : ''
                  }}
                  <span v-if="entry.elapsed"> · {{ entry.elapsed }}</span>
                </div>
              </div>
              <div v-else-if="entry.type === 'list'" class="output-list">
                <div
                  v-for="(item, li) in entry.items"
                  :key="li"
                  class="list-item"
                  :class="{ active: item.active }"
                >
                  <span v-if="item.active" class="active-marker">▸ </span>
                  <span>{{ item.text }}</span>
                </div>
              </div>
              <div v-else-if="entry.type === 'help'" class="output-help">
                <pre>{{ entry.text }}</pre>
              </div>
            </div>
          </div>

          <!-- Input -->
          <div class="terminal-input-line">
            <span class="prompt">{{ currentPrompt }}</span>
            <input
              ref="inputRef"
              v-model="currentInput"
              class="terminal-input"
              spellcheck="false"
              autocomplete="off"
              @keydown="handleKeydown"
              placeholder="Type a pam command... (try: help)"
            />
          </div>
        </div>
      </div>

      <!-- Quick Reference Card -->
      <div class="reference-card laminated-card">
        <div class="ref-header">
          <span class="ref-icon">📋</span>
          <h3>Quick Reference</h3>
          <span class="ref-badge">laminated for your convenience</span>
        </div>
        <div class="ref-categories">
          <div
            v-for="cat in commandCategories"
            :key="cat.label"
            class="ref-category"
          >
            <div class="ref-cat-label">{{ cat.label }}</div>
            <div class="ref-chips">
              <button
                v-for="cmd in cat.commands"
                :key="cmd.cmd"
                class="command-chip"
                @click="runQuickCommand(cmd.cmd)"
              >
                <code>{{ cmd.label }}</code>
                <span class="chip-desc">{{ cmd.desc }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Workflows -->
      <div class="workflows-section">
        <div class="workflows-header">
          <span class="memo-stamp stamp-approved">WORKFLOWS</span>
          <p class="workflows-sub">
            Click any workflow to run it step-by-step in the terminal
          </p>
        </div>
        <div class="workflow-grid">
          <div
            v-for="wf in workflows"
            :key="wf.num"
            class="workflow-card"
            @click="runWorkflow(wf)"
          >
            <div class="wf-header">
              <span class="wf-num">{{ String(wf.num).padStart(2, '0') }}</span>
              <span class="wf-title">{{ wf.title }}</span>
            </div>
            <p class="wf-desc">{{ wf.desc }}</p>
            <div class="wf-cmds">
              <code v-for="(cmd, i) in wf.preview" :key="i">{{ cmd }}</code>
            </div>
            <span class="wf-run">▸ Run Demo</span>
          </div>
        </div>
      </div>

      <!-- About the Data -->
      <div class="about-section">
        <h3 class="about-title">About the Sandbox Data</h3>
        <div class="about-grid">
          <div class="sticky-note about-note">
            <strong>👥 employees (15 rows)</strong><br />
            Pam, Michael, Dwight, Jim and the whole Scranton branch. Columns:
            <code>id</code>, <code>first_name</code>, <code>last_name</code>,
            <code>title</code>, <code>salary</code>, <code>department_id</code>,
            <code>hire_date</code>.
          </div>
          <div class="sticky-note sticky-note-blue about-note">
            <strong>🏢 departments (5 rows)</strong><br />
            Engineering, Sales, Marketing, HR, Finance — each with a budget and
            location.
          </div>
          <div class="sticky-note sticky-note-green about-note">
            <strong>📁 projects (5 rows)</strong><br />
            Active, completed, and planned projects with department ownership.
          </div>
          <div class="sticky-note sticky-note-pink about-note">
            <strong>🔗 employee_projects</strong><br />
            Junction table linking employees to projects with a role column.
            Great for <code>explain</code> and JOIN demos.
          </div>
        </div>
        <p class="about-hint annotation">
          “Wikipedia is the best thing ever. Anyone in the world can write
          anything they want about any subject — so you know you are getting the
          best possible information.” — Michael Scott
        </p>
      </div>
    </MemoPage>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import MemoPage from '../components/MemoPage.vue'

// ── Reactive state ──────────────────────────────────────────────────────────
const outputHistory = ref([])
const currentInput = ref('')
const inputRef = ref(null)
const terminalBody = ref(null)
const commandHistory = ref([])
const historyIndex = ref(-1)

// ── PAM state ───────────────────────────────────────────────────────────────
let db = null
let dbReady = false

const connections = ref({})
const activeConnection = ref('')
const savedQueries = ref({}) // { connName: { queryName: { id, sql } } }
let queryIdCounter = 0

const currentPrompt = ref('pam ▸ ')

// ── Command categories (for reference card) ────────────────────────────────────────────
const commandCategories = [
  {
    label: '📡 Connections',
    commands: [
      { cmd: 'init mydb sqlite ""', label: 'init', desc: 'Create connection' },
      { cmd: 'switch playground', label: 'switch', desc: 'Change active' },
      { cmd: 'status', label: 'status', desc: 'Active connection' },
      { cmd: 'ls', label: 'ls', desc: 'List connections' },
      { cmd: 'disconnect', label: 'disconnect', desc: 'Drop connection' }
    ]
  },
  {
    label: '🗂 Schema',
    commands: [
      { cmd: 'tables', label: 'tables', desc: 'List all tables' },
      { cmd: 'tv employees', label: 'tv employees', desc: 'Table structure' },
      { cmd: 'info tables', label: 'info tables', desc: 'Schema info' },
      { cmd: 'explain employees', label: 'explain', desc: 'Relationships' }
    ]
  },
  {
    label: '💾 Saved Queries',
    commands: [
      { cmd: 'list queries', label: 'list queries', desc: 'See saved queries' },
      {
        cmd: 'add top5 "SELECT first_name, salary FROM employees ORDER BY salary DESC LIMIT 5"',
        label: 'add',
        desc: 'Save a query'
      },
      {
        cmd: 'run top_earners',
        label: 'run top_earners',
        desc: 'Run pre-loaded'
      },
      { cmd: 'run by_dept', label: 'run by_dept', desc: 'Dept summary' }
    ]
  },
  {
    label: '🔍 Data',
    commands: [
      {
        cmd: 'tables employees',
        label: 'tables employees',
        desc: 'Browse table'
      },
      {
        cmd: 'query --table=employees "SELECT first_name, title, salary FROM employees WHERE salary > 60000 ORDER BY salary DESC"',
        label: 'query --table',
        desc: 'Filtered query'
      },
      {
        cmd: 'run active_projects',
        label: 'run active_projects',
        desc: 'Active projects'
      },
      {
        cmd: 'run recent_hires',
        label: 'run recent_hires',
        desc: 'Recent hires'
      }
    ]
  }
]

// ── Workflows ──────────────────────────────────────────────────────────────────────────────
const workflows = [
  {
    num: 1,
    title: 'Explore the Schema',
    desc: 'List tables, view structure, and understand relationships.',
    preview: ['tables', 'tv employees', 'explain employees'],
    commands: [
      'tables',
      'tv employees',
      'tv departments',
      'explain employees',
      'explain projects'
    ]
  },
  {
    num: 2,
    title: 'Build a Query Library',
    desc: 'Save reusable named queries and run them on demand.',
    preview: ['add top5 "..."', 'list queries', 'run top5'],
    commands: [
      'add high_earners "SELECT first_name, last_name, title, salary FROM employees WHERE salary > 60000 ORDER BY salary DESC"',
      'add dept_budgets "SELECT name, budget, location FROM departments ORDER BY budget DESC"',
      'list queries',
      'run high_earners',
      'run dept_budgets'
    ]
  },
  {
    num: 3,
    title: 'Cross-Table Analysis',
    desc: 'Use the pre-loaded queries to explore joins and aggregations.',
    preview: ['run by_dept', 'run active_projects', 'run recent_hires'],
    commands: [
      'run by_dept',
      'run active_projects',
      'run recent_hires',
      'run top_earners'
    ]
  }
]

// ── Helpers ─────────────────────────────────────────────────────────────────

function pushOutput(entry) {
  outputHistory.value.push(entry)
  scrollToBottom()
}

function pushCommand(text) {
  pushOutput({ type: 'command', prompt: currentPrompt.value, text })
}

function pushInfo(text) {
  pushOutput({ type: 'info', text })
}

function pushSuccess(text) {
  pushOutput({ type: 'success', text })
}

function pushError(text) {
  pushOutput({ type: 'error', text })
}

function pushTable(columns, rows, elapsed) {
  pushOutput({ type: 'table', columns, rows, elapsed })
}

function pushList(items) {
  pushOutput({ type: 'list', items })
}

function pushHelp(text) {
  pushOutput({ type: 'help', text })
}

async function scrollToBottom() {
  await nextTick()
  if (terminalBody.value) {
    terminalBody.value.scrollTop = terminalBody.value.scrollHeight
  }
}

function execSQL(sql) {
  if (!db) throw new Error('Database not initialized')
  const t0 = performance.now()
  const result = db.exec(sql)
  const t1 = performance.now()
  const elapsed = ((t1 - t0) / 1000).toFixed(3) + 's'
  return { result, elapsed }
}

function execRun(sql) {
  if (!db) throw new Error('Database not initialized')
  db.run(sql)
}

function resultToTable(result, elapsed) {
  if (!result || result.length === 0) return null
  const res = result[0]
  if (!res || !res.columns) return null
  const rows = (res.values || []).map(row => {
    const obj = {}
    res.columns.forEach((col, i) => {
      obj[col] = row[i]
    })
    return obj
  })
  return { columns: res.columns, rows, elapsed }
}

function updatePrompt() {
  if (activeConnection.value) {
    currentPrompt.value = `pam (${activeConnection.value}) ▸ `
  } else {
    currentPrompt.value = 'pam ▸ '
  }
}

// ── Command parser ──────────────────────────────────────────────────────────

function parseArgs(input) {
  const args = []
  let current = ''
  let inQuote = false
  let quoteChar = ''

  for (let i = 0; i < input.length; i++) {
    const ch = input[i]
    if (inQuote) {
      if (ch === quoteChar) {
        inQuote = false
      } else {
        current += ch
      }
    } else if (ch === '"' || ch === "'") {
      inQuote = true
      quoteChar = ch
    } else if (ch === ' ') {
      if (current) {
        args.push(current)
        current = ''
      }
    } else {
      current += ch
    }
  }
  if (current) args.push(current)
  return args
}

// ── Command handlers ────────────────────────────────────────────────────────

function handleCommand(input) {
  const trimmed = input.trim()
  if (!trimmed) return

  // Store in history
  commandHistory.value.push(trimmed)
  historyIndex.value = commandHistory.value.length

  pushCommand(trimmed)

  // Strip leading "pam " if user types it
  let cmd = trimmed
  if (cmd.startsWith('pam ')) {
    cmd = cmd.slice(4)
  }

  const args = parseArgs(cmd)
  if (args.length === 0) return

  const command = args[0].toLowerCase()
  const rest = args.slice(1)

  try {
    switch (command) {
      case 'help':
        cmdHelp(rest)
        break
      case 'init':
      case 'create':
        cmdInit(rest)
        break
      case 'switch':
      case 'use':
        cmdSwitch(rest)
        break
      case 'status':
      case 'test':
        cmdStatus()
        break
      case 'add':
      case 'save':
        cmdAdd(rest)
        break
      case 'remove':
      case 'delete':
        cmdRemove(rest)
        break
      case 'run':
        cmdRun(rest)
        break
      case 'query':
        cmdQuery(rest)
        break
      case 'list':
        cmdList(rest)
        break
      case 'ls':
        cmdListConnections()
        break
      case 'tables':
      case 't':
      case 'explore':
        cmdTables(rest)
        break
      case 'table-view':
      case 'tv':
        cmdTableView(rest)
        break
      case 'info':
        cmdInfo(rest)
        break
      case 'explain':
        cmdExplain(rest)
        break
      case 'disconnect':
      case 'unset':
        cmdDisconnect()
        break
      case 'edit':
        cmdEdit(rest)
        break
      case 'completion':
        cmdCompletion()
        break
      case 'history':
        cmdHistory()
        break
      case 'clear':
        outputHistory.value = []
        break
      case 'dm':
        cmdDM()
        break
      case '-v':
      case '--version':
        pushInfo("Pam's database drawer\nversion: v0.3.0-playground")
        break
      default:
        pushError(
          `Unknown command: ${command}\nType "help" to see available commands.`
        )
    }
  } catch (e) {
    pushError(e.message)
  }
}

function cmdHelp(rest) {
  if (rest.length > 0) {
    const sub = rest[0].toLowerCase()
    const helpTexts = {
      init: `Command: init
Create and validate a new database connection.

Usage:
  pam init <name> <type> <connection-string>

In this playground, type is always "sqlite" and connection string is ignored.
A new in-memory SQLite database is created.

Examples:
  pam init mydb sqlite ":memory:"
  pam init dev sqlite ""`,

      switch: `Command: switch
Switch the active connection used by other commands.

Usage:
  pam switch <connection-name>
  pam use <connection-name>

Examples:
  pam switch dev
  pam use playground`,

      add: `Command: add
Save a new named query under the current connection.

Usage:
  pam add <query-name> <sql>

Examples:
  pam add list_users "SELECT * FROM employees"
  pam add top_earners "SELECT * FROM employees WHERE salary > 70000"`,

      run: `Command: run
Execute a saved query or inline SQL.

Usage:
  pam run <query-name-or-id>
  pam run "<inline-sql>"

Examples:
  pam run list_users
  pam run 1
  pam run "SELECT * FROM departments"`,

      query: `Command: query
Run a SQL query against a specific table.

Usage:
  pam query --table=<table> [sql]
  pam query -t <table> [sql]

If no SQL is provided, defaults to SELECT * FROM <table>.

Examples:
  pam query --table=employees
  pam query --table=employees "SELECT first_name, salary FROM employees WHERE salary > 60000"
  pam query -t departments`,

      tables: `Command: tables
List all tables or query a specific table.

Usage:
  pam tables
  pam tables <table-name>

Examples:
  pam tables
  pam tables employees`,

      'table-view': `Command: table-view
View the structure of a table (columns, types, constraints).

Usage:
  pam table-view <table-name>
  pam tv <table-name>

Examples:
  pam tv employees
  pam table-view departments`,

      tv: `Command: table-view (alias: tv)
View the structure of a table (columns, types, constraints).

Usage:
  pam tv <table-name>

Examples:
  pam tv employees`,

      list: `Command: list
List connections or queries.

Usage:
  pam list connections
  pam list queries [search-term]

Examples:
  pam list connections
  pam list queries
  pam list queries emp`,

      info: `Command: info
Show tables or views in the current connection.

Usage:
  pam info tables
  pam info views

Examples:
  pam info tables`,

      explain: `Command: explain
Show foreign key relationships for a table.

Usage:
  pam explain <table>

Examples:
  pam explain employees
  pam explain projects`,

      remove: `Command: remove
Remove a saved query by name or ID.

Usage:
  pam remove <query-name-or-id>

Examples:
  pam remove list_users
  pam remove 1`,

      status: `Command: status
Show the current active connection.

Usage:
  pam status`,

      disconnect: `Command: disconnect
Clear the active connection.

Usage:
  pam disconnect`,

      completion: `Command: completion
Generate shell completion scripts.

Usage:
  pam completion <bash|zsh|fish> [--install]

(Not available in playground)`,

      edit: `Command: edit
Open config or queries in your editor.

Usage:
  pam edit config
  pam edit queries

(Not available in playground — use "pam list queries" instead)`
    }

    if (helpTexts[sub]) {
      pushHelp(helpTexts[sub])
    } else {
      pushError(`No help available for "${sub}"`)
    }
    return
  }

  pushHelp(`Pam's database drawer — query manager for your databases

Usage:
  pam <command> [arguments]

Commands:
  init          Create a new database connection
  switch/use    Switch the active connection
  disconnect    Disconnect from current database
  status        Show current connection
  add/save      Save a new named query
  remove/delete Remove a saved query
  run           Execute a saved query or inline SQL
  query         Run SQL against a specific table
  tables/t      List or query database tables
  table-view/tv View table structure (columns, types)
  list          List connections or queries
  ls            List connections (shortcut)
  info          Show tables or views
  explain       Show table relationships
  edit          Open config or queries in editor
  completion    Generate shell completion script
  help          Show this help

Examples:
  pam init mydb sqlite ":memory:"
  pam switch playground
  pam add list_users "SELECT * FROM employees"
  pam run list_users
  pam tables employees
  pam query --table=employees "SELECT * FROM employees WHERE salary > 70000"

Type "help <command>" for detailed help on a specific command.`)
}

function cmdInit(rest) {
  if (rest.length < 1) {
    pushError(
      'Usage: pam init <name> [type] [connection-string]\n\nIn the playground, a new SQLite database is created automatically.'
    )
    return
  }

  const name = rest[0]

  if (connections.value[name]) {
    pushError(
      `Connection "${name}" already exists. Use "pam switch ${name}" to activate it.`
    )
    return
  }

  connections.value[name] = { type: 'sqlite', conn: ':memory:' }
  savedQueries.value[name] = {}
  activeConnection.value = name
  updatePrompt()
  pushSuccess(
    `✓ Connection "${name}" created and activated (SQLite in-browser)`
  )
}

function cmdSwitch(rest) {
  if (rest.length < 1) {
    pushError('Usage: pam switch <connection-name>')
    return
  }

  const name = rest[0]
  if (!connections.value[name]) {
    pushError(
      `Connection "${name}" not found.\nAvailable: ${Object.keys(connections.value).join(', ') || '(none)'}`
    )
    return
  }

  activeConnection.value = name
  updatePrompt()
  pushSuccess(`✓ Switched to "${name}"`)
}

function cmdStatus() {
  if (!activeConnection.value) {
    pushInfo(
      'No active connection. Use "pam switch <name>" or "pam init <name>" to connect.'
    )
    return
  }

  const conn = connections.value[activeConnection.value]
  pushInfo(
    `Active connection: ${activeConnection.value}\nType: ${conn.type}\nConnection: ${conn.conn}`
  )
}

function cmdAdd(rest) {
  requireConnection()

  if (rest.length < 2) {
    pushError(
      'Usage: pam add <query-name> "<sql>"\n\nExample: pam add list_users "SELECT * FROM employees"'
    )
    return
  }

  const name = rest[0]
  const sql = rest.slice(1).join(' ')
  const connQueries = savedQueries.value[activeConnection.value]

  if (connQueries[name]) {
    pushError(`Query "${name}" already exists. Use "pam remove ${name}" first.`)
    return
  }

  queryIdCounter++
  connQueries[name] = { id: queryIdCounter, sql }
  pushSuccess(`✓ Query saved: ${name} (id: ${queryIdCounter})\n  SQL: ${sql}`)
}

function cmdRemove(rest) {
  requireConnection()

  if (rest.length < 1) {
    pushError('Usage: pam remove <query-name-or-id>')
    return
  }

  const selector = rest[0]
  const connQueries = savedQueries.value[activeConnection.value]

  // Try by name
  if (connQueries[selector]) {
    delete connQueries[selector]
    pushSuccess(`✓ Removed query "${selector}"`)
    return
  }

  // Try by id
  const id = parseInt(selector)
  if (!isNaN(id)) {
    for (const [name, q] of Object.entries(connQueries)) {
      if (q.id === id) {
        delete connQueries[name]
        pushSuccess(`✓ Removed query "${name}" (id: ${id})`)
        return
      }
    }
  }

  pushError(`Query "${selector}" not found.`)
}

function cmdRun(rest) {
  requireConnection()

  if (rest.length === 0) {
    pushError(
      'Usage: pam run <query-name-or-id>\n       pam run "<inline-sql>"'
    )
    return
  }

  const selector = rest.join(' ')
  let sql = null

  // Check if it looks like inline SQL
  const upper = selector.trimStart().toUpperCase()
  if (
    upper.startsWith('SELECT') ||
    upper.startsWith('INSERT') ||
    upper.startsWith('UPDATE') ||
    upper.startsWith('DELETE') ||
    upper.startsWith('CREATE') ||
    upper.startsWith('DROP') ||
    upper.startsWith('ALTER') ||
    upper.startsWith('WITH') ||
    upper.startsWith('PRAGMA')
  ) {
    sql = selector
  } else {
    // Look up saved query
    const connQueries = savedQueries.value[activeConnection.value]

    // By name
    if (connQueries[selector]) {
      sql = connQueries[selector].sql
      pushInfo(`Running query "${selector}": ${sql}`)
    } else {
      // By id
      const id = parseInt(selector)
      if (!isNaN(id)) {
        for (const [name, q] of Object.entries(connQueries)) {
          if (q.id === id) {
            sql = q.sql
            pushInfo(`Running query "${name}" (id: ${id}): ${sql}`)
            break
          }
        }
      }
    }

    if (!sql) {
      pushError(
        `Query "${selector}" not found. Use "pam list queries" to see saved queries.`
      )
      return
    }
  }

  executeAndDisplay(sql)
}

function cmdQuery(rest) {
  requireConnection()

  let tableName = ''
  let sqlParts = []

  for (let i = 0; i < rest.length; i++) {
    const arg = rest[i]
    if (arg.startsWith('--table=')) {
      tableName = arg.slice(8)
    } else if (arg === '--table' || arg === '-t') {
      tableName = rest[i + 1] || ''
      i++
    } else {
      sqlParts.push(arg)
    }
  }

  if (!tableName && sqlParts.length === 0) {
    pushError(
      'Usage: pam query --table=<table> [sql]\n\nExample: pam query --table=employees "SELECT * FROM employees WHERE salary > 60000"'
    )
    return
  }

  let sql = sqlParts.join(' ')
  if (!sql && tableName) {
    sql = `SELECT * FROM ${tableName}`
  }

  if (tableName) {
    pushInfo(`Table: ${tableName}`)
  }

  executeAndDisplay(sql)
}

function cmdTables(rest) {
  requireConnection()

  if (rest.length > 0) {
    // Query specific table
    const table = rest[0]
    const sql = `SELECT * FROM ${table}`
    pushInfo(`SELECT * FROM ${table}`)
    executeAndDisplay(sql)
    return
  }

  // List all tables
  const { result } = execSQL(
    "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
  )
  if (result.length > 0 && result[0].values.length > 0) {
    const tables = result[0].values.map(r => r[0])
    pushInfo('Tables in current database:')
    pushList(tables.map(t => ({ text: t, active: false })))
  } else {
    pushInfo('No tables found.')
  }
}

function cmdTableView(rest) {
  requireConnection()

  if (rest.length < 1) {
    pushError('Usage: pam table-view <table-name>\n       pam tv <table-name>')
    return
  }

  const table = rest[0]
  const { result, elapsed } = execSQL(`PRAGMA table_info(${table})`)

  if (!result || result.length === 0 || result[0].values.length === 0) {
    pushError(`Table "${table}" not found.`)
    return
  }

  const res = result[0]
  const columns = ['Column', 'Type', 'Nullable', 'Default', 'PK']
  const rows = res.values.map(row => ({
    Column: row[1],
    Type: row[2] || 'TEXT',
    Nullable: row[3] ? 'NOT NULL' : 'NULL',
    Default: row[4] !== null ? String(row[4]) : '',
    PK: row[5] ? '⚿ PK' : ''
  }))

  pushInfo(`Structure of "${table}":`)
  pushTable(columns, rows, elapsed)
}

function cmdList(rest) {
  const sub = (rest[0] || 'queries').toLowerCase()

  if (sub === 'connections') {
    cmdListConnections()
    return
  }

  if (sub === 'queries') {
    requireConnection()
    const connQueries = savedQueries.value[activeConnection.value]
    const entries = Object.entries(connQueries)
    const searchTerm = rest[1] || ''

    if (entries.length === 0) {
      pushInfo(
        `No saved queries for "${activeConnection.value}".\nUse "pam add <name> <sql>" to save one.`
      )
      return
    }

    const filtered = searchTerm
      ? entries.filter(
          ([name, q]) =>
            name.includes(searchTerm) ||
            q.sql.toLowerCase().includes(searchTerm.toLowerCase())
        )
      : entries

    if (filtered.length === 0) {
      pushInfo(`No queries matching "${searchTerm}".`)
      return
    }

    pushInfo(`Saved queries for "${activeConnection.value}":`)
    for (const [name, q] of filtered) {
      pushInfo(`  [${q.id}] ${name}\n      ${q.sql}`)
    }
    return
  }

  pushError('Usage: pam list [connections | queries] [search]')
}

function cmdListConnections() {
  const names = Object.keys(connections.value)
  if (names.length === 0) {
    pushInfo('No connections configured.\nUse "pam init <name>" to create one.')
    return
  }

  pushInfo('Connections:')
  pushList(
    names.map(n => ({
      text: `${n} (${connections.value[n].type})`,
      active: n === activeConnection.value
    }))
  )
}

function cmdInfo(rest) {
  requireConnection()

  const sub = (rest[0] || 'tables').toLowerCase()

  if (sub === 'tables') {
    const { result, elapsed } = execSQL(
      "SELECT name, type FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
    )
    if (result.length > 0 && result[0].values.length > 0) {
      const columns = ['Name', 'Type']
      const rows = result[0].values.map(r => ({ Name: r[0], Type: r[1] }))
      pushTable(columns, rows, elapsed)
    } else {
      pushInfo('No tables found.')
    }
  } else if (sub === 'views') {
    const { result, elapsed } = execSQL(
      "SELECT name, type FROM sqlite_master WHERE type='view' ORDER BY name"
    )
    if (result.length > 0 && result[0].values.length > 0) {
      const columns = ['Name', 'Type']
      const rows = result[0].values.map(r => ({ Name: r[0], Type: r[1] }))
      pushTable(columns, rows, elapsed)
    } else {
      pushInfo('No views found.')
    }
  } else {
    pushError('Usage: pam info <tables|views>')
  }
}

function cmdExplain(rest) {
  requireConnection()

  if (rest.length < 1) {
    pushError('Usage: pam explain <table-name>')
    return
  }

  const table = rest[0]

  // Get foreign keys for this table (belongs to)
  const { result: fkResult } = execSQL(`PRAGMA foreign_key_list(${table})`)

  // Get tables that reference this table (has many)
  const { result: allTables } = execSQL(
    "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
  )

  let output = `Relationships for "${table}":\n`

  // Belongs to
  if (fkResult && fkResult.length > 0 && fkResult[0].values.length > 0) {
    output += '\n  belongs to [N:1]:\n'
    for (const row of fkResult[0].values) {
      const refTable = row[2]
      const from = row[3]
      const to = row[4]
      output += `    ${table}.${from} → ${refTable}.${to}\n`
    }
  } else {
    output += '\n  belongs to [N:1]: (none)\n'
  }

  // Has many
  let hasMany = []
  if (allTables && allTables.length > 0) {
    for (const trow of allTables[0].values) {
      const otherTable = trow[0]
      if (otherTable === table) continue
      try {
        const { result: otherFk } = execSQL(
          `PRAGMA foreign_key_list(${otherTable})`
        )
        if (otherFk && otherFk.length > 0) {
          for (const fk of otherFk[0].values) {
            if (fk[2] === table) {
              hasMany.push(`    ${otherTable}.${fk[3]} → ${table}.${fk[4]}`)
            }
          }
        }
      } catch {
        /* skip */
      }
    }
  }

  if (hasMany.length > 0) {
    output += '\n  has many [1:N]:\n' + hasMany.join('\n') + '\n'
  } else {
    output += '\n  has many [1:N]: (none)\n'
  }

  pushHelp(output)
}

function cmdDisconnect() {
  if (!activeConnection.value) {
    pushInfo('No active connection.')
    return
  }

  const prev = activeConnection.value
  activeConnection.value = ''
  updatePrompt()
  pushSuccess(`✓ Disconnected from "${prev}"`)
}

function cmdEdit(rest) {
  pushInfo(
    'Editor not available in the playground.\nUse "pam list queries" to view saved queries, or "pam add" / "pam remove" to manage them.'
  )
}

function cmdCompletion() {
  pushInfo(
    'Shell completion is not available in the playground.\nIn a real terminal, run:\n  pam completion bash --install\n  pam completion zsh --install\n  pam completion fish --install'
  )
}

function cmdHistory() {
  if (commandHistory.value.length === 0) {
    pushInfo('No command history yet.')
    return
  }

  pushInfo('Command history:')
  pushList(
    commandHistory.value.map((cmd, i) => ({
      text: `[${i + 1}] ${cmd}`,
      active: false
    }))
  )
}

function cmdDM() {
  const quotes = [
    '"I am running away from my responsibilities. And it feels good." — Michael Scott',
    '"Would I rather be feared or loved? Easy. Both. I want people to be afraid of how much they love me." — Michael Scott',
    '"Bears. Beets. Battlestar Galactica." — Dwight Schrute',
    '"I am fast. To give you a reference point I am somewhere between a snake and a mongoose. And a panther." — Dwight Schrute',
    '"I talk a lot, so I\'ve learned to tune myself out." — Kelly Kapoor',
    '"I want people to fear how much they love me." — Michael Scott',
    '"The worst thing about prison was the dementors." — Michael Scott',
    '"Perfectenschlag." — Michael Scott'
  ]
  const q = quotes[Math.floor(Math.random() * quotes.length)]
  pushInfo('📎 Dunder Mifflin Wisdom:\n\n  ' + q)
}

// ── SQL execution ───────────────────────────────────────────────────────────────────

function executeAndDisplay(sql) {
  const t0 = performance.now()

  // Handle multiple statements
  const statements = sql
    .split(';')
    .map(s => s.trim())
    .filter(s => s.length > 0)
  let lastResult = null

  for (const stmt of statements) {
    const upper = stmt.trimStart().toUpperCase()
    if (
      upper.startsWith('SELECT') ||
      upper.startsWith('WITH') ||
      upper.startsWith('PRAGMA')
    ) {
      lastResult = db.exec(stmt)
    } else {
      db.run(stmt)
      lastResult = null
    }
  }

  const t1 = performance.now()
  const elapsed = ((t1 - t0) / 1000).toFixed(3) + 's'

  if (lastResult && lastResult.length > 0) {
    const tableData = resultToTable(lastResult, elapsed)
    if (tableData) {
      pushTable(tableData.columns, tableData.rows, elapsed)
    } else {
      pushSuccess(`✓ Query executed (${elapsed})`)
    }
  } else {
    pushSuccess(`✓ Query executed (${elapsed})`)
  }
}

function requireConnection() {
  if (!activeConnection.value) {
    throw new Error(
      'No active connection. Use "pam switch <name>" or "pam init <name>" first.'
    )
  }
  if (!dbReady) {
    throw new Error('Database not ready. Please wait for initialization.')
  }
}

// ── Input handling ───────────────────────────────────────────────────────────────────

const allCommands = [
  'help',
  'init',
  'switch',
  'use',
  'status',
  'test',
  'add',
  'save',
  'remove',
  'delete',
  'run',
  'query',
  'list',
  'ls',
  'tables',
  't',
  'explore',
  'table-view',
  'tv',
  'info',
  'explain',
  'disconnect',
  'clear',
  'unset',
  'edit',
  'completion',
  'history',
  'dm'
]

function handleKeydown(e) {
  if (e.key === 'Tab') {
    e.preventDefault()
    const input = currentInput.value
    const parts = input.split(' ')
    if (parts.length === 1 && parts[0]) {
      const partial = parts[0].toLowerCase()
      const match = allCommands.find(
        c => c.startsWith(partial) && c !== partial
      )
      if (match) currentInput.value = match + ' '
    }
    return
  }
  if (e.key === 'Enter') {
    e.preventDefault()
    const input = currentInput.value.trim()
    currentInput.value = ''
    if (input) {
      handleCommand(input)
    }
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (historyIndex.value > 0) {
      historyIndex.value--
      currentInput.value = commandHistory.value[historyIndex.value]
    }
  } else if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (historyIndex.value < commandHistory.value.length - 1) {
      historyIndex.value++
      currentInput.value = commandHistory.value[historyIndex.value]
    } else {
      historyIndex.value = commandHistory.value.length
      currentInput.value = ''
    }
  } else if (e.key === 'l' && e.ctrlKey) {
    e.preventDefault()
    outputHistory.value = []
  }
}

function runQuickCommand(cmd) {
  currentInput.value = cmd
  handleCommand(cmd)
  currentInput.value = ''
  if (inputRef.value) inputRef.value.focus()
}

async function runWorkflow(wf) {
  if (inputRef.value) inputRef.value.focus()
  pushInfo(`\n\u25b8 Running workflow: ${wf.title}`)
  for (const cmd of wf.commands) {
    await new Promise(resolve => setTimeout(resolve, 120))
    handleCommand(cmd)
  }
}

// ── Initialization ──────────────────────────────────────────────────────────

onMounted(async () => {
  pushInfo(
    'Welcome to the PAM Playground 🗄️\n' +
      'Dunder Mifflin Paper Company \u2014 Scranton Branch\n' +
      '\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n' +
      'Simulates real PAM commands with SQLite WASM.\n' +
      'Data: employees \u00b7 departments \u00b7 projects \u00b7 employee_projects\n'
  )

  try {
    const initSqlJs = (await import('sql.js')).default
    const SQL = await initSqlJs({
      locateFile: file => import.meta.env.BASE_URL + file
    })
    db = new SQL.Database()

    // Seed sample data
    db.exec(`
      CREATE TABLE departments (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        budget REAL,
        location TEXT
      );

      INSERT INTO departments VALUES (1, 'Engineering', 500000, 'Floor 3');
      INSERT INTO departments VALUES (2, 'Sales', 300000, 'Floor 1');
      INSERT INTO departments VALUES (3, 'Marketing', 200000, 'Floor 2');
      INSERT INTO departments VALUES (4, 'HR', 150000, 'Floor 1');
      INSERT INTO departments VALUES (5, 'Finance', 250000, 'Floor 2');

      CREATE TABLE employees (
        id INTEGER PRIMARY KEY,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        email TEXT UNIQUE,
        title TEXT,
        salary REAL,
        department_id INTEGER REFERENCES departments(id),
        hire_date TEXT
      );

      INSERT INTO employees VALUES (1,  'Michael', 'Scott',    'michael@dundermifflin.com', 'Regional Manager',       95000,  2, '2000-03-15');
      INSERT INTO employees VALUES (2,  'Dwight',  'Schrute',  'dwight@dundermifflin.com',  'Assistant Regional Mgr', 78000,  2, '2001-06-01');
      INSERT INTO employees VALUES (3,  'Jim',     'Halpert',  'jim@dundermifflin.com',     'Sales Rep',              65000,  2, '2002-10-01');
      INSERT INTO employees VALUES (4,  'Pam',     'Beesly',   'pam@dundermifflin.com',     'Office Administrator',   42000,  4, '2000-03-15');
      INSERT INTO employees VALUES (5,  'Ryan',    'Howard',   'ryan@dundermifflin.com',    'Temp',                   30000,  2, '2004-05-20');
      INSERT INTO employees VALUES (6,  'Stanley', 'Hudson',   'stanley@dundermifflin.com', 'Sales Rep',              68000,  2, '1998-02-01');
      INSERT INTO employees VALUES (7,  'Kevin',   'Malone',   'kevin@dundermifflin.com',   'Accountant',             52000,  5, '2000-06-01');
      INSERT INTO employees VALUES (8,  'Angela',  'Martin',   'angela@dundermifflin.com',  'Senior Accountant',      62000,  5, '1999-11-01');
      INSERT INTO employees VALUES (9,  'Oscar',   'Martinez', 'oscar@dundermifflin.com',   'Accountant',             58000,  5, '2000-01-15');
      INSERT INTO employees VALUES (10, 'Toby',    'Flenderson','toby@dundermifflin.com',   'HR Representative',      55000,  4, '1999-06-01');
      INSERT INTO employees VALUES (11, 'Kelly',   'Kapoor',   'kelly@dundermifflin.com',   'Customer Service Rep',   44000,  3, '2001-09-01');
      INSERT INTO employees VALUES (12, 'Meredith','Palmer',   'meredith@dundermifflin.com','Supplier Relations',     48000,  2, '1997-03-01');
      INSERT INTO employees VALUES (13, 'Creed',   'Bratton',  'creed@dundermifflin.com',   'Quality Assurance',      46000,  1, '1996-01-01');
      INSERT INTO employees VALUES (14, 'Darryl',  'Philbin',  'darryl@dundermifflin.com',  'Warehouse Foreman',      54000,  1, '1999-08-01');
      INSERT INTO employees VALUES (15, 'Andy',    'Bernard',  'andy@dundermifflin.com',     'Sales Rep',              63000,  2, '2006-01-15');

      CREATE TABLE projects (
        id INTEGER PRIMARY KEY,
        name TEXT NOT NULL,
        status TEXT CHECK(status IN ('active', 'completed', 'planned')),
        department_id INTEGER REFERENCES departments(id),
        start_date TEXT
      );

      INSERT INTO projects VALUES (1, 'Website Redesign',    'active',    1, '2024-01-10');
      INSERT INTO projects VALUES (2, 'Q4 Sales Campaign',   'completed', 2, '2023-10-01');
      INSERT INTO projects VALUES (3, 'Brand Refresh',       'active',    3, '2024-03-01');
      INSERT INTO projects VALUES (4, 'Employee Onboarding', 'planned',   4, '2024-06-01');
      INSERT INTO projects VALUES (5, 'Budget Audit',        'active',    5, '2024-02-15');

      CREATE TABLE employee_projects (
        employee_id INTEGER REFERENCES employees(id),
        project_id INTEGER REFERENCES projects(id),
        role TEXT,
        PRIMARY KEY (employee_id, project_id)
      );

      INSERT INTO employee_projects VALUES (1, 2, 'Sponsor');
      INSERT INTO employee_projects VALUES (2, 2, 'Lead');
      INSERT INTO employee_projects VALUES (3, 2, 'Member');
      INSERT INTO employee_projects VALUES (13, 1, 'Lead');
      INSERT INTO employee_projects VALUES (14, 1, 'Member');
      INSERT INTO employee_projects VALUES (11, 3, 'Lead');
      INSERT INTO employee_projects VALUES (4, 4, 'Lead');
      INSERT INTO employee_projects VALUES (10, 4, 'Member');
      INSERT INTO employee_projects VALUES (7, 5, 'Member');
      INSERT INTO employee_projects VALUES (8, 5, 'Lead');
      INSERT INTO employee_projects VALUES (9, 5, 'Member');
    `)

    dbReady = true

    // Create default "playground" connection
    connections.value['playground'] = { type: 'sqlite', conn: ':memory:' }
    queryIdCounter = 4
    savedQueries.value['playground'] = {
      top_earners: {
        id: 1,
        sql: 'SELECT first_name, last_name, title, salary FROM employees ORDER BY salary DESC LIMIT 5'
      },
      by_dept: {
        id: 2,
        sql: 'SELECT d.name AS department, COUNT(e.id) AS headcount, ROUND(AVG(e.salary), 0) AS avg_salary FROM departments d JOIN employees e ON e.department_id = d.id GROUP BY d.id ORDER BY avg_salary DESC'
      },
      active_projects: {
        id: 3,
        sql: "SELECT p.name AS project, p.status, d.name AS department, p.start_date FROM projects p JOIN departments d ON p.department_id = d.id WHERE p.status = 'active' ORDER BY p.start_date"
      },
      recent_hires: {
        id: 4,
        sql: 'SELECT first_name, last_name, title, hire_date FROM employees ORDER BY hire_date DESC LIMIT 6'
      }
    }
    activeConnection.value = 'playground'
    updatePrompt()

    pushSuccess(
      '\u2713 SQLite WASM loaded \u2014 connection "playground" ready (4 tables, 4 pre-saved queries).'
    )
    pushInfo(
      'Pre-loaded queries: top_earners \u00b7 by_dept \u00b7 active_projects \u00b7 recent_hires\n' +
        'Try: run top_earners   or   tables   or   help\n' +
        'Tip: press Tab to autocomplete commands  |  Ctrl+L or "clear" to reset'
    )
  } catch (e) {
    pushError('Failed to load SQLite: ' + e.message)
  }

  // Focus input
  await nextTick()
  if (inputRef.value) inputRef.value.focus()
})
</script>

<style scoped>
.playground-page {
  padding-top: 80px;
}
.page-inner {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 3rem 1.5rem 4rem;
}

.page-header {
  text-align: center;
  margin-bottom: 2.5rem;
}
.page-header h1 {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
}
.gradient-text {
  background: var(--gradient-hero);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.lead {
  color: var(--text-secondary);
  font-size: 1rem;
}

/* Terminal frame */
.terminal-frame {
  border: 2px dashed var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
  margin-bottom: 2.5rem;
  background: var(--bg-card);
}
.terminal-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: var(--paper-accent);
  border-bottom: 1px dashed var(--border);
}
.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}
.dot.red {
  background: #ff5f57;
}
.dot.yellow {
  background: #febc2e;
}
.dot.green {
  background: #28c840;
}
.terminal-title {
  flex: 1;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.75rem;
  font-family: var(--font-mono);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}
.terminal-body {
  background: var(--paper-muted);
  padding: 1rem;
  max-height: 600px;
  overflow-y: auto;
  cursor: text;
}
.terminal-body::-webkit-scrollbar {
  width: 6px;
}
.terminal-body::-webkit-scrollbar-track {
  background: transparent;
}
.terminal-body::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 3px;
}

/* Output entries */
.output-entry {
  margin-bottom: 0.5rem;
}
.output-command {
  font-family: var(--font-mono);
  font-size: 0.875rem;
  line-height: 1.6;
}
.prompt {
  color: var(--green);
  font-family: var(--font-mono);
  font-size: 0.875rem;
  white-space: nowrap;
  user-select: none;
}
.command-text {
  color: var(--text);
}
.output-info {
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 0.85rem;
  white-space: pre-wrap;
  line-height: 1.5;
  padding: 0.15rem 0;
}
.output-success {
  color: var(--green);
  font-family: var(--font-mono);
  font-size: 0.85rem;
  white-space: pre-wrap;
  line-height: 1.5;
}
.output-error {
  color: var(--red);
  font-family: var(--font-mono);
  font-size: 0.85rem;
  white-space: pre-wrap;
  line-height: 1.5;
  padding: 0.4rem 0.75rem;
  background: rgba(197, 48, 48, 0.08);
  border-left: 3px solid var(--red);
  border-radius: 0 4px 4px 0;
  margin: 0.25rem 0;
}
.output-help {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  line-height: 1.5;
  color: var(--text-secondary);
}
.output-help pre {
  margin: 0;
  background: none;
  border: none;
  padding: 0;
  font-size: inherit;
  line-height: inherit;
  color: inherit;
  white-space: pre-wrap;
}

/* Tables */
.output-table-wrap {
  margin: 0.25rem 0;
  overflow-x: auto;
  border: 1px dashed var(--border);
  border-radius: 6px;
  background: var(--paper-muted);
}
.output-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8rem;
  font-family: var(--font-mono);
}
.output-table th {
  background: var(--paper-accent);
  padding: 0.4rem 0.65rem;
  text-align: left;
  font-weight: 700;
  color: var(--accent);
  border-bottom: 1px dashed var(--border);
  white-space: nowrap;
  position: sticky;
  top: 0;
}
.output-table td {
  padding: 0.3rem 0.65rem;
  border-bottom: 1px dashed rgba(214, 201, 184, 0.7);
  color: var(--text-secondary);
  white-space: nowrap;
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.output-table tbody tr:hover {
  background: var(--accent-soft);
}
.row-num {
  color: var(--text-muted);
  font-size: 0.7rem;
  width: 35px;
  text-align: center;
}
.null-val {
  color: var(--text-muted);
  font-style: italic;
}
.table-meta {
  padding: 0.3rem 0.65rem;
  font-size: 0.75rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
  border-top: 1px dashed var(--border);
  background: var(--paper-accent);
}

/* Lists */
.output-list {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  padding: 0.25rem 0;
}
.list-item {
  padding: 0.15rem 0 0.15rem 1rem;
  color: var(--text-secondary);
}
.list-item.active {
  color: var(--green);
  font-weight: 600;
}
.active-marker {
  color: var(--green);
}

/* Input line */
.terminal-input-line {
  display: flex;
  align-items: center;
  gap: 0;
  margin-top: 0.25rem;
  padding-top: 0.25rem;
}
.terminal-input {
  flex: 1;
  background: none;
  border: none;
  color: var(--text);
  font-family: var(--font-mono);
  font-size: 0.875rem;
  outline: none;
  caret-color: var(--green);
  line-height: 1.6;
}
.terminal-input::placeholder {
  color: var(--text-muted);
  opacity: 0.5;
}

/* Reference card */
.reference-card {
  margin-bottom: 2.5rem;
}
.ref-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.25rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px dashed var(--border);
}
.ref-icon {
  font-size: 1.2rem;
}
.ref-header h3 {
  font-size: 1rem;
  font-weight: 700;
  font-family: var(--font-mono);
  flex: 1;
}
.ref-badge {
  font-family: var(--font-mono);
  font-size: 0.62rem;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  color: var(--text-muted);
  border: 1px dashed var(--border);
  border-radius: 4px;
  padding: 2px 8px;
}
.ref-categories {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.25rem;
}
.ref-category {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.ref-cat-label {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: var(--text-muted);
  padding-bottom: 0.35rem;
  border-bottom: 1px dashed var(--border);
}
.ref-chips {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}
.command-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.command-chip {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.2rem;
  padding: 0.5rem 0.9rem;
  background: var(--bg-card);
  border: 1px dashed var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition:
    border-color 0.15s,
    transform 0.15s;
  text-align: left;
  font-family: var(--font-mono);
  box-shadow: var(--shadow-soft);
}
.command-chip:hover {
  border-color: var(--accent);
  transform: translateY(-1px);
}
.command-chip code {
  font-size: 0.75rem;
  background: none;
  padding: 0;
  color: var(--accent);
}
.chip-desc {
  font-size: 0.7rem;
  color: var(--text-muted);
}

/* Info cards */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.25rem;
}
.info-card {
  background: var(--bg-card);
  border: 1px dashed var(--border);
  border-radius: var(--radius);
  padding: 1.5rem;
  transition: border-color 0.25s;
  box-shadow: var(--shadow-soft);
}
.info-card:hover {
  border-color: var(--accent);
}
.info-icon {
  font-size: 1.5rem;
  display: block;
  margin-bottom: 0.5rem;
}
.info-card h3 {
  font-size: 1rem;
  margin-bottom: 0.3rem;
}
.info-card p {
  font-size: 0.9rem;
  color: var(--text-secondary);
}
.info-card code {
  font-size: 0.8rem;
}

@media (max-width: 640px) {
  .terminal-body {
    max-height: 450px;
    padding: 0.75rem;
  }
  .command-chip {
    width: 100%;
  }
}

/* Workflows */
.workflows-section {
  margin-bottom: 2.5rem;
}
.workflows-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.25rem;
}
.workflows-sub {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  color: var(--text-muted);
}
.workflow-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}
.workflow-card {
  background: var(--bg-card);
  border: 1px dashed var(--border);
  border-radius: var(--radius);
  padding: 1.25rem;
  cursor: pointer;
  transition:
    border-color 0.2s,
    box-shadow 0.2s,
    transform 0.15s;
  box-shadow: var(--shadow-soft);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.workflow-card:hover {
  border-color: var(--accent);
  box-shadow: var(--shadow);
  transform: translateY(-2px);
}
.wf-header {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}
.wf-num {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  font-weight: 700;
  color: var(--text-muted);
  background: var(--bg-code);
  border: 1px dashed var(--border);
  border-radius: 4px;
  padding: 2px 6px;
}
.wf-title {
  font-weight: 700;
  font-size: 0.9rem;
  font-family: var(--font-mono);
}
.wf-desc {
  font-size: 0.82rem;
  color: var(--text-secondary);
  line-height: 1.45;
}
.wf-cmds {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
  margin-top: 0.25rem;
}
.wf-cmds code {
  font-size: 0.68rem;
  padding: 1px 6px;
}
.wf-run {
  font-family: var(--font-mono);
  font-size: 0.72rem;
  color: var(--accent);
  margin-top: auto;
  padding-top: 0.25rem;
}

/* About section */
.about-section {
  margin-bottom: 2.5rem;
}
.about-title {
  font-size: 1rem;
  font-weight: 700;
  font-family: var(--font-mono);
  margin-bottom: 1rem;
}
.about-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.25rem;
  margin-bottom: 1rem;
}
.about-note {
  font-size: 0.82rem;
  line-height: 1.55;
  padding: 1rem 1rem 0.9rem;
  margin-top: 0.5rem;
}
.about-note code {
  font-size: 0.75rem;
}
.about-hint {
  margin-top: 0.5rem;
  font-size: 0.75rem;
}
</style>
