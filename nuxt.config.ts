// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: true,
  nitro: {
    output: {
      publicDir: 'dist'
    },
    prerender: {
      failOnError: false,
      concurrency: 4
    }
  },
  app: {
    baseURL: process.env.NUXT_APP_BASE_URL || '/'
  }
})
