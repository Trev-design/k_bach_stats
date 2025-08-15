<template>
  <router-link
  :to="dest"
  :class="[`flex justify-center items-center rounded-full md:px-8 md:py-2 md:mr-4 sm:mt-2 md:mt-0 no-underline cursor-pointer sm:font-bold
  ${primary ? 'md:bg-primary' : 'md:bg-black'} ${primary ? 'dark:md:bg-light' : 'dark:md:bg-secondary'}
  text-secondary md:text-white md:dark:text-dark`]"
  >{{label}}</router-link>
</template>

<script setup>
defineProps({
  dest: {
    type:String,
    required:true
  },
  primary: {
    type:Boolean,
    default: true
  },
  label: {
    type:String,
    required:true
  }
})
</script>