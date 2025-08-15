<template>
  <section class="min-h-screen w-full flex justify-center items-center">
    <UserForm header="register" :large="true">
      <FormInput 
      v-for="(item, index) in inputs" 
      :key="index"
      :id="item.id"
      :label="item.label"
      :type="item.type"
      :placeholder="item.placeholder"
      :rows="item.rows"
      v-model="item.modelValue"
      />
      <div class="w-[70%]">
        <SubmitButton text="submit" :func="submitFunc"/>
      </div>
    </UserForm>
  </section>
</template>

<script setup>
import FormInput from '../UI/FormInput.vue'
import UserForm from '../UI/UserForm.vue'
import SubmitButton from '../UI/SubmitButton.vue'
import { ref } from 'vue'

const inputs = ref([
  {id:'username', label:'username', type:'text', placeholder:'username', rows:undefined, modelValue:''},
  {id:'email', label:'email', type:'email', placeholder:'email', rows:undefined, modelValue:''},
  {id:'password', label:'password', type:'password', placeholder:'password', rows:undefined, modelValue:''},
  {id:'confirmation', label:'confirmation', type:'password', placeholder:'confirmation', rows:undefined, modelValue:''}
])

const submitFunc = () => {
  inputs.value.map((obj) => console.log(obj.modelValue))
}
</script>