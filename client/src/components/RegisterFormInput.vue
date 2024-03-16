<template>
  <section class="form-input-container">
    <div class="register-label">
      <p class="register-label-text">register</p>
    </div>
    <form class="input-form" @submit.prevent="handleSubmit">
      <label class="form-label" for="name">Name:</label>
      <input 
        class="input-area" 
        id="name" 
        type="text" 
        required 
        v-model="userName"
        @focus="nameFocus = true"
        @blur="nameFocus = false">

      <label class="form-label" for="email">Email:</label>
      <input 
        class="input-area" 
        id="email" 
        type="email" 
        required 
        v-model="email"
        @focus="emailFocus = true"
        @blur="emailFocus = false">

      <label class="form-label" for="password">Password:</label>
      <input 
        class="input-area" 
        id="password" 
        type="password" 
        required 
        v-model="password"
        @focus="passwordFocus = true"
        @blur="passwordFocus = false">

      <label class="form-label" for="confirmation">Confirmation:</label>
      <input 
        class="input-area" 
        id="confirmation" 
        type="password" 
        required 
        v-model="confirmation"
        @focus="confirmationFocus = true"
        @blur="confirmationFocus = false">

        <div class="error-message-wrapper">
          <div class="error-message-container">
            <p class="error-message" v-if="nameFocus">{{ nameErrorMessage }}</p>
            <p class="error-message" v-else-if="emailFocus">{{ emailErrorMessage }}</p>
            <p class="error-message" v-else-if="passwordFocus">{{ passwordErrorMessage }}</p>
            <p class="error-message" v-else>{{ confirmationErrorMessage }}</p>
          </div>
        </div>

      <div class="submit">
        <button class="submit-button">Create Account</button>
      </div>
    </form>
  </section>
</template>


<script>
export default {
  name: 'RegisterFormInput',
  data: () => (
    {
      userName: '',
      email: '',
      password: '',
      confirmation: '',
      nameFocus: false,
      emailFocus: false,
      passwordFocus: false,
      confirmationFocus: false,
      nameErrorMessage: '',
      emailErrorMessage: '',
      passwordErrorMessage: '',
      confirmationErrorMessage: ''
    }
  ),
  methods: {
    handleSubmit() {
      this.$store.dispatch(
        'registerRequest',
        {
          name: this.userName,
          email: this.email,
          password: this.password,
          confirmation: this.confirmation
        }
      )
      .then((_ok) => {this.$router.push('/verify')})
      .catch((_error) => {this.$router.push('/register')})
    }
  },
  watch: {
    userName() {
      if (this.userName.length < 4) {
        this.nameErrorMessage = 'name must have mor than one character'
      } else {
        this.nameErrorMessage = ''
      }
    },

    email() {
      if (!/^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[a-z]{2,4}$/.test(this.email)) {
        this.emailErrorMessage = 'please enter a regular email address'
      } else {
        this.emailErrorMessage = ''
      }
    },

    password() {
      if (!/^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!?&%$§#@€]).{10,}$/.test(this.password)) {
        this.passwordErrorMessage = 'password needs to have et least 10 characters including lowercase letters, uppercase letters, digits and special characters'
      } else {
        this.passwordErrorMessage = ''
      }
    },

    confirmation() {
      if (this.confirmation !== this.password) {
        this.confirmationErrorMessage = 'confirmation does not match'
      } else {
        this.confirmationErrorMessage = ''
      }
    }
  }
}
</script>


<style scoped>
.register-label {
  position: absolute;
  top: -4rem;
  width: 600px;
}

.register-label-text {
  padding-left: .5rem;
  font-size: 1.2rem;
  font-weight: 700;
}
.form-input-container {
  width: 400px;
  height: 600px;
  position: relative;
  margin-top: 3rem;
  background-color: rgb(3, 6, 18);
  border-radius: 6px;
  &::before {
    content: '';
    position: absolute;
    top: -4rem;
    left: -5px;
    right: .2rem;
    width: 410px;
    height: 669px;
    z-index: -10;
    background-color: rgb(110, 170, 250);
    border-radius: 10px;
  }
}

.form-label {
  margin: 2rem 0 .2rem 2rem;
  color: rgb(110, 170, 250);
}

.input-form {
  display: flex;
  flex-direction: column;
}

.input-area {
  margin: 0 2rem 2rem 2rem;
  padding: .1rem .5rem;
  height: 20px;
  background-color: rgb(3, 6, 18);
  border-radius: 6px;
  border: 1px solid rgb(110, 170, 250);
  color: rgb(110, 170, 250);
  &:focus {outline: none;}
}

.error-message-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
}

.error-message-container {
  position: absolute;
  width: 80%;
}

.error-message {
  color: rgb(165, 50, 80);
}

.submit {
  margin: 2rem 0;
  text-align: center;
}

.submit-button {
  padding: .6rem 2.5rem;
  background-color: rgb(3, 6, 18);
  color: rgb(110, 170, 250);
  border-radius: 100vh;
}
</style>