<template>
  <section :class="`sidebar-container ${isUnfold ? 'active' : 'inactive'}`">
    <div class="search-option-container">

    </div>
    <div class="menu-items-container">
      <ProtectedSidebarMenuLink 
        :isUnfold="isUnfold" 
        optionName="profile" 
        optionIcon="mdi:account-circle-outline"
      />
      <ProtectedSidebarMenuLink 
        :isUnfold="isUnfold" 
        optionName="contacts" 
        optionIcon="mdi:account-group-outline"
      />
      <ProtectedSidebarMenuLink 
        :isUnfold="isUnfold" 
        optionName="contacts" 
        optionIcon="eos-icons:project-outlined"
      />
      <ProtectedSidebarMenuLink
        :isUnfold="isUnfold" 
        optionName="contacts" 
        optionIcon="eos-icons:project-outlined"
      />
    </div>
    <div class="menu-setings-container">
      <ProtectedSidebarMenuLink
        :isUnfold="isUnfold" 
        optionName="contacts" 
        optionIcon="eos-icons:project-outlined"
      />
    </div>
  </section>
</template>

<script>
import ProtectedSidebarMenuLink from './ProtectedSidebarMenuLink.vue';

export default {
  name: 'ProtectedSidebar',

  components: {
    ProtectedSidebarMenuLink
  },

  props: {
    isUnfold: {
      type: Boolean,
      required: true
    }
  }
}
</script>

<style scoped>
.sidebar-container {
  position: fixed;
  top: 5rem;
  bottom: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  width: 4rem;
  transition: all 0.3s ease;
  overflow: hidden;
}

.sidebar-container.active {width: 175px;}
.sidebar-container.inactive {width: 4rem;}

.menu-items-container {
  position: relative;
  display: flex;
  flex-direction: column;
  margin: 10rem 0 0;
}

.menu-setings-container {
  margin-top: auto;
  margin-bottom: 2rem;
}
</style>