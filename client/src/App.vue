<script>
import BackToHomeButton from './components/BackToHomeButton.vue'
import SignoutButton from './components/SignoutButton.vue'
import ProtectedTopBar from './components/ProtectedTopBar.vue'
import ProtectedSidebar from './components/ProtectedSidebar.vue'

export default {
  name: 'App',
  components: {
    ProtectedTopBar,
    BackToHomeButton,
    SignoutButton,
    ProtectedSidebar
  },

  data() {
    return {
      refreshInterval: null,
      isUnfold: false
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
    <ProtectedTopBar 
      v-if="$store.getters.isAuthenticated" 
      v-model:isUnfold="isUnfold"
      @unfold="isUnfold = !isUnfold">
      <BackToHomeButton/>
      <SignoutButton/>
    </ProtectedTopBar>
    <ProtectedSidebar 
      v-if="$store.getters.isAuthenticated" 
      v-model:isUnfold="isUnfold"
    />
    <router-view></router-view>
  </div>
</template>

<style scoped>
</style>
