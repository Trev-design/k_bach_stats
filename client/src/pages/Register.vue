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
      position: relative;
    }
</style>