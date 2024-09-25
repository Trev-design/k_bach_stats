<template>
  <div class="input-container">
    <input 
      :type="inputType"
      :id="inputID"
      class="form-input"
      @input="$emit('update:inputValue', $event.target.value)"
      @focus="$emit('focussed')"
      @blur="$emit('blurred')"
      required
    >
    <label :for="inputID" class="form-label">{{inputName}}</label>
  </div>
</template>


<script scoped>
export default {
  name: 'BaseInput',

  props: {
    inputType: {
      type: String,
      default: 'text'
    },
    inputID: {
      type: String,
      required: true
    },
    inputName: {
      type: String,
      required: true
    },
    inputvalue: {
      type: String,
      default: ''
    }
  }
}
</script>


<style scoped>
.input-container {
  position: relative;
  width: 65%;
  height: 4rem;
}

.form-input {
  position: relative;
  width: 100%;
  outline: none;
  line-height: 2.5rem;
  font-size: 1rem;
  background: transparent;
  transition: 0.2s ease;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  border: none;
  border-radius: 10px;
  padding: 0 1rem;
  margin: 0 -1rem;
  border: 2px solid #f0eeee;
}

.form-label {
  position: absolute;
  top: .45rem;
  left: .5rem;
  padding: 0 .3rem;
  margin: 0 -0.5rem;
  color: #f0eeee;
  background: #191f4d;
  font-size: 1.2rem;
  transition: 0.4s ease;
}

.form-input:focus,
.form-input:valid {
  color: #5eeb49;
  border: 2px solid #5eeb49;
}


.form-input:focus + .form-label,
.form-input:valid + .form-label {
  color: #5eeb49;
  transform: translate(-0.6rem, -1.4rem) scale(0.8);
  z-index: 100;
}
</style>