<template>
  <section class="verify-page-container"> 
    <VerifyFormInput/>
  </section>
</template>


<script>
import VerifyFormInput from '../components/VerifyFormInput.vue';
export default {
  name: 'verify',
  components: {VerifyFormInput}
}
</script>


<style >
.verify-page-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>