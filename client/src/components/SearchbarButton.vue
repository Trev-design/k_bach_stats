<template>
  <button class="searchbar-button">
    <Icon 
      :class="`searchbar-menu-icon ${searchEnabled ? 'active' : ''}`"
      icon="mingcute:search-3-line"
    />
  </button>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  name: 'SearchbarButton',

  components: {
    Icon
  },

  props: {
    searchEnabled: {
      type: Boolean,
      required: true
    },
  }
}
</script>

<style>
.searchbar-button {
  background: inherit;
  outline: none;
  border: none;
  margin-right: .5rem;
  cursor: pointer;
  z-index: 10;
}

.searchbar-menu-icon {
  width: 1.25rem;
  height: 1.25rem;
  color: aqua;
  transition: all .3s ease;
}

.searchbar-menu-icon.active {
  color: rgb(125, 255, 255);
}

.searchbar-menu-icon:hover {
  color: rgb(125, 255, 255);
}
</style>