<template>
  <div class="characters-list">
    <h1>Characters</h1>
    
    <div v-if="loading" class="loading">Loading characters...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <div class="characters-grid">
        <div
          v-for="character in paginatedCharacters"
          :key="character._id || character.id"
          class="character-card"
          @click="goToCharacter(character._id || character.id)"
        >
          <h3>{{ character.name || 'Unnamed Character' }}</h3>
          <p class="edition">{{ character.edition || 'Unknown Edition' }}</p>
          <p class="description">{{ truncateDescription(character.description) }}</p>
        </div>
      </div>
      
      <div class="pagination" v-if="totalPages > 1">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="pagination-btn"
        >
          Previous
        </button>
        <span class="page-info">
          Page {{ currentPage }} of {{ totalPages }}
        </span>
        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="pagination-btn"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { characterService } from '../services/api'

const router = useRouter()
const characters = ref([])
const loading = ref(true)
const error = ref(null)
const currentPage = ref(1)
const itemsPerPage = 10

const totalPages = computed(() => {
  return Math.ceil(characters.value.length / itemsPerPage)
})

const paginatedCharacters = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return characters.value.slice(start, end)
})

const truncateDescription = (description) => {
  if (!description) return 'No description'
  return description.length > 150 ? description.substring(0, 150) + '...' : description
}

const goToCharacter = (id) => {
  router.push(`/characters/${id}`)
}

const loadCharacters = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await characterService.getCharacters()
    characters.value = data
  } catch (err) {
    error.value = 'Failed to load characters. Please make sure the backend is running.'
    console.error('Error loading characters:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCharacters()
})
</script>

<style scoped>
.characters-list {
  width: 100%;
}

h1 {
  margin-bottom: 2rem;
  color: #2c3e50;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
  font-size: 1.1rem;
}

.error {
  color: #e74c3c;
}

.characters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.character-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.character-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.character-card h3 {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.edition {
  font-size: 0.9rem;
  color: #7f8c8d;
  margin-bottom: 0.5rem;
  text-transform: capitalize;
}

.description {
  color: #555;
  line-height: 1.5;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.pagination-btn {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #2980b9;
}

.pagination-btn:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.page-info {
  font-size: 1rem;
  color: #555;
}
</style>


