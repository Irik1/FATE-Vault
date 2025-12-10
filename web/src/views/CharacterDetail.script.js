import { ref, computed, onMounted, nextTick, watch, onUpdated } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { characterService } from '../services/api'
import { storageService } from '../services/storage'

export function useCharacterDetail() {
  const route = useRoute()
  const router = useRouter()
  
  const character = ref(null)
  const editedCharacter = ref({})
  const loading = ref(true)
  const error = ref(null)
  const saving = ref(false)
  const saveMessage = ref('')
  const saveSuccess = ref(false)
  const uploading = ref(false)
  const uploadError = ref('')
  const characterImages = ref([])

  // Skills management - array of skill groups (allows duplicate levels)
  const skills = ref([])

  // Aspects management - now always an array
  const aspects = ref([])

  // Stunts management
  const stunts = ref([])

  // Consequences management - now always an array
  const consequences = ref([])
  const activeTab = ref('main')
  const isLocked = ref(false)

  const loadCharacter = async () => {
    loading.value = true
    error.value = null
    try {
      const data = await characterService.getCharacters()
      const found = data.find(c => (c._id || c.id) === route.params.id)
      if (found) {
        character.value = found
        editedCharacter.value = JSON.parse(JSON.stringify(found))

        // --- NEW: assume aspects is always an array of {type, value} ---
        if (Array.isArray(editedCharacter.value.aspects)) {
          aspects.value = editedCharacter.value.aspects.map((a, i) => ({
            id: a.id ?? i,
            type: a.type || '',
            value: a.value || ''
          }))
        } else if (editedCharacter.value.aspects && typeof editedCharacter.value.aspects === 'object') {
          // Convert old format to new array format for backward compatibility
          const oldAspects = editedCharacter.value.aspects
          aspects.value = []
          
          if (oldAspects.highConcept) {
            aspects.value.push({
              id: 0,
              type: 'High Concept',
              value: oldAspects.highConcept
            })
          }
          if (oldAspects.trouble) {
            aspects.value.push({
              id: 1,
              type: 'Trouble',
              value: oldAspects.trouble
            })
          }
          if (Array.isArray(oldAspects.others)) {
            oldAspects.others.forEach((other, idx) => {
              if (other) {
                aspects.value.push({
                  id: aspects.value.length,
                  type: 'Other',
                  value: other
                })
              }
            })
          }
        } else {
          aspects.value = []
        }
        
        // --- NEW: assume stress is always an array of {type, boxes} ---
        if (!Array.isArray(editedCharacter.value.stress)) {
          // default array if absent or wrong type
          editedCharacter.value.stress = [
            { type: 'physical', boxes: [] },
            { type: 'mental', boxes: [] }
          ]
        } else {
          // normalize boxes entries (ensure isFilled property exists)
          editedCharacter.value.stress = editedCharacter.value.stress.map(s => {
            return {
              type: s.type || 'unknown',
              boxes: Array.isArray(s.boxes)
                ? s.boxes.map(box => {
                  if ('current' in box && !('isFilled' in box)) {
                    return { size: box.size, isFilled: box.current > 0 }
                  }
                  return box
                })
                : []
            }
          })
        }
        
        // --- NEW: assume skills is already an array of groups ---
        if (Array.isArray(editedCharacter.value.skills)) {
          skills.value = editedCharacter.value.skills.map((group, index) => ({
            id: group.id ?? index,
            level: group.level ?? '+0',
            skills: Array.isArray(group.skills) ? [...group.skills] : []
          }))
        } else {
          // fallback empty array
          skills.value = []
        }
        
        // --- NEW: assume stunts is always an array of {name,description} ---
        if (Array.isArray(editedCharacter.value.stunts)) {
          stunts.value = editedCharacter.value.stunts.map((s, i) => ({
            id: i,
            name: s.name || '',
            description: s.description || ''
          }))
        } else {
          stunts.value = []
        }
        
        // --- NEW: consequences expected as array already ---
        if (Array.isArray(editedCharacter.value.consequences)) {
          consequences.value = editedCharacter.value.consequences.map((c, i) => ({
            id: c.id ?? i,
            type: c.type || 'minor',
            size: c.size || 2,
            description: c.description || '',
            status: c.status || 'none',
            ...c
          }))
        } else {
          consequences.value = []
        }
        
        // Initialize images
        if (Array.isArray(editedCharacter.value.images)) {
          characterImages.value = [...editedCharacter.value.images]
        } else {
          characterImages.value = []
        }
        
        // Initialize refresh - handle both old number format and new object format
        if (editedCharacter.value.refresh && typeof editedCharacter.value.refresh === 'object') {
          // Already in object format, ensure both fields exist
          if (typeof editedCharacter.value.refresh.current !== 'number') {
            editedCharacter.value.refresh.current = editedCharacter.value.refresh.max || 3
          }
          if (typeof editedCharacter.value.refresh.max !== 'number') {
            editedCharacter.value.refresh.max = editedCharacter.value.refresh.current || 3
          }
        } else if (typeof editedCharacter.value.refresh === 'number') {
          // Convert old number format to object format
          editedCharacter.value.refresh = {
            current: editedCharacter.value.refresh,
            max: editedCharacter.value.refresh
          }
        } else {
          // Default
          editedCharacter.value.refresh = {
            current: 3,
            max: 3
          }
        }

        // Initialize locked mode based on playMode
        isLocked.value = editedCharacter.value.playMode === true
      } else {
        error.value = 'Character not found'
      }
    } catch (err) {
      error.value = 'Failed to load character. Please make sure the backend is running.'
      console.error('Error loading character:', err)
    } finally {
      loading.value = false
      // Resize stunt textareas after loading is complete and DOM is updated
      nextTick(() => {
        resizeAllStuntTextareas()
      })
    }
  }

  const saveCharacter = async () => {
    saving.value = true
    saveMessage.value = ''

    try {
      // --- NEW: persist skills as an array (no conversion to object) ---
      editedCharacter.value.skills = skills.value.map(group => ({
        level: group.level,
        skills: Array.isArray(group.skills) ? group.skills.filter(s => typeof s === 'string' && s.trim()) : []
      }))

      // --- NEW: persist aspects as an array of {type, value} ---
      editedCharacter.value.aspects = aspects.value.map(a => ({
        type: a.type || '',
        value: a.value || ''
      }))

      // --- NEW: persist stunts as an array of {name,description} ---
      editedCharacter.value.stunts = stunts.value.map(s => ({
        name: s.name || '',
        description: s.description || ''
      }))

      // consequences is already an array; persist the cleaned version
      editedCharacter.value.consequences = consequences.value.map(c => ({
        type: c.type || 'minor',
        size: c.size || 2,
        description: c.description || '',
        status: c.status || 'none'
      }))

      // Save images
      editedCharacter.value.images = characterImages.value

      // Ensure stress stays as array
      if (!Array.isArray(editedCharacter.value.stress)) {
        editedCharacter.value.stress = [
          { type: 'physical', boxes: [] },
          { type: 'mental', boxes: [] }
        ]
      }

      // Save character
      await characterService.updateCharacter(route.params.id, editedCharacter.value)
      
      saveMessage.value = 'Character saved successfully!'
      saveSuccess.value = true
      
      setTimeout(() => {
        saveMessage.value = ''
      }, 3000)
    } catch (err) {
      saveMessage.value = 'Failed to save character: ' + (err.response?.data?.error || err.message)
      saveSuccess.value = false
      console.error('Error saving character:', err)
    } finally {
      saving.value = false
    }
  }

  // Skills management - array of skill groups (allows duplicate levels)
  const addSkill = (groupId) => {
    const group = skills.value.find(g => g.id === groupId)
    if (group) {
      group.skills.push('')
    }
  }

  const removeSkill = (groupId, index) => {
    const group = skills.value.find(g => g.id === groupId)
    if (group) {
      group.skills.splice(index, 1)
    }
  }

  const addSkillLevel = () => {
    const newId = skills.value.length > 0 
      ? Math.max(...skills.value.map(g => g.id || 0)) + 1 
      : 0
    skills.value.push({
      id: newId,
      level: '+0',
      skills: []
    })
  }

  const removeSkillLevel = (groupId) => {
    const index = skills.value.findIndex(g => g.id === groupId)
    if (index !== -1) {
      skills.value.splice(index, 1)
    }
  }

  const updateSkillLevel = (groupId, newLevel) => {
    if (!newLevel) return
    // Validate new level format (should be like +4, -2, etc.)
    if (!/^[+-]?\d+$/.test(newLevel)) {
      return // Invalid format, don't update
    }
    const group = skills.value.find(g => g.id === groupId)
    if (group) {
      group.level = newLevel
    }
  }

  // Get skill level description
  const getSkillLevelDescription = (level) => {
    const num = parseInt(level.replace(/[+-]/, ''))
    const isNegative = level.startsWith('-')
    const value = isNegative ? -num : num
    
    if (value >= 8) return 'Legendary'
    if (value === 7) return 'Epic'
    if (value === 6) return 'Fantastic'
    if (value === 5) return 'Superb'
    if (value === 4) return 'Great'
    if (value === 3) return 'Good'
    if (value === 2) return 'Fair'
    if (value === 1) return 'Average'
    if (value === 0) return 'Mediocre'
    if (value === -1) return 'Poor'
    return 'Terrible'
  }

  // Get sorted skills (biggest to lowest)
  const sortedSkills = computed(() => {
    return [...skills.value].sort((a, b) => {
      const numA = parseInt(a.level.replace(/[+-]/, ''))
      const numB = parseInt(b.level.replace(/[+-]/, ''))
      const isNegA = a.level.startsWith('-')
      const isNegB = b.level.startsWith('-')
      const valueA = isNegA ? -numA : numA
      const valueB = isNegB ? -numB : numB
      return valueB - valueA // Biggest to lowest
    })
  })

  // Stunts management
  const addStunt = () => {
    const newId = stunts.value.length > 0 
      ? Math.max(...stunts.value.map(s => s.id)) + 1 
      : 0
    stunts.value.push({ id: newId, name: '', description: '' })
    // Resize textareas after adding new stunt
    nextTick(() => {
      resizeAllStuntTextareas()
    })
  }

  const removeStunt = (index) => {
    stunts.value.splice(index, 1)
  }

  // Consequences management (array)
  const addConsequence = () => {
    const newId = consequences.value.length > 0
      ? Math.max(...consequences.value.map(c => c.id || 0)) + 1
      : 0
    consequences.value.push({
      id: newId,
      type: 'minor',
      size: 2,
      description: '',
      status: 'none'
    })
  }

  const removeConsequence = (index) => {
    consequences.value.splice(index, 1)
  }

  // Stress management (array of {type, boxes})
  const addStressBox = (stressIndex) => {
    if (!editedCharacter.value.stress) editedCharacter.value.stress = []
    const s = editedCharacter.value.stress[stressIndex]
    if (s) {
      s.boxes = s.boxes || []
      s.boxes.push({ size: 1, isFilled: false })
    } else {
      editedCharacter.value.stress.push({ type: 'new', boxes: [{ size: 1, isFilled: false }] })
    }
  }

  const removeStressBox = (stressIndex, boxIndex) => {
    const s = editedCharacter.value.stress && editedCharacter.value.stress[stressIndex]
    if (s && Array.isArray(s.boxes)) s.boxes.splice(boxIndex, 1)
  }

  const addStressType = () => {
    if (!editedCharacter.value.stress) editedCharacter.value.stress = []
    let base = 'new'
    let idx = 1
    while (editedCharacter.value.stress.find(s => s.type === base + (idx === 1 ? '' : idx))) idx++
    const name = base + (idx === 1 ? '' : idx)
    editedCharacter.value.stress.push({ type: name, boxes: [] })
  }

  const removeStressType = (stressIndex) => {
    if (editedCharacter.value.stress && editedCharacter.value.stress[stressIndex]) {
      editedCharacter.value.stress.splice(stressIndex, 1)
    }
  }

  const updateStressTypeName = (stressIndex, newType) => {
    if (!newType) return
    const s = editedCharacter.value.stress && editedCharacter.value.stress[stressIndex]
    if (!s) return
    if (s.type === newType) return
    const existing = editedCharacter.value.stress.find((st, i) => st.type === newType && i !== stressIndex)
    if (existing) {
      existing.boxes = [...(existing.boxes || []), ...(s.boxes || [])]
      editedCharacter.value.stress.splice(stressIndex, 1)
    } else {
      s.type = newType
    }
  }

  // Aspects management (array)
  const addAspect = () => {
    const newId = aspects.value.length > 0 
      ? Math.max(...aspects.value.map(a => a.id || 0)) + 1 
      : 0
    aspects.value.push({ id: newId, type: '', value: '' })
  }

  const removeAspect = (index) => {
    aspects.value.splice(index, 1)
  }

  // Drag and drop handlers for reordering
  const draggedIndex = ref(null)
  const draggedItemType = ref(null) // 'aspect', 'stunt', or 'consequence'

  const handleDragStart = (event, index, itemType) => {
    draggedIndex.value = index
    draggedItemType.value = itemType
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/html', event.target)
    const row = event.currentTarget
    row.style.opacity = '0.5'
  }

  const handleDragEnd = (event) => {
    const row = event.currentTarget
    row.style.opacity = '1'
    draggedIndex.value = null
    draggedItemType.value = null
  }

  const handleDragOver = (event) => {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
    
    const targetRow = event.currentTarget
    const allRows = Array.from(targetRow.closest('tbody').querySelectorAll('tr[draggable="true"]'))
    const targetIndex = allRows.indexOf(targetRow)
    
    if (targetIndex === -1) return
    
    // Remove previous drag-over class from all rows
    allRows.forEach(row => row.classList.remove('drag-over'))
    
    // Add drag-over class to target row
    if (draggedIndex.value !== null && draggedIndex.value !== targetIndex) {
      targetRow.classList.add('drag-over')
    }
  }

  const handleDragLeave = (event) => {
    event.currentTarget.classList.remove('drag-over')
  }

  const handleDrop = (event, dropIndex, itemType) => {
    event.preventDefault()
    event.currentTarget.classList.remove('drag-over')
    
    if (draggedIndex.value === null || draggedItemType.value !== itemType) {
      return
    }
    
    const sourceIndex = draggedIndex.value
    
    if (sourceIndex === dropIndex) {
      return
    }
    
    // Reorder based on item type
    if (itemType === 'aspect') {
      const item = aspects.value.splice(sourceIndex, 1)[0]
      aspects.value.splice(dropIndex, 0, item)
    } else if (itemType === 'stunt') {
      const item = stunts.value.splice(sourceIndex, 1)[0]
      stunts.value.splice(dropIndex, 0, item)
    } else if (itemType === 'consequence') {
      const item = consequences.value.splice(sourceIndex, 1)[0]
      consequences.value.splice(dropIndex, 0, item)
    }
    
    draggedIndex.value = null
    draggedItemType.value = null
  }

  // Image management
  const handleFileSelect = async (event) => {
    const file = event.target.files[0]
    if (!file) return
    
    if (!file.type.startsWith('image/')) {
      uploadError.value = 'Please select an image file'
      setTimeout(() => { uploadError.value = '' }, 3000)
      return
    }
    
    if (file.size > 10 * 1024 * 1024) {
      uploadError.value = 'File size must be less than 10MB'
      setTimeout(() => { uploadError.value = '' }, 3000)
      return
    }
    
    uploading.value = true
    uploadError.value = ''
    
    try {
      const folder = `images/characters/${route.params.id || 'temp'}`
      const result = await storageService.uploadFile(file, folder)
      characterImages.value.push(result.filename)
      editedCharacter.value.images = [...characterImages.value]
      event.target.value = ''
    } catch (err) {
      uploadError.value = 'Failed to upload image: ' + (err.response?.data?.error || err.message)
      console.error('Error uploading image:', err)
    } finally {
      uploading.value = false
    }
  }

  const removeImage = async (index) => {
    const imageUrl = characterImages.value[index]
    try {
      await storageService.deleteFile(imageUrl)
    } catch (err) {
      console.warn('Could not delete image from storage:', err)
    }
    characterImages.value.splice(index, 1)
    editedCharacter.value.images = [...characterImages.value]
  }

  const getImageUrl = (filename) => {
    if (filename.startsWith('http://') || filename.startsWith('https://')) {
      return filename
    }
    return storageService.getFileUrl(filename)
  }

  const goBack = () => {
    router.push('/characters')
  }

  const updateRefresh = (field, delta) => {
    if (!editedCharacter.value.refresh) {
      editedCharacter.value.refresh = { current: 0, max: 3 }
    }
    if (field === 'current') {
      const newValue = (editedCharacter.value.refresh.current || 0) + delta
      if (newValue >= 0 && newValue <= (editedCharacter.value.refresh.max || 0)) {
        editedCharacter.value.refresh.current = newValue
      }
    } else if (field === 'max') {
      const newValue = (editedCharacter.value.refresh.max || 0) + delta
      if (newValue >= 1) {
        editedCharacter.value.refresh.max = newValue
        // Ensure current doesn't exceed max
        if (editedCharacter.value.refresh.current > newValue) {
          editedCharacter.value.refresh.current = newValue
        }
      }
    }
  }

  const toggleLock = () => {
    isLocked.value = !isLocked.value
    // When unlocking, set playMode to false
    // When locking, set playMode to true
    editedCharacter.value.playMode = isLocked.value
  }

  // Auto-resize textarea based on content
  const autoResizeTextarea = (event) => {
    const textarea = event.target
    textarea.style.height = 'auto'
    textarea.style.height = textarea.scrollHeight + 'px'
  }

  // Resize all stunt textareas (used after loading)
  const resizeAllStuntTextareas = () => {
    // Use setTimeout with nextTick to ensure DOM is fully updated and rendered
    nextTick(() => {
      setTimeout(() => {
        // Find all textareas - they should be in the stunts table
        // We'll check all textareas and resize them
        const allTextareas = document.querySelectorAll('textarea.form-textarea')
        
        allTextareas.forEach(textarea => {
          // Only resize if it's in a table (stunts table)
          if (textarea.closest('table.data-table')) {
            // Reset height to auto to get accurate scrollHeight
            textarea.style.height = 'auto'
            const scrollHeight = textarea.scrollHeight
            // Set minimum height to at least 2.5rem (40px) or scrollHeight, whichever is larger
            const minHeight = 40 // 2.5rem in pixels
            textarea.style.height = Math.max(scrollHeight, minHeight) + 'px'
          }
        })
      }, 100) // Small delay to ensure rendering is complete
    })
  }

  onMounted(() => {
    loadCharacter()
  })

  // Watch for stunts changes and resize textareas
  watch(stunts, () => {
    nextTick(() => {
      resizeAllStuntTextareas()
    })
  }, { deep: true })

  return {
    character,
    editedCharacter,
    loading,
    error,
    saving,
    saveMessage,
    saveSuccess,
    uploading,
    uploadError,
    characterImages,
    aspects,
    skills,
    stunts,
    consequences,
    saveCharacter,
    addSkill,
    removeSkill,
    addSkillLevel,
    removeSkillLevel,
    updateSkillLevel,
    getSkillLevelDescription,
    sortedSkills,
    addStunt,
    removeStunt,
    addConsequence,
    removeConsequence,
    addStressBox,
    removeStressBox,
    addStressType,
    removeStressType,
    updateStressTypeName,
    activeTab,
    addAspect,
    removeAspect,
    handleFileSelect,
    removeImage,
    getImageUrl,
    goBack,
    isLocked,
    updateRefresh,
    toggleLock,
    autoResizeTextarea,
    handleDragStart,
    handleDragEnd,
    handleDragOver,
    handleDragLeave,
    handleDrop
  }
}

