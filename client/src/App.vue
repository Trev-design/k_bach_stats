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

  created() {
    if (localStorage.getItem('userId') != null) {
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
        <h1 class="brand">KBach</h1>
      <ul class="nav-items">
        <li><router-link class="nav-link" to="/">{{ guest }}</router-link></li>
        <li><router-link to="/signin" class="nav-link" @click="logout()">Logout</router-link></li>
      </ul>
    </div>
    <router-view/>
  </div>
</template>

<style scoped>
.navbar-container {
  position: fixed;
  top: 0;
  z-index: 2;
  width: 100%;
  height: 4rem;
  display: flex;
  justify-content: end;
  align-items: center;
  background-color: rgb(3, 6, 18);
}

.nav-items {
  display: flex;
  flex-direction: row;

  li {
    list-style: none;
  }
}

.nav-link {
  margin: 0 1.5rem 0 0;
  text-decoration: none;
  font-size: 1.1rem;
  color: rgb(160, 180, 200);
  position: relative;
  &::before {
    content: '';
    width: 120%;
    left: -10%;
    bottom: -3px;
    height: 2px;
    background-color: rgb(110, 170, 250);
    position: absolute;
    transform: scale(0, 1);
    transition: all 0.3s ease;
  }

  &:hover {&::before {transform: scale(1, 1);}}
}

.brand {
  background: rgb(17,158,186);
  background: linear-gradient(90deg, rgba(17,158,186,1) 16%, rgba(43,43,222,1) 78%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-right: auto;
  margin-left: 1rem;
}
</style>
