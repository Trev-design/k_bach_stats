<template>
  <section class="forgotten-password-page-container">
    <ChangePasswordFormInput/>
  </section>
</template>


<script>
import ChangePasswordFormInput from '../components/ChangePasswordFormInput.vue'
export default {
  name: 'ForgottenPassword',
  components: {
    ChangePasswordFormInput
  }
}
</script>


<style scoped>
.forgotten-password-page-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>