<template>
  <section class="form-input-container">
    <form action="" class="input-form">
      <label for="" class="form-input-label">Email: </label>
      <input type="email" class="input-area" required v-model="email">

      <label for="" class="form-input-label">Password: </label>
      <input type="password" class="input-area" required v-model="password">

      <div class="submit">
        <button class="submit-button">Submit</button>
      </div>
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
      console.log(this.email)
    }
  }
}
</script>


<style scoped>
.form-input-container {
  position: absolute;
  top: 50%;
  left: 50%;
  border: 1px solid black;
  width: 350px;
  height: 400px;
  margin: -200px 0 0 -175px;
} 

.input-form {
  display: flex;
  flex-direction: column;
}

.form-input-label {
  margin: .4rem 2rem;
}

.input-area {
  height: 25px;
  margin: 0 2rem;
}

.submit {
  margin: 3rem 0;
  text-align: center;
}

.submit-button {
  padding: .4rem 2rem;
}
</style>