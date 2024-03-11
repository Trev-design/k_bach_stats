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
              this.$router.push('signin')
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
    },

    logout() {
      this.$store.dispatch('signoutRequest')
      this.$router.push('')
    }
  },

  computed: {
    jwt() {return this.$store.state.jwt},
    guest() {return localStorage.getItem('guest')}
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
    <div class="navbar-container" v-if="jwt != ''">
      <ul class="nav-items">
        <li><router-link class="nav-link" to="/">{{ guest }}</router-link></li>
        <li><router-link to="/signin" class="nav-link" @click="logout()">Logout</router-link></li>
      </ul>
    </div>
    <router-view></router-view>
  </div>
</template>

<style>
.navbar-container {
  position: fixed;
  top: 0;
  z-index: 2;
  width: 100%;
  display: flex;
  justify-content: end;
  align-items: center;
}

.nav-items {
  display: flex;
  flex-direction: row;

  li {
    list-style: none;
  }
}

.nav-link {
  margin: 0 1.2rem 0 0;
  text-decoration: none;
}
</style>
