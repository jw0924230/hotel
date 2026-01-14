import { defineNuxtConfig } from 'nuxt/config'
import fs from 'fs'
import path from 'path'

// Helper to load JSON
const loadJSON = (filePath: fs.PathOrFileDescriptor) => {
  try {
    return JSON.parse(fs.readFileSync(filePath, 'utf-8'))
  } catch (e) {
    return []
  }
}

// Helper to generate dynamic routes
const getDynamicRoutes = () => {
  const routes: string[] = []

  // 1. Hotel Routes
  const detailsDir = path.resolve(__dirname, 'data_json/hotels/details')
  if (fs.existsSync(detailsDir)) {
    const files = fs.readdirSync(detailsDir)
    files.filter(file => file.endsWith('.json'))
      .forEach(file => routes.push(`/detail/${file.replace('.json', '')}`))
  }

  // 2. Area Routes
  const cityIds = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]
  for (const id of cityIds) {
    const areaPath = path.resolve(__dirname, `app/data/areas/area_${id}.json`)
    const areaData = loadJSON(areaPath)
    if (areaData.length > 0) {
      const totalPages = Math.ceil(areaData.length / 20)
      for (let p = 1; p <= totalPages; p++) {
        routes.push(`/area/${id}/${p}`)
      }
    }
  }

  // 3. Blog Routes
  const blogPath = path.resolve(__dirname, 'data_json/blog/articles.json')
  routes.push('/blog')
  if (fs.existsSync(blogPath)) {
    const articles = JSON.parse(fs.readFileSync(blogPath, 'utf-8'))
    articles.forEach((a: { id: any }) => routes.push(`/blog/${a.id}`))
  }

  return routes
}

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  ssr: true,
  modules: ['@nuxtjs/sitemap'],
  site: {
    url: 'https://jw0924230.github.io/',
    name: '休息3小時'
  },
  sitemap: {
    urls: getDynamicRoutes
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
    baseURL: '/'
  },
  hooks: {
    async 'nitro:config'(nitroConfig) {
      if (process.env.NODE_ENV === 'development') {
        return
      }

      const routes = getDynamicRoutes()

      // Add to prerender routes
      // Add to prerender routes
      if (nitroConfig.prerender) {
        nitroConfig.prerender.routes = nitroConfig.prerender.routes || []
        nitroConfig.prerender.routes.push(...routes)
      }
    }
  },
  vite: {
    server: {
      watch: {
        ignored: ['**/data_json/**']
      }
    }
  }
})
