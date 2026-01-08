<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="page-grid">
        <!-- LEFT SIDEBAR (Desktop) -->
        <aside class="sidebar d-none d-md-block">
          <!-- City Guide -->
          <div class="side-panel">
            <div class="panel-head">縣市導覽</div>
            <div class="city-grid">
              <NuxtLink :to="`/area?city=${city.name}`" v-for="city in cities" :key="city.id">{{ city.name }}</NuxtLink>
            </div>
          </div>
          <!-- Service Guide -->
          <div class="side-panel mt-10">
            <div class="panel-head">服務導覽</div>
            <ul class="service-list">
              <li><a href="#">會員註冊</a></li>
              <li><a href="#">會員登入</a></li>
              <li><a href="#">討論版</a></li>
              <li><a href="#">評價版</a></li>
              <li><a href="#">最新消息</a></li>
              <li><a href="#">優惠下載</a></li>
            </ul>
          </div>
        </aside>

        <!-- MAIN CONTENT -->
        <main class="main-content">
          <!-- Hero Banner (Mobile shows smaller version) -->
          <div class="hero-banner mb-15">
             <div class="banner-box">
               <img src="https://placehold.co/800x250/333333/ffffff?text=全台毛小孩住宿懶人包+即刻預訂" alt="Hero" class="img-responsive">
             </div>
          </div>

          <!-- Discount Booking (Icon + Scroll) -->
          <section class="section-block">
            <div class="sec-header red-bg">
              <h3>優惠訂房</h3>
              <div class="sec-tags">
                <a href="#">熱門</a> <a href="#">基隆</a> <a href="#">台北</a> <a href="#">台中</a> <a href="#">高雄</a>
              </div>
            </div>
            <div class="discount-scroll">
              <div class="d-card" v-for="h in discountHotels" :key="h.id">
                <div class="d-img">
                   <img :src="h.image" :alt="h.name">
                   <span class="city-label">{{ h.city }}</span>
                </div>
                <div class="d-info">
                   <div class="d-name">{{ h.name }}</div>
                   <div class="d-price">{{ h.price }}</div>
                </div>
              </div>
            </div>
          </section>

          <!-- Partner Stores -->
          <section class="section-block mt-15">
            <div class="sec-header brown-bg">
              <h3>特約店家 - 持QK休閒卡消費可享優惠</h3>
            </div>
            <div class="partner-grid">
               <div class="p-card" v-for="i in 4" :key="i">
                 <img :src="`https://placehold.co/180x100/ecf0f1/2c3e50?text=Partner+${i}`" alt="Partner">
                 <div class="p-text">住宿折100元</div>
               </div>
            </div>
          </section>

          <!-- Regional Recommendations (3 Columns) -->
          <section class="section-block mt-15">
             <div class="region-row">
               <!-- North -->
               <div class="region-col">
                 <div class="r-head">北部Motel訂房推薦</div>
                 <ul class="r-list">
                   <li v-for="(h, i) in regionHotels.north" :key="i">
                     <span class="c-name">[{{h.city}}]</span>
                     <span class="h-name">{{h.name}}</span>
                     <span class="h-price">{{h.price}}</span>
                   </li>
                 </ul>
               </div>
               <!-- Central -->
               <div class="region-col">
                 <div class="r-head">中部Motel訂房推薦</div>
                 <ul class="r-list">
                   <li v-for="(h, i) in regionHotels.central" :key="i">
                     <span class="c-name">[{{h.city}}]</span>
                     <span class="h-name">{{h.name}}</span>
                     <span class="h-price">{{h.price}}</span>
                   </li>
                 </ul>
               </div>
               <!-- South -->
               <div class="region-col">
                 <div class="r-head">南部Motel訂房推薦</div>
                 <ul class="r-list">
                   <li v-for="(h, i) in regionHotels.south" :key="i">
                     <span class="c-name">[{{h.city}}]</span>
                     <span class="h-name">{{h.name}}</span>
                     <span class="h-price">{{h.price}}</span>
                   </li>
                 </ul>
               </div>
             </div>
          </section>

          <!-- News -->
          <section class="section-block mt-15">
            <div class="sec-header green-bg">
               <h3>業者消息</h3>
            </div>
            <ul class="news-rows">
               <li v-for="(n, i) in news" :key="i">
                 <span class="n-tag">{{n.tag}}</span>
                 <span class="n-city">[{{n.city}}]</span>
                 <span class="n-store">{{n.store}}</span>
                 <span class="n-title">{{n.title}}</span>
               </li>
            </ul>
          </section>
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { cities, discountHotels, regionHotels, news } = useHotelData()
</script>

<style scoped>
/* Grid Layout */
.page-grid { display: flex; gap: 10px; padding-top: 10px; }
.sidebar { width: 180px; flex-shrink: 0; }
.main-content { flex: 1; min-width: 0; } /* min-width 0 prevents overflow in flex */

/* Utilities */
.mt-10 { margin-top: 10px; }
.mt-15 { margin-top: 15px; }
.mb-15 { margin-bottom: 15px; }
.img-responsive { width: 100%; height: auto; display: block; }
a { text-decoration: none; color: inherit; }

/* Sidebar */
.side-panel { border: 1px solid #7B7B7B; background: #fff; }
.panel-head { background: #FFFACD; color: #000080; font-weight: bold; padding: 5px; text-align: center; border-bottom: 1px solid #ccc; font-size: 14px; }
.city-grid { display: flex; flex-wrap: wrap; padding: 5px; }
.city-grid a { width: 50%; font-size: 12px; color: blue; padding: 3px 0; text-align: center; }
.city-grid a:hover { background: #eee; text-decoration: underline; }
.service-list li { border-bottom: 1px dotted #ccc; }
.service-list a { display: block; padding: 6px 10px; color: blue; font-size: 13px; font-weight: bold; background: #FFFFF0; text-align: center; }
.service-list a:hover { background: #ffc; }

/* Sections */
.section-block { border: 1px solid #ccc; background: #fff; border-radius: 4px; overflow: hidden; }
.sec-header { padding: 6px 10px; display: flex; justify-content: space-between; align-items: center; color: white; }
.sec-header h3 { margin: 0; font-size: 14px; font-weight: bold; }
.sec-tags a { color: white; font-size: 12px; margin-left: 8px; }

.red-bg { background: #CC0000; }
.brown-bg { background: #8B4513; }
.green-bg { background: #6B8E23; }

/* Discount Scroll */
.discount-scroll { display: flex; overflow-x: auto; padding: 10px; gap: 10px; }
.d-card { width: 140px; flex-shrink: 0; border: 1px solid #eee; }
.d-img { position: relative; height: 90px; }
.d-img img { width: 100%; height: 100%; object-fit: cover; }
.city-label { position: absolute; top: 0; right: 0; background: #CC0000; color: white; font-size: 10px; padding: 2px 4px; }
.d-info { padding: 5px; text-align: center; }
.d-name { font-size: 12px; color: #333; overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
.d-price { font-size: 12px; color: #CC0000; font-weight: bold; }

/* Partner Grid */
.partner-grid { display: flex; padding: 10px; gap: 10px; overflow-x: auto; }
.p-card { text-align: center; width: 160px; flex-shrink: 0; }
.p-card img { width: 100%; border: 1px solid #ccc; }
.p-text { font-size: 12px; color: #CC0000; margin-top: 3px; }

/* Region List (3 Cols) */
.region-row { display: flex; border-top: 1px solid #000080; }
.region-col { flex: 1; border-right: 1px solid #ccc; }
.region-col:last-child { border-right: none; }
.r-head { background: #FFE4E1; color: #800000; font-weight: bold; font-size: 13px; padding: 5px; border-bottom: 1px solid #ccc; }
.r-list { padding: 5px; list-style: none; margin: 0; }
.r-list li { font-size: 12px; border-bottom: 1px dotted #ccc; padding: 3px 0; margin-bottom: 2px; }
.c-name { color: #000080; margin-right: 2px; }
.h-name { color: #0000FF; margin-right: 3px; }
.h-price { color: #CC0000; float: right; }

/* News Rows */
.news-rows { padding: 5px 10px; background: #F0FFF0; }
.news-rows li { font-size: 13px; padding: 5px 0; border-bottom: 1px dotted #999; }
.n-tag { background: gold; font-size: 11px; padding: 1px 3px; border-radius: 2px; margin-right: 5px; }
.n-city { color: #CC0000; font-weight: bold; margin-right: 5px; }
.n-store { color: #000080; font-weight: bold; margin-right: 5px; }
.n-title { color: #333; }

/* Mobile RWD */
@media (max-width: 768px) {
  .page-grid { display: block; }
  .sidebar { display: none; }
  .region-row { display: block; }
  .region-col { border-right: none; margin-bottom: 15px; }
  
  /* Mobile: Show simple city grid instead of sidebar */
  .mobile-city-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 5px; padding: 10px; background: #f0f0f0; margin-bottom: 10px; }
}
</style>
