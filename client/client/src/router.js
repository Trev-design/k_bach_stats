import { createRouter, createWebHistory } from "vue-router"
import { useAuthStore } from "./authStore"

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
    if (to.meta.requiresAuth) {
        const auth = useAuthStore()
        const isAuthenticated = auth.isAuthenticated()

        if (!isAuthenticated) {
            next({ name: 'Signin'})
        }
    }

    next()
})

export default router