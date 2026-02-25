<template>
  <div class="app">
    <NavBar @open-search="searchRef?.toggle()" />
    <SearchPalette ref="searchRef" />
    <main>
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <FooterSection />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import NavBar from './components/NavBar.vue'
import FooterSection from './components/FooterSection.vue'
import SearchPalette from './components/SearchPalette.vue'

const searchRef = ref(null)
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
