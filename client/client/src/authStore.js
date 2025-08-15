import { defineStore } from "pinia"

export const useAuthStore = defineStore('auth', {
    state: () => {
        return {jwt: ''}
    },

    actions: {
        login(newToken) { this.jwt = newToken },

        logout() { this.jwt = ''}
    },

    getters: {
        isAuthenticated() { return this.jwt != ''},
        token() { return this.jwt }
    }
})