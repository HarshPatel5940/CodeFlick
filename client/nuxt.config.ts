// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    "@vueuse/nuxt",
    "@nuxt/icon",
    "@unocss/nuxt",
    "@pinia/nuxt",
    "@pinia-plugin-persistedstate/nuxt",
    "@nuxt/eslint",
  ],

  // ssr: false,

  imports: {
    dirs: ["./layouts", "./components", "./pages", "./store", "./utils"],
  },

  devtools: { enabled: true },
  css: ["./app/app.css"],

  srcDir: "app",
  serverDir: "server",
  compatibilityDate: "2024-10-24",

  eslint: {
    config: {
      standalone: false,
      nuxt: {
        sortConfigKeys: true,
      },
    },
  },

  unocss: {
    nuxtLayers: true,
  },
});
