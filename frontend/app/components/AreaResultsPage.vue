<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="sec-header">
        <div>
          <h1 class="city-title">{{ heading }} (第 {{ page }} 頁)</h1>
          <div class="breadcrumbs">
            <NuxtLink to="/">首頁</NuxtLink> &gt;
            <button v-if="township" type="button" @click="goCity">{{ city.name }}</button>
            <span v-else class="active">{{ city.name }}</span>
            <template v-if="township"> &gt; <span class="active">{{ township.name }}</span></template>
          </div>
        </div>
      </div>

      <TownshipSingleFilter
        :city="city"
        :selected-township-id="township ? Number(township.id) : null"
        @select="goTownship"
      />

      <KlookDealsWidget
        v-if="klookCityId"
        :key="klookCityId"
        :city-id="klookCityId"
      />

      <div v-if="loading" class="loading-state">正在取得最新資料...</div>
      <div v-else-if="hotels.length" class="hotel-grid">
        <div v-for="hotel in hotels" :key="hotel.id" class="h-card">
          <NuxtLink :to="`/detail/${hotel.id}`" class="card-link h-image-link">
            <div class="h-img-wrapper">
              <img :src="hotel.image" :alt="hotel.name" loading="lazy" @error="handleImageError">
              <div v-if="hotel.price" class="price-tag">{{ hotel.price }}</div>
            </div>
          </NuxtLink>
          <div class="h-info">
            <NuxtLink :to="`/detail/${hotel.id}`" class="hotel-name-link"><h2 class="h-name">{{ hotel.name }}</h2></NuxtLink>
            <div class="h-address">{{ hotel.address }}</div>
            <div v-if="hotel.tags?.length" class="hotel-card-tags" aria-label="旅館標籤">
              <NuxtLink v-for="tag in hotel.tags" :key="tag.id" :to="`/tag/${tag.id}/1`">{{ tag.name }}</NuxtLink>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="not-found"><p>此分類目前沒有旅館資料</p></div>

      <div v-if="totalPages > 1" class="pagination">
        <button v-if="page > 1" class="page-btn" type="button" @click="goPage(page - 1)">上一頁</button>
        <span class="page-info">第 {{ page }} / {{ totalPages }} 頁</span>
        <button v-if="page < totalPages" class="page-btn" type="button" @click="goPage(page + 1)">下一頁</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { joinURL } from "ufo";
import { resolveKlookCityId } from "~/utils/klook";

const props = defineProps<{
  city: any;
  township?: any | null;
  page: number;
  initialHotels: any[];
  initialTotal: number;
}>();

const config = useRuntimeConfig();
const baseURL = config.app.baseURL;
const { handleImageError } = useHotelData();
const hotels = ref(props.initialHotels);
const total = ref(props.initialTotal);
const loading = ref(false);
const pageSize = 20;
const liveCache = useState<any>("area-live-result", () => null);

const klookCityId = computed(() =>
  resolveKlookCityId(props.city.name, props.township?.name),
);

const heading = computed(() =>
  props.township
    ? `${props.city.name}${props.township.name}住宿與休息推薦`
    : `${props.city.name}住宿與休息推薦`,
);
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)));

const mapHotels = (items: any[]) => items.map((hotel) => ({
  id: hotel.id,
  name: hotel.name,
  image: hotel.images?.length
    ? (hotel.images[0].startsWith("http") ? hotel.images[0] : joinURL(baseURL, `data/images/${hotel.images[0]}`))
    : joinURL(baseURL, `data/images/${hotel.id}.jpg`),
  price: hotel.price || "",
  address: hotel.address,
  tags: hotel.tags || [],
}));

const fetchLiveAndNavigate = async (township: any | null, targetPage: number) => {
  loading.value = true;
  try {
    const params = new URLSearchParams({
      page: String(targetPage),
      limit: String(pageSize),
      area: props.city.name,
    });
    if (township) params.set("township_ids", String(township.id));
    const result = await $fetch<any>(`${config.public.backendApiUrl}/api/hotels?${params}`);
    const key = township
      ? `area-${props.city.id}-${township.id}-${targetPage}`
      : `area-${props.city.id}-${targetPage}`;
    liveCache.value = { key, hotels: mapHotels(result.data || []), total: result.total || 0 };
    const path = township
      ? `/area/${props.city.id}/${township.id}/${targetPage}`
      : `/area/${props.city.id}/${targetPage}`;
    await navigateTo(path);
  } finally {
    loading.value = false;
  }
};

const goTownship = (township: any | null) => {
  if (!township && !props.township && props.page === 1) return;
  void fetchLiveAndNavigate(township, 1);
};
const goCity = () => {
  if (!props.township && props.page === 1) return;
  void fetchLiveAndNavigate(null, 1);
};
const goPage = (targetPage: number) => void fetchLiveAndNavigate(props.township || null, targetPage);
</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }
.sec-header { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 18px; border-bottom: 1px solid #ddd; padding-bottom: 15px; }
.city-title { margin: 0 0 8px; font-size: 24px; font-weight: 700; color: #2c3e50; }
.breadcrumbs { color: #7f8c8d; font-size: 14px; }
.breadcrumbs a, .breadcrumbs button { color: #2c3e50; text-decoration: none; background: none; border: 0; padding: 0; cursor: pointer; font: inherit; }
.breadcrumbs .active { color: #e74c3c; }
.hotel-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
.h-card { background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 5px rgba(0,0,0,.05); transition: transform .2s; }
.h-card:hover { transform: translateY(-5px); box-shadow: 0 5px 15px rgba(0,0,0,.1); }
.card-link { text-decoration: none; color: inherit; display: block; }
.hotel-name-link { color: inherit; text-decoration: none; }
.h-img-wrapper { position: relative; height: 160px; overflow: hidden; }
.h-img-wrapper img { width: 100%; height: 100%; object-fit: cover; }
.price-tag { position: absolute; bottom: 0; right: 0; background: rgba(231,76,60,.9); color: white; padding: 5px 10px; font-size: 13px; font-weight: 700; border-top-left-radius: 8px; }
.h-info { padding: 15px; }
.h-name { font-weight: 700; font-size: 16px; margin: 0 0 5px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.h-address { font-size: 13px; color: #7f8c8d; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.hotel-card-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; }
.hotel-card-tags a { border-radius: 999px; background: #f1f5f9; padding: 4px 9px; color: #64748b; font-size: 11px; font-weight: 700; text-decoration: none; }
.hotel-card-tags a:hover { background: #fee2e2; color: #c0392b; }
.not-found, .loading-state { text-align: center; padding: 40px; color: #7f8c8d; }
.pagination { display: flex; justify-content: center; align-items: center; margin-top: 40px; gap: 15px; }
.page-btn { padding: 8px 16px; background: white; border: 1px solid #ddd; border-radius: 4px; color: #2c3e50; cursor: pointer; }
.page-info { font-weight: 700; color: #555; }
@media (max-width: 992px) { .hotel-grid { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 768px) { .hotel-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 576px) { .hotel-grid { grid-template-columns: 1fr; } .sec-header { align-items: flex-start; } }
</style>
