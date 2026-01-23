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
import { computed } from 'vue'

// Import all area JSON files
const areaFiles = import.meta.glob('../data/areas/area_*.json', { eager: true })

// Region mapping configuration
const regionMap = {
  '基隆': '北部', '台北': '北部', '新北': '北部', '桃園': '北部', '新竹': '北部', '宜蘭': '北部',
  '苗栗': '中部', '台中': '中部', '彰化': '中部', '南投': '中部', '雲林': '中部',
  '嘉義': '南部', '台南': '南部', '高雄': '南部', '屏東': '南部',
  '花蓮': '東部', '台東': '東部',
  '澎湖': '東部', '金門': '東部', '馬祖': '東部' // Group islands with East or separate if needed
}

// Process data
const processAreas = () => {
  const regions = {
    '北部': [],
    '中部': [],
    '南部': [],
    '東部': [] // Includes East and Islands
  }
  
  for (const path in areaFiles) {
    const mod = areaFiles[path]
    const data = mod.default || mod
    
    if (Array.isArray(data) && data.length > 0) {
      const areaName = data[0].area
      
      // Extract ID from filename
      const match = path.match(/area_(\d+)\.json$/)
      const id = match ? parseInt(match[1]) : 999
      
      let regionName = regionMap[areaName]
      if (!regionName) {
        // Fallback for unknown areas, verify if they belong to existing groups or default to '北部' etc?
        // Let's check common ones. If missing from map, maybe put in '其他' or just skip grouping logic (put in North? or iterate keys)
        // For now, assume map covers standard Taiwan areas.
        // If "屏東縣" vs "屏東", our data seems to use short names "基隆", "台北".
        // Use '東部' as fallback or create '離島'
         regionName = '東部' 
      }

      regions[regionName].push({
        id,
        name: areaName,
        displayName: `${areaName}住宿與休息推薦`,
        link: `/area/${id}/1`
      })
    }
  }

  // Convert to array for rendering, ensuring order: North, Central, South, East
  const result = [
    { name: '北部', areas: regions['北部'].sort((a, b) => a.id - b.id) },
    { name: '中部', areas: regions['中部'].sort((a, b) => a.id - b.id) },
    { name: '南部', areas: regions['南部'].sort((a, b) => a.id - b.id) },
    { name: '東部', areas: regions['東部'].sort((a, b) => a.id - b.id) }
  ]
  
  return result
}

const regionGroups = processAreas()
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
