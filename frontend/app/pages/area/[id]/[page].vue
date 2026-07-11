<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="sec-header">
        <h1 class="city-title">
           <svg class="title-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
           {{ currentCityName }} 住宿與休息推薦 (第 {{ currentPage }} 頁)
        </h1>
        <div class="breadcrumbs">
           <NuxtLink to="/">首頁</NuxtLink> &gt; <span class="active">{{ currentCityName }}</span>
        </div>
      </div>

      <div class="hotel-grid" v-if="finalHotels.length > 0">
        <div v-for="h in finalHotels" :key="h.id" class="h-card">
           <NuxtLink :to="`/detail/${h.id}`" class="card-link">
             <div class="h-img-wrapper">
               <img :src="h.image" :alt="h.name" loading="lazy" @error="handleImageError">
               <div class="price-tag">{{ h.price }}</div>
             </div>
             <div class="h-info">
               <h2 class="h-name">{{ h.name }}</h2>
               <div class="h-address">{{ h.address }}</div>
             </div>
           </NuxtLink>
        </div>
      </div>
      <div v-else class="not-found">
        <p>此頁面無更多資料</p>
      </div>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPages > 1">
         <NuxtLink 
            v-if="currentPage > 1" 
            :to="`/area/${areaId}/${currentPage - 1}`" 
            class="page-btn"
         >
            上一頁
         </NuxtLink>
         <span class="page-info"> 第 {{ currentPage }} / {{ totalPages }} 頁 </span>
         <NuxtLink 
            v-if="currentPage < totalPages" 
            :to="`/area/${areaId}/${currentPage + 1}`" 
            class="page-btn"
         >
            下一頁
         </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { joinURL } from 'ufo'

const route = useRoute()
const config = useRuntimeConfig()
const baseURL = config.app.baseURL
const { defaultImage, handleImageError } = useHotelData()

const areaId = computed(() => route.params.id as string)
const currentPage = computed(() => parseInt(route.params.page as string) || 1)
const pageSize = 20

const getImageUrl = (imgName: string) => {
  if (!imgName) return defaultImage
  if (imgName.startsWith('http://') || imgName.startsWith('https://')) {
    return imgName
  }
  return joinURL(baseURL, `data/images/${imgName}`)
}

const { data: cities } = await useAsyncData('city-categories', () =>
  $fetch<any[]>(`${config.public.backendApiUrl}/api/categories?type=city`)
)

const currentCityName = computed(() => {
  const c = (cities.value || []).find(city => String(city.sort_order || city.id) === areaId.value)
  return c ? c.name : '未知地區'
})

// Fetch city data exclusively from the API.
const { data: fetchResult } = await useAsyncData(`area-${areaId.value}-${currentPage.value}`, async () => {
  let citiesList = cities.value
  if (!citiesList) {
    citiesList = await $fetch<any[]>(`${config.public.backendApiUrl}/api/categories?type=city`)
  }
  const c = (citiesList || []).find(city => String(city.sort_order || city.id) === areaId.value)
  const cityName = c ? c.name : ''
  
  if (cityName && cityName !== '其他') {
    const result = await $fetch<any>(`${config.public.backendApiUrl}/api/hotels?page=${currentPage.value}&limit=${pageSize}&area=${encodeURIComponent(cityName)}`)
    return {
      hotels: (result.data || []).map((hotel: any) => ({
        id: hotel.id,
        name: hotel.name,
        image: hotel.images?.length ? getImageUrl(hotel.images[0]) : joinURL(baseURL, `data/images/${hotel.id}.jpg`),
        price: hotel.price || '',
        address: hotel.address
      })),
      total: result.total || 0
    }
  }
  
  return { hotels: [], total: 0 }
}, {
  watch: [areaId, currentPage]
})

const finalHotels = computed(() => fetchResult.value?.hotels || [])
const totalHotels = computed(() => fetchResult.value?.total || 0)
const totalPages = computed(() => Math.ceil(totalHotels.value / pageSize))

useSeoMeta({
  title: computed(() => `${currentCityName.value}飯店、商旅、汽車旅館住宿與休息推薦`),
  description: computed(() => `${currentCityName.value} 臨時需要假日休息、平日休息，還是規劃一趟輕旅行的假日住宿、平日住宿，這裡一次整理${currentCityName.value} 熱門的飯店、商旅與汽車旅館推薦清單，滿足不同族群與使用情境需求。從適合短暫放鬆的2小時、3小時休息方案，到高 CP 值的過夜住宿選擇，完整比較地點、價格與彈性時段，協助你快速找到最適合的住宿或休息空間，無論情侶約會、商務出差或臨時歇腳，都能安心入住、輕鬆選擇。`)
})

</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

.sec-header { 
  display: flex; 
  justify-content: space-between; 
  align-items: flex-end; 
  margin-bottom: 25px; 
  border-bottom: 1px solid #ddd;
  padding-bottom: 15px;
}
.city-title { margin: 0; font-size: 24px; font-weight: 700; color: #2C3E50; display: flex; align-items: center; }
.title-icon { width: 22px; height: 22px; color: #E74C3C; margin-right: 8px; }

.breadcrumbs { color: #7f8c8d; font-size: 14px; margin-bottom: 5px; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }
.breadcrumbs a:hover { text-decoration: underline; }
.breadcrumbs .active { color: #E74C3C; }

/* Grid */
.hotel-grid { 
  display: grid; 
  grid-template-columns: repeat(4, 1fr); 
  gap: 20px; 
}
@media (max-width: 992px) { .hotel-grid { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 768px) { .hotel-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 576px) { .hotel-grid { grid-template-columns: 1fr; } }

/* Card */
.h-card { background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 5px rgba(0,0,0,0.05); transition: transform 0.2s; }
.h-card:hover { transform: translateY(-5px); box-shadow: 0 5px 15px rgba(0,0,0,0.1); }
.card-link { text-decoration: none; color: inherit; display: block; }

.h-img-wrapper { position: relative; height: 160px; overflow: hidden; }
.h-img-wrapper img { width: 100%; height: 100%; object-fit: cover; }
.price-tag { position: absolute; bottom: 0; right: 0; background: rgba(231, 76, 60, 0.9); color: white; padding: 5px 10px; font-size: 13px; font-weight: bold; border-top-left-radius: 8px; }

.h-info { padding: 15px; }
.h-name { font-weight: bold; font-size: 16px; margin: 0 0 5px 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.h-address { font-size: 13px; color: #7f8c8d; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.not-found { text-align: center; padding: 40px; color: #7f8c8d; }

/* Pagination */
.pagination { display: flex; justify-content: center; align-items: center; margin-top: 40px; gap: 15px; }
.page-btn { padding: 8px 16px; background: white; border: 1px solid #ddd; border-radius: 4px; text-decoration: none; color: #2C3E50; transition: all 0.2s; }
.page-btn:hover { background: #f0f0f0; border-color: #ccc; }
.page-info { font-weight: bold; color: #555; }
</style>
