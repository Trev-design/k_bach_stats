<script>
export default {
  name: 'App',
  data() {
    return {
      refreshInterval: null
    }
  },

  created() {this.startRefreshInterval()},

  methods: {
    startRefreshInterval() {
      this.refreshInterval = setInterval(() => {
        this.$store.dispatch('refreshSession').catch(error => console.log(error))
      }, 600000)
    }
  },

  beforeUnmount() {clearInterval(this.refreshInterval)}
}
</script>

<template>
  <router-view></router-view>
</template>

<style scoped>
</style>
