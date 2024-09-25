<template>
  <div class="error-message-container">
    <p class="error-message"> {{ errorMessage }} </p>
  </div>
</template>

<script>
export default {
  name: 'BaseInputErrorMessage', 

  props: {
    errorMessage: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.error-message-container {
  margin: 1.5rem 0;
  width: 70%;
  height: 3.5rem;
  display: flex;
  justify-content: left;
  align-items: left;
  box-sizing: border-box;
  word-wrap: break-word;
}

.error-message {
  font-size: 0.85rem;
  color: #a82c38;
}
</style>