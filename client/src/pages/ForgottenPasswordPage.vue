<template>
  <section class="forgotten_password_page_container">
    <BaseInputHeader headerText="New Password"/>
    <BaseInputForm>
      <BaseInput
        inputID="verifyinput"
        inputType="text"
        inputName="verification code"
        v-model:inputValue="verifyInput"
      />

      <BaseInput
        inputID="passwordinput"
        inputType="password"
        inputName="password"
        v-model:inputValue="passwordInput"
      />

      <BaseInput
        inputID="confirmationinput"
        inputType="password"
        inputName="confirmation"
        v-model:inputValue="confirmationInput" 
      />
    </BaseInputForm>
  </section>
</template>

<script>
import BaseInput from '../components/BaseInput.vue'
import BaseInputForm from '../components/BaseInputForm.vue'
import BaseInputHeader from '../components/BaseInputHeader.vue'
import BaseInputSubmit from '../components/BaseInputSubmit.vue'

export default {
  name: 'ForgottenPasswordPage',
  
  components: {
    BaseInput,
    BaseInputForm,
    BaseInputHeader,
    BaseInputSubmit
  },

  data() {
    return {
      verifyInput: '',
      passwordInput: '',
      confirmationInput: ''
    }
  },

  methods: {
    newPasswordRequest() {
      const payload = {
        verify: this.verifyInput,
        password: this.passwordInput,
        confirmation: this.confirmationInput
      }

      fetch('http://localhost:4000/forgotten_password', {
        method: 'POST',
        credentials: 'include',
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
        this.$router.push(`/home/${localStorage.getItem('account')}`)        
      })
    }
  }
}
</script>

<style scoped>
.forgotten_password_page_container {
  width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>