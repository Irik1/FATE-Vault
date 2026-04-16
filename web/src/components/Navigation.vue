<template>
  <nav class="navigation">
    <div class="nav-brand">
      <router-link to="/" class="brand-link">FATE Vault</router-link>
    </div>
    <div class="nav-links">
      <router-link to="/characters" class="nav-link">Characters</router-link>
      <router-link to="/stunts" class="nav-link">Stunts</router-link>
      <div class="nav-auth">
        <template v-if="isAuthenticated">
          <router-link to="/user" class="nav-link">Account</router-link>
          <button type="button" class="nav-btn" @click="doLogout">Sign out</button>
        </template>
        <template v-else>
          <router-link to="/login" class="nav-link">Sign in</router-link>
          <router-link to="/register" class="nav-link nav-cta">Register</router-link>
        </template>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { isAuthenticated, logout } = useAuth()

async function doLogout() {
  await logout()
  router.push({ name: 'Characters' })
}
</script>

<style scoped>
.navigation {
  background-color: #2c3e50;
  color: white;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.nav-brand {
  font-size: 1.5rem;
  font-weight: bold;
}

.brand-link {
  color: white;
  text-decoration: none;
  transition: opacity 0.2s;
}

.brand-link:hover {
  opacity: 0.8;
}

.nav-links {
  display: flex;
  gap: 2rem;
  align-items: center;
  flex: 1;
  justify-content: flex-end;
}

.nav-auth {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-left: 2rem;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.nav-link:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.nav-link.router-link-active {
  background-color: rgba(255, 255, 255, 0.2);
}

.nav-cta {
  background-color: rgba(255, 255, 255, 0.15);
}

.nav-btn {
  margin-left: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.45);
  background: transparent;
  color: white;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.nav-btn:hover {
  background-color: rgba(255, 255, 255, 0.12);
}
</style>


