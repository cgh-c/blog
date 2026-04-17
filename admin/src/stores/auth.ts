import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMe } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(false)
  const username = ref('')

  async function checkAuth() {
    try {
      const res: any = await getMe()
      authenticated.value = res.data.authenticated
      username.value = res.data.username || ''
      return authenticated.value
    } catch {
      authenticated.value = false
      return false
    }
  }

  function setAuth(user: { username: string }) {
    authenticated.value = true
    username.value = user.username
  }

  function clearAuth() {
    authenticated.value = false
    username.value = ''
  }

  return { authenticated, username, checkAuth, setAuth, clearAuth }
})
