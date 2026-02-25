<template>
  <div class="code-block">
    <div class="code-header" v-if="title">
      <span class="code-lang">{{ title }}</span>
      <button class="copy-btn" @click="copy" :class="{ copied }">
        {{ copied ? '✓ Copied' : 'Copy' }}
      </button>
    </div>
    <pre><code><slot /></code></pre>
  </div>
</template>

<script setup>
import { ref, useSlots } from 'vue'

defineProps({ title: String })
const copied = ref(false)
const slots = useSlots()

function copy() {
  const text = slots.default?.()[0]?.children || ''
  navigator.clipboard.writeText(text.toString().trim())
  copied.value = true
  setTimeout(() => (copied.value = false), 2000)
}
</script>

<style scoped>
.code-block {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
  margin: 1rem 0;
}
.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  background: rgba(110, 118, 129, 0.1);
  border-bottom: 1px solid var(--border);
}
.code-lang {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.copy-btn {
  background: none;
  border: 1px solid var(--border);
  color: var(--text-secondary);
  font-size: 0.75rem;
  padding: 2px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-family: var(--font-sans);
  transition: all 0.2s;
}
.copy-btn:hover {
  color: var(--text);
  border-color: var(--text-secondary);
}
.copy-btn.copied {
  color: var(--green);
  border-color: var(--green);
}
.code-block pre {
  margin: 0;
  border: none;
  border-radius: 0;
}
</style>
