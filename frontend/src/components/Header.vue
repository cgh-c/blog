<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useSiteStore } from '../stores/site'

const site = useSiteStore()
</script>

<template>
  <header class="site-header">
    <div class="container">
      <RouterLink to="/" class="site-title">
        {{ site.settings.site_name || 'My Blog' }}
      </RouterLink>
      <nav class="site-nav">
        <RouterLink to="/">Home</RouterLink>
        <RouterLink to="/archive">Archive</RouterLink>
        <RouterLink to="/about">About</RouterLink>
        <button class="dark-toggle" @click="site.toggleDarkMode">
          {{ site.darkMode ? 'Light' : 'Dark' }}
        </button>
      </nav>
    </div>
  </header>
</template>
