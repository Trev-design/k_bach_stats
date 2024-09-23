import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createRouter, createWebHashHistory} from 'vue-router'
import RegisterPage from './pages/RegisterPage.vue'
import SigninPage from './pages/SigninPage.vue'
import VerifyPage from './pages/VerifyPage.vue'
import NewVerifyPage from './pages/NewVerifyPage.vue'
import HomePage from './pages/HomePage.vue'
import Home from './pages/Home.vue'
import {createStore} from 'vuex'


const app = createApp(App)

const routes = [
  {path: '/', component: HomePage},
  {path: '/register', component: RegisterPage},
  {path: '/signin', component: SigninPage},
  {path: '/verify', component: VerifyPage},
  {path: '/new-verify', component: NewVerifyPage, props: {action: 'new_verify'}},
  {path: '/account/:id', component: Home, meta: {requiredAuth: true}}
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
  linkActiveClass: 'active'
})

const store = createStore({
  state: {
    jwt: null
  },

  mutations: {
    setJWT(state, jwt) {
      state.jwt = jwt
    }
  },

  actions: {
    setJWT({commit}, token) {
      commit('setJWT', token)
    },

    removeJWT({commit}) {
      commit('setJWT', null)
    },

    refreshSession({commit}) {
      return new Promise((resolve, reject) => {
        fetch('http://localhost:4000/session/refresh', {
          method: 'GET',
          headers: {'Content-Type': 'application/json'}        
        })
        .then(response => {
          if (!response.ok) {
            data = response.json()
            reject(data.message)
          }

          return response.json()
        })
        .then(data => {
          if (!data.jwt) {
            reject('something went wrong')
          }

          commit('setJWT', data.jwt)
          resolve()
        })
        .catch(error => reject(error));
      })
    }
  },

  getters: {
    isAuthenticated(state) {return state.jwt != null}
  }
})

router.beforeEach((to, from, next) => {
  console.log(to.meta)
  if (to.meta.requiredAuth && !store.getters.isAuthenticated) {
    console.log("try to check authentication")
    next('/signin')
  } else {
    console.log("no reason to check auth")
    next()
  }
})

app.use(router)
app.use(store)

app.mount('#app')

