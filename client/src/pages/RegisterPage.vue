<template>
  <section class="register-page-container">
    <BaseInputHeader headerText="Register"/>
    <BaseInputForm>
      <BaseInput 
        inputID="nameInput"
        inputType="text"
        inputName="name"
        v-model:inputValue="nameInput"
        @focussed="inputType='NAME'"
        @blurred="inputType='NONE'"
      />
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
      <BaseInput 
        inputID="confirmationInput"
        inputType="password"
        inputName="confirm"
        v-model:inputValue="confirmationInput"
        @focussed="inputType='CONFIRMATION'"
        @blurred="inputType='NONE'"
      />
      <BaseProviderSelect/>
      <BaseInputErrorMessage :errorMessage="errorMessage"/>
      <BaseInputSubmit @click="registerRequest()"/>
    </BaseInputForm>
  </section>
</template>

<script>
import BaseInput from '../components/BaseInput.vue'
import BaseInputSubmit from '../components/BaseInputSubmit.vue'
import BaseInputForm from '../components/BaseInputForm.vue'
import BaseInputHeader from '../components/BaseInputHeader.vue'
import BaseInputErrorMessage from '../components/BaseInputErrorMessage.vue'
import BaseProviderSelect from '../components/BaseProviderSelect.vue'

export default {
  name: 'RegisterPage',

  components: {
    BaseInput,
    BaseInputSubmit,
    BaseInputForm,
    BaseInputHeader,
    BaseInputErrorMessage,
    BaseProviderSelect
  },

  data(){
    return {
      nameInput: '',
      emailInput: '',
      passwordInput: '',
      confirmationInput: '',
      inputType: ''
    }
  },

  computed: {
    errorMessage() {
      const inputType = this.inputType

      return inputType === 'NAME' && this.nameInput !== '' && this.nameInput.length < 4 ? 'Invalid Name: A name must have at least 4 characters':
             inputType === 'EMAIL' && this.emailInput !== '' && !/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(this.emailInput) ? 'Invalid Email: This email does not follow regular email conventions':
             inputType === 'PASSWORD' && this.passwordInput !== '' && !/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/.test(this.passwordInput) ? 'Invalid or weak Password: Password must have at least one uppercase letter, one uppercase letter, one digit, one special icon like @$!%*?& and at least 8 characters at all':
             inputType === 'CONFIRMATION' && !/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/.test(this.passwordInput) && this.confirmationInput !== this.passwordInput ? 'password does not match': ''
    }
  }, 

  methods: {
    registerRequest() {
      const request = {
        name: this.nameInput,
        email: this.emailInput,
        password: this.passwordInput,
        confirm: this.confirmationInput
      }

      fetch('http://localhost:4000/account/signup', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(request)
      })
      .then(response => {
        if (!response.ok) {
          throw new Error('invalid credentials')
        }
        return response.json()
      })
      .then(data => {
        if(!data.name) {
          throw new Error("invalid credentials")
        }

        if(!data.id) {
          throw new Error("invalid credentials")
        }

        localStorage.setItem('account', data.id)
        localStorage.setItem('username', data.name)
        
        this.$router.push('/verify')
      })
      .catch(err => console.log(err))
    }
  }
}
</script>

<style scoped>
.register-page-container {
  width: 500px;
  height: 700px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
</style>