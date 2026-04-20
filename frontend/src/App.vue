<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSiteStore } from './stores/site'
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'
import BackToTop from './components/BackToTop.vue'

const route = useRoute()
const site = useSiteStore()
const isWelcome = computed(() => route.name === 'Welcome')

onMounted(() => {
  site.initDarkMode()
  site.loadSettings()
})
</script>

<template>
  <Header v-if="!isWelcome" />
  <main :class="{ 'main-content': !isWelcome }">
    <RouterView v-slot="{ Component }">
      <Transition name="page" mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>
  </main>
  <Footer v-if="!isWelcome" />
  <BackToTop v-if="!isWelcome" />
</template>
