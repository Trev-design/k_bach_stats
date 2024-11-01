<template>
  <div class="experience-adder-container">
    <div class="added-experience-container">
      
    </div>

    <div class="experience-adder">
      <div class="input-container">
        <input type="text" class="add-experinece-input">
        <Icon/>
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
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  
  name: 'ExperienceAdder',

  components: {
    Icon
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

</style>