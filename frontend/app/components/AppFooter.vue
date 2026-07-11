<template>
  <footer class="site-footer">
    <div class="container container-wide">
      <div class="footer-grid">
        <div v-for="region in regionGroups" :key="region.name" class="region-column">
          <h3 class="region-title">{{ region.name }}</h3>
          <div class="region-links">
            <NuxtLink 
              v-for="area in region.areas" 
              :key="area.id" 
              :to="area.link" 
              class="area-link"
            >
              {{ area.displayName }}
            </NuxtLink>
          </div>
        </div>
      </div>
      
      <div class="footer-bottom">
        <p>Copyright © 2026 喜助數位. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>
<script setup>
const regionGroups = [
  {
    name: '北部',
    areas: [
      { id: 1, displayName: '基隆住宿與休息推薦', link: '/area/1/1' },
      { id: 2, displayName: '台北住宿與休息推薦', link: '/area/2/1' },
      { id: 3, displayName: '新北住宿與休息推薦', link: '/area/3/1' },
      { id: 4, displayName: '桃園住宿與休息推薦', link: '/area/4/1' },
      { id: 5, displayName: '新竹住宿與休息推薦', link: '/area/5/1' },
      { id: 6, displayName: '宜蘭住宿與休息推薦', link: '/area/6/1' }
    ]
  },
  {
    name: '中部',
    areas: [
      { id: 7, displayName: '苗栗住宿與休息推薦', link: '/area/7/1' },
      { id: 8, displayName: '台中住宿與休息推薦', link: '/area/8/1' },
      { id: 10, displayName: '彰化住宿與休息推薦', link: '/area/10/1' },
      { id: 11, displayName: '南投住宿與休息推薦', link: '/area/11/1' },
      { id: 12, displayName: '雲林住宿與休息推薦', link: '/area/12/1' }
    ]
  },
  {
    name: '南部',
    areas: [
      { id: 13, displayName: '嘉義住宿與休息推薦', link: '/area/13/1' },
      { id: 14, displayName: '台南住宿與休息推薦', link: '/area/14/1' },
      { id: 16, displayName: '高雄住宿與休息推薦', link: '/area/16/1' },
      { id: 18, displayName: '屏東住宿與休息推薦', link: '/area/18/1' }
    ]
  },
  {
    name: '東部',
    areas: [
      { id: 19, displayName: '花蓮住宿與休息推薦', link: '/area/19/1' },
      { id: 20, displayName: '台東住宿與休息推薦', link: '/area/20/1' },
      { id: 21, displayName: '澎湖住宿與休息推薦', link: '/area/21/1' },
      { id: 22, displayName: '金門住宿與休息推薦', link: '/area/22/1' },
      { id: 23, displayName: '馬祖住宿與休息推薦', link: '/area/23/1' },
      { id: 24, displayName: '其他住宿與休息推薦', link: '/area/24/1' }
    ]
  }
]
</script>

<style scoped>
*, *::before, *::after { box-sizing: border-box; }

/* Container styles to match default.vue */
.container { width: 100%; padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

.site-footer {
  background: #2C3E50;
  color: #ECF0F1;
  padding: 50px 0 20px;
  margin-top: 60px;
  font-size: 14px;
}

.footer-grid {
  display: flex;
  flex-wrap: wrap;
  /* Use negative margin to counteract column padding if needed, 
     but here we use gap which is easier. 
     Just ensure justification is correct. */
  margin-right: -15px;
  margin-left: -15px;
}

.region-column {
  flex: 0 0 100%;
  max-width: 100%;
  padding-left: 15px;
  padding-right: 15px;
  margin-bottom: 30px;
}

/* Tablet & Desktop: 4 columns */
@media (min-width: 768px) {
  .region-column {
    flex: 0 0 25%;
    max-width: 25%;
    /* Use padding to create spacing, but smaller */
    padding-left: 8px;
    padding-right: 8px;
  }
}

.region-title {
  font-size: 18px; /* Larger title */
  font-weight: 700;
  color: #fff;
  margin-bottom: 15px; /* More space below title */
  padding-bottom: 8px;
  border-bottom: 2px solid #E74C3C;
  display: inline-block;
  min-width: 60px; 
}

.region-links {
  display: flex;
  flex-direction: column;
  gap: 10px; /* Increased vertical spacing */
}

.area-link {
  color: #BDC3C7;
  text-decoration: none;
  font-size: 15px; /* Larger link text */
  transition: all 0.2s;
  display: block;
  width: fit-content;
  line-height: 1.5;
}

.area-link:hover {
  color: #E74C3C;
  transform: translateX(3px);
}

.footer-bottom {
  border-top: 1px solid #34495e;
  padding-top: 20px;
  text-align: center;
  color: #7F8C8D;
  margin-top: 10px;
}
</style>
