<script>
import BaseInput from './components/BaseInput.vue';
export default {
  name: 'App',
  components: {
    BaseInput
  },
  data() {
    return {
      nameInput: ''
    }
  }
}
</script>

<template>
  <router-view></router-view>
</template>

<style scoped>
</style>
