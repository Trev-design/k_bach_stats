<template>
  <section class="form-input-container">
    <div class="verify-label">
      <p class="verify-label-text">verify</p>
    </div>
    <form action="" class="input-form">
      <div class="input">
        <label for="" class="form-input-label">enter Code:</label>
        <input type="text" class="input-area" required v-model="verifyCode">
      </div>

      <button class="submit-button" @click="handleSubmit()">Submit</button>
    </form>
  </section>
</template>


<script>
export default {
  name: 'VerifyFormInput',
  data: () => (
    {
      verifyCode: '',
      verifyError: ''
    }
  ),
  methods: {
    handleSubmit() {
      this.$store.dispatch('verifyRequest', {verification: this.verifyCode})
        .then((_ok) => {this.$router.push('/')})
        .catch((error) => {this.verifyError = error; console.log(error)})
    }
  }
}
</script>


<style scoped>

.verify-label {
  position: absolute;
  top: -3.8rem;
  left: .2rem;
}

.verify-label-text {
  font-size: 1.2rem;
  font-weight: 700;
}

.form-input-container {
  width: 350px;
  height: 400px;
  position: relative;

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
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  align-items: center;
  background-color: rgb(3, 6, 18);
  border-radius: 6px;
}

.input {
  display: flex;
  flex-direction: column;
}

.form-input-label {
  margin: .4rem 0;
  color: rgb(110, 170, 250);
}

.input-area {
  width: 260px;
  height: 25px;
  background-color: rgb(3, 6, 18);
  color: rgb(110, 170, 250);
  border: 1px solid rgb(110, 170, 250);
  border-radius: 6px;
  &:focus {outline: none;}
  padding: .1rem .5rem;
}

.submit-button {
  padding: .6rem 2rem;
  border-radius: 100vh;
  color: rgb(110, 170, 250);
  background-color: rgb(3, 6, 18);
  border: 1px solid rgb(110, 170, 250);
}
</style>