<template>
  <section class="form-input-container">
    <div class="verify-request-label">
      <p class="verify-request-label-text">register</p>
    </div>
    <form action="" class="input-form">
      <label for="email" class="input-label">Email:</label>
      <input type="text" class="verify-request-input-area" id="email" v-model="email">
    </form>
    <div class="error-message-container">
      <p class="error-message"></p>
    </div>
    <div class="error-message-wrapper">
      <div class="error-message-container">
        <p class="error-message">{{ errorMessage }}</p>
      </div>
    </div>
    <div class="submit">
      <button class="submit-button" @click="handleSubmit()">Submit</button>
    </div>
  </section>
</template>


<script>
export default {
  name: 'VerifyRequestForm',
  data() {
    return {
      email: '',
      errorMessage: ''
    }
  },

  methods: {
    handleSubmit() {
      this.$store.dispatch(
        'newVerifyRequest',
        {email: this.email}
      )
      .then((_ok) => {this.$router.push('/verify')})
      .catch((err) => {this.errorMessage = err})
    }
  },
}
</script>


<style scoped>
  .verify-request-label {
    position: absolute;
    top: -4rem;
    width: 600px;
  }

  .verify-request-label-text {
    padding-left: .5rem;
    font-size: 1.2rem;
    font-weight: 700;
  }

  .form-input-container {
    width: 350px;
    height: 300px;
    position: relative;
    background-color: rgb(3, 6, 18);
    border-radius: 6px;

    &::before {
      content: '';
      position: absolute;
      top: -4rem;
      left: -5px;
      width: 360px;
      height: 369px;
      z-index: -10;
      background-color: rgb(110, 170, 250);
      border-radius: 10px;
    }
  }

  .input-form {
    display: flex;
    flex-direction: column;
  }

  .input-label {
    margin: 1.5rem 0 .4rem 2rem;
    color: rgb(110, 170, 250);
  }

  .verify-request-input-area {
    height: 25px;
    margin: 0 2rem 2rem 2rem;
    padding: .1rem .5rem;
    color: rgb(110, 170, 250);
    background-color: rgb(3, 6, 18);
    border: 1px solid rgb(110, 170, 250);
    border-radius: 6px;
    &:focus {outline: none;}
  }

  .submit {
    margin: 3rem 0;
    text-align: center;
  }

  .submit-button {
    padding: .6rem 2.5rem;
    background-color: rgb(3, 6, 18);
    border-radius: 100vh;
    color: rgb(110, 170, 250);
    font-size: 1rem;
    cursor: pointer;
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
</style>