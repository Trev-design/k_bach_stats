<template>
  <div>
    <Navbar/>
  </div>
</template>

<script setup>
import { defineAsyncComponent } from 'vue'

const Navbar = defineAsyncComponent(() => import('../layout/Navbar.vue'))
</script>

<style>

</style>