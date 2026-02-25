<template>
  <nav class="navbar">
    <div class="nav-inner">
      <router-link to="/" class="nav-brand">
        <span class="brand-icon">🗄️</span>
        <span class="brand-text">pam</span>
      </router-link>
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
import { ref } from 'vue'

defineEmits(['openSearch'])

const menuOpen = ref(false)
const metaKey = navigator.platform.includes('Mac') ? '⌘' : 'Ctrl'
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(13, 17, 23, 0.85);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border);
}
.nav-inner {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 1.5rem;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.nav-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 700;
  font-size: 1.2rem;
  color: var(--text);
}
.brand-icon {
  font-size: 1.5rem;
}
.nav-links {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}
.nav-links a {
  color: var(--text-secondary);
  font-size: 0.9rem;
  font-weight: 500;
  transition: color 0.2s;
}
.nav-links a:hover,
.nav-links a.router-link-active {
  color: var(--text);
}
.nav-github {
  display: flex;
  align-items: center;
}

/* Search trigger */
.search-trigger {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--bg-code);
  border: 1px solid var(--border);
  color: var(--text-muted);
  padding: 5px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-family: var(--font-sans);
  font-size: 0.8rem;
  transition: all 0.15s;
}
.search-trigger:hover {
  border-color: var(--text-secondary);
  color: var(--text-secondary);
}
.search-trigger kbd {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  background: rgba(110, 118, 129, 0.15);
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
  .nav-links {
    position: fixed;
    top: 64px;
    left: 0;
    right: 0;
    background: rgba(13, 17, 23, 0.97);
    backdrop-filter: blur(16px);
    flex-direction: column;
    padding: 1.5rem;
    gap: 1.25rem;
    border-bottom: 1px solid var(--border);
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
