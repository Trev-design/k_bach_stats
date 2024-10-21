<template>
  <secttion class="protected-topbar-container">
    <div class="interaction-container">
      <ProtectedSidebarUnfoldButton :isUnfold="unfold" @click="$emit('unfold')"/>
      <div class="brandlogo-container">
        <Brandlogo/>
      </div>
    </div>
    <div class="navlink-container">
      <slot></slot>
    </div>
  </secttion>
</template>

<script>
import ProtectedSidebarUnfoldButton from './ProtectedSidebarUnfoldButton.vue'
import Brandlogo from './Brandlogo.vue'

export default {
  name: 'ProtectedTopBar',

  components: {
    ProtectedSidebarUnfoldButton,
    Brandlogo
  },

  props: {
    isUnfold: {
        type: Boolean,
        Required: true
    }
  },

  computed: {
    unfold() {
      return this.isUnfold
    }
  },
}
</script>

<style scoped>
.protected-topbar-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 5rem;
  width: 100vw;
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: end;
  background: #191f4d;
  z-index: 10;
}

.navlink-container {
  height: inherit;
  display: flex;
  align-items: center;
  justify-content: right;
}

.interaction-container {
  height: inherit;
  display: flex;
  align-items: center;
  margin-right: auto;
}

.brandlogo-container {
  height: inherit;
  display: flex;
  align-items: center;
}
</style>