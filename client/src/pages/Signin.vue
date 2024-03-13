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


<style>
.signin-page-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>