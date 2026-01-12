<template>
  <div class="portal-page">
    <div class="container container-wide">
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
            <h3 class="city-title">
               <span class="icon">📍</span> {{ cityData.name }}嚴選住宿
            </h3>
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
                     <div class="h-name">{{ h.name }}</div>
                     <div class="h-address">{{ h.address }}</div>
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
import { ref } from 'vue'

const { selectedCitiesData, regionCities, defaultImage, handleImageError } = useHotelData()
const scrollContainer = ref<HTMLElement | null>(null)

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
.h-address { font-size: 14px; color: #95A5A6; display: flex; align-items: center; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.h-address::before { content: "📍"; margin-right: 5px; font-size: 12px; }

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
