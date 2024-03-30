<template>
  <div class="signin-label">
    <p class="signin-label-text">signin</p>
  </div>
  <section class="form-input-container">
    <form action="" class="input-form">
      <label for="" class="form-input-label">Email: </label>
      <input type="email" class="input-area" required v-model="email">

      <label for="" class="form-input-label">Password: </label>
      <input type="password" class="input-area" required v-model="password">

      <div class="submit">
        <button class="submit-button" @click="handleSubmit()">Submit</button>
      </div>
      <router-link to="/forgot-password">forgot password</router-link>
    </form>
  </section>
</template>


<script>
export default {
  name: 'SigninFormInput',
  data: () => (
    {
      email: '',
      password: ''
    }
  ),
  methods: {
    handleSubmit() {
      this.$store.dispatch(
        'signinRequest',
        {
          email: this.email,
          password: this.password
        }
      )
      .then((_ok) => {this.$router.push('/')})
      .catch((_error) => {this.$router.push('/signin')})
    }
  }
}
</script>


<style scoped>
.signin-label {
  width: 350px;
  height: 4rem;
}

.signin-label-text {
  padding-left: .5rem;
  font-size: 1.2rem;
  font-weight: 700;
}

.form-input-container {
  width: 350px;
  height: 400px;
  position: relative;
  background-color: rgb(3, 6, 18);
  border-radius: 6px;

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

.input-form {
  display: flex;
  flex-direction: column;
}

.form-input-label {
  margin: 1.5rem 0 .4rem 2rem;
  color: rgb(110, 170, 250);
}

.input-area {
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
</style>