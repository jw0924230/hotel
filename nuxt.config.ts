import { defineNuxtConfig } from 'nuxt/config'
import fs from 'fs'
import path from 'path'

// Helper to load JSON
const loadJSON = (filePath) => {
  try {
    return JSON.parse(fs.readFileSync(filePath, 'utf-8'))
  } catch (e) {
    return []
  }
}

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
      concurrency: 4,
      routes: ['/'] // explicit root
    }
  },
  app: {
    // If deploying to https://<username>.github.io/<repo-name>/
    // You MUST set baseURL to '/<repo-name>/'
    // The user's repo is 'hotel', so '/hotel/'
    baseURL: process.env.NODE_ENV === 'production' ? '/hotel/' : '/'
  },
  hooks: {
    async 'nitro:config'(nitroConfig) {
      if (process.env.NODE_ENV === 'development') {
        return
      }
      // 1. Generate Hotel Routes
      const detailsDir = path.resolve(__dirname, 'data_json/hotels/details')
      let hotelRoutes: string[] = []

      if (fs.existsSync(detailsDir)) {
        const files = fs.readdirSync(detailsDir)
        hotelRoutes = files
          .filter(file => file.endsWith('.json'))
          .map(file => `/detail/${file.replace('.json', '')}`)
      }

      // 2. Generate Area Routes
      const cityIds = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21]
      const areaRoutes = []

      for (const id of cityIds) {
        const areaPath = path.resolve(__dirname, `app/data/areas/area_${id}.json`)
        const areaData = loadJSON(areaPath)
        if (areaData.length > 0) {
          const totalPages = Math.ceil(areaData.length / 20)
          for (let p = 1; p <= totalPages; p++) {
            areaRoutes.push(`/area/${id}/${p}`)
          }
        }
      }

      // 3. Generate Blog Routes
      const blogPath = path.resolve(__dirname, 'data_json/blog/articles.json')
      const blogRoutes = ['/blog']
      if (fs.existsSync(blogPath)) {
        const articles = JSON.parse(fs.readFileSync(blogPath, 'utf-8'))
        articles.forEach(a => blogRoutes.push(`/blog/${a.id}`))
      }

      // Add to prerender routes
      nitroConfig.prerender.routes = nitroConfig.prerender.routes || []
      nitroConfig.prerender.routes.push(...hotelRoutes, ...areaRoutes, ...blogRoutes)
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
