<template>
  <div class="user-page">
    <template v-if="!isAuthenticated">
      <h1>Account</h1>
      <div class="card guest">
        <p>You are not signed in.</p>
        <div class="actions">
          <router-link to="/login" class="btn primary">Sign in</router-link>
          <router-link to="/register" class="btn secondary">Register</router-link>
        </div>
      </div>
    </template>

    <template v-else>
      <h1>Account</h1>
      <div class="card">
        <h2>Profile</h2>
        <dl class="profile-dl">
          <div>
            <dt>Username</dt>
            <dd>{{ user.username }}</dd>
          </div>
          <div>
            <dt>Role</dt>
            <dd>{{ user.role }}</dd>
          </div>
          <div v-if="user._id">
            <dt>User ID</dt>
            <dd class="mono">{{ user._id }}</dd>
          </div>
        </dl>

        <h2 class="section-title">Update username</h2>
        <form class="inline-form" @submit.prevent="saveUsername">
          <input
            v-model.trim="editUsername"
            type="text"
            :disabled="saving"
            autocomplete="username"
          />
          <button type="submit" class="btn small" :disabled="saving || !editUsername">
            {{ saving ? 'Saving…' : 'Save' }}
          </button>
        </form>
        <p v-if="profileError" class="error-msg">{{ profileError }}</p>
        <p v-if="profileOk" class="ok-msg">{{ profileOk }}</p>

        <button type="button" class="btn danger outline" @click="doLogout">Sign out</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { user, isAuthenticated, logout, updateProfile } = useAuth()

const editUsername = ref('')
const saving = ref(false)
const profileError = ref(null)
const profileOk = ref(null)

watch(
  user,
  (u) => {
    if (u?.username) editUsername.value = u.username
  },
  { immediate: true }
)

function apiErrorMessage(err) {
  return err?.response?.data?.error || err?.message || 'Request failed'
}

async function saveUsername() {
  profileError.value = null
  profileOk.value = null
  if (!user.value?._id || editUsername.value === user.value.username) return
  saving.value = true
  try {
    await updateProfile(user.value._id, { username: editUsername.value })
    profileOk.value = 'Username updated.'
  } catch (err) {
    profileError.value = apiErrorMessage(err)
  } finally {
    saving.value = false
  }
}

async function doLogout() {
  await logout()
  profileOk.value = null
  profileError.value = null
  router.push({ name: 'Characters' })
}
</script>

<style scoped>
.user-page {
  width: 100%;
  max-width: 560px;
}

h1 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
}

h2 {
  margin: 0 0 1rem;
  font-size: 1.15rem;
  color: #2c3e50;
}

.section-title {
  margin-top: 1.75rem;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 1.75rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card.guest {
  text-align: center;
}

.card.guest p {
  margin-bottom: 1.25rem;
  color: #555;
  font-size: 1.05rem;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.profile-dl {
  margin: 0;
}

.profile-dl > div {
  display: grid;
  grid-template-columns: 100px 1fr;
  gap: 0.5rem 1rem;
  margin-bottom: 0.75rem;
  align-items: baseline;
}

dt {
  margin: 0;
  font-size: 0.85rem;
  color: #7f8c8d;
  font-weight: 600;
}

dd {
  margin: 0;
  color: #333;
}

.mono {
  font-family: ui-monospace, monospace;
  font-size: 0.85rem;
  word-break: break-all;
}

.inline-form {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 0.5rem;
}

.inline-form input {
  flex: 1;
  min-width: 180px;
  padding: 0.55rem 0.65rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.inline-form input:focus {
  outline: none;
  border-color: #3498db;
}

.btn {
  display: inline-block;
  padding: 0.65rem 1.25rem;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  border: none;
  cursor: pointer;
  transition: background-color 0.2s, color 0.2s;
}

.btn.primary {
  background-color: #3498db;
  color: white;
}

.btn.primary:hover {
  background-color: #2980b9;
}

.btn.secondary {
  background-color: #ecf0f1;
  color: #2c3e50;
}

.btn.secondary:hover {
  background-color: #dde4e6;
}

.btn.small {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
}

.btn.small:hover:not(:disabled) {
  background-color: #2980b9;
}

.btn.small:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.btn.danger.outline {
  margin-top: 1.5rem;
  background: transparent;
  color: #c0392b;
  border: 1px solid #c0392b;
}

.btn.danger.outline:hover {
  background: rgba(192, 57, 43, 0.08);
}

.error-msg {
  color: #e74c3c;
  font-size: 0.9rem;
  margin: 0.25rem 0 0;
}

.ok-msg {
  color: #27ae60;
  font-size: 0.9rem;
  margin: 0.25rem 0 0;
}
</style>
