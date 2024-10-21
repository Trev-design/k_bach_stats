<template>
  <button :class="`button ${isUnfold ? 'active' : 'inactive'}`"/>
</template>

<script>
export default {
  name: 'ProtectedSidebarUnfoldButton',
  props: {
    isUnfold: {
      type: String,
      required: true
    }
  }
}
</script>

<style>
.button {
  position: relative;
  border: none;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  margin-left: 1rem;
  justify-content: center;
  align-items: center;
  background: black;
  cursor: pointer;
}

.button::before,
.button::after {
  content: '';
  position: absolute;
  width: 1rem;
  height: 2px;
  background: white;
  transition: all 0.3s ease;
}

.button::before {top: 1.5rem;}
.button::after { top: 1rem; }

.button.active::before {
  top: 1.2rem;
  transform: rotate(135deg);
}

.button.active::after {
  top: 1.2rem;
  transform: rotate(-135deg);
}

.button.inactive::before {
  transform: rotate(0);
}

.button.inavtive::after {
  transform: rotate(0);
}
</style>