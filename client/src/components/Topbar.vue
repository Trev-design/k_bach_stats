<template>
  <section class="topbar-container">
    <div class="brandlogo-container">
      <Brandlogo/>
    </div>
    <div class="navlink-container">
      <slot></slot>
    </div>
  </section>
</template>

<script>
import Brandlogo from './Brandlogo.vue'
export default {
  name: 'Topbar',
  components: {
    Brandlogo
  }
}
</script>

<style scoped>
.topbar-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 5rem;
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

.brandlogo-container {
  margin-right: auto;
  height: inherit;
  display: flex;
  align-items: center;
}
</style>