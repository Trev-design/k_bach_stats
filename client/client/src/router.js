import { createRouter, createWebHistory } from "vue-router"
import { useAuthStore, pinia } from "./store"

const routes = [
    {
        path: '/',
        name: 'Landing',
        component: () => import('./components/sites/Landing.vue')
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('./components/sites/Register.vue')
    },
    {
        path: '/signin',
        name: 'Signin',
        component: () => import('./components/sites/Signin.vue')
    },
    {
        path: '/verify',
        name: '/Verify',
        component: () => import('./components/sites/Verify.vue')
    },
    {
        path: '/user/:id',
        name: 'Home',
        component: () => import('./components/sites/Home.vue'),
        meta: { requiresAuth: true}
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, _from, next) => {
    const auth = useAuthStore()

    if (to.meta.requiresAuth && !auth.isAuthenticated) {
        next({ name: 'Signin'})
    } else {
        next()
    }
})

export default router