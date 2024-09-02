<template>
  <section class="verify-page-container">
    <BaseInputHeader headerText="Verification"/>
    <BaseInputForm>
      <BaseInput 
        inputID="vertifyinput"
        inputType="text"
        inputName="verification code"
        v-model:inputValue="verifyInput"
        @focussed="inputType='EMAIL'"
        @blurred="inputType='NONE'"
      />
      <BaseInputSubmit/>
    </BaseInputForm>
  </section>
</template>

<script>
import BaseInput from '../components/BaseInput.vue';
import BaseInputForm from '../components/BaseInputForm.vue';
import BaseInputHeader from '../components/BaseInputHeader.vue';
import BaseInputSubmit from '../components/BaseInputSubmit.vue';

export default {
  name: 'VerifyPage',
  
  components: {
    BaseInput,
    BaseInputForm,
    BaseInputHeader,
    BaseInputSubmit
  }, 

  data() {
    return {
      verifyInput: ''
    }
  }
}
</script>

<style scoped>
.verify-page-container {
  width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>