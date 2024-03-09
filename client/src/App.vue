<script>
export default {
  name: 'App',

  data(){
    return {
      refreshInterval: null
    }
  },

  methods: {
    refresh(status) {
      if (status) {
        const interval = setInterval(() => {
          this.$store.dispatch('refreshRequest')
            .catch((_error) => {
              localStorage.removeItem('id')
              localStorage.removeItem('guest')
              this.$router.push('/signin')
            })
        }, 1000 * 60 * 10)
        this.setRefreshInterval(interval)
      } else {
        this.setRefreshInterval(null)
        clearInterval(this.refreshInterval)
      }
    },

    setRefreshInterval(interval) {
      this.$data.refreshInterval = interval
    }
  },

  computed: {
    jwt() {return this.$store.state.jwt}
  },

  watch: {
    jwt: {
      handler() {this.refresh(true)},
      once: true
    }
  },

  mounted() {
    if (localStorage.getItem('id') != null) {
      this.$store.dispatch('refreshRequest')
        .catch((_error) => {
          localStorage.removeItem('id')
          localStorage.removeItem('guest')
          this.$router.push('/signin')
        })
    }
  },
}
</script>


<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>
