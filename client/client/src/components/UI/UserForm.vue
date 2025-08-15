<template>
  <form :class="[`w-full ${large ? 'min-h-[500px]' : 'min-h-[300px]'} md:w-[450px] flex flex-col justify-between items-center`]" @submit.prevent="$emit('submit')">
    <div class="w-full">
      <h1 class="sm:ml-12 text-dark dark:text-light uppercase font-bold text-lg md:text-xl lg:text-2xl">{{header}}</h1>
    </div>
    <slot></slot>
  </form>
</template>

<script setup>
defineProps({
  header: {type:String,required:true,},
  large: {type:Boolean,default:false}
}) 
</script>