<template>
  <section class="form-input-container">
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
          .then((_ok) => {this.$router.push('/')})
          .catch((_error) => {this.$router.push('/register')})
      }
    }
}
</script>


<style scoped>
.form-input-container {
  width: 400px;
  height: 600px;
  border: 1px solid black;
  position: absolute;
  top: 50%;
  left: 50%;
  margin: 70px 0 0 -200px;
}

.form-label {
  margin: 2rem 0 .2rem 2rem;
}

.input-form {
  display: flex;
  flex-direction: column;
}

.input-area {
  margin: 0 2rem;
  height: 20px;
}

.submit {
  margin: 3rem 0;
  text-align: center;
}

.submit-button {
  padding: .4rem 1rem;
}
</style>