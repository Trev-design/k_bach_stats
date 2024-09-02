<template>
  <section class="new-verify-page-container">
    <BaseInputHeader headerText="new verification"/> 
    <BaseInputForm>
      <BaseInput 
        inputID="emailInput"
        inputType="text"
        inputName="email"
        v-model:inputValue="emailInput"
        @focussed="inputType='EMAIL'"
        @blurred="inputType='NONE'"
      />
      <BaseInputSubmit/>
    </BaseInputForm>
  </section>
</template>

<script>
import BaseInput from '../components/BaseInput.vue'
import BaseInputForm from '../components/BaseInputForm.vue'
import BaseInputSubmit from '../components/BaseInputSubmit.vue'
import BaseInputHeader from '../components/BaseInputHeader.vue'

export default {
  name: 'NewVerifyPage',

  components: {
    BaseInput,
    BaseInputForm,
    BaseInputSubmit,
    BaseInputHeader
  },

  data() {
    return {
      emailInput: ''
    }
  }
}
</script>

<style scoped>
.new-verify-page-container {
  width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>