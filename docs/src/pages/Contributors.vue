<template>
  <div class="contributors-page">
    <div class="page-inner">
      <div class="page-header">
        <h1>Contributors</h1>
        <p class="lead">
          Pam is open-source and built by the community.
          <br />Thanks to everyone who has contributed! 💜
        </p>
      </div>

      <!-- Maintainer -->
      <section class="section">
        <h2>Current Maintainer</h2>
        <div class="maintainer-card">
          <img
            :src="maintainer.avatar"
            :alt="maintainer.name"
            class="avatar avatar-lg"
          />
          <div>
            <h3>
              <a :href="maintainer.url" target="_blank">{{
                maintainer.name
              }}</a>
            </h3>
            <p class="handle">@{{ maintainer.login }}</p>
            <p class="bio">{{ maintainer.bio }}</p>
          </div>
        </div>
      </section>

      <!-- Original Creator -->
      <section class="section">
        <h2>Original Creator</h2>
        <div class="maintainer-card">
          <img
            src="https://github.com/eduardofuncao.png"
            alt="Eduardo Funcao"
            class="avatar avatar-lg"
          />
          <div>
            <h3>
              <a href="https://github.com/eduardofuncao" target="_blank"
                >Eduardo Funcao</a
              >
            </h3>
            <p class="handle">@eduardofuncao</p>
            <p class="bio">Original creator of Pam's Database Drawer.</p>
          </div>
        </div>
      </section>

      <!-- Contributors -->
      <section class="section">
        <h2>Contributors</h2>
        <p class="section-desc" v-if="loading">
          Loading contributors from GitHub...
        </p>
        <p class="section-desc" v-else-if="error">{{ error }}</p>
        <div class="contrib-grid" v-else>
          <a
            v-for="c in contributors"
            :key="c.login"
            :href="c.html_url"
            target="_blank"
            class="contrib-card"
          >
            <img :src="c.avatar_url" :alt="c.login" class="avatar" />
            <span class="contrib-name">@{{ c.login }}</span>
            <span class="contrib-count"
              >{{ c.contributions }} commit{{
                c.contributions !== 1 ? 's' : ''
              }}</span
            >
          </a>
        </div>
      </section>

      <!-- How to Contribute -->
      <section class="section">
        <h2>How to Contribute</h2>
        <div class="how-grid">
          <div class="how-card" v-for="h in howTo" :key="h.title">
            <span class="how-icon">{{ h.icon }}</span>
            <h3>{{ h.title }}</h3>
            <p>{{ h.desc }}</p>
          </div>
        </div>

        <div class="contribute-steps">
          <h3>Development Workflow</h3>
          <ol>
            <li>Fork the repository</li>
            <li>
              Create a feature branch: <code>git checkout -b my-feature</code>
            </li>
            <li>
              Make your changes and test with
              <a href="https://github.com/caiolandgraf/dbeesly" target="_blank"
                >dbeesly</a
              >
              databases
            </li>
            <li>Commit: <code>git commit -m "Add my feature"</code></li>
            <li>Push: <code>git push origin my-feature</code></li>
            <li>Open a Pull Request</li>
          </ol>
        </div>
      </section>

      <!-- Acknowledgments -->
      <section class="section">
        <h2>Acknowledgments</h2>
        <p class="section-desc">
          Pam wouldn't exist without these fantastic projects:
        </p>
        <div class="ack-grid">
          <a
            v-for="a in acknowledgments"
            :key="a.name"
            :href="a.url"
            target="_blank"
            class="ack-card"
          >
            <h4>{{ a.name }}</h4>
            <p>{{ a.desc }}</p>
          </a>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const maintainer = {
  name: 'Caio Landgraf',
  login: 'caiolandgraf',
  avatar: 'https://github.com/caiolandgraf.png',
  url: 'https://github.com/caiolandgraf',
  bio: "Current maintainer of Pam's Database Drawer. Full Stack Developer at EiCode, passionate about CLI tools and developer experience."
}

const contributors = ref([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const res = await fetch(
      'https://api.github.com/repos/caiolandgraf/pam/contributors?per_page=100'
    )
    if (!res.ok)
      throw new Error('GitHub API rate limit reached. Please try again later.')
    const data = await res.json()
    contributors.value = data.filter(c => c.type === 'User')
  } catch (e) {
    error.value = e.message
    // Fallback static list
    contributors.value = [
      {
        login: 'eduardofuncao',
        avatar_url: 'https://github.com/eduardofuncao.png',
        html_url: 'https://github.com/eduardofuncao',
        contributions: 0
      },
      {
        login: 'DeprecatedLuar',
        avatar_url: 'https://github.com/DeprecatedLuar.png',
        html_url: 'https://github.com/DeprecatedLuar',
        contributions: 0
      },
      {
        login: 'caiolandgraf',
        avatar_url: 'https://github.com/caiolandgraf.png',
        html_url: 'https://github.com/caiolandgraf',
        contributions: 0
      }
    ]
  } finally {
    loading.value = false
  }
})

const howTo = [
  {
    icon: '🐛',
    title: 'Report Bugs',
    desc: 'Found a bug? Open an issue on GitHub with steps to reproduce it.'
  },
  {
    icon: '💡',
    title: 'Suggest Features',
    desc: "Have an idea? Open an issue describing what you'd like to see."
  },
  {
    icon: '🔧',
    title: 'Submit PRs',
    desc: 'Fix bugs, add features, or improve docs — all contributions are welcome.'
  },
  {
    icon: '🧪',
    title: 'Test with Databases',
    desc: 'Test Pam with different databases using the dbeesly project.'
  }
]

const acknowledgments = [
  {
    name: 'naggie/dstask',
    url: 'https://github.com/naggie/dstask',
    desc: 'Elegant CLI design patterns and file-based data storage approach.'
  },
  {
    name: 'DeprecatedLuar/better-curl-saul',
    url: 'https://github.com/DeprecatedLuar/better-curl-saul',
    desc: 'Simple and genius approach to making a CLI tool.'
  },
  {
    name: 'DBeaver',
    url: 'https://github.com/dbeaver/dbeaver',
    desc: 'The OG database management tool — inspiration for many features.'
  },
  {
    name: 'Bubble Tea',
    url: 'https://github.com/charmbracelet/bubbletea',
    desc: "The powerful TUI framework that makes Pam's interface possible."
  }
]
</script>

<style scoped>
.contributors-page {
  padding-top: 64px;
}
.page-inner {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 3rem 1.5rem 4rem;
}
.page-header {
  text-align: center;
  margin-bottom: 3rem;
}
.page-header h1 {
  font-size: 2.2rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
}
.lead {
  color: var(--text-secondary);
  font-size: 1.1rem;
}

.section {
  margin-bottom: 3.5rem;
}
.section h2 {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
}
.section-desc {
  color: var(--text-secondary);
  margin-bottom: 1.25rem;
}

/* Maintainer */
.maintainer-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 1.5rem;
}
.avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 2px solid var(--border);
}
.avatar-lg {
  width: 80px;
  height: 80px;
}
.maintainer-card h3 {
  font-size: 1.15rem;
  margin-bottom: 0.15rem;
}
.handle {
  color: var(--text-muted);
  font-size: 0.85rem;
  font-family: var(--font-mono);
  margin-bottom: 0.4rem;
}
.bio {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

/* Contributor grid */
.contrib-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}
.contrib-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1.25rem;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  transition:
    border-color 0.2s,
    transform 0.2s;
  text-align: center;
}
.contrib-card:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
}
.contrib-name {
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--text);
}
.contrib-count {
  font-size: 0.8rem;
  color: var(--text-muted);
}

/* How to */
.how-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}
.how-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 1.25rem;
}
.how-icon {
  font-size: 1.5rem;
  display: block;
  margin-bottom: 0.5rem;
}
.how-card h3 {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.3rem;
}
.how-card p {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.contribute-steps {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 1.5rem;
}
.contribute-steps h3 {
  font-size: 1rem;
  margin-bottom: 0.75rem;
}
.contribute-steps ol {
  padding-left: 1.5rem;
  color: var(--text-secondary);
  font-size: 0.9rem;
}
.contribute-steps li {
  margin-bottom: 0.4rem;
}

/* Acknowledgments */
.ack-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1rem;
}
.ack-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 1.25rem;
  transition:
    border-color 0.2s,
    transform 0.2s;
}
.ack-card:hover {
  border-color: var(--accent);
  transform: translateY(-2px);
}
.ack-card h4 {
  font-size: 0.95rem;
  font-weight: 600;
  margin-bottom: 0.3rem;
  color: var(--text);
}
.ack-card p {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

@media (max-width: 600px) {
  .maintainer-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>
