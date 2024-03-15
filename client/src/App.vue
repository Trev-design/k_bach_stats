<script>
export default {
  name: 'App',

  data() {
    return {
      refreshDelay: null
    }
  },

  methods: {
    refresh(status) {
      if (status) {
        const delay = setTimeout(() => {
        this.$store.dispatch('refreshRequest')
          .catch((_error) => {
            localStorage.removeItem('id')
            localStorage.removeItem('guest')
            this.$router.push('signin')
          })
        }, 1000 *  10)
        this.setRefreshDelay(delay)
      } else {
        clearTimeout(this.refreshDelay)
        this.setRefreshDelay(null)
      }
    },

    setRefreshDelay(delay) {
      this.refreshDelay = delay
    },

    logout() {
      this.refresh(false)
      this.$store.dispatch('signoutRequest')
      this.$router.push('')
    },
    userName() {
      return localStorage.getItem('guest')
    }
  },

  computed: {
    jwt() {return this.$store.state.jwt}
  },

  watch: {
    jwt(token) {
      if (token != '') {
        this.refresh(true)
      }
    }
  },

  created() {
    if (!!localStorage.getItem('userId')) {
      console.log(this.userID != null)
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
        <li><router-link class="nav-link" to="/">{{ this.userName() }}</router-link></li>
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
    bottom: -4px;
    height: 2px;
    background-color: rgb(110, 170, 250);
    position: absolute;
    transform: scale(0, 1);
    transition: all 0.3s ease;
  }

  &:hover {&::before {transform: scale(1, 1);}}

  @media screen and (min-width: 850px) {
    margin: 0 1.5rem 0 0;
    font-size: 1.25rem;
  }
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
