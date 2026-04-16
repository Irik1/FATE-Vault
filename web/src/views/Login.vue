<template>
  <div class="auth-page">
    <div class="auth-card">
      <h1>Sign in</h1>
      <p class="subtitle">Use your FATE Vault account.</p>

      <form @submit.prevent="submit">
        <label class="field">
          <span>Username</span>
          <input
            v-model.trim="username"
            type="text"
            autocomplete="username"
            required
            :disabled="loading"
          />
        </label>
        <label class="field">
          <span>Password</span>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            required
            :disabled="loading"
          />
        </label>

        <p v-if="error" class="error-msg">{{ error }}</p>

        <button type="submit" class="submit-btn" :disabled="loading">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </button>
      </form>

      <p class="footer-link">
        No account?
        <router-link to="/register">Create one</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { login } = useAuth()

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref(null)

function apiErrorMessage(err) {
  return err?.response?.data?.error || err?.message || 'Request failed'
}

async function submit() {
  error.value = null
  loading.value = true
  try {
    await login(username.value, password.value)
    await router.replace({ name: 'User' })
  } catch (err) {
    error.value = apiErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 1rem;
}

.auth-card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  max-width: 420px;
  width: 100%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h1 {
  margin: 0 0 0.5rem;
  color: #2c3e50;
  font-size: 1.75rem;
}

.subtitle {
  margin: 0 0 1.5rem;
  color: #7f8c8d;
  font-size: 0.95rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  margin-bottom: 1rem;
}

.field span {
  font-size: 0.9rem;
  color: #555;
  font-weight: 500;
}

.field input {
  padding: 0.65rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.field input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.field input:disabled {
  background: #f4f4f4;
}

.error-msg {
  color: #e74c3c;
  font-size: 0.9rem;
  margin: 0 0 1rem;
}

.submit-btn {
  width: 100%;
  padding: 0.75rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background-color: #2980b9;
}

.submit-btn:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.footer-link {
  margin-top: 1.5rem;
  text-align: center;
  color: #555;
  font-size: 0.95rem;
}

.footer-link a {
  color: #3498db;
  font-weight: 500;
}
</style>
