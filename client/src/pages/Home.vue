<template>
  <div>
    <LandingNavbar v-if="jwt === ''"/>
  </div>
</template>


<script>
import LandingNavbar from '../components/LandingNavbar.vue';

export default {
  name: 'Home',
  components: {LandingNavbar},
  computed: {
    jwt() { return this.$store.state.jwt }
  }
}
</script>


<style>

</style>