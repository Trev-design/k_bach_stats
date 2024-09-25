<template>
  <div class="input-header-container">
    <div class="input-header-text-container">
      <h1 class="base-input-header"> {{ headerText }} </h1>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BaseInputHeader',

  props: {
    headerText: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.input-header-container {
  width: 80%;
  height: 6rem;
  display: flex;
  flex: column;
  justify-content: left;
  align-items: center;
}

.input-header-text-container {
  margin: 0 1rem;
}

.base-input-header {
  background: linear-gradient(90deg, #00dbde, #fc00ff);
  font-size: 2rem;
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
</style>