<template>
  <section class="homepage-con">
    <Topbar/>
    <Sidebar/>
    <section class="homepage-content-container">

    </section>
  </section>
</template>

<script>
import Sidebar from '../components/Sidebar.vue';
import Topbar from '../components/Topbar.vue';
export default {
  name: 'HomePage',
  components: {
    Sidebar,
    Topbar
  }
}
</script>

<style scoped>

.homepage-content-container {
  position: relative;
  left: 2rem;
  width: 85vw;
  height: 100vh;
  margin: 7rem;
  border: 1px solid white;
}
</style>