<template>
  <nav class="navbar">
    <div class="nav-inner">
      <div class="nav-left">
        <router-link to="/" class="nav-brand">
          <span class="brand-row">
            <span class="brand-icon">🗄️</span>
            <span class="brand-text">Pam's Database Drawer</span>
          </span>
          <span class="brand-sub">Scranton Branch</span>
        </router-link>
        <div class="nav-meta">
          <span><span class="meta-label">To:</span> The Office</span>
          <span><span class="meta-label">From:</span> Reception</span>
          <span><span class="meta-label">Subject:</span> Queries & Docs</span>
          <span><span class="meta-label">Form:</span> DM-00</span>
        </div>
      </div>
      <button
        class="mobile-toggle"
        @click="menuOpen = !menuOpen"
        :class="{ active: menuOpen }"
      >
        <span></span><span></span><span></span>
      </button>
      <div class="nav-links" :class="{ open: menuOpen }">
        <router-link to="/" @click="menuOpen = false">Home</router-link>
        <router-link to="/docs" @click="menuOpen = false">Docs</router-link>
        <router-link to="/playground" @click="menuOpen = false"
          >Playground</router-link
        >
        <router-link to="/contributors" @click="menuOpen = false"
          >Contributors</router-link
        >
        <button
          class="theme-toggle"
          @click="toggleTheme"
          :aria-label="`Switch to ${themeLabel} mode`"
        >
          <span class="theme-icon">{{ themeIcon }}</span>
          <span class="theme-text">{{ themeLabel }}</span>
        </button>
        <button class="search-trigger" @click="$emit('openSearch')">
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
          <span class="search-text">Search</span>
          <kbd>{{ metaKey }}K</kbd>
        </button>
        <a
          href="https://github.com/caiolandgraf/pam"
          target="_blank"
          class="nav-github"
        >
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path
              d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"
            />
          </svg>
        </a>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

defineEmits(['openSearch'])

const menuOpen = ref(false)
const metaKey = navigator.platform.includes('Mac') ? '⌘' : 'Ctrl'
const theme = ref('light')

const themeIcon = computed(() => (theme.value === 'dark' ? '☀️' : '🌙'))
const themeLabel = computed(() => (theme.value === 'dark' ? 'Light' : 'Dark'))

function applyTheme(value) {
  document.documentElement.dataset.theme = value
  localStorage.setItem('pam-docs-theme', value)
}

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  applyTheme(theme.value)
}

onMounted(() => {
  const stored = localStorage.getItem('pam-docs-theme')
  if (stored) {
    theme.value = stored
  } else if (
    window.matchMedia &&
    window.matchMedia('(prefers-color-scheme: dark)').matches
  ) {
    theme.value = 'dark'
  }
  applyTheme(theme.value)
})
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--bg-card);
  border-bottom: 2px solid var(--border);
  box-shadow: var(--shadow-soft);
}
.nav-inner {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0.625rem 1.5rem;
  min-height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: nowrap;
}
.nav-left {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}
.nav-brand {
  display: flex;
  flex-direction: column;
  gap: 0.1rem;
  font-weight: 700;
  font-size: 1.1rem;
  color: var(--text);
}
.brand-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  white-space: nowrap;
}
.brand-icon {
  font-size: 1.35rem;
}
.brand-sub {
  font-size: 0.7rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: var(--text-muted);
  font-family: var(--font-mono);
}
.nav-meta {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.15rem 1rem;
  font-size: 0.68rem;
  font-family: var(--font-mono);
  color: var(--text-secondary);
  flex-shrink: 0;
}
.meta-label {
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  margin-right: 0.3rem;
}
.nav-links {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}
.nav-links a {
  color: var(--text-secondary);
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.2s;
  padding: 6px 10px;
  border: 1px solid var(--border);
  border-bottom-width: 2px;
  border-radius: 6px 6px 0 0;
  background: var(--bg-card);
  font-family: var(--font-mono);
}
.nav-links a:hover,
.nav-links a.router-link-active {
  color: var(--text);
  border-color: var(--accent);
  box-shadow: var(--shadow-soft);
}
.nav-github {
  display: flex;
  align-items: center;
}

/* Search trigger */
.theme-toggle {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: var(--bg-code);
  border: 1px dashed var(--border);
  color: var(--text-secondary);
  padding: 5px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: 0.7rem;
  transition: all 0.15s;
}
.theme-toggle:hover {
  border-color: var(--accent);
  color: var(--text);
}
.theme-icon {
  font-size: 0.9rem;
}

.search-trigger {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--bg-code);
  border: 1px dashed var(--border);
  color: var(--text-muted);
  padding: 5px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-family: var(--font-mono);
  font-size: 0.75rem;
  transition: all 0.15s;
}
.search-trigger:hover {
  border-color: var(--accent);
  color: var(--text-secondary);
}
.search-trigger kbd {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  background: rgba(214, 201, 184, 0.35);
  border: 1px solid var(--border);
  padding: 1px 5px;
  border-radius: 3px;
  color: var(--text-muted);
}

.mobile-toggle {
  display: none;
  flex-direction: column;
  gap: 5px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
}
.mobile-toggle span {
  display: block;
  width: 22px;
  height: 2px;
  background: var(--text-secondary);
  border-radius: 2px;
  transition: all 0.3s;
}
.mobile-toggle.active span:nth-child(1) {
  transform: rotate(45deg) translate(5px, 5px);
}
.mobile-toggle.active span:nth-child(2) {
  opacity: 0;
}
.mobile-toggle.active span:nth-child(3) {
  transform: rotate(-45deg) translate(5px, -5px);
}

@media (max-width: 1020px) {
  .nav-meta {
    display: none;
  }
}

@media (max-width: 860px) {
  .brand-sub {
    display: none;
  }
}

@media (max-width: 768px) {
  .mobile-toggle {
    display: flex;
  }
  .search-text {
    display: none;
  }
  .search-trigger kbd {
    display: none;
  }
  .theme-text {
    display: none;
  }
  .nav-links {
    position: fixed;
    top: 64px;
    left: 0;
    right: 0;
    background: var(--bg-card);
    flex-direction: column;
    padding: 1.25rem;
    gap: 0.75rem;
    border-bottom: 2px solid var(--border);
    transform: translateY(-120%);
    opacity: 0;
    transition: all 0.3s ease;
  }
  .nav-links.open {
    transform: translateY(0);
    opacity: 1;
  }
}
</style>
