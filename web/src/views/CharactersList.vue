<template>
  <div class="characters-list">
    <div class="header-section">
      <h1>Characters</h1>
      <button @click="showEditionModal = true" class="create-btn">+ Create Character</button>
    </div>
    
    <!-- Edition Selection Modal -->
    <div v-if="showEditionModal" class="modal-overlay" @click="showEditionModal = false">
      <div class="modal-content" @click.stop>
        <h2>Select Edition</h2>
        <div v-if="templatesLoading" class="loading">Loading templates...</div>
        <div v-else-if="templates.length === 0" class="error">No templates available</div>
        <div v-else class="edition-options">
          <button 
            v-for="template in templates" 
            :key="template.edition"
            @click="createCharacter(template.edition)"
            class="edition-btn"
          >
            {{ template.edition ? (template.edition.charAt(0).toUpperCase() + template.edition.slice(1)) : 'Unknown' }}
          </button>
        </div>
        <button @click="showEditionModal = false" class="cancel-btn">Cancel</button>
      </div>
    </div>
    
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
const showEditionModal = ref(false)
const templates = ref([])
const templatesLoading = ref(false)

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

const loadTemplates = async () => {
  templatesLoading.value = true
  try {
    const data = await characterService.getTemplates()
    templates.value = data
  } catch (err) {
    console.error('Error loading templates:', err)
    // Fallback to empty array if templates can't be loaded
    templates.value = []
  } finally {
    templatesLoading.value = false
  }
}

const createCharacter = (edition) => {
  showEditionModal.value = false
  router.push(`/characters/new?edition=${edition}`)
}

onMounted(() => {
  loadCharacters()
  loadTemplates()
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

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.create-btn {
  padding: 0.75rem 1.5rem;
  background-color: #27ae60;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background-color: #229954;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-content h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #2c3e50;
}

.edition-options {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.edition-btn {
  padding: 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
  text-transform: capitalize;
}

.edition-btn:hover {
  background-color: #2980b9;
}

.cancel-btn {
  width: 100%;
  padding: 0.75rem;
  background-color: #95a5a6;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.cancel-btn:hover {
  background-color: #7f8c8d;
}
</style>


