import { defineNuxtConfig } from 'nuxt/config'

const getDynamicRoutes = async () => {
  const routes = ['/blog']
  const backendUrl = process.env.BACKEND_API_URL || 'http://localhost:8080'
  const buildToken = process.env.NEXT_PUBLIC_SSG_BUILD_TOKEN
  const headers: Record<string, string> = {}
  if (buildToken) {
    headers['x-github-build-token'] = buildToken
  }

  try {
    const [hotelsResponse, citiesResponse, postsResponse] = await Promise.all([
      fetch(`${backendUrl}/api/hotels?limit=10000`, { headers }),
      fetch(`${backendUrl}/api/categories?type=city`, { headers }),
      fetch(`${backendUrl}/api/posts?limit=10000`, { headers })
    ])

    if (hotelsResponse.ok) {
      const hotels = await hotelsResponse.json()
      for (const hotel of hotels.data || []) {
        routes.push(`/detail/${hotel.id}`)
      }
    }

    if (citiesResponse.ok) {
      const cities = await citiesResponse.json()
      for (const city of cities) {
        const cityId = city.sort_order || city.id
        const response = await fetch(
          `${backendUrl}/api/hotels?limit=1&area=${encodeURIComponent(city.name)}`,
          { headers }
        )
        if (!response.ok) continue
        const result = await response.json()
        const totalPages = Math.max(1, Math.ceil((result.total || 0) / 20))
        for (let page = 1; page <= totalPages; page++) {
          routes.push(`/area/${cityId}/${page}`)
        }
      }
    }

    if (postsResponse.ok) {
      const posts = await postsResponse.json()
      for (const post of posts.data || []) {
        routes.push(`/blog/${post.id}`)
      }
    }
  } catch (e) {
    console.warn('Dynamic API routes unavailable during build:', e)
  }

  return routes
}

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: true,
  runtimeConfig: {
    public: {
      backendApiUrl: process.env.BACKEND_API_URL || 'http://localhost:8080',
      ssgBuildToken: process.env.NEXT_PUBLIC_SSG_BUILD_TOKEN || ''
    }
  },
  css: ['~/assets/css/reset.css'],
  modules: ['@nuxtjs/sitemap'],
  site: {
    url: 'https://www.qk3houronline.com/',
    name: '休息3小時'
  },
  sitemap: {
    urls: ['/blog']
  },
  nitro: {
    output: {
      publicDir: 'dist'
    },
    prerender: {
      failOnError: false,
      concurrency: 4,
      routes: ['/'] // explicit root
    }
  },
  app: {
    // If deploying to https://<username>.github.io/<repo-name>/
    // You MUST set baseURL to '/<repo-name>/'
    // The user's repo is 'hotel', so '/hotel/'
    baseURL: '/',
    head: {
      script: [
        {
          innerHTML: `(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','GTM-NTR8MG73');`
        }
      ]
    }
  },
  hooks: {
    async 'nitro:config'(nitroConfig) {
      if (process.env.NODE_ENV === 'development') {
        return
      }

      const routes = await getDynamicRoutes()

      // Add to prerender routes
      // Add to prerender routes
      if (nitroConfig.prerender) {
        nitroConfig.prerender.routes = nitroConfig.prerender.routes || []
        nitroConfig.prerender.routes.push(...routes)
      }
    }
  }
})
