<template>
  <div class="experience-adder-container">
    <div class="added-experience-container">
      <ExperienceratingList/>
    </div>

    <div class="adder-container">
      <div class="experience-adder">
        <div class="input-container">
          <input id="experienceInputLabel" type="text" class="add-experinece-input">
          <label for="experienceInputLabel" class="add-experience-input-label">add experience</label>
        </div>
        
        <div class="add-self-assessment-container">
          <Icon 
            v-for="star in 5"
            :key="star"
            :class="`star-icon ${star <= ratingInput ? 'active' : ''}`"
            icon="fluent:star-24-filled"
          />
        </div>
      </div>

      <button class="add-experience-button">
        <Icon
          class="experience-adder-button-icon"
          icon="mdi:add"
        />
      </button>
    </div>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'
import ExperienceRatingList from './ExperienceRatingList.vue'

export default {
  
  name: 'ExperienceAdder',

  components: {
    Icon,
    ExperienceRatingList,
  },

  data() {
    return {
      newExperiences: [],
      ratingInput: 0
    }
  },

  props: {
    isUfold: {
      type: Boolean,
      required: true
    }
  }
}
</script>

<style>
.experience-adder-container {
  width: 90%;
}

.added-experience-container {
  width: 100%;
  min-height: 120px;
  border: 2px solid aqua;
  border-radius: 5px;
}

.adder-container {
  display: flex;
  flex-direction: row;
  justify-content: end;
  align-items: center;
}

.experience-adder {
  display: flex;
  flex-direction: column;
  justify-content: end;
  align-items: end;
  height: 4rem;
  padding: 1rem 0;
}
</style>