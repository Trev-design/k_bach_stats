<template>
  <section class="home-container">
    <Sidebar />
    <Topbar/>
    <section class="home-account-content-container">

    </section>
  </section>
</template>

<script>
import Sidebar from '../components/Sidebar.vue'
import Topbar from '../components/Topbar.vue'

export default {
  name: 'Home',
  components: {
    Sidebar,
    Topbar
  }
}
</script>

<style scoped>
.home-container {
  width: 100%;
  height: 100%;
}

.home-account-content-container {
  position: relative;
  left: 2rem;
  width: 85vw;
  height: 100vh;
  margin: 7rem;
  border: 1px solid white;
}
</style>