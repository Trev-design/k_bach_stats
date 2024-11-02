<template>
  <section :class="`add-experience-section-container ${addingReady ? 'active' : ''}`">
    <div :class="`searchbar-container ${searchEnabled ? 'active' : ''}`">
      <SearchbarButton
        :searchEnabled="searchEnabled"
        @click="$emit('handleSearch')"
      />

      <div class="search-input-container">
        <input 
        id="searchInput"
        type="text" 
        :disabled="!searchEnabled"
        :class="`searchbar-input ${searchEnabled ? 'active' : ''}`"
        @input="$emit('update:inputValue', $event.target.value)"
        @focus="searchbarFocused = true"
        @blur="searchbarFocused = false"
      >

      <label 
        for="searchInput" 
        :class="`searchbar-input-label ${searchEnabled ? 'enabled' : ''}`"
      >search</label>
      </div>
    </div>

    <ExperienceAdder/>

    <div class="submit-button-container">
      <button 
        class="add-experience-submit-button"
        @click="$emit('disable')"
      >submit</button>
    </div>
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
  },
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
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: left;
}

.add-experience-section-container.active {
  transform: scale(1, 1);
}

.searchbar-container {
  position: relative;
  display: flex;
  flex-direction: row;
  width: 85%;
  padding: .4rem 1rem;
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
  animation: .3s draw-border linear forwards;
}

@keyframes draw-border {
  from {--angle: 360deg;}
  to {--angle: 0deg;}
}

.searchbar-input {
  position: relative;
  background:#191f4d;
  outline: none;
  border: none;
  color: aqua;
  width: 100%;
}

.search-input-container {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 87.5%;
}

.searchbar-input-label {
  position: absolute;
  background: #191f4d;
  color: aqua;
  top: -3px;
  transform: translateX(-20px);
  opacity: 0;
  transition: all .3s ease;
  padding: .15rem;
}

.searchbar-input-label.enabled {
  transform: translateX(0);
  opacity: 1;
}

.searchbar-input:focus + .searchbar-input-label,
.searchbar-input:valid + .searchbar-input-label{
  top: -1.5rem;
}

.submit-button-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.add-experience-submit-button {
  width: 120px;
  height: 40px;
  color: rgb(15, 25, 40);
  background: rgb(50, 160, 180);
  outline: none;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: all .3s ease;
}

.add-experience-submit-button:hover {
  background: rgb(80, 200, 220);
}
</style>