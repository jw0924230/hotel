<template>
  <div class="portal-page">
    <div class="container container-wide">
      <header class="sec-header">
        <h1>{{ tag.name }}旅館推薦（第 {{ page }} 頁）</h1>
        <div class="breadcrumbs">
          <NuxtLink to="/">首頁</NuxtLink> &gt;
          <span>特色分類</span> &gt;
          <span class="active">{{ tag.name }}</span>
        </div>
      </header>

      <nav class="feature-category-nav" aria-label="特色分類">
        <div class="feature-category-heading">
          <h2>特色分類</h2>
          <span>選擇想找的住宿特色</span>
        </div>
        <div class="feature-category-list">
          <NuxtLink
            v-for="item in allTags"
            :key="item.id"
            :to="`/tag/${item.id}/1`"
            class="feature-category-pill"
            :class="{ active: Number(item.id) === Number(tag.id) }"
            :aria-current="Number(item.id) === Number(tag.id) ? 'page' : undefined"
          >
            {{ item.name }}
          </NuxtLink>
        </div>
      </nav>

      <div v-if="loading" class="loading-state">正在取得最新資料...</div>
      <div v-else-if="hotels.length" class="hotel-grid">
        <article v-for="hotel in hotels" :key="hotel.id" class="h-card">
          <NuxtLink :to="`/detail/${hotel.id}`" class="card-link">
            <div class="h-img-wrapper">
              <img :src="hotel.image" :alt="hotel.name" loading="lazy" @error="handleImageError" />
              <div v-if="hotel.price" class="price-tag">{{ hotel.price }}</div>
            </div>
          </NuxtLink>
          <div class="h-info">
            <NuxtLink :to="`/detail/${hotel.id}`" class="hotel-name-link"><h2>{{ hotel.name }}</h2></NuxtLink>
            <div class="h-address">{{ hotel.address }}</div>
            <div v-if="hotel.tags?.length" class="hotel-card-tags" aria-label="旅館標籤">
              <NuxtLink v-for="item in hotel.tags" :key="item.id" :to="`/tag/${item.id}/1`" :class="{ active: item.id === tag.id }">
                {{ item.name }}
              </NuxtLink>
            </div>
          </div>
        </article>
      </div>
      <div v-else class="not-found">此標籤目前沒有旅館資料</div>

      <nav v-if="totalPages > 1" class="pagination" aria-label="搜尋結果分頁">
        <button v-if="page > 1" type="button" @click="goPage(page - 1)">上一頁</button>
        <span>第 {{ page }} / {{ totalPages }} 頁</span>
        <button v-if="page < totalPages" type="button" @click="goPage(page + 1)">下一頁</button>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { joinURL } from "ufo";

const props = defineProps<{
  tag: any;
  page: number;
  initialHotels: any[];
  initialTotal: number;
  allTags: any[];
}>();

const config = useRuntimeConfig();
const baseURL = config.app.baseURL;
const { handleImageError } = useHotelData();
const hotels = ref(props.initialHotels);
const total = ref(props.initialTotal);
const loading = ref(false);
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / 20)));
const liveCache = useState<any>("tag-live-result", () => null);
const mapHotels = (items: any[]) => items.map((hotel: any) => ({
  ...hotel,
  image: hotel.images?.length
    ? (hotel.images[0].startsWith("http") ? hotel.images[0] : joinURL(baseURL, `data/images/${hotel.images[0]}`))
    : joinURL(baseURL, `data/images/${hotel.id}.jpg`),
  tags: hotel.tags || [],
}));

const goPage = async (targetPage: number) => {
  loading.value = true;
  try {
    const result = await $fetch<any>(`${config.public.backendApiUrl}/api/hotels`, {
      query: { page: targetPage, limit: 20, tag_id: props.tag.id },
    });
    liveCache.value = {
      key: `tag-${props.tag.id}-${targetPage}`,
      hotels: mapHotels(result.data || []),
      total: result.total || 0,
    };
    await navigateTo(`/tag/${props.tag.id}/${targetPage}`);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }
.sec-header { margin-bottom: 24px; border-bottom: 1px solid #ddd; padding-bottom: 15px; }
.sec-header h1 { margin: 0 0 8px; color: #2c3e50; font-size: 24px; }
.breadcrumbs { color: #7f8c8d; font-size: 14px; }
.breadcrumbs a { color: #2c3e50; text-decoration: none; }
.breadcrumbs .active { color: #e74c3c; }
.feature-category-nav { margin-bottom: 22px; border-radius: 10px; background: #eef1f4; padding: 18px 20px 20px; }
.feature-category-heading { display: flex; align-items: baseline; gap: 10px; margin-bottom: 14px; }
.feature-category-heading h2 { margin: 0; color: #2c3e50; font-size: 18px; }
.feature-category-heading span { color: #7f8c8d; font-size: 13px; }
.feature-category-list { display: flex; flex-wrap: wrap; gap: 9px; }
.feature-category-pill { border: 1px solid #d7dde3; border-radius: 999px; background: #fff; padding: 7px 14px; color: #536273; font-size: 13px; font-weight: 700; line-height: 1.4; text-decoration: none; transition: border-color .2s, background .2s, color .2s, transform .2s; }
.feature-category-pill:hover { border-color: #e74c3c; color: #c0392b; transform: translateY(-1px); }
.feature-category-pill:focus-visible { outline: 3px solid rgba(231,76,60,.25); outline-offset: 2px; }
.feature-category-pill.active { border-color: #e74c3c; background: #e74c3c; color: #fff; }
.hotel-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
.h-card { overflow: hidden; border-radius: 8px; background: white; box-shadow: 0 2px 5px rgba(0,0,0,.05); transition: transform .2s; }
.h-card:hover { transform: translateY(-5px); box-shadow: 0 5px 15px rgba(0,0,0,.1); }
.card-link, .hotel-name-link { display: block; color: inherit; text-decoration: none; }
.h-img-wrapper { position: relative; height: 160px; overflow: hidden; }
.h-img-wrapper img { width: 100%; height: 100%; object-fit: cover; }
.price-tag { position: absolute; right: 0; bottom: 0; border-top-left-radius: 8px; background: rgba(231,76,60,.9); padding: 5px 10px; color: white; font-size: 13px; font-weight: 700; }
.h-info { padding: 15px; }
.h-info h2 { margin: 0 0 5px; overflow: hidden; color: #2c3e50; font-size: 16px; text-overflow: ellipsis; white-space: nowrap; }
.h-address { overflow: hidden; color: #7f8c8d; font-size: 13px; text-overflow: ellipsis; white-space: nowrap; }
.hotel-card-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; }
.hotel-card-tags a { border-radius: 999px; background: #f1f5f9; padding: 4px 9px; color: #64748b; font-size: 11px; font-weight: 700; text-decoration: none; }
.hotel-card-tags a:hover, .hotel-card-tags a.active { background: #fee2e2; color: #c0392b; }
.not-found, .loading-state { padding: 50px 20px; color: #7f8c8d; text-align: center; }
.pagination { display: flex; align-items: center; justify-content: center; gap: 15px; margin-top: 40px; }
.pagination button { border: 1px solid #ddd; border-radius: 4px; background: white; padding: 8px 16px; color: #2c3e50; cursor: pointer; }
.pagination span { color: #555; font-weight: 700; }
@media (max-width: 992px) { .hotel-grid { grid-template-columns: repeat(3, 1fr); } }
@media (max-width: 768px) { .hotel-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 576px) {
  .hotel-grid { grid-template-columns: 1fr; }
  .feature-category-nav { padding: 16px; }
  .feature-category-heading { display: block; }
  .feature-category-heading span { display: block; margin-top: 4px; }
  .feature-category-list { gap: 7px; }
  .feature-category-pill { padding: 7px 12px; }
}
</style>
