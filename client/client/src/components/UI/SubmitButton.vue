<template>
  <button 
  @click="func"
  class="w-[120px] block bg-primary dark:bg-light hover:bg-secondary text-light dark:text-dark px-2 py-0.5 rounded-xl">
    {{ text }}
  </button>
</template>

<script setup>
defineProps({
  text: {type:String, required:true},
  func: {type:Function, default:() => {}}
})
</script>

<style>

</style>