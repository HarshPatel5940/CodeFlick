// https://nuxt.com/docs/api/configuration/nuxt-config
// eslint-disable-next-line no-undef
export default defineNuxtConfig({
  modules: [
    '@nuxt/eslint',
    '@nuxt/icon',
    '@nuxt/ui',
    '@pinia/nuxt',
    'pinia-plugin-persistedstate/nuxt',
    '@vueuse/nuxt',
    '@vueuse/motion/nuxt',
    'nuxt-time',
  ],

  ssr: false,

  imports: {
      dirs: [
        './store',
        './components',
        './components/ui',
        './components/home',
        './components/dashboard',
      ],
    },
  devtools: {
    enabled: false,

    timeline: {
      enabled: true,
    },
  },

  css:['~/assets/css/onedark.css'],

  srcDir: 'app',
  serverDir: 'nuxt-server',
  compatibilityDate: '2024-10-24',

  eslint: {
    config: {
      standalone: false,
      nuxt: {
        sortConfigKeys: true,
      },
    },
  },

  components: [
      {
        path: '~/components',
        pathPrefix: false,
      },
    ],

})