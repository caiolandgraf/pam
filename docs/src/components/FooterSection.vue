<template>
  <footer class="footer">
    <div class="footer-inner">
      <div class="footer-brand">
        <span class="brand-icon">🗄️</span>
        <span>Pam's Database Drawer</span>
      </div>
      <p class="footer-quote">
        <em
          >"I don't think it would be the worst thing if it didn't work out...
          Wait, can I say that?"</em
        >
        <span class="quote-attr">— Pam Beesly</span>
      </p>
      <div class="footer-links">
        <a href="https://github.com/caiolandgraf/pam" target="_blank">GitHub</a>
        <a href="https://github.com/caiolandgraf/pam/releases" target="_blank"
          >Releases</a
        >
        <a
          href="https://github.com/caiolandgraf/pam/blob/main/LICENSE"
          target="_blank"
          >MIT License</a
        >
      </div>
      <p class="footer-copy">
        Made with 👚 by
        <a href="https://github.com/caiolandgraf" target="_blank"
          >@caiolandgraf</a
        >
        · originally by
        <a href="https://github.com/eduardofuncao" target="_blank"
          >@eduardofuncao</a
        >
      </p>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  border-top: 1px solid var(--border);
  margin-top: 6rem;
  padding: 3rem 1.5rem;
  text-align: center;
}
.footer-inner {
  max-width: var(--max-width);
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}
.footer-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 1.1rem;
}
.brand-icon {
  font-size: 1.4rem;
}
.footer-quote {
  color: var(--text-secondary);
  font-size: 0.9rem;
  max-width: 480px;
}
.quote-attr {
  display: block;
  color: var(--text-muted);
  font-size: 0.8rem;
  margin-top: 0.25rem;
}
.footer-links {
  display: flex;
  gap: 1.5rem;
}
.footer-links a {
  color: var(--text-muted);
  font-size: 0.85rem;
}
.footer-links a:hover {
  color: var(--accent);
}
.footer-copy {
  color: var(--text-muted);
  font-size: 0.8rem;
}
</style>
