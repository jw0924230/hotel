<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div v-if="hotel" class="hotel-detail">
        <!-- Header -->
        <div class="sec-header">
            <h1 class="hotel-title">{{ hotel.name }}</h1>
            <div class="breadcrumbs">
               <NuxtLink to="/">首頁</NuxtLink> &gt; 
               <NuxtLink v-if="cityId" :to="`/area/${cityId}/1`">{{ cityName }}</NuxtLink>
               <span v-else>{{ cityName || '地區' }}</span> &gt; 
               <span class="active">{{ hotel.name }}</span>
            </div>
        </div>

        <div class="detail-grid">
            <!-- Left Column: Image -->
            <div class="detail-left">
                <div class="main-img">
                    <img :src="processedImage" :alt="hotel.name" @error="handleImageError">
                </div>
            </div>

            <!-- Right Column: Basic Info -->
            <div class="detail-right">
                <div class="info-card">
                    <div class="info-row" v-if="hotel.address">
                        <span class="label">地址：</span>
                        <span class="text">{{ hotel.address }}</span>
                    </div>
                    <div class="info-row" v-if="hotel.phone">
                        <span class="label">電話：</span>
                        <span class="text">{{ hotel.phone }}</span>
                    </div>
                    <div class="info-row" v-if="hotel.website">
                        <span class="label">網站：</span>
                        <a :href="hotel.website" target="_blank" class="text-link">{{ hotel.website }}</a>
                    </div>
                    <div class="info-row" v-if="hotel.price">
                        <span class="label">價格：</span>
                        <span class="price-text">{{ hotel.price }}</span>
                    </div>
                    <div class="info-row" v-if="hotel.rest_info">
                        <span class="label">休息：</span>
                        <span class="text">{{ hotel.rest_info }}</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Tabs Section -->
        <div class="tabs-container">
            <div class="tabs-nav">
                <button 
                    :class="['tab-btn', { active: currentTab === 'intro' }]" 
                    @click="currentTab = 'intro'"
                >
                    簡介
                </button>
                <button 
                    :class="['tab-btn', { active: currentTab === 'rules' }]" 
                    @click="currentTab = 'rules'"
                >
                    須知與規定
                </button>
                 <button 
                    :class="['tab-btn', { active: currentTab === 'map' }]" 
                    @click="currentTab = 'map'"
                >
                    地圖導航
                </button>
            </div>

            <div class="tab-content">
                <!-- Intro -->
                <div v-show="currentTab === 'intro'" class="content-pane">
                    <div class="long-text" v-html="formattedDesc"></div>
                </div>

                <!-- Rules -->
                <div v-show="currentTab === 'rules'" class="content-pane">
                    <div class="rule-block" v-if="hotel.housing_rules">
                        <h3>住房規定</h3>
                        <div class="long-text" v-html="formatText(hotel.housing_rules)"></div>
                    </div>
                    <div class="rule-block" v-if="hotel.precautions">
                        <h3>注意事項</h3>
                        <div class="long-text" v-html="formatText(hotel.precautions)"></div>
                    </div>
                </div>

                <!-- Map -->
                <div v-show="currentTab === 'map'" class="content-pane">
                    <div class="map-container">
                        <iframe 
                            width="100%" 
                            height="450" 
                            style="border:0" 
                            loading="lazy" 
                            allowfullscreen 
                            :src="googleMapUrl">
                        </iframe>
                    </div>
                </div>
            </div>
        </div>

      </div>
      <div v-else class="not-found">
        <p>找不到此旅館資料 (ID: {{ route.params.id }})</p>
        <NuxtLink to="/" class="btn-back">回首頁</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

import { joinURL } from 'ufo'

const route = useRoute()
const config = useRuntimeConfig()
const baseURL = config.app.baseURL
const { defaultImage, cities } = useHotelData()

const hotelId = route.params.id as string

const { data: hotel } = await useFetch<any>(`/api/detail/${hotelId}`)

const currentTab = ref('intro')

// Image Logic
const processedImage = computed(() => {
    if (!hotelId) return defaultImage
    return joinURL(baseURL, `data/images/${hotelId}.jpg`)
})

const handleImageError = (e: Event) => {
    (e.target as HTMLImageElement).src = defaultImage
}

// City Logic
const cityData = computed(() => {
    if(!hotel.value?.address) return null
    // Extract first 2-3 chars as city name usually
    const addr = hotel.value.address
    const name = addr.substring(0, 3).replace('台', '臺') // normalize if needed, but data seems to use 台/臺 mixed?
    // Actually cities list has '台北', '新北'. Address might be '台北市...'
    // Lets try to match
    const found = cities.find(c => addr.includes(c.name))
    return found
})

const cityName = computed(() => cityData.value?.name || hotel.value?.address?.substring(0, 3) || '')
const cityId = computed(() => cityData.value?.id)


// Format helpers
const formatText = (text: string) => {
    if(!text) return ''
    let cleaned = text
    
    // Remove adsbygoogle
    cleaned = cleaned.replace(/\(adsbygoogle\s*=\s*window\.adsbygoogle\s*\|\|\s*\[\]\)\.push\(\{\}\);/g, '')
    
    // Remove HTML comments
    cleaned = cleaned.replace(/<!--[\s\S]*?-->/g, '')

    // Remove window.___gcfg
    cleaned = cleaned.replace(/window\.___gcfg\s*=\s*\{.*?\};/g, '')
    cleaned = cleaned.replace(/\(function\(\)\s*\{[\s\S]*?\}\)\(\);/g, '')

    // Remove repeated header-like text often found in scraping
    // e.g. "最新消息", "優惠活動" repeated
    const rubbish = ['簡　　介', '最新消息', '優惠活動', '住房介紹', '附近景點', '交通指南', '相關連結', '連鎖分館', '精選相簿', '店家資訊 QRCode', '本站聲明：本網站上所刊登之業者圖片及文字皆為業者自行上傳，本站不負任何責任，如有侵權請來信告知。']
    rubbish.forEach(r => {
        cleaned = cleaned.replace(new RegExp(r, 'g'), '')
    })

    // Trim and normalize newlines
    cleaned = cleaned.replace(/\n\s*\n/g, '\n') // collapse multiple newlines
    cleaned = cleaned.trim()

    return cleaned.replace(/\n/g, '<br>')
}

const formattedDesc = computed(() => formatText(hotel.value?.description || hotel.value?.stay_info || '尚無簡介')) // Use stay_info if description is empty, as per user's json it seems stay_info has the 'intro'

const googleMapUrl = computed(() => {
    if (!hotel.value?.address) return ''
    const query = encodeURIComponent(hotel.value.address)
    // Using simple output=embed
    return `https://maps.google.com/maps?q=${query}&output=embed&z=16`
})

</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

.sec-header { margin-bottom: 25px; border-bottom: 1px solid #ddd; padding-bottom: 15px; }
.hotel-title { font-size: 32px; color: #2C3E50; margin: 0 0 10px 0; font-weight: 700; }
.breadcrumbs { color: #7f8c8d; font-size: 14px; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }
.breadcrumbs a:hover { text-decoration: underline; }
.breadcrumbs .active { color: #E74C3C; }

/* Adjusted Grid */
.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 40px; align-items: stretch; }
@media(max-width: 768px) { .detail-grid { grid-template-columns: 1fr; } }

.main-img { 
    width: 100%; 
    height: 100%;
    min-height: 400px;
    border-radius: 12px; overflow: hidden; 
    box-shadow: 0 5px 15px rgba(0,0,0,0.1); 
    background: #eee;
    position: relative;
}
.main-img img { 
    width: 100%; 
    height: 100%; 
    object-fit: cover; 
    position: absolute;
    top: 0;
    left: 0;
}

.detail-right {
    display: flex;
    flex-direction: column;
}

.info-card { 
    background: white; padding: 30px; border-radius: 12px; 
    box-shadow: 0 5px 15px rgba(0,0,0,0.05); 
    flex: 1; /* Fill height */
    display: flex; flex-direction: column; justify-content: center;
}

.info-row { margin-bottom: 18px; display: flex; align-items: baseline; }
.info-row:last-child { margin-bottom: 0; }
.info-row .label { width: 70px; font-weight: bold; color: #7f8c8d; flex-shrink: 0; }
.info-row .text { color: #2C3E50; font-size: 16px; line-height: 1.5; }
.info-row .text-link { color: #3498DB; text-decoration: none; word-break: break-all; }
.info-row .text-link:hover { text-decoration: underline; }
.price-text { color: #E74C3C; font-size: 28px; font-weight: 700; }

/* Tabs */
.tabs-container { background: white; border-radius: 12px; box-shadow: 0 5px 15px rgba(0,0,0,0.05); overflow: hidden; }
.tabs-nav { display: flex; border-bottom: 1px solid #eee; background: #fdfdfd; }
.tab-btn { 
    flex: 1; padding: 15px; border: none; background: none; 
    font-size: 16px; cursor: pointer; color: #7f8c8d; font-weight: 600;
    transition: all 0.2s;
    border-bottom: 3px solid transparent;
}
.tab-btn:hover { background: #f0f0f0; color: #2C3E50; }
.tab-btn.active { color: #E74C3C; border-bottom-color: #E74C3C; background: white; }

.tab-content { padding: 40px; }
.content-pane { animation: fadeIn 0.3s; }
.long-text { line-height: 1.8; color: #444; font-size: 15px; }
.rule-block { margin-bottom: 30px; border-bottom: 1px dashed #eee; padding-bottom: 20px; }
.rule-block:last-child { border-bottom: none; }
.rule-block h3 { color: #2C3E50; margin-bottom: 15px; border-left: 4px solid #E74C3C; padding-left: 10px; font-size: 18px; }

.map-container { border-radius: 8px; overflow: hidden; border: 1px solid #ddd; }

.not-found { text-align: center; padding: 50px; font-size: 18px; color: #7f8c8d; }
.btn-back { display: inline-block; margin-top: 15px; padding: 10px 20px; background: #3498DB; color: white; text-decoration: none; border-radius: 5px; }

@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
