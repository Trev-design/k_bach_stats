<template lang="">
  <section class="signin-page-container">  
    <SigninFormInput/>   
  </section>
</template>


<script>
import SigninFormInput from '../components/SigninFormInput.vue';
export default {
  name: 'Signin',
  components: {SigninFormInput}
}
</script>


<style lang="">
.signin-page-container {
  position: relative;
}
</style>