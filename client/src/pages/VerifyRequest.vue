<template>
  <section class="new-verify-page-container">
    <VerifyRequestForm/>  
  </section>
</template>


<script>
import VerifyRequestForm from '../components/VerifyRequestForm.vue';
export default {
  name: 'VerifyRequest',
  components: {
    VerifyRequestForm
  }
}
</script>


<style scoped>
.new-verify-page-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>