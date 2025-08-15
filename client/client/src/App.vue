<script setup>
import { defineAsyncComponent } from 'vue'
const Navbar = defineAsyncComponent(() => import('./components/layout/Navbar.vue'))
</script>

<template>
<div class="bg-light dark:bg-dark min-h-screen">
  <router-view></router-view>
</div>
</template>

<style scoped>

</style>
