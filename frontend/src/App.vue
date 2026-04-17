<script setup lang="ts">
import { onMounted } from 'vue'
import { useSiteStore } from './stores/site'
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'

const site = useSiteStore()

onMounted(() => {
  site.initDarkMode()
  site.loadSettings()
})
</script>

<template>
  <Header />
  <main class="main-content">
    <RouterView />
  </main>
  <Footer />
</template>
