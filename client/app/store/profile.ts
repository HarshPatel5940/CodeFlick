import { defineStore } from 'pinia'

export const useProfile = defineStore('profile', {
  state: () => ({
    UserID: '',
    name: '',
    email: '',
    isAdmin: false,
    isDeleted: false,
    isPremium: false,
  }),
  getters: {
    data: state => state,
  },
  actions: {
    set(
      UserID: string,
      name: string,
      email: string,
      isAdmin: boolean,
      isDeleted: boolean,
      isPremium: boolean,
    ) {
      this.UserID = UserID
      this.name = name
      this.email = email
      this.isAdmin = isAdmin
      this.isDeleted = isDeleted
      this.isPremium = isPremium
    },
    reset() {
      this.UserID = ''
      this.name = ''
      this.email = ''
      this.isAdmin = false
      this.isDeleted = false
      this.isPremium = false
    },
  },
  persist: {
    storage: piniaPluginPersistedstate.cookies(),
  },
})
