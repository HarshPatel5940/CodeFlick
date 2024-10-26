export const useCount = defineStore({
  id: 'count',
  state: () => ({ count: 0 }),
  actions: {
    increment() {
      this.count++
    },
    decrement() {
      this.count--
    },
  },
  persist: true,
})
