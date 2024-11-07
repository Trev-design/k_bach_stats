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

<style scoped>
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
  width: 100%;
}

.experience-adder {
  display: flex;
  flex-direction: column;
  justify-content: end;
  align-items: end;
  height: 4rem;
  padding: 1rem 0;
  width: inherit;
}

.input-container {
  position: relative;
  display: flex;
  flex-direction: column;
  height: 1.6rem;
  border: 2px solid aqua;
  border-radius: 5px;
  width: 100%;
}

.add-experinece-input {
  height: inherit;
  background: inherit;
  color: aquamarine;
  border: none;
  outline: none;
  margin: 0 .4rem;
}

.add-experience-input-label {
  position: absolute;
  color: aqua;
  background: #191f4d;
  padding: 0 .15rem;
  left: .5rem;
  transform: scale(1);
  transition: all 0.3s ease;
}

.add-experinece-input:focus + .add-experience-input-label,
.add-experinece-input:valid + .add-experience-input-label {
  transform: translate(0, -0.8rem) scale(0.75);
}

.add-self-assessment-container {
  width: 100%;
  display: flex;
  flex-direction: row;
  justify-content: start;
  align-items: center;
}

.star-icon {
  width: 1.1rem;
  height: 1.1rem;
  padding-top: .4rem;
  padding-left: .2rem;
  color: rgb(80, 80, 140);
  transition: all .3s ease;
  cursor: pointer;
}

.star-icon.active{
  color: rgb(236, 236, 77);
}

.add-experience-button {
  background: #191f4d;
  outline: none;
  border: none;
  cursor: pointer;
}

.experience-adder-button-icon {
  color: aquamarine;
  width: 1.6rem;
  height: 1.6rem;
  margin-top: -0.3rem;
}
</style>