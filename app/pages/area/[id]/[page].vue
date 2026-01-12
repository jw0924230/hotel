<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="sec-header">
        <h1 class="city-title">
           <span class="icon">📍</span> {{ currentCityName }} 住宿推薦 (第 {{ currentPage }} 頁)
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
               <div class="h-name">{{ h.name }}</div>
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
const { getCityDataById, defaultImage, cities, handleImageError } = useHotelData()


const areaId = computed(() => route.params.id as string)
const currentPage = computed(() => parseInt(route.params.page as string) || 1)
const pageSize = 20

// Fetch city data
const rawData = computed(() => getCityDataById(areaId.value))
const currentCityName = computed(() => {
    const c = cities.find(c => String(c.id) === areaId.value)
    return c ? c.name : '未知地區'
})

const totalHotels = computed(() => rawData.value.length)
const totalPages = computed(() => Math.ceil(totalHotels.value / pageSize))

const finalHotels = computed(() => {
    const start = (currentPage.value - 1) * pageSize
    const end = start + pageSize
    
    return rawData.value.slice(start, end).map(item => {
        if (!item || !item.link) return null
        const id = item.link.split('n=')[1] || ''
        return {
            id,
            name: item.name,
            image: joinURL(baseURL, `data/images/${id}.jpg`),
            price: item.price,
            address: item.address,
            link: item.link
        }
    }).filter(i => i !== null)
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
.city-title .icon { margin-right: 8px; }

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
.h-name { font-weight: bold; font-size: 16px; margin-bottom: 5px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.h-address { font-size: 13px; color: #7f8c8d; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.not-found { text-align: center; padding: 40px; color: #7f8c8d; }

/* Pagination */
.pagination { display: flex; justify-content: center; align-items: center; margin-top: 40px; gap: 15px; }
.page-btn { padding: 8px 16px; background: white; border: 1px solid #ddd; border-radius: 4px; text-decoration: none; color: #2C3E50; transition: all 0.2s; }
.page-btn:hover { background: #f0f0f0; border-color: #ccc; }
.page-info { font-weight: bold; color: #555; }
</style>
