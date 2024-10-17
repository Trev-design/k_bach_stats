<template>
  <div class="back-to-home-button-container">
    <button class="back-to-home-button" @click="redirectToHome()"> {{ user }} </button>
  </div>
</template>

<script>
export default {
  name: 'BackToHomeButton',

  methods: {
    redirectToHome() {
      this.$router.push(`/${this.account}/${this.user}`)
    }
  },

  computed: {
    user() {
      return localStorage.getItem('username')
    },

    account() {
      return localStorage.getItem('account')
    }
  }
}
</script>

<style scoped>

.back-to-home-button-container {
  background: #191f4d;
  margin-right: 1.8rem;
  border-radius: 5px;
  transition: all 0.5s ease;
}

.back-to-home-button {
  display: inline-block;
  text-decoration: none;
  border: none;
  font-size: 1.4rem;
  padding: 0.6rem 1rem;
  border-radius: 4px;
  color: transparent;
  cursor: pointer;
  background: linear-gradient(to right, rgb(117, 179, 245) 30%, rgb(199, 43, 251));
  background-size: 200% 100%;
  background-position: left bottom;
  background-clip: text;
  -webkit-background-clip: text;
  transition: all 0.5s ease;
}

.back-to-home-button:hover {
  background-position: right bottom;
}

.back-to-home-button-container:hover {
  background: #212c82;
}
</style>