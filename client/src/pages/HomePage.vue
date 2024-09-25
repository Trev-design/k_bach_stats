<template>
  <section class="homepage-container">
    <Topbar>
      <Navlink dest="signin"/>
      <Navlink dest="signup"/>
    </Topbar>
    <section class="homepage-content-container">

    </section>
  </section>
</template>

<script>
import Topbar from '../components/Topbar.vue'
import Navlink from '../components/Navlink.vue'

export default {
  name: 'HomePage',
  components: {
    Topbar,
    Navlink
  }
}
</script>

<style scoped>
.homepage-container {
  width: 100%;
  height: 100%;
}

.homepage-content-container {
  position: relative;
  left: 2rem;
  width: 85vw;
  height: 100vh;
  margin: 7rem;
  border: 1px solid white;
}
</style>