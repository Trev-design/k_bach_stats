<template>
  <div>
    <section v-if="!canChange" class="forgot-password-form-container">
      <div class="new-password-label">
        <p class="new-password-labeltext">Forgot Password</p>
      </div>
      <form action="" class="new-password-input-form">
        <label for="email" class="form-label" required>Email</label>
        <input 
          type="email" 
          class="input-area"
          id="email" required 
          v-model="email">
        
          <div class="error-message-wrapper">
            <div class="error-message-container">
              <p class="error-message">{{ serverErrorMessage }}</p>
            </div>
          </div>

          <div class="submit">
            <button class="submit-button" @click="handleNewPasswordSubmit">New password</button>
          </div>
      </form>
    </section>

    <section v-else class="change-password-form-container">
      <div class="new-password-label">
        <p class="new-password-labeltext">Create Password</p>
      </div>
      <form action="" class="new-password-input-form">
        <label for="password" class="form-label">Password</label>
        <input 
          type="password" 
          class="input-area" 
          id="password" 
          required
          v-model="password"
          @focus="passwordFocus = true"
          @blur="passwordFocus = false">

        <label for="confirmation" class="form-label">Confirm</label>
        <input 
          type="password" 
          class="input-area" 
          id="confirmation" 
          required
          v-model="confirmation"
          @focus="confirmationFocus = true"
          @blur="confirmationFocus = false">

        <label for="verification" class="form-label">Verify Code</label>
        <input 
          type="text"
          class="input-area"
          id="verification"
          required
          v-model="verify"
          @focus="verifyFocus = true"
          @blur="verifyFocus = false">

        <div class="error-message-wrapper">
          <div class="error-message-container">
            <p class="error-message" v-if="passwordFocus">{{ passwordErrorMessage }}</p>
            <p class="error-message" v-else-if="confirmationFocus">{{ confirmationErrorMessage }}</p>
            <p class="error-message" v-else-if="verifyErrorMessage">{{ verifyErrorMessage }}</p>
            <p class="error-message" v-else>{{ serverErrorMessage }}</p>
          </div>
        </div>

        <div class="submit">
          <button class="submit-button" @click="handleCreatePasswordSubmit">Change Password</button>
        </div>
      </form>
    </section>
  </div>
</template>


<script>
export default {
  name: 'ChangePasswordFormInput',
  data : () => (
    {
      canChange: false,
      email: '',
      password: '',
      confirmation: '',
      verify: '',
      passwordFocus: false,
      confirmationFocus: false,
      verifyFocus: false,
      passwordErrorMessage: '',
      confirmationErrorMessage: '',
      verifyErrorMessage: '',
      serverErrorMessage: ''
    }
  ),

  methods: {
    handleNewPasswordSubmit() {
      this.$store.dispatch(
        'requestNewPassword',
        {email: this.email}
      )
      .then((_ok) => {
        this.canChange = true
        this.serverErrorMessage = ''
      })
      .catch((error) => {this.serverErrorMessage = error})
    },

    handleCreatePasswordSubmit() {
      this.$store.dispatch(
        'requestChangePassword',
        {
          password: this.password,
          confirmation: this.confirmation,
          verify: this.verify
        }
      )
      .then((_ok) => {this.$router.push('/signin')})
      .catch((error) => {this.serverErrorMessage = error})
    }
  },

  watch: {
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
    },

    verify() {
      if (!/^\d+$/.test(this.verify) && this.verify.length > 0) {
        this.verifyErrorMessage = 'input must be a string'
      } else {
        this.verifyErrorMessage = ''
      }
    }
  }
}
</script>


<style scoped>
.new-password-label {
  position: absolute;
  top: -4rem;
  width: 600px;
}

.new-password-labeltext {
  padding-left: .5rem;
  font-size: 1.2rem;
  font-weight: 700;
}

.forgot-password-form-container {
  width: 350px;
  height: 400px;
  position: relative;
  background-color: rgb(3, 6, 18);

  &::before {
    content: '';
    position: absolute;
    top: -4rem;
    left: -5px;
    right: .2rem;
    width: 360px;
    height: 469px;
    z-index: -10;
    background-color: rgb(110, 170, 250);
    border-radius: 10px;
  }
}

.change-password-form-container {
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

.new-password-input-form {
  display: flex;
  flex-direction: column;
}

.form-label {
  margin: 2rem 0 .2rem 2rem;
  color: rgb(110, 170, 250);
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