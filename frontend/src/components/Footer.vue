<script setup lang="ts">
import { useSiteStore } from '../stores/site'

const site = useSiteStore()
</script>

<template>
  <footer class="site-footer">
    <div class="container">
      <p>&copy; {{ new Date().getFullYear() }} {{ site.settings.site_name || 'My Blog' }}</p>
      <p v-if="site.settings.icp_number">{{ site.settings.icp_number }}</p>
    </div>
  </footer>
</template>
