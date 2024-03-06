<template>
  <section class="register-page-container">
    <FormInput></FormInput>  
  </section>
</template>


<script>
import FormInput from '../components/FormInput.vue';
export default {
  name: 'Register',
  components: {FormInput}
}
</script>


<style>
    .register-page-container {
      position: relative;
    }
</style>