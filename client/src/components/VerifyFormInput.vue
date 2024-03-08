<template>
  <section class="form-input-container">
    <form action="" class="input-form">
      <div class="input">
        <label for="" class="form-input-label">enter Code:</label>
        <input type="text" class="input-area">
      </div>

      <button class="submit-button">Submit</button>
    </form>
  </section>
</template>


<script>
export default {
  name: 'VerifyFormInput'
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
  margin: 200px 0 0 -175px;
}

.input-form {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  align-items: center;
}

.input {
  display: flex;
  flex-direction: column;
}

.form-input-label {
  margin: .4rem 0;
}

.input-area {
  width: 260px;
  height: 25px;
}

.submit-button {
  padding: .4rem 1.5rem;
}

</style>