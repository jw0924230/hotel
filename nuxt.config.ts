// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: true,
  nitro: {
    output: {
      publicDir: 'dist'
    }
  },
  app: {
    baseURL: process.env.NUXT_APP_BASE_URL || '/'
  }
})
