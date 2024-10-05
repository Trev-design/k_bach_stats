<template>
  <button class="signout-button" @click="handleSignout()">signout</button>
</template>

<script>
export default {
  name: 'SignoutButton',
  methods: {
    handleSignout() {
      this.$store.dispatch('signout').catch(error => console.log(error))
      this.$router.push('/signin')
    }
  },

  created() {

  }
}
</script>

<style>
.signout-button {
  position: relative;
  text-decoration: none;
  border: none;
  font-size: 1.4rem;
  margin-right: 1.8rem;
  color: rgb(180, 200, 225);
  background: inherit;
  cursor: pointer;
}

.signout-button::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 2px;
  bottom: -4px;
  left: 0;
  right: 0;
  background: rgb(180, 200, 225);
  transform: scale(0, 1);
  transition: transform 0.3s ease;
}

.signout-button:hover::after {
  transform: scale(1, 1);
}
</style>