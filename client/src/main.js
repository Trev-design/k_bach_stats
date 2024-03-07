import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createStore } from 'vuex'
import './style.css'
import App from './App.vue'
import Home from './pages/Home.vue'
import Register from './pages/Register.vue'
import Signin from './pages/Signin.vue'
import Verify from './pages/Verify.vue'

const app = createApp(App)

const routes = [
  {path: "/", component: Home},
  {path: "/register", component: Register},
  {path: "/signin", component: Signin},
  {path: "/verify", component: Verify}
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
        jwt: ''
      }
    },
    mutations: {

    },
    actions: {
      
    }
  }
)

app.use(store)
app.use(router)
app.mount('#app')
