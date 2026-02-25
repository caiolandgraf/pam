<template>
  <div class="playground-page">
    <div class="page-inner">
      <div class="page-header">
        <h1>
          <span class="gradient-text">Playground</span>
        </h1>
        <p class="lead">
          Try SQL queries in your browser — powered by SQLite (WASM).
          <br />No server needed. Everything runs locally.
        </p>
      </div>

      <!-- Terminal-style frame -->
      <div class="terminal-frame">
        <div class="terminal-bar">
          <span class="dot red"></span>
          <span class="dot yellow"></span>
          <span class="dot green"></span>
          <span class="terminal-title">pam playground — SQLite in-browser</span>
        </div>

        <div class="terminal-body">
          <!-- Status -->
          <div class="status-bar">
            <span
              class="status-indicator"
              :class="dbReady ? 'ready' : 'loading'"
            >
              {{ dbReady ? '● connected' : '○ loading...' }}
            </span>
            <span class="status-info" v-if="dbReady"
              >SQLite WASM · {{ tableCount }} tables · sample data loaded</span
            >
          </div>

          <!-- SQL Editor -->
          <div class="editor-area">
            <div class="editor-header">
              <span class="editor-label">SQL</span>
              <div class="editor-actions">
                <select
                  v-model="selectedExample"
                  @change="loadExample"
                  class="example-select"
                >
                  <option value="">-- examples --</option>
                  <option
                    v-for="ex in examples"
                    :key="ex.name"
                    :value="ex.name"
                  >
                    {{ ex.name }}
                  </option>
                </select>
                <button
                  class="run-btn"
                  @click="runQuery"
                  :disabled="!dbReady || running"
                >
                  <span v-if="running" class="spinner"></span>
                  <span v-else>▶</span>
                  {{ running ? 'Running...' : 'Run' }}
                  <kbd>{{ metaKey }}+Enter</kbd>
                </button>
              </div>
            </div>
            <textarea
              ref="editorRef"
              v-model="sql"
              class="sql-editor"
              spellcheck="false"
              :rows="editorRows"
              @keydown="handleEditorKey"
              placeholder="-- Write your SQL here..."
            ></textarea>
          </div>

          <!-- Results -->
          <div class="results-area" v-if="queryRan">
            <!-- Error -->
            <div v-if="errorMsg" class="result-error">
              <span class="error-icon">✗</span>
              {{ errorMsg }}
            </div>

            <!-- Success non-select -->
            <div
              v-else-if="!resultColumns.length && !errorMsg"
              class="result-success"
            >
              <span class="success-icon">✓</span>
              Command executed successfully
              <span class="result-time">{{ elapsed }}</span>
            </div>

            <!-- Table -->
            <div v-else class="result-table-wrap">
              <div class="result-meta">
                <span
                  >{{ resultRows.length }} row{{
                    resultRows.length !== 1 ? 's' : ''
                  }}</span
                >
                <span
                  >{{ resultColumns.length }} column{{
                    resultColumns.length !== 1 ? 's' : ''
                  }}</span
                >
                <span class="result-time">{{ elapsed }}</span>
              </div>
              <div class="table-scroll">
                <table class="result-table">
                  <thead>
                    <tr>
                      <th class="row-num">#</th>
                      <th v-for="col in resultColumns" :key="col">{{ col }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(row, i) in resultRows" :key="i">
                      <td class="row-num">{{ i + 1 }}</td>
                      <td
                        v-for="col in resultColumns"
                        :key="col"
                        :class="{ 'null-val': row[col] === null }"
                      >
                        {{ row[col] === null ? 'NULL' : row[col] }}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>

          <!-- History -->
          <div class="history-area" v-if="history.length > 0">
            <div class="history-header" @click="historyOpen = !historyOpen">
              <span>History ({{ history.length }})</span>
              <span class="chevron" :class="{ open: historyOpen }">›</span>
            </div>
            <div class="history-list" v-if="historyOpen">
              <button
                v-for="(h, i) in history"
                :key="i"
                class="history-item"
                @click="sql = h.sql"
              >
                <code>{{
                  h.sql.length > 80 ? h.sql.slice(0, 80) + '...' : h.sql
                }}</code>
                <span class="history-time">{{ h.time }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Info cards -->
      <div class="info-grid">
        <div class="info-card">
          <span class="info-icon">🖥️</span>
          <h3>In-Browser</h3>
          <p>
            Uses sql.js (SQLite compiled to WASM). No data leaves your browser.
          </p>
        </div>
        <div class="info-card">
          <span class="info-icon">📊</span>
          <h3>Sample Data</h3>
          <p>
            Pre-loaded with employees, departments and projects tables to
            experiment with.
          </p>
        </div>
        <div class="info-card">
          <span class="info-icon">⚡</span>
          <h3>Full SQL</h3>
          <p>
            Supports SELECT, INSERT, UPDATE, DELETE, CREATE TABLE, JOINs and
            more.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const dbReady = ref(false)
const running = ref(false)
const queryRan = ref(false)
const errorMsg = ref('')
const resultColumns = ref([])
const resultRows = ref([])
const elapsed = ref('')
const sql = ref('')
const selectedExample = ref('')
const history = ref([])
const historyOpen = ref(false)
const editorRef = ref(null)

const metaKey = navigator.platform.includes('Mac') ? '⌘' : 'Ctrl'

let db = null

const tableCount = ref(0)

const editorRows = computed(() => {
  const lines = sql.value.split('\n').length
  return Math.max(4, Math.min(lines + 1, 16))
})

const examples = [
  { name: 'All employees', sql: 'SELECT * FROM employees;' },
  {
    name: 'JOIN departments',
    sql: `SELECT e.first_name, e.last_name, d.name AS department, e.salary
FROM employees e
JOIN departments d ON e.department_id = d.id
ORDER BY d.name, e.last_name;`
  },
  {
    name: 'Salary stats',
    sql: `SELECT
  d.name AS department,
  COUNT(*) AS headcount,
  ROUND(AVG(e.salary), 2) AS avg_salary,
  MIN(e.salary) AS min_salary,
  MAX(e.salary) AS max_salary
FROM employees e
JOIN departments d ON e.department_id = d.id
GROUP BY d.name
ORDER BY avg_salary DESC;`
  },
  {
    name: 'Projects overview',
    sql: `SELECT
  p.name AS project,
  p.status,
  d.name AS department,
  COUNT(ep.employee_id) AS team_size
FROM projects p
JOIN departments d ON p.department_id = d.id
LEFT JOIN employee_projects ep ON ep.project_id = p.id
GROUP BY p.id
ORDER BY team_size DESC;`
  },
  {
    name: 'Top earners',
    sql: `SELECT first_name || ' ' || last_name AS name, title, salary
FROM employees
WHERE salary > 80000
ORDER BY salary DESC
LIMIT 10;`
  },
  {
    name: 'Create a table',
    sql: `CREATE TABLE notes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  body TEXT,
  created_at TEXT DEFAULT (datetime('now'))
);

INSERT INTO notes (title, body) VALUES ('Hello', 'My first note');
INSERT INTO notes (title, body) VALUES ('Pam', 'Database drawer is great');

SELECT * FROM notes;`
  },
  {
    name: 'Schema info',
    sql: `SELECT name, type, sql
FROM sqlite_master
WHERE type IN ('table', 'view')
ORDER BY type, name;`
  }
]

function loadExample() {
  const ex = examples.find(e => e.name === selectedExample.value)
  if (ex) sql.value = ex.sql
}

function handleEditorKey(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
    e.preventDefault()
    runQuery()
  }
  // Tab inserts spaces
  if (e.key === 'Tab') {
    e.preventDefault()
    const el = editorRef.value
    const start = el.selectionStart
    const end = el.selectionEnd
    sql.value = sql.value.substring(0, start) + '  ' + sql.value.substring(end)
    requestAnimationFrame(() => {
      el.selectionStart = el.selectionEnd = start + 2
    })
  }
}

async function runQuery() {
  if (!db || !sql.value.trim()) return
  running.value = true
  errorMsg.value = ''
  resultColumns.value = []
  resultRows.value = []

  const t0 = performance.now()

  try {
    // Split by ; to handle multi-statement
    const statements = sql.value
      .split(';')
      .map(s => s.trim())
      .filter(s => s.length > 0)

    let lastResult = null

    for (const stmt of statements) {
      const upper = stmt.toUpperCase()
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
    elapsed.value = ((t1 - t0) / 1000).toFixed(3) + 's'

    if (lastResult && lastResult.length > 0) {
      const res = lastResult[0]
      resultColumns.value = res.columns
      resultRows.value = res.values.map(row => {
        const obj = {}
        res.columns.forEach((col, i) => {
          obj[col] = row[i]
        })
        return obj
      })
    }

    queryRan.value = true

    // Update table count
    try {
      const tc = db.exec(
        "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
      )
      if (tc.length) tableCount.value = tc[0].values[0][0]
    } catch {
      /* ignore */
    }

    history.value.unshift({
      sql: sql.value.trim(),
      time: new Date().toLocaleTimeString()
    })
    if (history.value.length > 30) history.value.pop()
  } catch (e) {
    const t1 = performance.now()
    elapsed.value = ((t1 - t0) / 1000).toFixed(3) + 's'
    errorMsg.value = e.message
    queryRan.value = true
  } finally {
    running.value = false
  }
}

onMounted(async () => {
  try {
    const initSqlJs = (await import('sql.js')).default
    const SQL = await initSqlJs({
      locateFile: file => import.meta.env.BASE_URL + file
    })
    db = new SQL.Database()

    // Seed sample data
    db.run(`
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

    const tc = db.exec(
      "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
    )
    if (tc.length) tableCount.value = tc[0].values[0][0]

    dbReady.value = true
    sql.value = examples[0].sql
  } catch (e) {
    errorMsg.value = 'Failed to load SQLite: ' + e.message
    queryRan.value = true
  }
})
</script>

<style scoped>
.playground-page {
  padding-top: 64px;
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
  font-size: 2.2rem;
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
  font-size: 1.05rem;
}

/* Terminal frame */
.terminal-frame {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
  margin-bottom: 3rem;
}
.terminal-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 14px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border);
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
  font-size: 0.8rem;
  font-family: var(--font-mono);
}
.terminal-body {
  background: var(--bg);
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

/* Status */
.status-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.8rem;
  padding: 0.4rem 0;
}
.status-indicator {
  font-family: var(--font-mono);
}
.status-indicator.ready {
  color: var(--green);
}
.status-indicator.loading {
  color: var(--yellow);
}
.status-info {
  color: var(--text-muted);
}

/* Editor */
.editor-area {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
}
.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem;
  background: rgba(110, 118, 129, 0.08);
  border-bottom: 1px solid var(--border);
}
.editor-label {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.editor-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}
.example-select {
  background: var(--bg-code);
  border: 1px solid var(--border);
  color: var(--text-secondary);
  font-size: 0.8rem;
  font-family: var(--font-sans);
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}
.example-select:focus {
  outline: 1px solid var(--accent);
}
.run-btn {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  background: var(--green-dim);
  color: var(--text);
  border: none;
  padding: 5px 14px;
  border-radius: 6px;
  font-size: 0.85rem;
  font-family: var(--font-sans);
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s;
}
.run-btn:hover:not(:disabled) {
  background: var(--green);
}
.run-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.run-btn kbd {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  background: rgba(255, 255, 255, 0.1);
  padding: 1px 5px;
  border-radius: 3px;
  margin-left: 4px;
}
.spinner {
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.sql-editor {
  width: 100%;
  background: var(--bg-code);
  color: var(--text);
  border: none;
  padding: 1rem;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  line-height: 1.6;
  resize: vertical;
  outline: none;
  min-height: 100px;
}
.sql-editor::placeholder {
  color: var(--text-muted);
}

/* Results */
.results-area {
  margin-top: 0.25rem;
}
.result-error {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(248, 81, 73, 0.1);
  border: 1px solid rgba(248, 81, 73, 0.3);
  border-radius: var(--radius-sm);
  color: var(--red);
  font-size: 0.875rem;
  font-family: var(--font-mono);
}
.error-icon {
  font-weight: 700;
}
.result-success {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(63, 185, 80, 0.1);
  border: 1px solid rgba(63, 185, 80, 0.3);
  border-radius: var(--radius-sm);
  color: var(--green);
  font-size: 0.875rem;
}
.success-icon {
  font-weight: 700;
}
.result-time {
  margin-left: auto;
  font-size: 0.8rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
}

/* Table */
.result-meta {
  display: flex;
  gap: 1rem;
  padding: 0.5rem 0;
  font-size: 0.8rem;
  color: var(--text-muted);
}
.table-scroll {
  overflow-x: auto;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}
.result-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
  font-family: var(--font-mono);
}
.result-table th {
  background: var(--bg-card);
  padding: 0.5rem 0.75rem;
  text-align: left;
  font-weight: 600;
  color: var(--accent);
  border-bottom: 1px solid var(--border);
  white-space: nowrap;
  position: sticky;
  top: 0;
}
.result-table td {
  padding: 0.4rem 0.75rem;
  border-bottom: 1px solid rgba(48, 54, 61, 0.5);
  color: var(--text-secondary);
  white-space: nowrap;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.result-table tbody tr:hover {
  background: rgba(88, 166, 255, 0.05);
}
.row-num {
  color: var(--text-muted);
  font-size: 0.75rem;
  width: 40px;
  text-align: center;
}
.null-val {
  color: var(--text-muted);
  font-style: italic;
}

/* History */
.history-area {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
}
.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  font-size: 0.8rem;
  color: var(--text-muted);
  background: rgba(110, 118, 129, 0.05);
}
.history-header:hover {
  background: rgba(110, 118, 129, 0.1);
}
.chevron {
  transition: transform 0.2s;
  font-size: 1rem;
}
.chevron.open {
  transform: rotate(90deg);
}
.history-list {
  display: flex;
  flex-direction: column;
  max-height: 200px;
  overflow-y: auto;
}
.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 0.4rem 0.75rem;
  border: none;
  background: none;
  cursor: pointer;
  text-align: left;
  font-family: var(--font-sans);
  border-top: 1px solid rgba(48, 54, 61, 0.4);
  transition: background 0.1s;
}
.history-item:hover {
  background: rgba(88, 166, 255, 0.05);
}
.history-item code {
  font-size: 0.8rem;
  color: var(--text-secondary);
  background: none;
  padding: 0;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.history-time {
  font-size: 0.7rem;
  color: var(--text-muted);
  flex-shrink: 0;
}

/* Info cards */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.25rem;
}
.info-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 1.5rem;
  transition: border-color 0.25s;
}
.info-card:hover {
  border-color: var(--border-accent);
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
</style>
