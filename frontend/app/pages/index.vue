<template>
  <div class="portal-page">
    <div class="container container-wide">
      <!-- SEO Header -->
      <div class="home-intro">
          <h1 class="main-title">全台飯店、商旅、汽車旅館住宿與休息推薦</h1>
      </div>

      <!-- Quick Area Menu (Grouped + Drag Scroll) -->
      <div 
         class="area-menu-wrapper"
         ref="scrollContainer"
         @mousedown="startDrag"
         @mousemove="onDrag"
         @mouseup="stopDrag"
         @mouseleave="stopDrag"
      >
        <div class="area-group" v-for="region in regionCities" :key="region.name">
           <div class="group-label">{{ region.name }}</div>
           <div class="group-cities">
              <NuxtLink 
                v-for="city in region.cities" 
                :key="city.id" 
                :to="`/area/${city.id}/1`" 
                class="area-chip"
                draggable="false" 
              >
                {{ city.name }}
              </NuxtLink>
           </div>
        </div>
      </div>

      <div class="cities-container">
        <section v-for="cityData in selectedCitiesData" :key="cityData.name" class="city-section">
          <div class="sec-header">
            <h2 class="city-title">
               <svg class="title-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
               {{ cityData.name }}嚴選住宿
            </h2>
            <div class="sec-tags">
               <NuxtLink :to="`/area/${cityData.id}/1`" class="more-link">
                 查看全部 <span class="arrow">→</span>
               </NuxtLink>
            </div>
          </div>
          <div class="hotel-grid">
            <div class="h-card" v-for="h in cityData.hotels" :key="h.id">
               <NuxtLink :to="`/detail/${h.id}`" class="card-link">
                   <div class="h-img-wrapper">
                     <img :src="h.image" :alt="h.name" loading="lazy" @error="handleImageError">
                     <div class="price-tag" v-if="h.price">{{ h.price }}</div>
                   </div>
                    <div class="h-info">
                      <h3 class="h-name">{{ h.name }}</h3>
                      <div class="h-address">
                        <svg class="addr-icon" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
                        {{ h.address }}
                      </div>
                    </div>
               </NuxtLink>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { joinURL } from 'ufo'

const { defaultImage, handleImageError } = useHotelData()
const config = useRuntimeConfig()
const baseURL = config.app.baseURL

const getImageUrl = (imgName: string) => {
  if (!imgName) return defaultImage
  if (imgName.startsWith('http://') || imgName.startsWith('https://')) {
    return imgName
  }
  return joinURL(baseURL, `data/images/${imgName}`)
}

const backendUrl = process.env.BACKEND_API_URL || 'http://localhost:8080'

const { data: locationData } = await useAsyncData('locations', async () => {
  const [categories, regions] = await Promise.all([
    $fetch<any[]>(`${backendUrl}/api/categories?type=city`),
    $fetch<any[]>(`${backendUrl}/api/regions`)
  ])
  const cityMap = new Map(categories.map(city => [city.name, {
    id: city.sort_order || city.id,
    name: city.name
  }]))

  return {
    cities: categories.map(city => ({
      id: city.sort_order || city.id,
      name: city.name
    })),
    regions: regions.map(region => ({
      name: region.name,
      cities: region.cities.map((name: string) => cityMap.get(name)).filter(Boolean)
    }))
  }
})

const regionCities = computed(() => locationData.value?.regions || [])

const { data: selectedCitiesData } = await useAsyncData('home-cities', async () => {
  const featuredCityNames = ['台北', '新北', '桃園', '台中', '台南', '高雄']
  const cityByName = new Map((locationData.value?.cities || []).map(city => [city.name, city]))
  const citiesToFetch = featuredCityNames
    .map(name => cityByName.get(name))
    .filter(Boolean) as Array<{ id: number, name: string }>
  
  const promises = citiesToFetch.map(async (city) => {
    const result = await $fetch<any>(`${backendUrl}/api/hotels?limit=6&area=${encodeURIComponent(city.name)}`)
    return {
      id: city.id,
      name: city.name,
      hotels: (result.data || []).map((hotel: any) => ({
        id: hotel.id,
        name: hotel.name,
        image: hotel.images?.length ? getImageUrl(hotel.images[0]) : joinURL(baseURL, `data/images/${hotel.id}.jpg`),
        price: hotel.price || '',
        address: hotel.address
      }))
    }
  })
  
  return Promise.all(promises)
})

const scrollContainer = ref<HTMLElement | null>(null)

useSeoMeta({
  title: '全台飯店、商旅、汽車旅館住宿與休息推薦',
  description: '不論是臨時需要假日休息、平日休息，還是規劃一趟輕旅行的假日住宿、平日住宿，這裡一次整理全台熱門的飯店、商旅與汽車旅館推薦清單，滿足不同族群與使用情境需求。從適合短暫放鬆的2小時、3小時休息方案，到高 CP 值的過夜住宿選擇，完整比較地點、價格與彈性時段，協助你快速找到最適合的住宿或休息空間，無論情侶約會、商務出差或臨時歇腳，都能安心入住、輕鬆選擇。'
})

// Drag to Scroll Logic
let isDown = false
let startX = 0
let scrollLeft = 0

const startDrag = (e: MouseEvent) => {
  if (!scrollContainer.value) return
  isDown = true
  scrollContainer.value.classList.add('active')
  startX = e.pageX - scrollContainer.value.offsetLeft
  scrollLeft = scrollContainer.value.scrollLeft
}

const stopDrag = () => {
  if (!scrollContainer.value) return
  isDown = false
  scrollContainer.value.classList.remove('active')
}

const onDrag = (e: MouseEvent) => {
  if (!isDown || !scrollContainer.value) return
  e.preventDefault()
  const x = e.pageX - scrollContainer.value.offsetLeft
  const walk = (x - startX) * 2 // Scroll-fast
  scrollContainer.value.scrollLeft = scrollLeft - walk
}
</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

/* Intro */
.home-intro { margin-bottom: 30px; text-align: center; max-width: 800px; margin-left: auto; margin-right: auto; }
.main-title { font-size: 28px; font-weight: 800; color: #2C3E50; margin-bottom: 15px; }
.main-desc { font-size: 15px; color: #555; line-height: 1.6; }

/* Area Menu Wrapper */
.area-menu-wrapper {
  display: flex;
  overflow-x: auto;
  gap: 20px;
  padding: 15px 5px;
  margin-bottom: 30px;
  scrollbar-width: none;
  -ms-overflow-style: none;
  cursor: grab;
  user-select: none; /* Prevent text selection while dragging */
}
.area-menu-wrapper::-webkit-scrollbar { display: none; }
.area-menu-wrapper.active { cursor: grabbing; cursor: -webkit-grabbing; }

.area-group {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  background: white;
  padding: 5px 15px;
  border-radius: 50px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  border: 1px solid #eee;
}

.group-label {
  font-weight: 800;
  color: #E74C3C;
  margin-right: 15px;
  font-size: 14px;
  white-space: nowrap;
  border-right: 2px solid #eee;
  padding-right: 15px;
  height: 20px;
  line-height: 20px;
}

.group-cities {
  display: flex;
  gap: 8px;
}

.area-chip {
  padding: 6px 14px;
  background: #fdfdfd;
  border: 1px solid #f0f0f0;
  border-radius: 20px;
  color: #555;
  font-weight: 600;
  font-size: 14px;
  text-decoration: none;
  transition: all 0.2s ease;
  white-space: nowrap;
}
.area-chip:hover {
  background: #E74C3C;
  color: white;
  border-color: #E74C3C;
}

.city-section { margin-bottom: 50px; }

/* Header Styling */
.sec-header { 
  display: flex; 
  justify-content: space-between; 
  align-items: center; 
  margin-bottom: 25px; 
  border-left: 5px solid #E74C3C;
  padding-left: 15px;
}
.city-title { margin: 0; font-size: 24px; font-weight: 700; color: #2C3E50; display: flex; align-items: center; }
.city-title .icon { margin-right: 8px; font-size: 20px; }

.more-link { 
  color: #E74C3C; 
  font-weight: 600; 
  font-size: 14px; 
  text-decoration: none; 
  display: flex; 
  align-items: center; 
  transition: transform 0.2s;
}
.more-link:hover { transform: translateX(5px); }
.more-link .arrow { margin-left: 5px; font-size: 16px; }

/* Grid Styling */
.hotel-grid { 
  display: grid; 
  grid-template-columns: repeat(3, 1fr); 
  gap: 25px; 
}
@media (min-width: 992px) {
    .hotel-grid { grid-template-columns: repeat(3, 1fr); } 
}
@media (min-width: 1200px) {
    .hotel-grid { grid-template-columns: repeat(3, 1fr); }
}

/* Card Styling */
.h-card { 
  background: white; 
  border-radius: 12px; 
  overflow: hidden; 
  box-shadow: 0 10px 20px rgba(0,0,0,0.05); 
  transition: transform 0.3s, box-shadow 0.3s; 
  height: 100%;
}
.h-card:hover { 
  transform: translateY(-8px); 
  box-shadow: 0 15px 30px rgba(0,0,0,0.1); 
}
.card-link { text-decoration: none; color: inherit; display: block; height: 100%; }

.h-img-wrapper { position: relative; height: 200px; overflow: hidden; }
.h-img-wrapper img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.5s; }
.h-card:hover .h-img-wrapper img { transform: scale(1.1); }

.price-tag { 
  position: absolute; 
  bottom: 0; 
  right: 0; 
  background: rgba(231, 76, 60, 0.9); 
  color: white; 
  padding: 8px 15px; 
  font-weight: bold; 
  font-size: 14px; 
  border-top-left-radius: 12px;
}

.h-info { padding: 20px; }
.h-name { font-weight: 700; font-size: 18px; color: #2C3E50; margin-bottom: 8px; line-height: 1.4; display: -webkit-box; -webkit-line-clamp: 1; line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; }
.h-address { font-size: 14px; color: #95A5A6; display: flex; align-items: center; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; gap: 4px; }
.addr-icon { width: 14px; height: 14px; flex-shrink: 0; color: #95A5A6; }
.title-icon { width: 22px; height: 22px; vertical-align: middle; color: #E74C3C; display: inline-block; margin-right: 6px; }

/* Responsive Adjustments */
@media (max-width: 992px) {
  .hotel-grid { grid-template-columns: repeat(2, 1fr); }
}

@media (max-width: 576px) {
  .hotel-grid { grid-template-columns: 1fr; }
  .h-img-wrapper { height: 180px; }
  .city-title { font-size: 20px; }
}
</style>
