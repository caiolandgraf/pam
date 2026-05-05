<template>
  <section class="memo-page">
    <div class="dm-letterhead">
      <span>Dunder Mifflin Paper Company — Scranton Branch</span>
      <span>Est. 1949 · Scranton, PA 18503</span>
    </div>
    <div class="hole-punch-row">
      <span class="hole-punch"></span>
      <span class="hole-punch"></span>
      <span class="hole-punch"></span>
    </div>
    <header class="memo-header memo-header--page">
      <div class="memo-kv">
        <div><span class="memo-label">To:</span> {{ to }}</div>
        <div><span class="memo-label">From:</span> {{ from }}</div>
        <div v-if="subject">
          <span class="memo-label">Subject:</span> {{ subject }}
        </div>
        <div v-if="status">
          <span class="memo-label">Status:</span> {{ status }}
        </div>
      </div>
      <div class="memo-stamp">{{ stamp }}</div>
    </header>

    <div v-if="$slots.tags" class="memo-tags">
      <slot name="tags" />
    </div>

    <div v-if="$slots.meta" class="memo-meta">
      <slot name="meta" />
    </div>

    <div v-if="note || $slots.note" class="sticky-note memo-note">
      <slot name="note">{{ note }}</slot>
    </div>

    <div class="page-header">
      <div v-if="badge" class="page-badge">{{ badge }}</div>
      <h1 class="page-title">
        <slot name="title">{{ title }}</slot>
      </h1>
      <p v-if="subtitle || $slots.subtitle" class="lead">
        <slot name="subtitle">{{ subtitle }}</slot>
      </p>
    </div>

    <slot />

    <footer class="memo-page-footer">
      <span>{{ badge || 'Form DM-00' }}</span>
      <span>Page 1 of 1</span>
      <span>Dunder Mifflin Paper Company</span>
    </footer>
  </section>
</template>

<script setup>
defineProps({
  to: { type: String, default: 'Scranton Branch' },
  from: { type: String, default: 'Reception' },
  subject: { type: String, default: '' },
  status: { type: String, default: '' },
  stamp: { type: String, default: 'FILED' },
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  badge: { type: String, default: '' },
  note: { type: String, default: '' }
})
</script>

<style scoped>
.memo-page {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 2.5rem 1.5rem 3rem;
  display: flex;
  flex-direction: column;
  border-left: 3px solid rgba(184, 50, 43, 0.12);
}
.memo-header--page {
  margin-bottom: 0.75rem;
}
/* Override stamp rotation inside page headers for more character */
.memo-header--page :deep(.memo-stamp) {
  transform: rotate(-8deg);
  opacity: 0.78;
  font-size: 0.65rem;
  letter-spacing: 0.25em;
  border-width: 2.5px;
}
.memo-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}
.memo-meta {
  margin-bottom: 1rem;
}
.memo-note {
  margin-bottom: 2rem;
  display: inline-block;
}
.page-header {
  text-align: center;
  margin-bottom: 2.5rem;
}
.page-title {
  font-size: 2rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
  font-family: var(--font-mono);
  letter-spacing: 0.02em;
}
.page-badge {
  display: inline-block;
  font-family: var(--font-mono);
  font-size: 0.7rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  border: 1px dashed var(--border);
  border-radius: 6px;
  padding: 4px 10px;
  margin-bottom: 0.75rem;
  background: var(--bg-code);
  color: var(--text-secondary);
}
.lead {
  color: var(--text-secondary);
  font-size: 1rem;
}
</style>
