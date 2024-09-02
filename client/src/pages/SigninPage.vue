<template>
  <section class="signin-page-container">
    <BaseInputHeader headerText="Signin"/>
    <baseInputForm>
      <BaseInput 
        inputID="emailInput"
        inputType="text"
        inputName="email"
        v-model:inputValue="emailInput"
        @focussed="inputType='EMAIL'"
        @blurred="inputType='NONE'"
      />
      <BaseInput 
        inputID="passwordInput"
        inputType="password"
        inputName="password"
        v-model:inputValue="passwordInput"
        @focussed="inputType='PASSWORD'"
        @blurred="inputType='NONE'"
      />
      <BaseProviderSelect/>
      <BaseInputSubmit/>
    </baseInputForm>
  </section>
</template>

<script>
import BaseInput from '../components/BaseInput.vue'
import BaseInputForm from '../components/BaseInputForm.vue'
import BaseInputHeader from '../components/BaseInputHeader.vue'
import BaseInputSubmit from '../components/BaseInputSubmit.vue'
import BaseProviderSelect from '../components/BaseProviderSelect.vue'

export default {
  name: 'SigninPage',

  components: {
    BaseInput,
    BaseInputForm,
    BaseInputHeader,
    BaseInputSubmit,
    BaseProviderSelect
  },

  data() {
    return {
      emailInput: '',
      passwordInput: '',
      inputType: 'NONE'
    }
  }
}
</script>

<style scoped>
.signin-page-container {
  width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>