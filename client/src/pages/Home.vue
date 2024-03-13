<template>
  <div>
    <LandingNavbar v-if="jwt === ''"/>
    <About/>
    <Features/>
    <Testimonials/>
    <Pricing/>
  </div>
</template>


<script>
import LandingNavbar from '../components/LandingNavbar.vue';
import About from '../sections/About.vue';
import Features from '../sections/Features.vue';
import Testimonials from '../sections/Testimonials.vue';
import Pricing from '../sections/Pricing.vue'

export default {
  name: 'Home',
  components: {
    LandingNavbar,
    About,
    Features,
    Testimonials,
    Pricing
  },
  computed: {
    jwt() { return this.$store.state.jwt }
  }
}
</script>


<style>

</style>