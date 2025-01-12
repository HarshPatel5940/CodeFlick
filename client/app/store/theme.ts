import { defineStore } from 'pinia'

export const useTheme = defineStore({
  id: 'theme',
  state: () => ({ isDarkMode: true }),
  actions: {
    getThemeBool() {
      return this.isDarkMode
    },
    getTheme() {
      return this.isDarkMode ? 'dark' : 'light'
    },
    toggleTheme() {
      this.isDarkMode = !this.isDarkMode
      return this.isDarkMode ? 'dark' : 'light'
    },
  },
  persist: {
    storage: piniaPluginPersistedstate.cookies(),
  },
})
