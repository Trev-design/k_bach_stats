import { createApp } from 'vue'
import {createRouter, createWebHashHistory} from 'vue-router'
import './style.css'
import App from './App.vue'
import Home from './pages/Home.vue'
import Register from './pages/Register.vue'
import Signin from './pages/Signin.vue'

const app = createApp(App)

const routes = [
  {path: "/", component: Home},
  {path: "/register", component: Register},
  {path: "/signin", component: Signin}
]

const router = createRouter(
  {
    history: createWebHashHistory(),
    routes: routes
  }
)

app.use(router)
app.mount('#app')
