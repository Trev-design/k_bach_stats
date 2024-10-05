import { createApp, provide } from 'vue'
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
import { ApolloClient, createHttpLink, InMemoryCache } from '@apollo/client/core'
import { createApolloProvider } from '@vue/apollo-option'
import { setContext } from '@apollo/client/link/context'
import VueApolloComponents from '@vue/apollo-components'
import { provideApolloClient } from '@vue/apollo-composable'
import { h } from 'vue'


const routes = [
  {path: '/', component: HomePage},
  {path: '/signup', component: RegisterPage},
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
          credentials: 'include',
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
    },

    signout({commit}) {
      return new Promise((resolve, reject) => {
        fetch("http://localhost:4000/session/signout", {
          credentials: 'include',
          method: 'GET',
          headers: {'Content-Type': 'application/json'}
        })
        .then(response => {
          if (!response.ok) {
            reject('something went wrong')
          }

          return response.json()
        })
        .then(_data => {
          commit('setJWT', null)
          localStorage.removeItem('username')
          localStorage.removeItem('account')
          resolve()
        })
        .catch(_error => reject('something went wrong'))
      })
    }
  },

  getters: {
    isAuthenticated(state) {return state.jwt != null},
    token(state) {return state.jwt}
  }
})

router.beforeEach((to, from, next) => {
  console.log(to.meta)
  if (to.meta.requiredAuth && !store.getters.isAuthenticated) {
    console.log("try to check authentication")
    store.dispatch('refreshSession').catch(_error => next('/signin'))
  } else {
    console.log("no reason to check auth")
    next()
  }
})

const authLink = setContext((_, { headers }) => {
  const authToken = store.getters.token

  return {
    headers: {
      ...headers,
      Authorization: authToken ? `Baerer ${authToken}` : ''
    }
  }
})

const httpLink = createHttpLink({
  uri: "http://localhost:5148/graphql",
  credentials: 'include'
})

const apolloClient = new ApolloClient({
  cache: new InMemoryCache(),
  link: authLink.concat(httpLink)
})

//const apolloProvider = createApolloProvider({defaultClient: apolloClient})

const app = createApp({
  setup() {
    provideApolloClient(apolloClient)
  },

  render: () => h(App)
})

app.use(router)
app.use(store)
app.use(VueApolloComponents)

app.mount('#app')

