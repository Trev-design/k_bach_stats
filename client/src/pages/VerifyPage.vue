<template>
  <section class="verify-page-container">
    <BaseInputHeader headerText="Verification"/>
    <BaseInputForm>
      <BaseInput 
        inputID="vertifyinput"
        inputType="text"
        inputName="verification code"
        v-model:inputValue="verifyInput"
      />
      <BaseInputSubmit @click="verifyAccountRequest()"/>
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
  },

  methods: {
    verifyAccountRequest() {
      const payload = {
        verify: this.verifyInput
      }

      fetch('http://localhost:4000/verify/account', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'userid': localStorage.getItem('account')
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
        if (!data.token) {
          throw new Error('invalid token')
        }
        this.$store.dispatch('setJWT', data.token)
        localStorage.setItem('username', data.user)
        this.$router.push(`account/${localStorage.getItem('account')}`)
      })
      .catch(error => console.log(error))
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