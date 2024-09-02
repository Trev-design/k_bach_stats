import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createRouter, createWebHashHistory} from 'vue-router'
import RegisterPage from './pages/RegisterPage.vue'
import SigninPage from './pages/SigninPage.vue'
import VerifyPage from './pages/VerifyPage.vue'
import NewVerifyPage from './pages/NewVerifyPage.vue'
import HomePage from './pages/HomePage.vue'
import Vuex from 'vuex'


const app = createApp(App)

const routes = [
  {path: '/', component: HomePage},
  {path: '/register', component: RegisterPage},
  {path: '/signin', component: SigninPage},
  {path: '/verify', component: VerifyPage},
  {path: '/new-verify', component: NewVerifyPage}
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: routes,
  linkActiveClass: 'active'
})

app.use(router)

app.mount('#app')
