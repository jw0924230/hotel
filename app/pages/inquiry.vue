<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="page-grid">
        <!-- Sidebar -->
        <aside class="sidebar d-none d-md-block">
           <div class="side-panel">
            <div class="panel-head">縣市導覽</div>
            <div class="city-grid">
               <NuxtLink :to="`/area?city=${c.name}`" v-for="c in cities" :key="c.id">{{ c.name }}</NuxtLink>
            </div>
           </div>
        </aside>

        <!-- Main Content -->
        <main class="main-content">
          <!-- Breadcrumb -->
          <div class="breadcrumb mb-10">
            <NuxtLink to="/">首頁</NuxtLink> &gt; <span>{{ hotel?.city || '未知' }}</span> &gt; <span class="active">{{ hotel?.name || 'Loading...' }}</span>
          </div>

          <div v-if="hotel" class="hotel-detail">
            <h1 class="hotel-title">{{ hotel.name }}</h1>

            <!-- Tabs -->
            <div class="tabs">
               <div :class="['tab-item', { active: currentTab === 'intro' }]" @click="currentTab = 'intro'">簡　　介</div>
               <div :class="['tab-item', { active: currentTab === 'news' }]" @click="currentTab = 'news'">最新消息</div>
               <div :class="['tab-item', { active: currentTab === 'rooms' }]" @click="currentTab = 'rooms'">住房介紹</div>
               <div :class="['tab-item', { active: currentTab === 'map' }]" @click="currentTab = 'map'">交通指南</div>
            </div>

            <!-- Tab Content -->
            <div class="tab-content">
               <!-- INTRO TAB -->
               <div v-if="currentTab === 'intro'">
                   <p style="white-space: pre-line;">{{ hotel.desc }}</p>
                   <hr>
                   <div class="info-row"><strong>地址：</strong>{{ hotel.address }} <a href="#" class="map-link">[電子地圖]</a></div>
                   <div class="info-row"><strong>網站：</strong><a :href="hotel.website" target="_blank">{{ hotel.website }}</a></div>
                   <div class="info-row"><strong>信箱：</strong><a :href="`mailto:${hotel.email}`">{{ hotel.email }}</a></div>
                   <div class="mt-15">
                     <img :src="hotel.image" class="img-responsive rounded">
                   </div>
               </div>

               <!-- NEWS TAB -->
               <div v-if="currentTab === 'news'">
                 <ul class="news-list">
                   <li v-for="(n, i) in hotel.news" :key="i">{{ n }}</li>
                 </ul>
               </div>

               <!-- ROOMS TAB -->
               <div v-if="currentTab === 'rooms'">
                 <div class="room-grid">
                    <div v-for="(r, i) in hotel.roomTypes" :key="i" class="room-card">
                       <img :src="r.img" alt="Room">
                       <div class="r-info">
                         <div class="r-name">{{ r.name }}</div>
                         <div class="r-price">{{ r.price }}</div>
                       </div>
                    </div>
                 </div>
               </div>

               <!-- MAP TAB -->
               <div v-if="currentTab === 'map'">
                  <div class="map-box">
                    Google Maps Placeholder
                    <br>
                    {{ hotel.address }}
                  </div>
               </div>
            </div>
            
          </div>
          <div v-else class="not-found">
            查無此旅館資料
          </div>

        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const { cities, hotelDetails } = useHotelData()

const hotelId = computed(() => Number(route.query.n))
const hotel = computed(() => (hotelDetails as any)[hotelId.value] || null)

const currentTab = ref('intro')
if (route.query.s) {
  // Map 's' param to tab if needed
}
</script>

<style scoped>
/* Reuse layout basics */
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
.mt-10 { margin-top: 10px; }
.mb-10 { margin-bottom: 10px; }
.mt-15 { margin-top: 15px; }

/* Sidebar */
.side-panel { border: 1px solid #7B7B7B; background: #fff; margin-bottom: 10px; }
.panel-head { background: #FFFACD; color: #000080; font-weight: bold; padding: 5px; text-align: center; border-bottom: 1px solid #ccc; font-size: 14px; }
.city-grid { display: flex; flex-wrap: wrap; padding: 5px; }
.city-grid a { width: 50%; font-size: 12px; color: blue; padding: 3px 0; text-align: center; text-decoration: none; }
.city-grid a:hover { background: #eee; text-decoration: underline; }

/* Breadcrumb */
.breadcrumb { font-size: 13px; color: #666; }
.breadcrumb a { color: #000080; text-decoration: none; }
.breadcrumb .active { color: #333; }

/* Hotel Detail */
.hotel-title { color: #000080; font-size: 20px; font-weight: bold; margin-bottom: 10px; border-bottom: 1px dotted #ccc; padding-bottom: 5px; }
.info-row { font-size: 14px; margin-bottom: 5px; color: #333; }
.info-row a { color: blue; text-decoration: underline; }
.img-responsive { max-width: 100%; height: auto; }
.rounded { border-radius: 4px; }

/* Tabs */
.tabs { display: flex; border-bottom: 1px solid #CC0000; margin-bottom: 15px; background: #eee; }
.tab-item {
  padding: 8px 15px; cursor: pointer; font-size: 14px; font-weight: bold; color: #555;
  border-right: 1px solid #ccc;
}
.tab-item:hover { background: #ddd; }
.tab-item.active { background: #CC0000; color: white; }

.tab-content { padding: 10px; border: 1px solid #eee; background: white; min-height: 200px; }

/* Room Grid */
.room-grid { display: flex; flex-wrap: wrap; gap: 10px; }
.room-card { width: 150px; border: 1px solid #ccc; padding: 5px; text-align: center; }
.room-card img { width: 100%; height: 100px; object-fit: cover; }
.r-name { font-size: 13px; color: #000080; margin: 5px 0; }
.r-price { color: red; font-weight: bold; font-size: 12px; }

.news-list { padding-left: 20px; }
.news-list li { margin-bottom: 5px; font-size: 13px; color: #333; }

.map-box { background: #eee; height: 300px; display: flex; align-items: center; justify-content: center; color: #666; text-align: center; }
</style>
