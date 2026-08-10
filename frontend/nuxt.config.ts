import { defineNuxtConfig } from 'nuxt/config'

type RouteManifest = {
  static: string[]
  hotels: string[]
  areas: string[]
  articles: string[]
  tags: string[]
}

const getRouteManifest = async (): Promise<RouteManifest> => {
  const manifest: RouteManifest = {
    static: ['/', '/blog'],
    hotels: [],
    areas: [],
    articles: [],
    tags: [],
  }
  const backendUrl = process.env.BACKEND_API_URL || 'http://localhost:8080'
  const buildToken = process.env.NEXT_PUBLIC_SSG_BUILD_TOKEN
  const headers: Record<string, string> = {}
  if (buildToken) {
    headers['x-github-build-token'] = buildToken
  }

  try {
    const [hotelsResponse, locationsResponse, postsResponse, tagsResponse, articleTagsResponse] = await Promise.all([
      fetch(`${backendUrl}/api/hotels?limit=10000`, { headers }),
      fetch(`${backendUrl}/api/regions/combined`, { headers }),
      fetch(`${backendUrl}/api/posts?limit=10000`, { headers }),
      fetch(`${backendUrl}/api/hotel-tags`, { headers }),
      fetch(`${backendUrl}/api/article-tags`, { headers })
    ])

    const requiredResponses = [hotelsResponse, locationsResponse, postsResponse, tagsResponse, articleTagsResponse]
    if (requiredResponses.some(response => !response.ok)) {
      throw new Error(`Required API failed during SSG route discovery: ${requiredResponses.map(r => r.status).join(', ')}`)
    }

    const hotels = await hotelsResponse.json()
    for (const hotel of hotels.data || []) {
      manifest.hotels.push(`/detail/${hotel.id}`)
    }
    const locations = await locationsResponse.json()
    const townships = (locations.cities || []).flatMap((city: any) => city.townships || [])
    if (townships.length !== 368) {
      throw new Error(`Expected 368 township categories during SSG, received ${townships.length}`)
    }

    for (const city of locations.cities || []) {
      const cityResponse = await fetch(
        `${backendUrl}/api/hotels?limit=1&area=${encodeURIComponent(city.name)}`,
        { headers }
      )
      if (!cityResponse.ok) throw new Error(`Failed to count hotels for city ${city.name}`)
      const cityResult = await cityResponse.json()
      const cityPages = Math.max(1, Math.ceil((cityResult.total || 0) / 20))
      for (let page = 1; page <= cityPages; page++) {
        manifest.areas.push(`/area/${city.id}/${page}`)
      }

      for (const township of city.townships || []) {
        const townshipPages = Math.max(1, Math.ceil((township.hotel_count || 0) / 20))
        for (let page = 1; page <= townshipPages; page++) {
          manifest.areas.push(`/area/${city.id}/${township.id}/${page}`)
        }
      }
    }

    const posts = await postsResponse.json()
    for (const post of posts.data || []) {
      manifest.articles.push(`/blog/${post.id}`)
    }

    const tags = await tagsResponse.json()
    for (const tag of tags || []) {
      const tagPages = Math.max(1, Math.ceil((tag.enabled_hotel_count || 0) / 20))
      for (let page = 1; page <= tagPages; page++) {
        manifest.tags.push(`/tag/${tag.id}/${page}`)
      }
    }

    const articleTags = await articleTagsResponse.json()
    for (const tag of articleTags || []) {
      const tagPages = Math.max(1, Math.ceil((tag.post_count || 0) / 12))
      for (let page = 1; page <= tagPages; page++) manifest.tags.push(`/blog/tag/${tag.id}/${page}`)
    }
  } catch (e) {
    if (process.env.BACKEND_API_URL) throw e
    console.warn('Dynamic API routes unavailable during local build:', e)
  }

  for (const key of Object.keys(manifest) as Array<keyof RouteManifest>) {
    manifest[key] = [...new Set(manifest[key])]
  }
  return manifest
}

let routeManifestPromise: Promise<RouteManifest> | undefined
const resolveRouteManifest = () => routeManifestPromise ||= getRouteManifest()
const resolvePrerenderRoutes = async () => {
  const manifest = await resolveRouteManifest()
  return [...new Set(Object.values(manifest).flat())]
}

export default defineNuxtConfig({
  devServer: {
    port: 3002
  },
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
    sitemaps: {
      static: { includeAppSources: false, urls: async () => (await resolveRouteManifest()).static, chunks: 10000 },
      hotels: { includeAppSources: false, urls: async () => (await resolveRouteManifest()).hotels, chunks: 10000 },
      areas: { includeAppSources: false, urls: async () => (await resolveRouteManifest()).areas, chunks: 10000 },
      articles: { includeAppSources: false, urls: async () => (await resolveRouteManifest()).articles, chunks: 10000 },
      tags: { includeAppSources: false, urls: async () => (await resolveRouteManifest()).tags, chunks: 10000 },
    },
  },
  nitro: {
    output: {
      publicDir: 'dist'
    },
    prerender: {
      failOnError: true,
      concurrency: 4,
      routes: ['/', '/sitemap_index.xml']
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

      const routes = await resolvePrerenderRoutes()

      // Add to prerender routes
      // Add to prerender routes
      if (nitroConfig.prerender) {
        nitroConfig.prerender.routes = nitroConfig.prerender.routes || []
        nitroConfig.prerender.routes.push(...routes)
      }
    }
  }
})
