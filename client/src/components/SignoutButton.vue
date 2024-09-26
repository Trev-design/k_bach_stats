<template>
  <a class="signout-button" @click="handleSignout()">signout</a>
</template>

<script>
export default {
  name: 'SignoutButton',
  methods: {
    handleSignout() {
      this.$store.dispatch('signout').catch(error => console.log(error))
      this.$router.push('/signin')
    }
  }
}
</script>

<style>
.signout-button {
  position: relative;
  border: none;
  font-size: 1.4rem;
  margin-right: 1.8rem;
  background: inherit;
  color: rgb(180, 200, 225);
  cursor: pointer;
}

.signout-button::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  top: -4px;
  left: 0;
  right: 0;
  color: rgb(180, 200, 225);
  transform: scale(0, 1);
  transition: transform 0.3s ease;
}

.signout-button:hover::after {
  transform: scale(1, 1);
}
</style>