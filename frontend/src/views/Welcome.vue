<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSiteStore } from '../stores/site'

const router = useRouter()
const site = useSiteStore()
const ready = ref(false)

onMounted(() => {
  requestAnimationFrame(() => {
    ready.value = true
  })
})

function enter() {
  router.push('/home')
}
</script>

<template>
  <div class="welcome-page" :class="{ ready }">
    <!-- Animated background -->
    <div class="welcome-bg">
      <div class="welcome-orb orb-1"></div>
      <div class="welcome-orb orb-2"></div>
      <div class="welcome-orb orb-3"></div>
    </div>

    <div class="welcome-content">
      <div class="welcome-greeting">Hello, I'm</div>
      <h1 class="welcome-name">{{ site.settings.site_name || 'My Blog' }}</h1>
      <p class="welcome-desc">{{ site.settings.site_description || 'A personal blog about coding and life' }}</p>

      <div class="welcome-tags">
        <span class="welcome-tag">Go</span>
        <span class="welcome-tag">Vue</span>
        <span class="welcome-tag">Full Stack</span>
      </div>

      <button class="welcome-enter" @click="enter">
        Enter Blog
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="5" y1="12" x2="19" y2="12"/>
          <polyline points="12 5 19 12 12 19"/>
        </svg>
      </button>
    </div>

    <div class="welcome-footer">
      <a href="https://github.com" target="_blank" class="welcome-social">GitHub</a>
      <span class="welcome-dot"></span>
      <a href="#" class="welcome-social" @click.prevent="enter">Blog</a>
    </div>
  </div>
</template>
