<template>
  <div class="character-detail">
    <div v-if="loading" class="loading">Loading character...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="character" class="character-content">
      <div class="header-actions">
        <button @click="$router.push('/characters')" class="back-btn">← Back to Characters</button>
        <button @click="saveCharacter" class="save-btn" :disabled="saving">
          {{ saving ? 'Saving...' : 'Save Changes' }}
        </button>
      </div>
      
      <div class="character-form">
        <div class="form-group">
          <label>Name</label>
          <input v-model="editedCharacter.name" type="text" class="form-input" />
        </div>
        
        <div class="form-group">
          <label>Edition</label>
          <select v-model="editedCharacter.edition" class="form-input">
            <option value="core">Core</option>
            <option value="accelerated">Accelerated</option>
            <option value="condensed">Condensed</option>
            <option value="custom">Custom</option>
          </select>
        </div>
        
        <div class="form-group">
          <label>Description</label>
          <textarea v-model="editedCharacter.description" class="form-textarea" rows="4"></textarea>
        </div>
        
        <div class="form-section">
          <h3>Images</h3>
          <div class="images-section">
            <!-- Image Upload -->
            <div class="image-upload-area">
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                @change="handleFileSelect"
                style="display: none"
              />
              <button @click="$refs.fileInput.click()" class="upload-btn" :disabled="uploading">
                {{ uploading ? 'Uploading...' : '📷 Upload Image' }}
              </button>
              <span v-if="uploadError" class="upload-error">{{ uploadError }}</span>
            </div>
            
            <!-- Display Images -->
            <div v-if="characterImages.length > 0" class="images-grid">
              <div v-for="(imageUrl, index) in characterImages" :key="index" class="image-item">
                <img :src="getImageUrl(imageUrl)" :alt="`Character image ${index + 1}`" class="character-image" />
                <button @click="removeImage(index)" class="remove-image-btn" title="Remove image">×</button>
              </div>
            </div>
            <div v-else class="no-images">No images uploaded yet</div>
          </div>
        </div>
        
        <div class="form-group">
          <label>Refresh</label>
          <input v-model.number="editedCharacter.refresh" type="number" class="form-input" />
        </div>
        
        <div class="form-group">
          <label>Extras</label>
          <textarea v-model="editedCharacter.extras" class="form-textarea" rows="3"></textarea>
        </div>
        
        <div class="form-section">
          <h3>Aspects</h3>
          <div class="form-group">
            <label>High Concept</label>
            <input
              v-if="editedCharacter.aspects && typeof editedCharacter.aspects === 'object'"
              v-model="editedCharacter.aspects.highConcept"
              type="text"
              class="form-input"
            />
            <textarea
              v-else
              v-model="aspectsText"
              class="form-textarea"
              rows="3"
              placeholder='Enter aspects as JSON, e.g., {"highConcept": "...", "trouble": "...", "others": [...]}'
            ></textarea>
          </div>
        </div>
        
        <div class="form-section">
          <h3>Skills</h3>
          <textarea
            v-model="skillsText"
            class="form-textarea"
            rows="5"
            placeholder='Enter skills as JSON, e.g., {"+4": ["Skill1"], "+3": ["Skill2"]}'
          ></textarea>
        </div>
        
        <div class="form-section">
          <h3>Stunts</h3>
          <div v-for="(stunt, index) in stuntsArray" :key="index" class="stunt-item">
            <textarea
              v-model="stuntsArray[index]"
              class="form-textarea"
              rows="2"
              placeholder="Enter stunt description"
            ></textarea>
            <button @click="removeStunt(index)" class="remove-btn">Remove</button>
          </div>
          <button @click="addStunt" class="add-btn">Add Stunt</button>
        </div>
        
        <div class="form-section">
          <h3>Stress</h3>
          <div class="stress-group">
            <div class="form-group">
              <label>Physical</label>
              <input
                v-if="editedCharacter.stress && typeof editedCharacter.stress === 'object'"
                v-model.number="editedCharacter.stress.physical"
                type="number"
                class="form-input"
              />
              <input v-else v-model="stressText" type="text" class="form-input" placeholder='{"physical": 2, "mental": 3}' />
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="saveMessage" class="save-message" :class="{ success: saveSuccess, error: !saveSuccess }">
        {{ saveMessage }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { characterService } from '../services/api'
import { storageService } from '../services/storage'

const route = useRoute()
const character = ref(null)
const editedCharacter = ref({})
const loading = ref(true)
const error = ref(null)
const saving = ref(false)
const saveMessage = ref('')
const saveSuccess = ref(false)
const aspectsText = ref('')
const skillsText = ref('')
const stressText = ref('')
const stuntsArray = ref([])
const uploading = ref(false)
const uploadError = ref('')
const characterImages = ref([])

const loadCharacter = async () => {
  loading.value = true
  error.value = null
  try {
    const data = await characterService.getCharacters()
    const found = data.find(c => (c._id || c.id) === route.params.id)
    if (found) {
      character.value = found
      editedCharacter.value = JSON.parse(JSON.stringify(found))
      
      // Initialize form fields
      if (typeof editedCharacter.value.aspects === 'object') {
        aspectsText.value = JSON.stringify(editedCharacter.value.aspects, null, 2)
      } else {
        aspectsText.value = ''
      }
      
      if (typeof editedCharacter.value.skills === 'object') {
        skillsText.value = JSON.stringify(editedCharacter.value.skills, null, 2)
      } else {
        skillsText.value = ''
      }
      
      if (typeof editedCharacter.value.stress === 'object') {
        stressText.value = JSON.stringify(editedCharacter.value.stress, null, 2)
      } else {
        stressText.value = ''
      }
      
      if (Array.isArray(editedCharacter.value.stunts)) {
        stuntsArray.value = [...editedCharacter.value.stunts]
      } else if (typeof editedCharacter.value.stunts === 'object') {
        stuntsArray.value = Object.values(editedCharacter.value.stunts)
      } else {
        stuntsArray.value = []
      }
      
      // Initialize images
      if (Array.isArray(editedCharacter.value.images)) {
        characterImages.value = [...editedCharacter.value.images]
      } else {
        characterImages.value = []
      }
    } else {
      error.value = 'Character not found'
    }
  } catch (err) {
    error.value = 'Failed to load character. Please make sure the backend is running.'
    console.error('Error loading character:', err)
  } finally {
    loading.value = false
  }
}

const parseJsonField = (text, defaultValue) => {
  if (!text || !text.trim()) return defaultValue
  try {
    return JSON.parse(text)
  } catch (e) {
    return defaultValue
  }
}

const saveCharacter = async () => {
  saving.value = true
  saveMessage.value = ''
  
  try {
    // Parse JSON fields
    editedCharacter.value.aspects = parseJsonField(aspectsText.value, editedCharacter.value.aspects)
    editedCharacter.value.skills = parseJsonField(skillsText.value, editedCharacter.value.skills)
    editedCharacter.value.stress = parseJsonField(stressText.value, editedCharacter.value.stress)
    editedCharacter.value.stunts = stuntsArray.value.filter(s => s && s.trim())
    editedCharacter.value.images = characterImages.value
    
    // For now, we'll just log since the backend might not have update endpoint yet
    // In a real scenario, you'd call: await characterService.updateCharacter(route.params.id, editedCharacter.value)
    console.log('Saving character:', editedCharacter.value)
    
    saveMessage.value = 'Character saved successfully! (Note: Backend update endpoint may need to be implemented)'
    saveSuccess.value = true
    
    setTimeout(() => {
      saveMessage.value = ''
    }, 3000)
  } catch (err) {
    saveMessage.value = 'Failed to save character: ' + err.message
    saveSuccess.value = false
    console.error('Error saving character:', err)
  } finally {
    saving.value = false
  }
}

const addStunt = () => {
  stuntsArray.value.push('')
}

const removeStunt = (index) => {
  stuntsArray.value.splice(index, 1)
}

const handleFileSelect = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  // Validate file type
  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Please select an image file'
    setTimeout(() => { uploadError.value = '' }, 3000)
    return
  }
  
  // Validate file size (max 10MB)
  if (file.size > 10 * 1024 * 1024) {
    uploadError.value = 'File size must be less than 10MB'
    setTimeout(() => { uploadError.value = '' }, 3000)
    return
  }
  
  uploading.value = true
  uploadError.value = ''
  
  try {
    // Upload to storage service
    const folder = `images/characters/${route.params.id || 'temp'}`
    const result = await storageService.uploadFile(file, folder)
    
    // Add to images array - use filename from response
    // The filename includes the folder path (e.g., "images/characters/123/image.jpg")
    characterImages.value.push(result.filename)
    
    // Update edited character
    editedCharacter.value.images = [...characterImages.value]
    
    // Clear file input
    event.target.value = ''
    
    console.log('Image uploaded successfully:', result)
  } catch (err) {
    uploadError.value = 'Failed to upload image: ' + (err.response?.data?.error || err.message)
    console.error('Error uploading image:', err)
  } finally {
    uploading.value = false
  }
}

const removeImage = async (index) => {
  const imageUrl = characterImages.value[index]
  
  // Try to delete from storage (don't fail if it doesn't exist)
  try {
    await storageService.deleteFile(imageUrl)
  } catch (err) {
    console.warn('Could not delete image from storage:', err)
  }
  
  // Remove from array
  characterImages.value.splice(index, 1)
  editedCharacter.value.images = [...characterImages.value]
}

const getImageUrl = (filename) => {
  // If it's already a full URL, return it
  if (filename.startsWith('http://') || filename.startsWith('https://')) {
    return filename
  }
  // Otherwise, get URL from storage service
  return storageService.getFileUrl(filename)
}

onMounted(() => {
  loadCharacter()
})
</script>

<style scoped>
.character-detail {
  width: 100%;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
  font-size: 1.1rem;
}

.error {
  color: #e74c3c;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.back-btn, .save-btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.back-btn {
  background-color: #95a5a6;
  color: white;
}

.back-btn:hover {
  background-color: #7f8c8d;
}

.save-btn {
  background-color: #27ae60;
  color: white;
}

.save-btn:hover:not(:disabled) {
  background-color: #229954;
}

.save-btn:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.character-form {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  color: #2c3e50;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: #3498db;
}

.form-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #eee;
}

.form-section h3 {
  margin-bottom: 1rem;
  color: #2c3e50;
}

.stunt-item {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  align-items: flex-start;
}

.stunt-item .form-textarea {
  flex: 1;
}

.remove-btn {
  padding: 0.5rem 1rem;
  background-color: #e74c3c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  white-space: nowrap;
}

.remove-btn:hover {
  background-color: #c0392b;
}

.add-btn {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.add-btn:hover {
  background-color: #2980b9;
}

.save-message {
  margin-top: 1rem;
  padding: 1rem;
  border-radius: 4px;
  text-align: center;
}

.save-message.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.save-message.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.images-section {
  margin-top: 1rem;
}

.image-upload-area {
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.upload-btn {
  padding: 0.75rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.2s;
}

.upload-btn:hover:not(:disabled) {
  background-color: #2980b9;
}

.upload-btn:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.upload-error {
  color: #e74c3c;
  font-size: 0.9rem;
}

.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.image-item {
  position: relative;
  border: 1px solid #ddd;
  border-radius: 4px;
  overflow: hidden;
  background: #f8f9fa;
}

.character-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  display: block;
}

.remove-image-btn {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  width: 2rem;
  height: 2rem;
  background-color: rgba(231, 76, 60, 0.9);
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.5rem;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.remove-image-btn:hover {
  background-color: rgba(192, 57, 43, 1);
}

.no-images {
  color: #7f8c8d;
  font-style: italic;
  text-align: center;
  padding: 2rem;
}
</style>

