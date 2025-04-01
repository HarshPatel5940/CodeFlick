import { defineStore } from 'pinia'

export const useTheme = defineStore('theme', {
  state: () => ({ isDarkMode: true }),
  getters: {
    themeBool(state) {
      return state.isDarkMode
    },
    theme(state) {
      return state.isDarkMode ? 'dark' : 'light'
    },
  },
  actions: {
    toggleTheme() {
      this.isDarkMode = !this.isDarkMode
      return this.isDarkMode ? 'dark' : 'light'
    },
  },
  persist: {
    storage: piniaPluginPersistedstate.cookies(),
  },
})
