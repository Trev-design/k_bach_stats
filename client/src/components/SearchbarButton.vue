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
    searcEnabled: {
      type: Boolean,
      required: true
    },

    optionsUnfold: {
      type: Boolean,
      required: true
    }
  }
}
</script>

<style>

</style>