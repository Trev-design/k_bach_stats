<template>
  <section class="signin-page-container">
    <BaseInputHeader headerText="Signin"/>
    <baseInputForm>
      <BaseInput 
        inputID="emailInput"
        inputType="text"
        inputName="email"
        v-model:inputValue="emailInput"
      />
      <BaseInput 
        inputID="passwordInput"
        inputType="password"
        inputName="password"
        v-model:inputValue="passwordInput"
      />
      <BaseProviderSelect/>
      <BaseInputSubmit @click="signinRequest()"/>
    </baseInputForm>
  </section>
</template>

<script>
import { ssrModuleExportsKey } from 'vite/runtime'
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
      passwordInput: ''
    }
  },

  methods: {
    signinRequest() {
      const payload = {
        email: this.emailInput,
        password: this.passwordInput
      }

      fetch('http://localhost:4000/account/signin', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
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
        if (!data.jwt) {
          throw new Error('invalid token')
        }

        console.log(data.jwt)

        this.$store.dispatch('setJWT', data.jwt)
        localStorage.setItem('username', data.user)
        localStorage.setItem('account', data.account)
  
        this.$router.push(`account/${localStorage.getItem('account')}`)
      })
      .catch(error => console.log(error))
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