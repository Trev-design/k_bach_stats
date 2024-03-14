<template>
  <section class="register-page-container">
    <RegisterFormInput/>  
  </section>
</template>


<script>
import RegisterFormInput from '../components/RegisterFormInput.vue';
export default {
  name: 'Register',
  components: {RegisterFormInput}
}
</script>


<style>
  .register-page-container {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
</style>