export const useTheme = defineStore({
  id: "theme",
  state: () => {
    console.log("skill issue", persistedState.localStorage.getItem("theme"));
    if (persistedState.localStorage.getItem("theme")) {
      return {
        isDarkMode: persistedState.localStorage.getItem("theme") === "dark",
      };
    }
    return {
      isDarkMode: false,
    };
  },
  actions: {
    getThemeBool() {
      return this.isDarkMode;
    },
    getTheme() {
      return this.isDarkMode ? "dark" : "light";
    },
    toggleTheme() {
      this.isDarkMode = !this.isDarkMode;
      return this.isDarkMode ? "dark" : "light";
    },
  },
  persist: true,
});
