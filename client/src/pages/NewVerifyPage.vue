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
      <BaseInputSubmit @click="newVerifyRequest()"/>
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
  },

  props: ['action'],

  methods: {
    newVerifyRequest() {
      const payload = {
        email: this.emailInput,
        kind: this.action
      }

      fetch('http://localhost:4000/verify/account', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload)
      })
      .then(response => {
        if (!response.ok) {
          throw new Error('invalid request')
        }

        return response.json()
      })
      .then(data => {
        if (!data.id) {
          throw new Error('something went wrong')
        }
        localStorage.setItem('account', data.id)
        
        if (this.action === 'new_verify') {
          this.$router.push('/verify')
        } else {
          this.$router.push('/forgotten_password')
        }
      })
      .catch(error => console.log(error))
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