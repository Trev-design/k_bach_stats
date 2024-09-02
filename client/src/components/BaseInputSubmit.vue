<template>
  <div class="form-submit-button-container">
    <button 
      class="form-submit-button" 
      @click="$emit('onSubmit')"
    >
      Submit
    </button>
  </div>
</template>

<script>
export default {
  name: 'BaseInputSubmit'
}
</script>

<style scoped>
  .form-submit-button-container {
    margin: 2rem 0 0 0;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .form-submit-button {
    width: 50%;
    height: 4rem;
    border: none;
    outline: none;
    background: #611df2;
    border-radius: 5px;
    cursor: pointer;
  }
</style>