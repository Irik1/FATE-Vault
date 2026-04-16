import { ref, computed } from 'vue'
import { userService } from '../services/api'

const user = ref(null)
const initialized = ref(false)

export async function initAuth() {
  try {
    const data = await userService.me()
    if (data?.user) {
      user.value = data.user
    } else {
      user.value = null
    }
  } catch {
    user.value = null
  } finally {
    initialized.value = true
  }
}

export function useAuth() {
  const isAuthenticated = computed(() => !!user.value)

  async function login(username, password) {
    const data = await userService.auth({ username, password })
    if (data?.user) {
      user.value = data.user
    }
    return data
  }

  async function register(username, password) {
    const data = await userService.register({ username, password })
    if (data?.user) {
      user.value = data.user
    }
    return data
  }

  async function logout() {
    try {
      await userService.logout()
    } catch {
      // still clear local state
    }
    user.value = null
  }

  async function updateProfile(id, body) {
    const updated = await userService.update(id, body)
    user.value = updated
    return updated
  }

  return {
    user,
    initialized,
    isAuthenticated,
    login,
    register,
    logout,
    updateProfile
  }
}

export { user, initialized }
