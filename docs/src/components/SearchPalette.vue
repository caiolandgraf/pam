<template>
  <Teleport to="body">
    <Transition name="overlay">
      <div v-if="open" class="search-overlay" @click.self="close">
        <div class="search-modal">
          <div class="search-header">
            <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="2" stroke-linecap="round">
              <circle cx="11" cy="11" r="8" />
              <path d="m21 21-4.35-4.35" />
            </svg>
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="Search docs, commands, pages..."
              class="search-input"
              @keydown.escape="close"
              @keydown.down.prevent="moveDown"
              @keydown.up.prevent="moveUp"
              @keydown.enter.prevent="goToSelected"
            />
            <kbd class="search-esc">Esc</kbd>
          </div>

          <div class="search-body" v-if="query.length > 0">
            <div v-if="results.length === 0" class="search-empty">
              No results for "<strong>{{ query }}</strong>"
            </div>
            <div v-else class="search-results" ref="resultsRef">
              <button
                v-for="(r, i) in results"
                :key="r.id"
                class="search-result"
                :class="{ active: i === selectedIndex }"
                @click="goTo(r)"
                @mouseenter="selectedIndex = i"
              >
                <span class="result-icon">{{ r.icon }}</span>
                <div class="result-text">
                  <span class="result-title" v-html="highlight(r.title)"></span>
                  <span class="result-section">{{ r.section }}</span>
                </div>
                <span class="result-type">{{ r.type }}</span>
              </button>
            </div>
          </div>

          <div class="search-body search-hints" v-else>
            <div class="hint-group">
              <span class="hint-label">Pages</span>
              <button class="search-result" @click="goTo({ route: '/' })">
                <span class="result-icon">🏠</span>
                <span class="result-title">Home</span>
              </button>
              <button class="search-result" @click="goTo({ route: '/docs' })">
                <span class="result-icon">📖</span>
                <span class="result-title">Documentation</span>
              </button>
              <button class="search-result" @click="goTo({ route: '/playground' })">
                <span class="result-icon">🎮</span>
                <span class="result-title">Playground</span>
              </button>
              <button class="search-result" @click="goTo({ route: '/contributors' })">
                <span class="result-icon">👥</span>
                <span class="result-title">Contributors</span>
              </button>
            </div>
          </div>

          <div class="search-footer">
            <span><kbd>↑↓</kbd> navigate</span>
            <span><kbd>↵</kbd> select</span>
            <span><kbd>Esc</kbd> close</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const open = ref(false)
const query = ref('')
const selectedIndex = ref(0)
const inputRef = ref(null)
const resultsRef = ref(null)
const router = useRouter()

// Search index — all searchable content
const searchIndex = [
  // Pages
  { id: 'p-home', title: 'Home', section: 'Pages', type: 'page', icon: '🏠', route: '/', keywords: 'home landing start' },
  { id: 'p-docs', title: 'Documentation', section: 'Pages', type: 'page', icon: '📖', route: '/docs', keywords: 'docs documentation help guide' },
  { id: 'p-play', title: 'Playground', section: 'Pages', type: 'page', icon: '🎮', route: '/playground', keywords: 'playground demo try sql sqlite interactive' },
  { id: 'p-contrib', title: 'Contributors', section: 'Pages', type: 'page', icon: '👥', route: '/contributors', keywords: 'contributors team community' },

  // Installation
  { id: 'd-install', title: 'Installation', section: 'Getting Started', type: 'docs', icon: '📦', route: '/docs#installation', keywords: 'install setup download binary go nix' },
  { id: 'd-go-install', title: 'Go Install', section: 'Installation', type: 'docs', icon: '📦', route: '/docs#go-install', keywords: 'go install gobin' },
  { id: 'd-build', title: 'Build Manually', section: 'Installation', type: 'docs', icon: '🔨', route: '/docs#build-manually', keywords: 'build compile clone manual' },
  { id: 'd-nix', title: 'Nix / NixOS', section: 'Installation', type: 'docs', icon: '❄️', route: '/docs#nix', keywords: 'nix nixos flake' },
  { id: 'd-releases', title: 'Binary Releases', section: 'Installation', type: 'docs', icon: '📋', route: '/docs#releases', keywords: 'releases download binary' },

  // Connections
  { id: 'd-conn', title: 'Connections', section: 'Getting Started', type: 'docs', icon: '🔌', route: '/docs#connections', keywords: 'connection database connect' },
  { id: 'd-init', title: 'pam init', section: 'Connections', type: 'command', icon: '⚡', route: '/docs#init', keywords: 'init create connection new database' },
  { id: 'd-switch', title: 'pam switch', section: 'Connections', type: 'command', icon: '🔄', route: '/docs#switch', keywords: 'switch use connection change' },
  { id: 'd-status', title: 'pam status', section: 'Connections', type: 'command', icon: '📊', route: '/docs#status', keywords: 'status test ping connection active' },
  { id: 'd-disconnect', title: 'pam disconnect', section: 'Connections', type: 'command', icon: '🔌', route: '/docs#disconnect', keywords: 'disconnect close clear unset' },

  // Databases
  { id: 'd-db', title: 'Database Support', section: 'Getting Started', type: 'docs', icon: '🗄️', route: '/docs#databases', keywords: 'database support postgres mysql sqlite oracle sqlserver clickhouse firebird' },
  { id: 'd-postgres', title: 'PostgreSQL', section: 'Databases', type: 'docs', icon: '🐘', route: '/docs#db-postgresql', keywords: 'postgres postgresql pg' },
  { id: 'd-mysql', title: 'MySQL / MariaDB', section: 'Databases', type: 'docs', icon: '🐬', route: '/docs#db-mysqlmariadb', keywords: 'mysql mariadb' },
  { id: 'd-sqlite', title: 'SQLite', section: 'Databases', type: 'docs', icon: '📁', route: '/docs#db-sqlite', keywords: 'sqlite file local' },
  { id: 'd-oracle', title: 'Oracle', section: 'Databases', type: 'docs', icon: '🏛️', route: '/docs#db-oracle', keywords: 'oracle instant client' },
  { id: 'd-sqlserver', title: 'SQL Server', section: 'Databases', type: 'docs', icon: '🪟', route: '/docs#db-sqlserver', keywords: 'sql server mssql microsoft' },
  { id: 'd-clickhouse', title: 'ClickHouse', section: 'Databases', type: 'docs', icon: '🏠', route: '/docs#db-clickhouse', keywords: 'clickhouse analytics' },
  { id: 'd-firebird', title: 'Firebird', section: 'Databases', type: 'docs', icon: '🔥', route: '/docs#db-firebird', keywords: 'firebird' },

  // Query Management
  { id: 'd-queries', title: 'Query Management', section: 'Usage', type: 'docs', icon: '📂', route: '/docs#queries', keywords: 'query queries manage sql' },
  { id: 'd-add', title: 'pam add', section: 'Queries', type: 'command', icon: '➕', route: '/docs#add-query', keywords: 'add save query new create' },
  { id: 'd-run', title: 'pam run', section: 'Queries', type: 'command', icon: '▶️', route: '/docs#run-query', keywords: 'run execute query sql select' },
  { id: 'd-query-table', title: 'pam query --table', section: 'Queries', type: 'command', icon: '🎯', route: '/docs#query-table', keywords: 'query table select from where' },
  { id: 'd-list-q', title: 'pam list queries', section: 'Queries', type: 'command', icon: '📋', route: '/docs#list-queries', keywords: 'list queries search' },
  { id: 'd-remove', title: 'pam remove', section: 'Queries', type: 'command', icon: '🗑️', route: '/docs#remove-query', keywords: 'remove delete query' },

  // Tables & Exploration
  { id: 'd-tables', title: 'Database Exploration', section: 'Usage', type: 'docs', icon: '🔍', route: '/docs#tables', keywords: 'tables explore schema' },
  { id: 'd-tables-list', title: 'pam tables', section: 'Exploration', type: 'command', icon: '📊', route: '/docs#tables-list', keywords: 'tables list schema browse' },
  { id: 'd-tv', title: 'pam table-view', section: 'Exploration', type: 'command', icon: '🏗️', route: '/docs#table-view', keywords: 'table-view tv columns structure alter schema' },
  { id: 'd-explore', title: 'pam explore', section: 'Exploration', type: 'command', icon: '🗺️', route: '/docs#explore', keywords: 'explore tables views browse' },
  { id: 'd-explain', title: 'pam explain', section: 'Exploration', type: 'command', icon: '🔗', route: '/docs#explain', keywords: 'explain relationships foreign key fk' },
  { id: 'd-info', title: 'pam info', section: 'Exploration', type: 'command', icon: 'ℹ️', route: '/docs#info', keywords: 'info tables views schema' },

  // Config
  { id: 'd-config', title: 'Configuration', section: 'Configuration', type: 'docs', icon: '⚙️', route: '/docs#config', keywords: 'config configuration yaml settings' },
  { id: 'd-rowlimit', title: 'Row Limit', section: 'Config', type: 'docs', icon: '🔢', route: '/docs#row-limit', keywords: 'row limit default_row_limit' },
  { id: 'd-colwidth', title: 'Column Width', section: 'Config', type: 'docs', icon: '↔️', route: '/docs#column-width', keywords: 'column width dynamic responsive' },
  { id: 'd-colors', title: 'Color Schemes', section: 'Config', type: 'docs', icon: '🎨', route: '/docs#color-schemes', keywords: 'color scheme theme dracula gruvbox nord monokai catppuccin tokyo' },
  { id: 'd-editor', title: 'Editor Integration', section: 'Config', type: 'docs', icon: '✏️', route: '/docs#editor', keywords: 'editor vim nano vscode EDITOR' },
  { id: 'd-completion', title: 'Shell Completion', section: 'Config', type: 'docs', icon: '🐚', route: '/docs#completion', keywords: 'completion autocomplete bash zsh fish shell tab' },

  // Reference
  { id: 'd-commands', title: 'Commands Reference', section: 'Reference', type: 'docs', icon: '📚', route: '/docs#commands', keywords: 'commands reference all list' },
  { id: 'd-tui', title: 'TUI Keys', section: 'Reference', type: 'docs', icon: '⌨️', route: '/docs#tui', keywords: 'tui keys keyboard navigation vim' },

  // TUI specific keys
  { id: 'k-copy', title: 'Copy cell (y)', section: 'TUI Keys', type: 'key', icon: '📋', route: '/docs#tui', keywords: 'copy yank clipboard cell' },
  { id: 'k-visual', title: 'Visual mode (v)', section: 'TUI Keys', type: 'key', icon: '🔲', route: '/docs#tui', keywords: 'visual mode select multiple cells' },
  { id: 'k-export', title: 'Export (x)', section: 'TUI Keys', type: 'key', icon: '📤', route: '/docs#tui', keywords: 'export csv json sql markdown html' },
  { id: 'k-update', title: 'Update cell (u)', section: 'TUI Keys', type: 'key', icon: '✏️', route: '/docs#tui', keywords: 'update edit cell value' },
  { id: 'k-delete', title: 'Delete row (D)', section: 'TUI Keys', type: 'key', icon: '🗑️', route: '/docs#tui', keywords: 'delete row remove' },
  { id: 'k-sort', title: 'Sort column (f)', section: 'TUI Keys', type: 'key', icon: '🔽', route: '/docs#tui', keywords: 'sort column order asc desc' },
]

const results = computed(() => {
  if (!query.value) return []
  const q = query.value.toLowerCase().trim()
  const terms = q.split(/\s+/)

  return searchIndex
    .map(item => {
      const haystack = `${item.title} ${item.section} ${item.keywords}`.toLowerCase()
      let score = 0
      for (const term of terms) {
        if (haystack.includes(term)) {
          score += 1
          // Boost exact title match
          if (item.title.toLowerCase().includes(term)) score += 2
          // Boost prefix match
          if (item.title.toLowerCase().startsWith(term)) score += 3
        } else {
          return { ...item, score: -1 }
        }
      }
      return { ...item, score }
    })
    .filter(r => r.score > 0)
    .sort((a, b) => b.score - a.score)
    .slice(0, 15)
})

function highlight(text) {
  if (!query.value) return text
  const terms = query.value.trim().split(/\s+/)
  let result = text
  for (const term of terms) {
    const re = new RegExp(`(${term.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
    result = result.replace(re, '<mark>$1</mark>')
  }
  return result
}

watch(query, () => { selectedIndex.value = 0 })

function moveDown() {
  if (selectedIndex.value < results.value.length - 1) {
    selectedIndex.value++
    scrollToSelected()
  }
}
function moveUp() {
  if (selectedIndex.value > 0) {
    selectedIndex.value--
    scrollToSelected()
  }
}
function scrollToSelected() {
  nextTick(() => {
    const el = resultsRef.value?.children[selectedIndex.value]
    el?.scrollIntoView({ block: 'nearest' })
  })
}

function goToSelected() {
  const item = results.value[selectedIndex.value]
  if (item) goTo(item)
}

function goTo(item) {
  close()
  if (item.route) {
    const [path, hash] = item.route.split('#')
    router.push({ path: path || '/', hash: hash ? '#' + hash : undefined })
  }
}

function toggle() {
  open.value = !open.value
  if (open.value) {
    query.value = ''
    selectedIndex.value = 0
    nextTick(() => inputRef.value?.focus())
  }
}
function close() {
  open.value = false
}

function onKeydown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    toggle()
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))

defineExpose({ toggle })
</script>

<style scoped>
.search-overlay {
  position: fixed;
  inset: 0;
  z-index: 200;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  padding-top: min(20vh, 140px);
}
.search-modal {
  width: 100%;
  max-width: 600px;
  max-height: 480px;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.5);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: modal-in 0.15s ease;
}
@keyframes modal-in {
  from { transform: scale(0.98) translateY(-8px); opacity: 0; }
  to { transform: scale(1) translateY(0); opacity: 1; }
}

/* Header */
.search-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border);
}
.search-icon {
  color: var(--text-muted);
  flex-shrink: 0;
}
.search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  font-size: 1rem;
  font-family: var(--font-sans);
  color: var(--text);
}
.search-input::placeholder {
  color: var(--text-muted);
}
.search-esc {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  background: var(--bg-code);
  border: 1px solid var(--border);
  padding: 2px 8px;
  border-radius: 4px;
  color: var(--text-muted);
}

/* Body */
.search-body {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}
.search-empty {
  padding: 2rem;
  text-align: center;
  color: var(--text-muted);
  font-size: 0.9rem;
}
.search-empty strong {
  color: var(--text-secondary);
}

/* Results */
.search-results {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.search-result {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  border-radius: var(--radius-sm);
  border: none;
  background: none;
  cursor: pointer;
  width: 100%;
  text-align: left;
  font-family: var(--font-sans);
  transition: background 0.1s;
  color: var(--text-secondary);
  font-size: 0.9rem;
}
.search-result:hover,
.search-result.active {
  background: rgba(88, 166, 255, 0.1);
}
.search-result.active {
  outline: 1px solid var(--border-accent);
}
.result-icon {
  font-size: 1.1rem;
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}
.result-text {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}
.result-title {
  color: var(--text);
  font-weight: 500;
}
.result-title :deep(mark) {
  background: rgba(88, 166, 255, 0.3);
  color: var(--accent-hover);
  border-radius: 2px;
  padding: 0 2px;
}
.result-section {
  font-size: 0.75rem;
  color: var(--text-muted);
}
.result-type {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-muted);
  background: var(--bg-code);
  padding: 2px 8px;
  border-radius: 4px;
  flex-shrink: 0;
}

/* Hints */
.search-hints {
  padding: 0.75rem;
}
.hint-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.hint-label {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  padding: 0.4rem 0.75rem 0.25rem;
}

/* Footer */
.search-footer {
  display: flex;
  gap: 1.5rem;
  padding: 0.5rem 1rem;
  border-top: 1px solid var(--border);
  font-size: 0.75rem;
  color: var(--text-muted);
}
.search-footer kbd {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  background: var(--bg-code);
  border: 1px solid var(--border);
  padding: 1px 6px;
  border-radius: 3px;
  color: var(--text-secondary);
  margin-right: 4px;
}

/* Transitions */
.overlay-enter-active { transition: opacity 0.15s ease; }
.overlay-leave-active { transition: opacity 0.1s ease; }
.overlay-enter-from,
.overlay-leave-to { opacity: 0; }

@media (max-width: 640px) {
  .search-modal {
    max-width: calc(100vw - 2rem);
    max-height: 70vh;
  }
}
</style>
