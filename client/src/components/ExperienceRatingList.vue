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

  data() {
    return {
      experiences: []
    }
  }
}
</script>

<style>

</style>