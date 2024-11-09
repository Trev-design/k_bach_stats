<template>
  <div class="experience-rating-container">
    <div 
      class="experience-rating"
      v-for="experience in experiences"
      :key="experience.id"
    >
      <p class="experience-descriptor">{{experience.name}}</p>
      
      <Icon 
        v-for="star in 5"
        :key="star"
        :class="`rating-star ${star <= experience.rating ? 'active' : ''}`"
        icon="fluent:star-24-filled"
      />
      
      <button class="remove-experience-button">
        <Icon
          class="remove-experience-button-icon"
          icon="mdi:remove"
        />
      </button>
    </div>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  name: 'ExperienceratingList',

  components: {
    Icon
  },

  computed: {
    experiences() {
      return this.$store.getters.newExperienceList
    }
  }
}
</script>

<style>
.experience-rating-container {

}

.experience-rating {

}

.experience-descriptor {

}

.rating-star {

}

.rating-star.active {

}

.remove-experience-button {

}

.remove-experience-button-icon {

}
</style>