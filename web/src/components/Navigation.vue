<template>
  <nav class="navigation">
    <div class="nav-brand">
      <router-link to="/" class="brand-link">FATE Vault</router-link>
    </div>
    <div class="nav-links">
      <router-link to="/characters" class="nav-link">Characters</router-link>
      <router-link to="/stunts" class="nav-link">Stunts</router-link>
      <router-link to="/user" class="nav-link user-link">User</router-link>
    </div>
  </nav>
</template>

<script setup>
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

.user-link {
  margin-left: auto;
}
</style>


