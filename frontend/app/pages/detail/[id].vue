<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div v-if="hotel" class="hotel-detail">
        <!-- Header -->
        <div class="sec-header">
          <h1 class="hotel-title">{{ hotel.name }}</h1>
          <div class="breadcrumbs">
            <NuxtLink to="/">首頁</NuxtLink> &gt;
            <NuxtLink v-if="cityId" :to="`/area/${cityId}/1`">{{
              cityName
            }}</NuxtLink>
            <span v-else>{{ cityName || "地區" }}</span> &gt;
            <span class="active">{{ hotel.name }}</span>
          </div>
        </div>

        <div class="detail-grid">
          <!-- Left Column: Image or Slider -->
          <div class="detail-left">
            <div class="image-slider">
              <!-- Main Image Wrapper -->
              <div class="slider-wrapper">
                <img
                  :src="
                    hotelImages.length
                      ? hotelImages[activeImageIndex]
                      : processedImage
                  "
                  :alt="hotel.name"
                  class="slider-main-img"
                  @error="handleImageError"
                />
                <!-- Navigation buttons (only show if there are multiple images) -->
                <button
                  v-if="hotelImages.length > 1"
                  type="button"
                  class="slider-btn prev"
                  @click="prevImage"
                >
                  ‹
                </button>
                <button
                  v-if="hotelImages.length > 1"
                  type="button"
                  class="slider-btn next"
                  @click="nextImage"
                >
                  ›
                </button>
              </div>

              <!-- Thumbnails (only show if there are multiple images) -->
              <div class="slider-thumbs" v-if="hotelImages.length > 1">
                <div
                  v-for="(img, idx) in hotelImages"
                  :key="idx"
                  :class="['thumb-item', { active: activeImageIndex === idx }]"
                  @click="activeImageIndex = idx"
                >
                  <img :src="img" :alt="hotel.name" />
                </div>
              </div>
            </div>
          </div>

          <!-- Right Column: Basic Info -->
          <div class="detail-right">
            <div class="info-card">
              <div class="info-row" v-if="hotel.address">
                <h2 class="label">地址：</h2>
                <a
                  :href="`https://www.google.com/maps/search/?api=1&query=${hotel.address}`"
                  class="text-link"
                  target="_blank"
                  >{{ hotel.address }}</a
                >
              </div>
              <div class="info-row" v-if="hotel.phone">
                <h2 class="label">電話：</h2>
                <a :href="`tel:${hotel.phone}`" class="text-link text-phone">{{
                  hotel.phone
                }}</a>
              </div>
              <div class="info-row" v-if="hotel.website">
                <h2 class="label">網站：</h2>
                <a :href="hotel.website" target="_blank" class="text-link">{{
                  hotel.website
                }}</a>
              </div>
            </div>

            <!-- Independent Pricing Block -->
            <div class="pricing-card">
              <div class="pricing-row-single">
                <h2 class="pricing-label">住宿價格</h2>
                <span class="pricing-val">{{ formatStayPricing(hotel) }}</span>
              </div>

              <div class="pricing-row-single">
                <h2 class="pricing-label">休息價格</h2>
                <span class="pricing-val">{{ formatRestPricing(hotel) }}</span>
              </div>
            </div>

            <!-- Booking Action Button (Agoda/Booking Link) -->
            <div class="booking-action" v-if="hotel.booking_link">
              <a
                :href="hotel.booking_link"
                target="_blank"
                rel="noopener noreferrer"
                class="btn-booking"
              >
                前往訂房
              </a>
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
              <h3>簡介</h3>
            </button>
            <button
              :class="['tab-btn', { active: currentTab === 'rules' }]"
              @click="currentTab = 'rules'"
            >
              <h3>須知與規定</h3>
            </button>
          </div>

          <div class="tab-content">
            <!-- Intro -->
            <div v-show="currentTab === 'intro'" class="content-pane">
              <div class="long-text" v-html="formattedDesc"></div>
            </div>

            <!-- Rules -->
            <div v-show="currentTab === 'rules'" class="content-pane">
              <div
                v-for="(section, index) in ruleSections"
                :key="`${section.title}-${index}`"
                class="rule-block"
              >
                <h3>{{ section.title }}</h3>
                <div class="long-text" v-html="section.content"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="not-found">
        <div class="hotel-door" aria-hidden="true">
          <div class="door-number">404</div>
          <div class="door-line"></div>
          <div class="door-handle"></div>
          <div class="key-tag">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
            >
              <circle cx="8" cy="15" r="4"></circle>
              <path d="M11 12l8-8m-2 2 2 2m-5 1 2 2"></path>
            </svg>
            <span>NO VACANCY</span>
          </div>
        </div>
        <div class="not-found-copy">
          <span class="service-label">FRONT DESK · 24 HOURS</span>
          <h1>暫無這間旅館資料</h1>
          <p>
            找不到編號
            {{ route.params.id }} 的旅館資料，可能已下架或網址已經變更。
          </p>
          <div class="not-found-actions">
            <NuxtLink to="/" class="btn-back-primary">回到首頁</NuxtLink>
            <NuxtLink to="/area/2/1" class="btn-back-secondary"
              >瀏覽其他住宿</NuxtLink
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRoute } from "vue-router";

import { joinURL } from "ufo";

const route = useRoute();
const config = useRuntimeConfig();
const baseURL = config.app.baseURL;
const { defaultImage, handleImageError } = useHotelData();

const hotelId = route.params.id as string;

const { data: hotel } = await useAsyncData(`detail-hotel-${hotelId}`, () =>
  $fetch<any>(`${config.public.backendApiUrl}/api/hotels/${hotelId}`)
);
const { data: cities } = await useAsyncData("detail-city-categories", () =>
  $fetch<any[]>(`${config.public.backendApiUrl}/api/categories?type=city`),
);

const currentTab = ref("intro");
const activeImageIndex = ref(0);

// Image Logic
const processedImage = computed(() => {
  if (!hotelId) return defaultImage;
  // Primary/default image
  if (hotel.value?.images && hotel.value.images.length > 0) {
    const primaryImg = hotel.value.images[0];
    if (primaryImg.startsWith("http://") || primaryImg.startsWith("https://")) {
      return primaryImg;
    }
    return joinURL(baseURL, `data/images/${primaryImg}`);
  }
  return joinURL(baseURL, `data/images/${hotelId}.jpg`);
});

// All Images list
const hotelImages = computed(() => {
  if (!hotel.value?.images || hotel.value.images.length === 0) {
    return [];
  }
  return hotel.value.images.map((img: string) => {
    if (img.startsWith("http://") || img.startsWith("https://")) {
      return img;
    }
    return joinURL(baseURL, `data/images/${img}`);
  });
});

const prevImage = () => {
  if (activeImageIndex.value > 0) {
    activeImageIndex.value--;
  } else {
    activeImageIndex.value = hotelImages.value.length - 1;
  }
};

const nextImage = () => {
  if (activeImageIndex.value < hotelImages.value.length - 1) {
    activeImageIndex.value++;
  } else {
    activeImageIndex.value = 0;
  }
};

// Format Price Helper
const formatPrice = (price: any) => {
  if (!price) return "請電洽";
  const priceStr = String(price);
  // If it's pure numbers, append '起' for safety
  if (/^\d+$/.test(priceStr)) {
    return priceStr + "起";
  }
  return priceStr;
};

const formatStayPricing = (hotelData: any) => {
  if (!hotelData?.pricing) return formatPrice(hotelData?.price_accommodation);
  const parts = [];
  if (hotelData.pricing.weekday_stay)
    parts.push(`平日 ${hotelData.pricing.weekday_stay}`);
  if (hotelData.pricing.holiday_stay)
    parts.push(`假日 ${hotelData.pricing.holiday_stay}`);
  return parts.length > 0 ? parts.join(" / ") : "請電洽";
};

const formatRestPricing = (hotelData: any) => {
  if (!hotelData?.pricing) return formatPrice(hotelData?.price_rest);
  const parts = [];
  if (hotelData.pricing.weekday_rest) {
    const hours = hotelData.pricing.weekday_rest_hours
      ? `${hotelData.pricing.weekday_rest_hours} 小時 `
      : "";
    parts.push(`平日 ${hours}${hotelData.pricing.weekday_rest}`);
  }
  if (hotelData.pricing.holiday_rest) {
    const hours = hotelData.pricing.holiday_rest_hours
      ? `${hotelData.pricing.holiday_rest_hours} 小時 `
      : "";
    parts.push(`假日 ${hours}${hotelData.pricing.holiday_rest}`);
  }
  return parts.length > 0 ? parts.join(" / ") : "請電洽";
};

// City Logic
const cityData = computed(() => {
  if (!hotel.value?.address) return null;
  const addr = hotel.value.address;
  const found = (cities.value || []).find((city) => addr.includes(city.name));
  return found;
});

const cityName = computed(
  () => cityData.value?.name || hotel.value?.address?.substring(0, 3) || "",
);
const cityId = computed(() => cityData.value?.sort_order || cityData.value?.id);

// Format helpers
const formatText = (text: string) => {
  if (!text) return "";

  // If it contains rich HTML formatting, bypass simple raw-text formatting
  if (
    text.includes("<p>") ||
    text.includes("<h2>") ||
    text.includes("<h3>") ||
    text.includes("<b>") ||
    text.includes("<br>")
  ) {
    return text;
  }

  let cleaned = text;

  // Remove adsbygoogle
  cleaned = cleaned.replace(
    /\(adsbygoogle\s*=\s*window\.adsbygoogle\s*\|\|\s*\[\]\)\.push\(\{\}\);/g,
    "",
  );

  // Remove HTML comments
  cleaned = cleaned.replace(/<!--[\s\S]*?-->/g, "");

  // Remove window.___gcfg
  cleaned = cleaned.replace(/window\.___gcfg\s*=\s*\{.*?\};/g, "");
  cleaned = cleaned.replace(/\(function\(\)\s*\{[\s\S]*?\}\)\(\);/g, "");

  const rubbish = [
    "簡　　介",
    "最新消息",
    "優惠活動",
    "住房介紹",
    "附近景點",
    "交通指南",
    "相關連結",
    "連鎖分館",
    "精選相簿",
    "店家資訊 QRCode",
    "本站聲明：本網站上所刊登之業者圖片及文字皆為業者自行上傳，本站不負任何責任，如有侵權請來信告知。",
  ];
  rubbish.forEach((r) => {
    cleaned = cleaned.replace(new RegExp(r, "g"), "");
  });

  // Trim and normalize newlines
  cleaned = cleaned.replace(/\n\s*\n/g, "\n");
  cleaned = cleaned.trim();

  return cleaned.replace(/\n/g, "<br>");
};

const formattedDesc = computed(() =>
  formatText(hotel.value?.description || hotel.value?.stay_info || "尚無簡介"),
);

const ruleSections = computed(() => {
  const html = formatText(hotel.value?.housing_rules || "");
  if (!html) return [];

  const headingPattern = /<h3[^>]*>([\s\S]*?)<\/h3>/gi;
  const headings = [...html.matchAll(headingPattern)];
  if (headings.length === 0) {
    return [{ title: "住房規定", content: html }];
  }

  const sections: Array<{ title: string; content: string }> = [];
  const introduction = html.slice(0, headings[0].index).trim();
  if (introduction) {
    sections.push({ title: "住房規定", content: introduction });
  }

  headings.forEach((heading, index) => {
    const contentStart = (heading.index || 0) + heading[0].length;
    const contentEnd =
      index + 1 < headings.length ? headings[index + 1].index : html.length;
    const title = heading[1].replace(/<[^>]+>/g, "").trim() || "住房規定";
    sections.push({
      title,
      content: html.slice(contentStart, contentEnd).trim(),
    });
  });

  return sections;
});

const googleMapUrl = computed(() => {
  if (!hotel.value?.address) return "";
  const query = encodeURIComponent(hotel.value.address);
  return `https://www.google.com/maps/search/?api=1&query=${query}`;
});

useSeoMeta({
  title: computed(
    () =>
      `${hotel.value?.name || "旅館"}休息與住宿價格｜電話、地址、交通方式整理`,
  ),
  description: computed(
    () =>
      `${hotel.value?.name || "本旅館"}提供彈性多元的休息與住宿選擇，無論是2 小時、3 小時休息，或是平日、假日住宿需求，都能快速掌握最新價格方案。本頁整理 ${hotel.value?.name || "旅館"}休息與住宿價格、訂房電話、詳細地址與交通資訊，協助你在規劃行程或臨時住宿時，快速做出最合適的選擇。`,
  ),
});
</script>

<style scoped>
.portal-page {
  padding: 40px 0;
  background: #f8f9fa;
  min-height: 80vh;
}
.container {
  padding: 0 15px;
  margin: 0 auto;
}
.container-wide {
  max-width: 1200px;
}

.not-found {
  min-height: 650px;
  display: grid;
  grid-template-columns: minmax(280px, 420px) minmax(300px, 1fr);
  align-items: center;
  gap: clamp(40px, 8vw, 110px);
  padding: clamp(35px, 7vw, 80px);
  overflow: hidden;
  border-radius: 24px;
  background:
    radial-gradient(
      circle at 75% 20%,
      rgba(202, 164, 107, 0.16),
      transparent 30%
    ),
    linear-gradient(135deg, #17212b 0%, #23303b 52%, #111820 100%);
  box-shadow: 0 24px 70px rgba(15, 23, 42, 0.22);
}

.hotel-door {
  position: relative;
  width: min(100%, 340px);
  aspect-ratio: 0.68;
  justify-self: center;
  border: 10px solid #0d141a;
  border-bottom-width: 18px;
  border-radius: 4px 4px 0 0;
  background:
    linear-gradient(
      90deg,
      rgba(255, 255, 255, 0.03),
      transparent 20%,
      rgba(0, 0, 0, 0.12)
    ),
    #3d2f28;
  box-shadow:
    inset 0 0 0 2px rgba(202, 164, 107, 0.18),
    24px 28px 45px rgba(0, 0, 0, 0.32);
}

.door-number {
  position: absolute;
  top: 16%;
  left: 50%;
  transform: translateX(-50%);
  color: #d7b980;
  font-family: Georgia, "Times New Roman", serif;
  font-size: clamp(52px, 8vw, 82px);
  letter-spacing: 0.12em;
  text-indent: 0.12em;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.45);
}

.door-line {
  position: absolute;
  top: 35%;
  left: 18%;
  right: 18%;
  height: 1px;
  background: rgba(215, 185, 128, 0.35);
}

.door-handle {
  position: absolute;
  right: 15%;
  top: 56%;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #caa46b;
  box-shadow:
    0 0 0 6px #8d704b,
    0 5px 12px rgba(0, 0, 0, 0.45);
}

.key-tag {
  position: absolute;
  right: -42px;
  bottom: 10%;
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 10px 13px;
  border-radius: 5px;
  background: #d7b980;
  color: #17212b;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 0.12em;
  transform: rotate(-7deg);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.28);
}

.key-tag svg {
  width: 23px;
  height: 23px;
}
.not-found-copy {
  max-width: 560px;
  color: #f8fafc;
}
.service-label {
  color: #d7b980;
  font-size: 11px;
  font-weight: 800;
  letter-spacing: 0.22em;
}
.not-found-copy h1 {
  margin: 17px 0 18px;
  font-family: Georgia, "Times New Roman", serif;
  font-size: clamp(34px, 5vw, 58px);
  line-height: 1.12;
  font-weight: 500;
}
.not-found-copy p {
  max-width: 500px;
  margin: 0;
  color: #b8c2cc;
  font-size: 16px;
  line-height: 1.8;
}
.not-found-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 34px;
}
.btn-back-primary,
.btn-back-secondary {
  padding: 13px 22px;
  border-radius: 7px;
  font-size: 14px;
  font-weight: 750;
  text-decoration: none;
  transition:
    transform 0.2s,
    background-color 0.2s;
}
.btn-back-primary {
  background: #d7b980;
  color: #17212b;
}
.btn-back-secondary {
  background: rgba(255, 255, 255, 0.08);
  color: #f8fafc;
}
.btn-back-primary:hover,
.btn-back-secondary:hover {
  transform: translateY(-2px);
}
.btn-back-secondary:hover {
  background: rgba(255, 255, 255, 0.14);
}

@media (max-width: 760px) {
  .not-found {
    min-height: auto;
    grid-template-columns: 1fr;
    text-align: center;
  }
  .hotel-door {
    width: min(70vw, 260px);
  }
  .key-tag {
    right: -25px;
  }
  .not-found-copy {
    margin: 0 auto;
  }
  .not-found-actions {
    justify-content: center;
  }
}

.sec-header {
  margin-bottom: 25px;
  border-bottom: 1px solid #ddd;
  padding-bottom: 15px;
}
.hotel-title {
  font-size: 32px;
  color: #2c3e50;
  margin: 0 0 10px 0;
  font-weight: 700;
}
.breadcrumbs {
  color: #7f8c8d;
  font-size: 14px;
}
.breadcrumbs a {
  color: #2c3e50;
  text-decoration: none;
}
.breadcrumbs a:hover {
  text-decoration: underline;
}
.breadcrumbs .active {
  color: #e74c3c;
}

/* Adjusted Grid */
.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
  margin-bottom: 40px;
  align-items: stretch;
}
.detail-left {
  min-width: 0;
}
@media (max-width: 768px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

.main-img {
  width: 100%;
  height: 100%;
  min-height: 400px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
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

/* Image Slider styles */
.image-slider {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.slider-wrapper {
  position: relative;
  width: 100%;
  height: 400px;
  background: #eee;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}
.slider-main-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.slider-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border: none;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  font-size: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}
.slider-btn:hover {
  background: rgba(0, 0, 0, 0.8);
}
.slider-btn.prev {
  left: 15px;
}
.slider-btn.next {
  right: 15px;
}

.slider-thumbs {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 8px;
  scrollbar-width: thin;
  scrollbar-color: rgba(100, 116, 139, 0.4) transparent;
}
.slider-thumbs::-webkit-scrollbar {
  height: 6px;
}
.slider-thumbs::-webkit-scrollbar-track {
  background: transparent;
}
.slider-thumbs::-webkit-scrollbar-thumb {
  background: rgba(100, 116, 139, 0.4);
  border-radius: 10px;
}
.slider-thumbs::-webkit-scrollbar-thumb:hover {
  background: rgba(100, 116, 139, 0.6);
}
.slider-thumbs::-webkit-scrollbar-button {
  display: none;
}
.thumb-item {
  width: 60px;
  height: 60px;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
  flex-shrink: 0;
}
.thumb-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.thumb-item.active {
  border-color: #e74c3c;
}

.detail-right {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.info-card {
  background: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-bottom: 20px;
}

.info-row {
  margin-bottom: 18px;
  display: flex;
  align-items: baseline;
}
.info-row:last-child {
  margin-bottom: 0;
}
.info-row .label {
  width: 70px;
  font-weight: bold;
  color: #7f8c8d;
  flex-shrink: 0;
  margin: 0;
  font-size: 16px;
}
.info-row .text {
  color: #2c3e50;
  font-size: 16px;
  line-height: 1.5;
}
.info-row .text-link {
  color: #3498db;
  text-decoration: none;
  word-break: break-all;
}
.info-row .text-link:hover {
  text-decoration: underline;
}
.text-phone {
  font-weight: bold;
}

.price-text {
  color: #e74c3c;
  font-size: 28px;
  font-weight: 700;
}

/* Pricing Card Styles */
.pricing-card {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  border: 1px solid #f0f0f0;
  margin-bottom: 20px;
}
.pricing-row-single {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px dashed #eee;
}
.pricing-row-single:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.pricing-label {
  font-size: 18px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.pricing-val {
  color: #e74c3c;
  font-weight: 700;
  font-size: 18px;
  white-space: pre-line;
  text-align: right;
}

/* Booking Action Button */
.booking-action {
  margin-bottom: 20px;
}
.btn-booking {
  display: block;
  width: 100%;
  padding: 15px;
  text-align: center;
  background: linear-gradient(135deg, #e74c3c, #c0392b);
  color: white;
  font-size: 18px;
  font-weight: 700;
  text-decoration: none;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(231, 76, 60, 0.3);
  transition: all 0.3s;
}
.btn-booking:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(231, 76, 60, 0.4);
  background: linear-gradient(135deg, #ff6b6b, #e74c3c);
}

/* Tabs */
.tabs-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
  overflow: hidden;
}
.tabs-nav {
  display: flex;
  border-bottom: 1px solid #eee;
  background: #fdfdfd;
}
.tab-btn {
  flex: 1;
  padding: 15px;
  border: none;
  background: none;
  cursor: pointer;
  color: #7f8c8d;
  transition: all 0.2s;
  border-bottom: 3px solid transparent;
}
.tab-btn h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.tab-btn:hover {
  background: #f0f0f0;
  color: #2c3e50;
}
.tab-btn.active {
  color: #e74c3c;
  border-bottom-color: #e74c3c;
  background: white;
}

.tab-content {
  padding: 40px;
}
.content-pane {
  animation: fadeIn 0.3s;
}
.long-text {
  line-height: 1.8;
  color: #444;
  font-size: 15px;
}
.long-text :deep(h2) {
  font-size: 20px;
  color: #2c3e50;
  margin-top: 20px;
  margin-bottom: 10px;
  border-bottom: 1px solid #eee;
  padding-bottom: 5px;
}
.long-text :deep(h3) {
  font-size: 18px;
  color: #2c3e50;
  margin-top: 15px;
  margin-bottom: 8px;
}
.long-text :deep(ul) {
  padding-left: 20px;
  margin-bottom: 15px;
}
.long-text :deep(li) {
  margin-bottom: 5px;
  list-style-type: disc;
}

.rule-block {
  margin-bottom: 18px;
  padding: 22px;
  background: rgba(231, 76, 60, 0.05);
  border: 1px dashed #e74c3c;
  border-radius: 10px;
  box-shadow: none;
}
.rule-block:last-child {
  margin-bottom: 0;
}
.rule-block h3 {
  color: #2c3e50;
  margin-bottom: 15px;
  border-left: 4px solid #e74c3c;
  padding-left: 10px;
  font-size: 18px;
}

.not-found {
  text-align: center;
  padding: 50px;
  font-size: 18px;
  color: #7f8c8d;
}
.btn-back {
  display: inline-block;
  margin-top: 15px;
  padding: 10px 20px;
  background: #3498db;
  color: white;
  text-decoration: none;
  border-radius: 5px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
