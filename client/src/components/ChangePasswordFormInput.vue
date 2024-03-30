<template>
  <div>
    <section v-if="canChange" class="forgot-password-form-container">
      <div class="new-password-label">
        <p class="new-password-labeltext">Forgot Password</p>
      </div>
      <form action="" class="forgot-password-input-form">
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
      <form action="" class="change-password-input-form">
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
        .then((_ok) => {this.$router.push('/signin')})
        .catch((error) => {this.serverErrorMessage = error})
      )
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
    
</style>