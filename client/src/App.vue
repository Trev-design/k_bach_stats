<script>
import Topbar from './components/Topbar.vue'
import Navlink from './components/Navlink.vue'
import SignoutButton from './components/SignoutButton.vue'

export default {
  name: 'App',
  components: {
    Topbar,
    Navlink,
    SignoutButton
  },
  data() {
    return {
      refreshInterval: null
    }
  },

  computed: {
    username() {
      return localStorage.getItem('username')
    },

    account() {
      return localStorage.getItem('account')
    }
  },

  created() {this.startRefreshInterval()},

  methods: {
    startRefreshInterval() {
      this.refreshInterval = setInterval(() => {
        this.$store.dispatch('refreshSession').catch(error => console.log(error))
      }, 60000)
    }
  },

  beforeUnmount() {clearInterval(this.refreshInterval)}
}
</script>

<template>
  <div>
    <Topbar v-if="$store.getters.isAuthenticated">
      <Navlink :dest="`account/${account}`">{{ username }}</Navlink>
      <SignoutButton/>
    </Topbar>
    <router-view></router-view>
  </div>
</template>

<style scoped>
</style>
