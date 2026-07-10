export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig()
  const token = config.public.ssgBuildToken
  
  if (token) {
    globalThis.$fetch = $fetch.create({
      onRequest({ options }) {
        options.headers = options.headers || {}
        if (Array.isArray(options.headers)) {
          options.headers.push(['x-github-build-token', token])
        } else if (options.headers instanceof Headers) {
          options.headers.set('x-github-build-token', token)
        } else {
          options.headers['x-github-build-token'] = token
        }
      }
    })
  }
})
