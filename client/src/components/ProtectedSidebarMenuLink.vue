<template>
  <router-link class="router-link">
    <span class="menu-icon-span">
      <Icon class="menu-option-icon" :icon="optionIcon"/>
    </span>
    <span :class="`menu-option-name ${isUnfold ? 'active' : 'inactive'}`"> <p>{{ optionName }}</p> </span>
  </router-link>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  name: 'ProtectedSidebarMenuLink',

  components: {
    Icon
  },

  props: {
    isUnfold: {
      type: Boolean,
      required: true
    },

    optionName: {
      type: String,
      required: true
    },

    optionIcon: {
      type: String,
      required: true
    }
  }
}
</script>

<style>
.router-link {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin: .5rem 1.25rem;
  color: rgb(100, 175, 125);
  text-decoration: none;
}

.menu-option-icon {
  vertical-align: -0.4em;
  font-size: 1.5rem;
}

.menu-option-name {
  font-size: 1.2rem;
  opacity: 0;
  transition: all 0.3s ease;
}

.menu-option-name.active {
  opacity: 1;
  margin: 0 1.4rem;
}
</style>