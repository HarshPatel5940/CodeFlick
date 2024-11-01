// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    "@nuxt/eslint",
    "@nuxt/icon",
    "@nuxt/ui",
    "@pinia/nuxt",
    "@pinia-plugin-persistedstate/nuxt",
    "@vueuse/nuxt",
    "@vueuse/motion/nuxt",
  ],

  // ssr: false,

  imports: {
    dirs: ["./layouts", "./components", "./pages", "./store", "./utils"],
  },

  devtools: { enabled: true },

  srcDir: "app",
  serverDir: "nuxt-server",
  compatibilityDate: "2024-10-24",

  eslint: {
    config: {
      standalone: false,
      nuxt: {
        sortConfigKeys: true,
      },
    },
  },
});
