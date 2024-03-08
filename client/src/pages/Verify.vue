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
    position: relative;
}
</style>