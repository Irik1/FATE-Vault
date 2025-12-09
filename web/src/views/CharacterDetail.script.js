import { ref, computed, onMounted } from 'vue'
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
  
  // Stunts management
  const stunts = ref([])
  
  // Consequences management - object with keys like minor, moderate, severe
  const consequences = ref({})
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
        
        // Initialize aspects if they don't exist
        if (!editedCharacter.value.aspects) {
          editedCharacter.value.aspects = {
            highConcept: '',
            trouble: '',
            others: []
          }
        } else if (!editedCharacter.value.aspects.others) {
          editedCharacter.value.aspects.others = []
        }
        
        // Initialize stress if it doesn't exist
        if (!editedCharacter.value.stress) {
          editedCharacter.value.stress = {
            physical: { boxes: [] },
            mental: { boxes: [] }
          }
        } else {
          if (!editedCharacter.value.stress.physical) {
            editedCharacter.value.stress.physical = { boxes: [] }
          } else if (editedCharacter.value.stress.physical.boxes) {
            // Convert old format (with current) to new format (with isFilled)
            editedCharacter.value.stress.physical.boxes = editedCharacter.value.stress.physical.boxes.map(box => {
              if ('current' in box && !('isFilled' in box)) {
                return { size: box.size, isFilled: box.current > 0 }
              }
              return box
            })
          }
          if (!editedCharacter.value.stress.mental) {
            editedCharacter.value.stress.mental = { boxes: [] }
          } else if (editedCharacter.value.stress.mental.boxes) {
            // Convert old format (with current) to new format (with isFilled)
            editedCharacter.value.stress.mental.boxes = editedCharacter.value.stress.mental.boxes.map(box => {
              if ('current' in box && !('isFilled' in box)) {
                return { size: box.size, isFilled: box.current > 0 }
              }
              return box
            })
          }
        }
        
        // Initialize skills - convert from object to array format (allows duplicate levels)
        if (editedCharacter.value.skills && typeof editedCharacter.value.skills === 'object' && !Array.isArray(editedCharacter.value.skills)) {
          skills.value = []
          let idCounter = 0
          Object.entries(editedCharacter.value.skills).forEach(([level, skillList]) => {
            skills.value.push({
              id: idCounter++,
              level: level,
              skills: Array.isArray(skillList) ? [...skillList] : []
            })
          })
        } else if (Array.isArray(editedCharacter.value.skills)) {
          skills.value = editedCharacter.value.skills.map((group, index) => ({
            id: group.id || index,
            level: group.level || '+0',
            skills: group.skills || []
          }))
        } else {
          skills.value = []
        }
        
        // Initialize stunts
        if (Array.isArray(editedCharacter.value.stunts)) {
          stunts.value = editedCharacter.value.stunts.map((stunt, index) => ({
            id: index,
            name: typeof stunt === 'string' ? '' : Object.keys(stunt)[0] || '',
            description: typeof stunt === 'string' ? stunt : Object.values(stunt)[0] || ''
          }))
        } else if (editedCharacter.value.stunts && typeof editedCharacter.value.stunts === 'object') {
          stunts.value = Object.entries(editedCharacter.value.stunts).map(([name, description], index) => ({
            id: index,
            name,
            description: typeof description === 'string' ? description : ''
          }))
        } else {
          stunts.value = []
        }
        
        // Initialize consequences - convert from object format to array format
        if (editedCharacter.value.consequences && typeof editedCharacter.value.consequences === 'object' && !Array.isArray(editedCharacter.value.consequences)) {
          // Convert object format (minor/moderate/severe) to array format
          consequences.value = []
          Object.entries(editedCharacter.value.consequences).forEach(([type, data], index) => {
            consequences.value.push({
              id: index,
              type: type,
              size: data.size || 2,
              description: data.description || '',
              status: data.status || 'none'
            })
          })
        } else if (Array.isArray(editedCharacter.value.consequences)) {
          consequences.value = editedCharacter.value.consequences.map((consequence, index) => ({
            id: index,
            type: consequence.type || 'minor',
            size: consequence.size || 2,
            description: consequence.description || '',
            status: consequence.status || 'none',
            ...consequence
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

  const saveCharacter = async () => {
    saving.value = true
    saveMessage.value = ''
    
    try {
      // Convert skills array back to object format (group by level, merge duplicates)
      const skillsObj = {}
      skills.value.forEach(group => {
        if (group.level && Array.isArray(group.skills)) {
          if (!skillsObj[group.level]) {
            skillsObj[group.level] = []
          }
          group.skills.forEach(skill => {
            if (skill && skill.trim()) {
              skillsObj[group.level].push(skill)
            }
          })
        }
      })
      editedCharacter.value.skills = skillsObj
      
      // Convert stunts array back to object format
      const stuntsObj = {}
      stunts.value.forEach(stunt => {
        if (stunt.name && stunt.description) {
          stuntsObj[stunt.name] = stunt.description
        }
      })
      editedCharacter.value.stunts = stuntsObj
      
      // Save consequences - convert array to object format for backward compatibility
      const consequencesObj = {}
      consequences.value.forEach((consequence) => {
        if (consequence.type) {
          consequencesObj[consequence.type] = {
            size: consequence.size || 2,
            description: consequence.description || '',
            status: consequence.status || 'none'
          }
        }
      })
      editedCharacter.value.consequences = consequencesObj
      
      // Save images
      editedCharacter.value.images = characterImages.value
      
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
  }

  const removeStunt = (index) => {
    stunts.value.splice(index, 1)
  }

  // Consequences management - array of consequences with editable type
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

  // Stress management
  const addStressBox = (stressType) => {
    if (!editedCharacter.value.stress) {
      editedCharacter.value.stress = {}
    }
    if (!editedCharacter.value.stress[stressType]) {
      editedCharacter.value.stress[stressType] = { boxes: [] }
    }
    editedCharacter.value.stress[stressType].boxes.push({
      size: 1,
      isFilled: false
    })
  }

  const removeStressBox = (stressType, index) => {
    if (editedCharacter.value.stress && editedCharacter.value.stress[stressType] && editedCharacter.value.stress[stressType].boxes) {
      editedCharacter.value.stress[stressType].boxes.splice(index, 1)
    }
  }

  const addStressType = () => {
    if (!editedCharacter.value.stress) {
      editedCharacter.value.stress = {}
    }
    const newType = 'new'
    editedCharacter.value.stress[newType] = { boxes: [] }
  }

  const removeStressType = (stressType) => {
    if (editedCharacter.value.stress && editedCharacter.value.stress[stressType]) {
      delete editedCharacter.value.stress[stressType]
    }
  }

  const updateStressTypeName = (oldType, newType) => {
    if (oldType === newType || !newType) return
    if (editedCharacter.value.stress && editedCharacter.value.stress[oldType]) {
      editedCharacter.value.stress[newType] = editedCharacter.value.stress[oldType]
      delete editedCharacter.value.stress[oldType]
    }
  }

  // Aspects management
  const addAspect = () => {
    if (!editedCharacter.value.aspects) {
      editedCharacter.value.aspects = { highConcept: '', trouble: '', others: [] }
    }
    if (!editedCharacter.value.aspects.others) {
      editedCharacter.value.aspects.others = []
    }
    editedCharacter.value.aspects.others.push('')
  }

  const removeAspect = (index) => {
    if (editedCharacter.value.aspects && editedCharacter.value.aspects.others) {
      editedCharacter.value.aspects.others.splice(index, 1)
    }
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

  onMounted(() => {
    loadCharacter()
  })

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
    updateRefresh
  }
}

