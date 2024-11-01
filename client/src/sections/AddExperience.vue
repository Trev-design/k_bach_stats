<template>
  <section :class="`add-experience-section-container ${addingReady ? 'active' : ''}`">
    <div :class="`searchbar-container ${searchEnabled ? 'active' : ''}`">
      <SearchbarButton
        :searcEnabled="searchEnabled"
        :optionsUnfold="optionsUnfold"
        @click="$emit('handleSearch')"
      />

      <input 
        type="text" 
        :class="`searchbar-input ${searchEnabled ? 'active' : ''}`"
        @input="$emit('update:inputValue', $event.target.value)"
      >
    </div>

    <ExperienceAdder/>
 
    <button 
      class="add-experience-submit-button"
      @click="$emit('disable')"
    >submit</button>
  </section>
</template>

<script>
import ExperienceAdder from '../components/ExperienceAdder.vue'
import SearchbarButton from '../components/SearchbarButton.vue'

export default {
  name: 'AddExperience',

  components: {
    SearchbarButton,
    ExperienceAdder,
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
    },

    searchOptionUnfold: {
      type: Boolean,
      required: true,
    },

    addExperienceUnfold: {
      type: Boolean,
      required: true
    },

    submitUnfold: {
      type: Boolean
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

.searchbar-container {
  position: relative;
  display: flex;
  flex-direction: row;
  padding: 0.4rem 1rem;
  border-radius: 100vh;
  background: #191f4d;
  transition: all .5s ease;
}

@property --angle {
  syntax: "<angle>";
  inherits: false;
  initial-value: 360deg;
}

.searchbar-container::before {
  --angle: 360deg;
  content: '';
  position: absolute;
  top: -1.5px;
  left: -1.5px;
  bottom: -1.5px;
  right: -1.5px;
  border-radius: 100vh;
  background-image: conic-gradient(from 0deg, #0000 var(--angle), greenyellow 0deg);
  z-index: -1;
}

.searchbar-container.active::before {
  animation: 1s draw-border linear forwards;
}

@keyframes draw-border {
  from {--angle: 360deg;}
  to {--angle: 0deg;}
}

</style>