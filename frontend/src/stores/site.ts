import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSettings } from '../api'

export const useSiteStore = defineStore('site', () => {
  const settings = ref<Record<string, string>>({})
  const loaded = ref(false)
  const darkMode = ref(false)

  async function loadSettings() {
    try {
      const res: any = await getSettings()
      settings.value = res.data || {}
      loaded.value = true
    } catch {
      // silently fail
    }
  }

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
    document.documentElement.classList.toggle('dark', darkMode.value)
    localStorage.setItem('darkMode', darkMode.value ? '1' : '0')
  }

  function initDarkMode() {
    const saved = localStorage.getItem('darkMode')
    if (saved === '1') {
      darkMode.value = true
      document.documentElement.classList.add('dark')
    }
  }

  return { settings, loaded, darkMode, loadSettings, toggleDarkMode, initDarkMode }
})
