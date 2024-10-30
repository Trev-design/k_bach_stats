<template>
  <section :class="`add-experience-section-container ${addingReady ? 'active' : ''}`">
    <div :class="`searchbar-container ${searchEnabled ? 'active' : ''}`">
      <SearchbarButton
        :searcEnabled="searchEnabled"
        @click="$emit('handleSearch')"
      />

      <input 
        type="text" 
        :class="`searchbar-input ${searchEnabled ? 'active' : ''}`"
        @input="$emit('update:inputValue', $event.target.value)"
      >

      <button 
        class="add-experience-submit-button"
        @click="$emit('disable')"
      >submit</button>
    </div>
  </section>
</template>

<script>
import SearchbarButton from '../components/SearchbarButton.vue'

export default {
  name: 'AddExperience',

  components: {
    SearchbarButton,
  },

  props: {
    searchEnabled: {
      type: Boolean,
      required: true
    },

    addingReady: {
      type: Boolean,
      required: true
    },

    inputValue: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.add-experience-section-container {
  position: fixed;
  width: inherit;
  height: 350px;
  background: #191f4d;
  z-index: 1;
  transform: scale(0, 0);
  transition: transform 0.4s ease;
}

.add-experience-section-container.active {
  transform: scale(1, 1);
}
</style>