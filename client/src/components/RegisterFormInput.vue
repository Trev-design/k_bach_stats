<template>
  <section class="form-input-container">
    <div class="register-label">
      <p class="register-label-text">register</p>
    </div>
    <form class="input-form" @submit.prevent="handleSubmit">
      <label class="form-label" for="">Name:</label>
      <input class="input-area" type="text" required v-model="userName">

      <label class="form-label" for="">Email:</label>
      <input class="input-area" type="email" required v-model="email">

      <label class="form-label" for="">Password:</label>
      <input class="input-area" type="password" required v-model="password">

      <label class="form-label" for="">Confirmation;</label>
      <input class="input-area" type="password" required v-model="confirmation">

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

.submit {
  margin: 3rem 0;
  text-align: center;
}

.submit-button {
  padding: .6rem 2.5rem;
  background-color: rgb(3, 6, 18);
  color: rgb(110, 170, 250);
  border-radius: 100vh;
}
</style>