<template>
  <section class="home-container">
    <section class="home-account-content-container">

    </section>
  </section>
</template>

<script>
export default {
  name: 'Home',
  async created() {
    await this.$store.dispatch('initialData').catch(_err => console.log('could not fetch userid'))
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