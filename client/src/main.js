import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createStore } from 'vuex'
import './style.css'
import App from './App.vue'
import Home from './pages/Home.vue'
import Register from './pages/Register.vue'
import Signin from './pages/Signin.vue'
import Verify from './pages/Verify.vue'
import VerifyRequest from './pages/VerifyRequest.vue'
import ForgottenPassword from './pages/ForgottenPassword.vue'

const app = createApp(App)

const routes = [
  {path: "/", component: Home},
  {path: "/register", component: Register},
  {path: "/signin", component: Signin},
  {path: "/verify", component: Verify},
  {path: "/new-verify", component: VerifyRequest},
  {path: "/forgotten-password", component, ForgottenPassword}
]

const router = createRouter(
  {
    history: createWebHashHistory(),
    routes: routes
  }
)

const store = createStore({
  state() {
    return {
      jwt: '',
    }
  },

  mutations: {
    setAccessToken(state, token) {
      state.jwt = token
    },

    unsetAccessToken (state) {
      state.jwt = ''
    }
  },

  actions: {
    registerRequest({}, userdata) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include',
          body: JSON.stringify(userdata)
        }
        fetch('http://localhost:4000/account/create', requestOptions)
          .then((response) => {
            if (response.ok) {
              return response.json()
            } else {
              reject(response.json().then((data) => (data.message)))
            }
          })
          .then((data) => {
            localStorage.setItem('guest', data.guest)
            resolve('OK')
          })
          .catch((error) => {
            reject(error)
          })
      })
    },
    
    signinRequest({commit}, userdata) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include',
          body: JSON.stringify(userdata)
        }
        fetch('http://localhost:4000/account/signin', requestOptions)
          .then((response) => {
            if (response.ok) {
              return response.json()
            } else {
              reject(response.json().then((data) => (data.message)))
            }
          })
          .then((data) => {
            localStorage.setItem('guest', data.name)
            localStorage.setItem('userId', data.id)
            commit('setAccessToken', data.jwt)
            resolve('OK')
          })
          .catch((error) => {
            reject(error)
          })
      })
    },

    verifyRequest({commit}, userdata) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include',
          body: JSON.stringify(userdata)
        }
        fetch('http://localhost:4000/account/verify', requestOptions)
          .then((response) => {
            if (response.ok) {
              return response.json()
            } else {
              reject(response.json().then((data) => (data.message)))
            }
          })
          .then((data) => {
            localStorage.setItem('userId', data.id)
            commit('setAccessToken', data.jwt)
            resolve('OK')
          })
          .catch((error) => {
            reject(error)
          })
      })
    },

    refreshRequest({commit}) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'GET',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include'
        }
        fetch('http://localhost:4000/session/refresh_session', requestOptions)
          .then((response) => {
            if (response.ok) {
              return response.json()
            } else {
              reject(response.json().then((data) => (data.message)))
            }
          })
          .then((data) => {
            commit('setAccessToken', data.jwt)
            resolve('OK')
          })
          .catch((error) => {
            reject(error)
          })
      })
    },

    signoutRequest({commit}) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'GET',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include'
        }
        fetch('http://localhost:4000/session/signout', requestOptions)
          .then((response) => {
            if (response.ok) {
              return response.json()
            } else {
              reject(response.json().then((data) => (data.message)))
            }
          })
          .then((_data) => {
            localStorage.removeItem('guest')
            localStorage.removeItem('userId')
            commit('unsetAccessToken')
            resolve('OK')
          })
          .catch((error) => {
            reject(error)
          })
      })
    },

    newVerifyRequest({}, userdata) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include',
          body: JSON.stringify(userdata)
        }
        fetch('http://localhost:4000/account/new_verify', requestOptions)
          .then((response) => {
            if (response.ok) {
              return response.json()
            } else {
              reject(response.json().then((data) => (data.message)))
            }
          })
          .then((_data) => {resolve('OK')})
          .catch((error) => {
            reject(error)
          })
      })
    },

    requestNewPassword({}, userdata) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include',
          body: JSON.stringify(userdata)
        }

        fetch('http://localhost:4000/account/forgotten_password', requestOptions)
        .then((response) => {
          if (response.ok) {
            return response.json()
          }

          reject(response.json().then((data) => (data.message)))
        })
        .then((_data) => {resolve('OK')})
        .catch((error) => {reject(error.message)})
      })
    },

    requestChangePassword({}, userdata) {
      return new Promise((resolve, reject) => {
        const requestOptions = {
          method: 'PUT',
          headers: {'Content-Type': 'application/json'},
          credentials: 'include',
          body: JSON.stringify(userdata)
        }

        fetch('http://localhost:4000/account/change_password', requestOptions)
        .then((response) => {
          if (response.ok) {
            return response.json()
          }

          reject(response.json().then((data) => (data.message)))
        })
        .then((_data) => {resolve('OK')})
        .catch((error) => {reject(error.message)})
      })
    }
  }
})

app.use(store)
app.use(router)
app.mount('#app')
