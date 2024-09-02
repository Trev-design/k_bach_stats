<template>
  <form class="base-input-form">
    <slot></slot>
  </form>
</template>

<script>
export default {
    name: 'BaseInputForm'
}
</script>

<style>
.base-input-form {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
}
</style>