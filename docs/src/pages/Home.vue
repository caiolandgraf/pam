<template>
  <div class="home">
    <!-- Hero -->
    <section class="hero">
      <div class="hero-inner">
        <div class="memo-header hero-memo">
          <div class="memo-kv">
            <div><span class="memo-label">To:</span> Scranton Branch</div>
            <div><span class="memo-label">From:</span> Reception (Pam)</div>
            <div>
              <span class="memo-label">Subject:</span> Keep SQL organized
            </div>
          </div>
          <div class="memo-stamp">FILED</div>
        </div>
        <div class="hero-badge">interoffice memo • open-source CLI</div>
        <h1 class="hero-title">
          <span class="gradient-text">Pam's</span> Database Drawer
        </h1>
        <p class="hero-sub">
          A minimal CLI tool for managing and executing SQL queries across
          multiple databases. Written in Go, made beautiful with
          <a href="https://github.com/charmbracelet/bubbletea" target="_blank"
            >BubbleTea</a
          >.
        </p>
        <p class="hero-quote">
          <em>"Pam, the receptionist, has been doing a fantastic job."</em>
          — Michael Scott
        </p>
        <div class="hero-labels">
          <span class="label-chip">Form DM-01</span>
          <span class="label-chip">Status: Approved</span>
          <span class="label-chip">Filed: Reception</span>
        </div>
        <div class="hero-actions">
          <router-link to="/playground" class="btn btn-primary">
            Try the Playground
          </router-link>
          <router-link to="/docs" class="btn btn-secondary">
            Read the Docs
          </router-link>
        </div>
        <div class="hero-install sticky-note">
          <strong>Quick Install</strong>
          <code>go install github.com/caiolandgraf/pam/cmd/pam@latest</code>
        </div>
      </div>
    </section>

    <!-- Demo GIF -->
    <section class="demo section">
      <div class="section-inner">
        <div class="demo-frame">
          <div class="demo-dots">
            <span class="dot red"></span>
            <span class="dot yellow"></span>
            <span class="dot green"></span>
            <span class="demo-title-bar">pam — terminal</span>
          </div>
          <img
            src="https://github.com/user-attachments/assets/b62bec1d-2255-4d02-9b7f-1c99afbeb664"
            alt="Pam demo"
            class="demo-img"
            loading="lazy"
          />
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="features section">
      <div class="section-inner">
        <h2 class="section-title">Office Highlights</h2>
        <div class="features-grid">
          <div class="feature-card" v-for="f in features" :key="f.title">
            <span class="feature-icon">{{ f.icon }}</span>
            <h3>{{ f.title }}</h3>
            <p>{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- DB Support -->
    <section class="databases section">
      <div class="section-inner">
        <h2 class="section-title">Branch Support</h2>
        <p class="section-sub">
          Connect to all major databases with a single tool
        </p>
        <div class="db-grid">
          <div class="db-chip" v-for="db in databases" :key="db">{{ db }}</div>
        </div>
      </div>
    </section>

    <!-- Quick Start -->
    <section class="quickstart section">
      <div class="section-inner">
        <h2 class="section-title">Quick Start Checklist</h2>
        <div class="steps">
          <div class="step" v-for="(s, i) in steps" :key="i">
            <div class="step-number">{{ i + 1 }}</div>
            <div class="step-content">
              <h3>{{ s.title }}</h3>
              <CodeBlock :title="'bash'">{{ s.code }}</CodeBlock>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- TUI Keys -->
    <section class="tui section">
      <div class="section-inner">
        <h2 class="section-title">Conference Room Controls</h2>
        <p class="section-sub">
          Vim-style keybindings for blazing-fast navigation
        </p>
        <div class="keys-grid">
          <div class="key-group" v-for="g in keyGroups" :key="g.title">
            <h4>{{ g.title }}</h4>
            <ul>
              <li v-for="k in g.keys" :key="k.key">
                <kbd>{{ k.key }}</kbd>
                <span>{{ k.desc }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="cta section">
      <div class="section-inner cta-inner">
        <h2>Ready to organize your queries, Scranton-style?</h2>
        <p>Get started in seconds — one command to install, one to connect.</p>
        <div class="hero-actions">
          <a
            href="https://github.com/caiolandgraf/pam"
            target="_blank"
            class="btn btn-primary"
          >
            View on GitHub
          </a>
          <router-link to="/docs" class="btn btn-secondary">
            Documentation
          </router-link>
          <router-link to="/playground" class="btn btn-secondary">
            Try Playground
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import CodeBlock from '../components/CodeBlock.vue'

const features = [
  {
    icon: '📂',
    title: 'Query Library',
    desc: 'Save and organize your most-used queries with names and IDs.'
  },
  {
    icon: '⚡',
    title: 'Runs in the CLI',
    desc: 'Execute queries with minimal overhead, right from your terminal.'
  },
  {
    icon: '🔌',
    title: 'Multi-Database',
    desc: 'PostgreSQL, MySQL, SQLite, Oracle, SQL Server, ClickHouse & Firebird.'
  },
  {
    icon: '🖥️',
    title: 'TUI Table Viewer',
    desc: 'Interactive results with Vim-style bindings and keyboard-driven navigation.'
  },
  {
    icon: '✏️',
    title: 'In-Place Editing',
    desc: 'Update cells, delete rows and edit SQL directly from the results table.'
  },
  {
    icon: '📤',
    title: 'Export Data',
    desc: 'Export as CSV, JSON, SQL, Markdown or HTML tables with a single key.'
  }
]

const databases = [
  'PostgreSQL',
  'MySQL / MariaDB',
  'SQLite',
  'Oracle',
  'SQL Server',
  'ClickHouse',
  'Firebird'
]

const steps = [
  {
    title: 'Create a connection',
    code: 'pam init mydb postgres "postgresql://user:pass@localhost:5432/mydb"'
  },
  { title: 'Save a query', code: 'pam add list_users "SELECT * FROM users"' },
  { title: 'Run it', code: 'pam run list_users' }
]

const keyGroups = [
  {
    title: 'Navigation',
    keys: [
      { key: 'h j k l', desc: 'Move left / down / up / right' },
      { key: 'g / G', desc: 'Jump to first / last row' },
      { key: 'Ctrl+u / Ctrl+d', desc: 'Page up / down' }
    ]
  },
  {
    title: 'Data',
    keys: [
      { key: 'y', desc: 'Copy cell to clipboard' },
      { key: 'v', desc: 'Visual selection mode' },
      { key: 'x', desc: 'Export selected data' }
    ]
  },
  {
    title: 'Editing',
    keys: [
      { key: 'u', desc: 'Update current cell' },
      { key: 'D', desc: 'Delete current row' },
      { key: 'e', desc: 'Edit & re-run query' }
    ]
  },
  {
    title: 'Other',
    keys: [
      { key: 'f', desc: 'Toggle sort on column' },
      { key: 's', desc: 'Save current query' },
      { key: 'q / Esc', desc: 'Quit' }
    ]
  }
]
</script>

<style scoped>
/* ---------- HERO ---------- */
.hero {
  padding: 10rem 1.5rem 4rem;
  text-align: center;
}
.hero-memo {
  margin-bottom: 1.5rem;
}
.hero-inner {
  max-width: 720px;
  margin: 0 auto;
}
.hero-badge {
  display: inline-block;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.2em;
  font-weight: 600;
  color: var(--text-secondary);
  border: 1px dashed var(--border);
  padding: 4px 14px;
  border-radius: 6px;
  margin-bottom: 1.5rem;
  font-family: var(--font-mono);
  background: var(--bg-card);
}
.hero-title {
  font-size: clamp(2.4rem, 6vw, 3.6rem);
  font-weight: 800;
  line-height: 1.1;
  letter-spacing: -0.02em;
  margin-bottom: 1.25rem;
}
.gradient-text {
  background: var(--gradient-hero);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.hero-sub {
  font-size: 1.05rem;
  color: var(--text-secondary);
  max-width: 560px;
  margin: 0 auto 1rem;
}
.hero-quote {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 2rem;
  font-family: var(--font-mono);
}
.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
  margin-bottom: 1.5rem;
}
.hero-labels {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}
.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.7rem 1.6rem;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 0.9rem;
  transition: all 0.2s;
  border: 1px solid var(--border);
  font-family: var(--font-mono);
  background: var(--bg-card);
}
.btn-primary {
  background: var(--accent);
  color: var(--bg);
  border-color: var(--accent);
  box-shadow: var(--shadow-soft);
}
.btn-primary:hover {
  background: var(--accent-hover);
  color: var(--bg);
}
.btn-secondary {
  border-color: var(--border);
  color: var(--text);
  background: var(--bg-card);
}
.btn-secondary:hover {
  border-color: var(--accent);
  color: var(--text);
}
.hero-install {
  display: inline-flex;
  flex-direction: column;
  gap: 0.4rem;
  align-items: center;
}
.hero-install code {
  font-size: 0.85rem;
  padding: 0.35rem 0.75rem;
  background: var(--paper-muted);
  border: 1px dashed var(--border);
  border-radius: var(--radius-sm);
  color: var(--text);
}

/* ---------- SECTIONS ---------- */
.section {
  padding: 4rem 1.5rem;
}
.section-inner {
  max-width: var(--max-width);
  margin: 0 auto;
}
.section-title {
  font-size: 1.6rem;
  font-weight: 700;
  text-align: center;
  margin: 0 auto 0.5rem;
  padding-bottom: 0.4rem;
  border-bottom: 2px dashed var(--border);
  display: table;
}
.section-sub {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: 2.5rem;
}

/* ---------- DEMO ---------- */
.demo-frame {
  border: 2px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow: var(--shadow);
  background: var(--bg-card);
}
.demo-dots {
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
.demo-title-bar {
  flex: 1;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.75rem;
  font-family: var(--font-mono);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}
.demo-img {
  width: 100%;
  display: block;
}

/* ---------- FEATURES ---------- */
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.25rem;
  margin-top: 2.5rem;
}
.feature-card {
  background: var(--bg-card);
  border: 1px dashed var(--border);
  border-radius: var(--radius);
  padding: 1.5rem;
  transition:
    border-color 0.25s,
    transform 0.25s;
  box-shadow: var(--shadow-soft);
}
.feature-card:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
}
.feature-icon {
  font-size: 1.6rem;
  display: block;
  margin-bottom: 0.75rem;
}
.feature-card h3 {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.4rem;
}
.feature-card p {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

/* ---------- DATABASES ---------- */
.db-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.75rem;
}
.db-chip {
  background: var(--bg-card);
  border: 1px dashed var(--border);
  padding: 0.55rem 1.3rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  transition: border-color 0.2s;
  font-family: var(--font-mono);
}
.db-chip:hover {
  border-color: var(--accent);
}

/* ---------- QUICK START ---------- */
.steps {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 700px;
  margin: 2rem auto 0;
}
.step {
  display: flex;
  gap: 1.25rem;
  align-items: flex-start;
}
.step-number {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.9rem;
  background: var(--paper-muted);
  color: var(--text);
  margin-top: 2px;
  border: 1px dashed var(--border);
  font-family: var(--font-mono);
}
.step-content {
  flex: 1;
  min-width: 0;
}
.step-content h3 {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

/* ---------- TUI KEYS ---------- */
.keys-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-top: 1rem;
}
.key-group {
  background: var(--bg-card);
  border: 1px dashed var(--border);
  border-radius: var(--radius);
  padding: 1.25rem;
  box-shadow: var(--shadow-soft);
}
.key-group h4 {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.18em;
  color: var(--text-muted);
  margin-bottom: 0.75rem;
  font-family: var(--font-mono);
}
.key-group ul {
  list-style: none;
}
.key-group li {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.35rem 0;
  font-size: 0.9rem;
  color: var(--text-secondary);
}
kbd {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  background: var(--bg-code);
  border: 1px dashed var(--border);
  padding: 2px 8px;
  border-radius: 4px;
  color: var(--text);
  white-space: nowrap;
  min-width: 80px;
  text-align: center;
}

/* ---------- CTA ---------- */
.cta-inner {
  text-align: center;
  padding: 3rem;
  background: var(--bg-card);
  border: 2px dashed var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow-soft);
}
.cta-inner h2 {
  font-size: 1.6rem;
  margin-bottom: 0.5rem;
}
.cta-inner p {
  color: var(--text-secondary);
  margin-bottom: 1.5rem;
}
</style>
