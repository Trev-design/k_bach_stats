<template>
  <div class="relative z-10 w-[70%]">
    <label :for="id" class="block ml-2 mb-2 text-sm font-medium text-gray-900 dark:text-gray-300">
      {{ label }}        
    </label>
    <component 
    class="p-2 shadow-sm text-gray-900  dark:text-gray-300 text-sm rounded-xl block w-full focus:ring-primary focus:border-primary-500 border border-secondary"
    :is="type === 'textarea' ? 'textarea' : 'input'"
    :rows="type === 'textarea' ? rows : undefined"
    :placeholder="placeholder"
    :id="id"
    :type="type"
    :value="modelValue"
    @input="$emit('update:modelValue', $event.target.value)"/>
  </div>
</template>

<script setup>
defineProps({
  id: {type:String, required:true},
  label:{type:String, required:true},
  type:{type:String, default:'text'},
  placeholder:{type:String, default:''},
  rows:{type:Number, default:4},
  modelValue:{type:[String, Number],default:''}
})
</script>

<style>
</style>