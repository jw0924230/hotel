<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="sec-header">
        <h1 class="city-title">
           <span class="icon">📍</span> {{ cityName }} 住宿推薦 (第 {{ page }} 頁)
        </h1>
        <div class="breadcrumbs">
            <NuxtLink to="/">首頁</NuxtLink> &gt; <span>{{ cityName }}</span>
        </div>
      </div>

      <div class="hotel-grid">
        <div class="h-card" v-for="h in hotels" :key="h.id">
           <NuxtLink :to="`/hotel/${h.id}`" class="card-link">
               <div class="h-img-wrapper">
                 <img :src="h.image" :alt="h.name" loading="lazy" @error="(e) => (e.target as HTMLImageElement).src = defaultImage">
                 <div class="price-tag" v-if="h.price">{{ h.price }}</div>
               </div>
               <div class="h-info">
                 <div class="h-name">{{ h.name }}</div>
                 <div class="h-address">{{ h.address }}</div>
               </div>
           </NuxtLink>
        </div>
      </div>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPage > 1">
        <NuxtLink 
            v-if="page > 1" 
            :to="`/area/${routeId}/${page - 1}`" 
            class="page-link"
        >
            上一頁
        </NuxtLink>

        <div class="page-numbers">
            <span v-for="p in totalPage" :key="p">
                <NuxtLink 
                    v-if="p !== page"
                    :to="`/area/${routeId}/${p}`" 
                    class="page-num"
                >
                    {{ p }}
                </NuxtLink>
                <span v-else class="page-num active">{{ p }}</span>
            </span>
        </div>

        <NuxtLink 
            v-if="page < totalPage" 
            :to="`/area/${routeId}/${page + 1}`" 
            class="page-link"
        >
            下一頁
        </NuxtLink>
      </div>
      
      <div v-if="hotels.length === 0" class="no-data">
        沒有找到相關資料。
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const { getCityDataById, cities, defaultImage } = useHotelData()

const routeId = computed(() => route.params.id)
const page = computed(() => Number(route.params.page) || 1)
const pageSize = 20

// Fetch City Name
const cityName = computed(() => {
    const city = cities.find(c => String(c.id) === String(routeId.value))
    return city ? city.name : ''
})

// Fetch Data and Pagination
const rawData = computed(() => getCityDataById(routeId.value as string))
const totalHotels = computed(() => rawData.value.length)
const totalPage = computed(() => Math.ceil(totalHotels.value / pageSize))

const hotels = computed(() => {
    const start = (page.value - 1) * pageSize
    const end = start + pageSize
    
    // Use similar mapping logic as in useHotelData but we need to do it here
    // or expose the helper. Since we need local images, we need access to the glob or proper URLs.
    // The previous logic for `processCityData` contained the image mapping.
    // We should probably reuse that logic or duplicate it briefly here but correct.
    
    // Wait, getCityDataById returns raw JSON. We need to process it.
    // And to process it we need access to the image list which is inside the composable scope not exported.
    // Solution: Let's assume for now we duplicate the image lookup logic or ask useHotelData to export a processor?
    // Actually, `defaultImage` is exported. But the specific mapping is not.
    
    // Let's rely on the image mapping logic being duplicated slightly or refactored.
    // Ideally refactor useHotelData to export `getProcessedCityData(id)`
    
    // To allow this component to work immediately without refactoring useHotelData extensively again:
    // I will use a simple assumption that image path is predictable if I could import the glob here too?
    // No, `import.meta.glob` works in comps too.
    
    // Let's add the image logic here.
    
    return rawData.value.slice(start, end).map(item => {
        const id = item.link.split('n=')[1] || ''
        const imageUrl = `/data/images/${id}.jpg` // This assumes public path or we need the glob again?
        // Wait, local images in `data/images` are NOT served publicly unless they are in `public/`.
        // The previous implementation utilized Vite's import glob which resolves to assets.
        // If we want dynamic images in a production build, they should be in `public` or imported.
        // Since we have a LOT of potential images, putting them in `assets` and globbing is fine for small sets.
        // But for "all images", doing a full glob here is okay as long as it matches.
        
        // Actually, for the homepage we globbed ALL images. We can do the same here.
        return {
            id,
            name: item.name,
            image: item.image_url, // Placeholder, will fix with logic below
            price: item.price,
            address: item.address,
            link: item.link
        }
    })
})

// We need the proper image URL.
// Since we cannot easily import the same glob without duplicating code/performance hit,
// accessing `useHotelData`'s internal map would be better.
// Let's stick with the plan. I will duplicate the Glob for now as it is eager and cached by Vite effectively.
const hotelImages = import.meta.glob('~/data/images/*.jpg', { eager: true, import: 'default' })

const processedHotels = computed(() => {
    return hotels.value.map(h => {
        const imagePath = Object.keys(hotelImages).find(key => key.endsWith(`/${h.id}.jpg`))
        const imageUrl = imagePath ? (hotelImages[imagePath] as string) : defaultImage
        return { ...h, image: imageUrl }
    })
})

// Override hotels with processed
const finalHotels = processedHotels

</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

.sec-header { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  margin-bottom: 30px; 
}
.city-title { margin: 0; font-size: 28px; font-weight: 700; color: #2C3E50; display: flex; align-items: center; }
.city-title .icon { margin-right: 10px; }
.breadcrumbs { font-size: 14px; color: #7f8c8d; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }

.hotel-grid { 
  display: grid; 
  grid-template-columns: repeat(3, 1fr); 
  gap: 25px; 
  margin-bottom: 40px;
}
@media (max-width: 992px) { .hotel-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 576px) { .hotel-grid { grid-template-columns: 1fr; } }

.h-card { 
  background: white; 
  border-radius: 12px; 
  overflow: hidden; 
  box-shadow: 0 5px 15px rgba(0,0,0,0.05); 
  transition: transform 0.3s; 
}
.h-card:hover { transform: translateY(-5px); }
.card-link { text-decoration: none; color: inherit; display: block; }
.h-img-wrapper { position: relative; height: 200px; overflow: hidden; }
.h-img-wrapper img { width: 100%; height: 100%; object-fit: cover; }
.price-tag { 
  position: absolute; bottom: 0; right: 0; 
  background: rgba(231, 76, 60, 0.9); color: white; 
  padding: 5px 15px; font-weight: bold; 
  border-top-left-radius: 12px;
}
.h-info { padding: 15px; }
.h-name { font-weight: 700; font-size: 18px; color: #2C3E50; margin-bottom: 5px; line-height: 1.4; display: -webkit-box; -webkit-line-clamp: 1; line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; }
.h-address { font-size: 13px; color: #95A5A6; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* Pagination */
.pagination { display: flex; justify-content: center; align-items: center; gap: 10px; padding-top: 20px; }
.page-link, .page-num {
    padding: 8px 16px;
    background: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    color: #555;
    text-decoration: none;
    transition: all 0.2s;
}
.page-link:hover, .page-num:hover { background: #E74C3C; color: white; border-color: #E74C3C; }
.page-num.active { background: #E74C3C; color: white; border-color: #E74C3C; }

.no-data { text-align: center; color: #888; font-size: 18px; padding: 50px; }
</style>
