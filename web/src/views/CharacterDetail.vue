<template>
  <div class="character-detail">
    <div v-if="loading" class="loading">Loading character...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="character" class="character-wrapper">
      <div class="header-actions">
        <button @click="goBack" class="back-btn">← Back to Characters</button>
        <div class="header-right-actions">
          <button @click="isLocked = !isLocked" class="lock-btn" :class="{ locked: isLocked }">
            {{ isLocked ? '🔒 Locked' : '🔓 Unlocked' }}
          </button>
          <button @click="saveCharacter" class="save-btn" :disabled="saving">
            {{ saving ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="tabs">
        <button 
          @click="activeTab = 'main'" 
          :class="['tab-button', { active: activeTab === 'main' }]"
        >
          Main
        </button>
        <button 
          @click="activeTab = 'notes'" 
          :class="['tab-button', { active: activeTab === 'notes' }]"
        >
          Notes
        </button>
      </div>

      <!-- Main Tab -->
      <div v-if="activeTab === 'main'" class="character-form">
        <!-- Character Name and Edition at Top -->
        <div class="character-header-info">
          <h1 class="character-name">{{ editedCharacter.name || 'Unnamed Character' }}</h1>
          <span class="character-edition-label">System: {{ editedCharacter.edition || 'Unknown' }}</span>
        </div>

        <!-- Basic Information - Name, Description, Images inline -->
        <div class="form-section">
          <h2 class="section-header">Basic Information</h2>
          <div class="name-description-images-row">
            <div class="name-description-column">
              <div class="form-group">
                <label>Name</label>
                <input v-model="editedCharacter.name" type="text" class="form-input" :disabled="isLocked" />
              </div>
              <div class="form-group">
                <label>Description</label>
                <textarea v-model="editedCharacter.description" class="form-textarea" rows="4" :disabled="isLocked"></textarea>
              </div>
              <div class="form-group">
                <label>Refreshes</label>
                <div class="refresh-display">
                  <span class="refresh-text">Refreshes {{ editedCharacter.refresh?.current || 0 }}/{{ editedCharacter.refresh?.max || 0 }}</span>
                  <div class="refresh-controls">
                    <button @click="updateRefresh('current', -1)" class="refresh-btn" :disabled="(editedCharacter.refresh?.current || 0) <= 0">−</button>
                    <button @click="updateRefresh('current', 1)" class="refresh-btn" :disabled="(editedCharacter.refresh?.current || 0) >= (editedCharacter.refresh?.max || 0)">+</button>
                    <span v-if="!isLocked" class="refresh-separator">|</span>
                    <template v-if="!isLocked">
                      <button @click="updateRefresh('max', -1)" class="refresh-btn" :disabled="(editedCharacter.refresh?.max || 0) <= 1">−</button>
                      <button @click="updateRefresh('max', 1)" class="refresh-btn">+</button>
                    </template>
                  </div>
                </div>
              </div>
            </div>
            <div class="images-column">
              <div class="form-group">
                <label>Images</label>
                <div class="images-section-compact">
                  <div v-if="!isLocked" class="image-upload-area">
                    <input
                      ref="fileInput"
                      type="file"
                      accept="image/*"
                      @change="handleFileSelect"
                      style="display: none"
                    />
                    <button @click="$refs.fileInput.click()" class="upload-btn" :disabled="uploading">
                      {{ uploading ? 'Uploading...' : '📷 Upload' }}
                    </button>
                    <span v-if="uploadError" class="upload-error">{{ uploadError }}</span>
                  </div>
                  
                  <div v-if="characterImages.length > 0" class="images-grid-compact">
                    <div v-for="(imageUrl, index) in characterImages" :key="index" class="image-item-compact">
                      <img :src="getImageUrl(imageUrl)" :alt="`Character image ${index + 1}`" class="character-image-compact" />
                      <button v-if="!isLocked" @click="removeImage(index)" class="remove-image-btn-compact" title="Remove image">×</button>
                    </div>
                  </div>
                  <div v-else class="no-images-compact">No images</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Aspects -->
        <div class="form-section">
          <h2 class="section-header">Aspects</h2>
            <div class="form-group">
              <label>High Concept</label>
              <input
                v-model="editedCharacter.aspects.highConcept"
                type="text"
                class="form-input"
                placeholder="High Concept"
                :disabled="isLocked"
              />
            </div>
            <div class="form-group">
              <label>Trouble</label>
              <input
                v-model="editedCharacter.aspects.trouble"
                type="text"
                class="form-input"
                placeholder="Trouble"
                :disabled="isLocked"
              />
            </div>
          <div class="form-group">
            <label>Other Aspects</label>
            <div v-for="(aspect, index) in (editedCharacter.aspects.others || [])" :key="index" class="form-group" style="margin-bottom: 0.75rem; display: flex; gap: 0.5rem; align-items: center;">
              <input
                v-model="editedCharacter.aspects.others[index]"
                type="text"
                class="form-input"
                :placeholder="`Aspect ${index + 1}`"
                style="flex: 1;"
                :disabled="isLocked"
              />
              <button v-if="!isLocked" @click="removeAspect(index)" class="btn-icon btn-remove" title="Remove aspect">
                ×
              </button>
            </div>
            <button v-if="!isLocked" @click="addAspect" class="add-btn" style="margin-top: 0.5rem;">
              + Add Aspect
            </button>
          </div>
        </div>

        <!-- Skills -->
        <div class="form-section">
          <h2 class="section-header">Skills</h2>
          <div v-for="group in sortedSkills" :key="group.id" class="skill-group">
            <div class="skill-level-controls">
              <input
                :value="group.level"
                type="text"
                class="form-input skill-level-input"
                placeholder="+4"
                @blur="(e) => updateSkillLevel(group.id, e.target.value)"
                :disabled="isLocked"
              />
              <span class="skill-level-description">{{ getSkillLevelDescription(group.level) }}</span>
            </div>
            <div class="skill-list-compact">
              <div v-for="(skill, index) in group.skills" :key="`${group.id}-${index}`" class="skill-item-compact">
                <input
                  v-model="group.skills[index]"
                  type="text"
                  class="form-input skill-input-compact"
                  :placeholder="`Skill name`"
                  :disabled="isLocked"
                />
                <button v-if="!isLocked" @click="removeSkill(group.id, index)" class="btn-icon btn-remove" title="Remove skill">
                  ×
                </button>
              </div>
              <button v-if="!isLocked" @click="addSkill(group.id)" class="add-btn-small">+ Add Skill</button>
            </div>
            <button v-if="!isLocked" @click="removeSkillLevel(group.id)" class="btn-icon btn-remove" title="Remove skill level">
              ×
            </button>
          </div>
          <button v-if="!isLocked" @click="addSkillLevel" class="add-btn" style="margin-top: 1rem;">
            + Add Skill Level
          </button>
        </div>

        <!-- Stunts -->
        <div class="form-section">
          <h2 class="section-header">Stunts</h2>
          <table class="data-table">
            <thead>
              <tr>
                <th style="width: 200px;">Name</th>
                <th>Description</th>
                <th style="width: 80px;">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(stunt, index) in stunts" :key="stunt.id">
                <td>
                  <input v-model="stunt.name" type="text" class="form-input" placeholder="Stunt name" :disabled="isLocked" />
                </td>
                <td>
                  <textarea v-model="stunt.description" class="form-textarea" rows="2" placeholder="Stunt description" :disabled="isLocked"></textarea>
                </td>
                <td>
                  <div class="table-actions">
                    <button v-if="!isLocked" @click="removeStunt(index)" class="btn-icon btn-remove" title="Remove stunt">
                      ×
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="stunts.length === 0">
                <td colspan="3" style="text-align: center; color: #6c757d; padding: 2rem;">
                  No stunts added yet
                </td>
              </tr>
            </tbody>
          </table>
          <button v-if="!isLocked" @click="addStunt" class="add-btn">+ Add Stunt</button>
        </div>

        <!-- Stress -->
        <div class="form-section">
          <h2 class="section-header">Stress</h2>
          <div class="stress-line">
            <div v-for="(stressData, stressType) in editedCharacter.stress" :key="stressType" class="stress-type-line">
              <div class="stress-type-header">
                <input
                  :value="stressType"
                  type="text"
                  class="form-input stress-type-input"
                  @blur="(e) => updateStressTypeName(stressType, e.target.value)"
                  :disabled="isLocked"
                />
                <button v-if="!isLocked" @click="removeStressType(stressType)" class="btn-icon btn-remove" title="Remove stress type">
                  ×
                </button>
              </div>
              <div class="stress-boxes-line">
                <div v-for="(box, index) in stressData.boxes || []" :key="index" class="stress-box-checkbox">
                  <input
                    type="checkbox"
                    v-model="box.isFilled"
                    :id="`${stressType}-${index}`"
                    class="stress-checkbox"
                  />
                  <label :for="`${stressType}-${index}`" class="stress-box-label">
                    <input
                      v-model.number="box.size"
                      type="number"
                      min="1"
                      class="stress-box-size-input"
                      @click.stop
                      :disabled="isLocked"
                    />
                  </label>
                  <button 
                    v-if="!isLocked && index === (stressData.boxes || []).length - 1" 
                    @click="removeStressBox(stressType, index)" 
                    class="btn-icon btn-remove stress-box-remove" 
                    title="Remove box"
                  >
                    ×
                  </button>
                </div>
                <button v-if="!isLocked" @click="addStressBox(stressType)" class="add-btn-small">+</button>
              </div>
            </div>
          </div>
          <button v-if="!isLocked" @click="addStressType" class="add-btn" style="margin-top: 0.5rem;">
            + Add Stress Type
          </button>
        </div>

        <!-- Consequences -->
        <div class="form-section">
          <h2 class="section-header">Consequences</h2>
          <table class="data-table">
            <thead>
              <tr>
                <th style="width: 100px;">Type</th>
                <th style="width: 70px;">Size</th>
                <th>Description</th>
                <th style="width: 110px;">Status</th>
                <th style="width: 60px;">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(consequence, index) in consequences" :key="consequence.id || index">
                <td>
                  <input
                    v-model="consequence.type"
                    type="text"
                    class="form-input"
                    placeholder="Type"
                    :disabled="isLocked"
                  />
                </td>
                <td>
                  <input 
                    v-model.number="consequence.size" 
                    type="number" 
                    class="form-input consequence-size-input"
                    min="1"
                    max="9"
                    style="width: 60px;"
                    :disabled="isLocked"
                  />
                </td>
                <td>
                  <input 
                    v-model="consequence.description" 
                    type="text" 
                    class="form-input" 
                    placeholder="Consequence description"
                  />
                </td>
                <td>
                  <select v-model="consequence.status" class="form-select consequence-status-select" style="width: 100px;">
                    <option value="none">None</option>
                    <option value="active">Active</option>
                    <option value="healed">Healed</option>
                  </select>
                </td>
                <td>
                  <div class="table-actions">
                    <button v-if="!isLocked" @click="removeConsequence(index)" class="btn-icon btn-remove" title="Remove consequence">
                      ×
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="consequences.length === 0">
                <td colspan="5" style="text-align: center; color: #6c757d; padding: 2rem;">
                  No consequences added yet
                </td>
              </tr>
            </tbody>
          </table>
          <button v-if="!isLocked" @click="addConsequence" class="add-btn">+ Add Consequence</button>
        </div>

        <!-- Extras - Moved under Consequences -->
        <div class="form-section">
          <h2 class="section-header">Extras</h2>
          <div class="form-group">
            <textarea v-model="editedCharacter.extras" class="form-textarea" rows="3" :disabled="isLocked"></textarea>
          </div>
        </div>
      </div>

      <!-- Notes Tab -->
      <div v-if="activeTab === 'notes'" class="character-form">
        <div class="form-section">
          <h2 class="section-header">Character Notes</h2>
          <div class="form-group">
            <textarea 
              v-model="editedCharacter.notes" 
              class="form-textarea notes-textarea" 
              rows="20"
              placeholder="Enter your character notes here..."
            ></textarea>
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
import { useCharacterDetail } from './CharacterDetail.script.js'
import './CharacterDetail.style.css'

const {
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
  addAspect,
  removeAspect,
  handleFileSelect,
  removeImage,
  getImageUrl,
  goBack,
  activeTab,
  isLocked,
  updateRefresh
} = useCharacterDetail()
</script>
