<template>
  <div class="character-detail">
    <div v-if="loading" class="loading">Loading character...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="character || isCreating" class="character-wrapper">
      <div class="header-actions">
        <button @click="goBack" class="back-btn">← Back to Characters</button>
        <div class="header-right-actions">
          <button @click="toggleLock" class="lock-btn" :class="{ locked: isLocked }">
            {{ isLocked ? '🔒 Customization locked' : 'Customization unlocked' }}
          </button>
          <button @click="saveCharacter" class="save-btn" :disabled="saving">
            {{ saving ? (isCreating ? 'Creating...' : 'Saving...') : (isCreating ? 'Create' : 'Save Changes') }}
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
          <table class="data-table">
            <thead>
              <tr>
                <th style="width: 200px;">Type</th>
                <th>Value</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(aspect, index) in aspects" 
                :key="aspect.id || index"
                draggable="true"
                @dragstart="(e) => handleDragStart(e, index, 'aspect')"
                @dragend="handleDragEnd"
                @dragover="handleDragOver"
                @dragleave="handleDragLeave"
                @drop="(e) => handleDrop(e, index, 'aspect')"
                class="draggable-row"
              >
                <td>
                  <input v-model="aspect.type" type="text" class="form-input" placeholder="Aspect type" :disabled="isLocked" />
                </td>
                <td>
                  <input v-model="aspect.value" type="text" class="form-input" placeholder="Aspect value" :disabled="isLocked" />
                </td>
              </tr>
              <tr v-if="aspects.length === 0">
                <td colspan="2" style="text-align: center; color: #6c757d; padding: 2rem;">
                  No aspects added yet
                </td>
              </tr>
            </tbody>
          </table>
          <button v-if="!isLocked" @click="addAspect" class="add-btn">+ Add Aspect</button>
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
            <div 
              class="skill-list-compact"
              :data-group-id="group.id"
              @dragover="handleSkillDragOver"
              @dragleave="handleSkillDragLeave"
              @drop="handleSkillDrop($event, group.id)"
            >
              <div 
                v-for="(skill, index) in group.skills" 
                :key="`${group.id}-${index}`" 
                class="skill-item-compact"
                :draggable="!isLocked"
                @dragstart="handleSkillDragStart($event, group.id, index)"
                @dragend="handleSkillDragEnd"
              >
                <input
                  v-model="group.skills[index]"
                  type="text"
                  class="form-input skill-input-compact"
                  :placeholder="`Skill name`"
                  :disabled="isLocked"
                  @mousedown.stop
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
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(stunt, index) in stunts" 
                :key="stunt.id"
                draggable="true"
                @dragstart="(e) => handleDragStart(e, index, 'stunt')"
                @dragend="handleDragEnd"
                @dragover="handleDragOver"
                @dragleave="handleDragLeave"
                @drop="(e) => handleDrop(e, index, 'stunt')"
                class="draggable-row"
              >
                <td>
                  <input v-model="stunt.name" type="text" class="form-input" placeholder="Stunt name" :disabled="isLocked" />
                </td>
                <td>
                  <textarea 
                    v-model="stunt.description" 
                    class="form-textarea" 
                    rows="2" 
                    placeholder="Stunt description" 
                    :disabled="isLocked"
                    @input="autoResizeTextarea"
                    style="resize: none; overflow: hidden; min-height: 2.5rem;"
                  ></textarea>
                </td>
              </tr>
              <tr v-if="stunts.length === 0">
                <td colspan="2" style="text-align: center; color: #6c757d; padding: 2rem;">
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
            <div v-for="(stressData, idx) in (editedCharacter.stress || [])" :key="stressData.type || idx" class="stress-type-line">
              <div class="stress-type-header">
                <input
                  :value="stressData.type"
                  type="text"
                  class="form-input stress-type-input"
                  @blur="(e) => updateStressTypeName(idx, e.target.value)"
                  :disabled="isLocked"
                />
                <button v-if="!isLocked" @click="removeStressType(idx)" class="btn-icon btn-remove" title="Remove stress type">
                  ×
                </button>
              </div>
              <div class="stress-boxes-line">
                <div v-for="(box, boxIndex) in (stressData.boxes || [])" :key="boxIndex" class="stress-box-checkbox">
                  <input
                    type="checkbox"
                    v-model="box.isFilled"
                    :id="`${stressData.type}-${boxIndex}`"
                    class="stress-checkbox"
                  />
                  <label :for="`${stressData.type}-${boxIndex}`" class="stress-box-label">
                    <input
                      v-model.number="box.size"
                      type="number"
                      min="1"
                      class="stress-box-size-input"
                      :class="{ 'locked-input': isLocked }"
                      @click="!isLocked && $event.stopPropagation()"
                      :disabled="isLocked"
                    />
                  </label>
                  <button 
                    v-if="!isLocked && boxIndex === (stressData.boxes || []).length - 1" 
                    @click="removeStressBox(idx, boxIndex)" 
                    class="btn-icon btn-remove stress-box-remove" 
                    title="Remove box"
                  >
                    ×
                  </button>
                </div>
                <button v-if="!isLocked" @click="addStressBox(idx)" class="add-btn-small">+</button>
              </div>
            </div>

            <button v-if="!isLocked" @click="addStressType" class="add-btn" style="margin-top: 0.5rem;">
              + Add Stress Type
            </button>
          </div>
        </div>

        <!-- Consequences -->
        <div class="form-section">
          <h2 class="section-header">Consequences</h2>
          <table class="data-table">
            <thead>
              <tr>
                <template v-if="isLocked">
                  <th style="width: 150px;">Type</th>
                  <th>Description</th>
                  <th style="width: 110px;">Status</th>
                </template>
                <template v-else>
                  <th style="width: 100px;">Type</th>
                  <th style="width: 70px;">Size</th>
                  <th>Description</th>
                  <th style="width: 110px;">Status</th>
                  <th style="width: 60px;">Actions</th>
                </template>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(consequence, index) in consequences" 
                :key="consequence.id || index"
                draggable="true"
                @dragstart="(e) => handleDragStart(e, index, 'consequence')"
                @dragend="handleDragEnd"
                @dragover="handleDragOver"
                @dragleave="handleDragLeave"
                @drop="(e) => handleDrop(e, index, 'consequence')"
                class="draggable-row"
              >
                <template v-if="isLocked">
                  <td>
                    {{ consequence.type }} (-{{ consequence.size }})
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
                </template>
                <template v-else>
                  <td>
                    <input
                      v-model="consequence.type"
                      type="text"
                      class="form-input"
                      placeholder="Type"
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
                      <button @click="removeConsequence(index)" class="btn-icon btn-remove" title="Remove consequence">
                        ×
                      </button>
                    </div>
                  </td>
                </template>
              </tr>
              <tr v-if="consequences.length === 0">
                <td :colspan="isLocked ? 3 : 5" style="text-align: center; color: #6c757d; padding: 2rem;">
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
  addAspect,
  removeAspect,
  handleFileSelect,
  removeImage,
  getImageUrl,
  goBack,
  activeTab,
  isLocked,
  updateRefresh,
  toggleLock,
  autoResizeTextarea,
  handleDragStart,
  handleDragEnd,
  handleDragOver,
  handleDragLeave,
  handleDrop,
  handleSkillDragStart,
  handleSkillDragEnd,
  handleSkillDragOver,
  handleSkillDragLeave,
  handleSkillDrop,
  isCreating
} = useCharacterDetail()
</script>
