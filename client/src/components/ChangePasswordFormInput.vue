<template>
  <div>
    <section v-if="canChange" class="forgot-password-form-container">
      <div class="new-password-label">
        <p class="new-password-labeltext"></p>
      </div>
      <form action="" class="forgot-password-input-form">
        <label for="email" class="form-label" required>Email</label>
        <input 
          type="email" 
          class="input-area" 
          id="email" required 
          v-model="email">
      </form>
    </section>
    <section v-else class="change-password-form-container">
      <div class="new-password-label">
        <p class="new-password-labeltext"></p>
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
            <p class="error-message" v-else>{{ verifyErrorMessage }}</p>
          </div>
        </div>

        <div class="submit">
          <button class="submit-button">Change Password</button>
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
    }
  ),

  methods: {

  },

  watch: {

  }
}
</script>


<style scoped>
    
</style>