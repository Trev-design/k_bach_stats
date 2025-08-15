import { defineStore } from "pinia"
import { createPinia } from "pinia"

export const pinia = createPinia()

export const useAuthStore = defineStore('auth', {
    state: () => {
        return {jwt: ''}
    },

    actions: {
        login(newToken) { this.jwt = newToken },

        logout() { this.jwt = ''}
    },

    getters: {
        isAuthenticated: (state) => state.jwt != '',
        token: (state) => state.jwt
    }
})

export const useUserStore = defineStore('user', {
    state: () => {
        return {name: '', abo: '', roles: []}
    },
    getters: {
        user: (state) => state.name,
        type: (state) => state.abo,
        roles: (state) => state.roles
    }
})