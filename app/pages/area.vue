<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="page-grid">
        <!-- Sidebar (Reused layout structure) -->
        <aside class="sidebar d-none d-md-block">
          <div class="side-panel">
            <div class="panel-head">縣市導覽</div>
            <div class="city-grid">
               <NuxtLink :to="`/area?city=${c.name}`" v-for="c in cities" :key="c.id">{{ c.name }}</NuxtLink>
            </div>
          </div>
          <div class="side-panel mt-10">
            <div class="panel-head">服務導覽</div>
            <ul class="service-list">
              <li><a href="#">會員註冊</a></li>
              <li><a href="#">最新消息</a></li>
              <li><a href="#">特約店家</a></li>
            </ul>
          </div>
        </aside>

        <!-- Main Content -->
        <main class="main-content">
          <!-- Breadcrumb/Title Area -->
          <div class="area-header mb-15">
             <div class="breadcrumb">
               <NuxtLink to="/">首頁</NuxtLink> &gt; <span>{{ currentCity }}</span>
             </div>
             <div class="filter-bar">
               <span class="badgem active">全部</span>
               <span class="badgem">汽車旅館</span>
               <span class="badgem">飯店/旅館</span>
               <span class="badgem">民宿</span>
             </div>
          </div>

          <!-- Result List -->
          <div class="result-list">
            <div v-if="filteredHotels.length === 0" class="no-result">
               查無資料，請選擇其他地區 (試試 "彰化")
            </div>

            <div v-for="(hotel, index) in filteredHotels" :key="hotel.id" class="result-item">
               <div class="r-index">{{ index + 1 }}</div>
               <div class="r-thumb">
                 <img :src="hotel.image" :alt="hotel.name">
               </div>
               <div class="r-content">
                 <h3 class="r-title">
                   <NuxtLink to="#">{{ hotel.name }}</NuxtLink>
                 </h3>
                 <div class="r-desc">{{ hotel.desc }}</div>
                 <div class="r-meta">
                   <div class="r-price">{{ hotel.price }}</div>
                   <div class="r-loc">{{ hotel.city }}</div>
                 </div>
                 <div class="r-action">
                    <NuxtLink :to="`/inquiry?n=${hotel.id}`" class="btn-detail">詳全文</NuxtLink>
                 </div>
               </div>
            </div>
          </div>

        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const { cities, allHotels } = useHotelData()

const currentCity = computed(() => {
  return (route.query.city as string) || '全部'
})

const filteredHotels = computed(() => {
  const city = route.query.city
  if (!city) return allHotels
  return allHotels.filter(h => h.city === city)
  // Simple filter for mock purposes
})
</script>

<style scoped>
/* Reusing Grid from index.vue (Scoped, so needs copy or global) 
   Ideally this should be in a shared CSS file, but strictly following single-file habits for the task.
*/
.portal-page { padding: 10px 0; font-family: Arial, sans-serif; }
.container-wide { max-width: 1000px; margin: 0 auto; padding: 0 10px; }
.page-grid { display: flex; gap: 10px; }
.sidebar { width: 180px; flex-shrink: 0; }
.main-content { flex: 1; min-width: 0; }
.d-none { display: none; }
.d-md-block { display: block; }

@media (min-width: 768px) {
  .d-none { display: block !important; }
  .d-md-block { display: block !important; }
}

/* Sidebar Styles */
.side-panel { border: 1px solid #7B7B7B; background: #fff; margin-bottom: 10px; }
.panel-head { background: #FFFACD; color: #000080; font-weight: bold; padding: 5px; text-align: center; border-bottom: 1px solid #ccc; font-size: 14px; }
.city-grid { display: flex; flex-wrap: wrap; padding: 5px; }
.city-grid a { width: 50%; font-size: 12px; color: blue; padding: 3px 0; text-align: center; text-decoration: none; }
.city-grid a:hover { background: #eee; text-decoration: underline; }
.service-list { list-style: none; padding: 0; margin: 0; }
.service-list li { border-bottom: 1px dotted #ccc; }
.service-list a { display: block; padding: 6px 10px; color: blue; font-size: 13px; font-weight: bold; background: #FFFFF0; text-align: center; text-decoration: none; }

/* Area Header */
.area-header { background: #f9f9f9; padding: 10px; border: 1px solid #eee; border-radius: 4px; }
.breadcrumb { font-size: 13px; color: #666; margin-bottom: 10px; }
.breadcrumb a { color: #000080; text-decoration: none; }
.filter-bar { display: flex; gap: 5px; }
.badgem { font-size: 12px; padding: 3px 8px; border: 1px solid #ccc; background: white; cursor: pointer; border-radius: 3px; }
.badgem.active { background: #000080; color: white; border-color: #000080; }

/* Result List */
.result-item { display: flex; border: 1px solid #ddd; padding: 10px; margin-bottom: 10px; background: white; position: relative; }
.r-index { position: absolute; top: 0; left: 0; background: #ccc; color: white; font-size: 10px; padding: 2px 5px; }
.r-thumb { width: 140px; height: 100px; flex-shrink: 0; margin-right: 15px; }
.r-thumb img { width: 100%; height: 100%; object-fit: cover; border: 1px solid #eee; }
.r-content { flex: 1; display: flex; flex-direction: column; justify-content: space-between; }
.r-title { margin: 0 0 5px; font-size: 16px; }
.r-title a { color: #000080; text-decoration: none; }
.r-title a:hover { text-decoration: underline; }
.r-desc { font-size: 13px; color: #666; margin-bottom: 5px; }
.r-meta { display: flex; justify-content: space-between; align-items: center; font-size: 13px; }
.r-price { color: #CC0000; font-weight: bold; }
.r-loc { color: #888; }
.r-action { text-align: right; margin-top: 5px; }
.btn-detail { font-size: 12px; color: #fff; background: #666; padding: 3px 8px; text-decoration: none; border-radius: 2px; }
.btn-detail:hover { background: #CC0000; }

.no-result { padding: 20px; text-align: center; color: #666; }
</style>
