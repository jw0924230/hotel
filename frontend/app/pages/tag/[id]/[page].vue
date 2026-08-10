<template>
  <TagResultsPage
    :tag="tag"
    :page="currentPage"
    :initial-hotels="result.hotels"
    :initial-total="result.total"
    :all-tags="allTags || []"
  />
</template>

<script setup lang="ts">
import { computed } from "vue";
import { joinURL } from "ufo";

definePageMeta({ key: (route) => route.fullPath });

const route = useRoute();
const config = useRuntimeConfig();
const baseURL = config.app.baseURL;
const tagId = computed(() => String(route.params.id));
const currentPage = computed(() => Number(route.params.page) || 1);
const liveCache = useState<any>("tag-live-result", () => null);

const { data: fetchedTag } = await useAsyncData(`hotel-tag-${tagId.value}`, () =>
  $fetch<any>(`${config.public.backendApiUrl}/api/hotel-tags/${tagId.value}`),
);
const { data: allTags } = await useAsyncData("hotel-tags-public", () =>
  $fetch<any[]>(`${config.public.backendApiUrl}/api/hotel-tags`),
);
const tag = computed(() => fetchedTag.value || { id: Number(tagId.value), name: "未知標籤" });
const key = computed(() => `tag-${tagId.value}-${currentPage.value}`);
const mapHotels = (items: any[]) => items.map((hotel: any) => ({
  id: hotel.id,
  name: hotel.name,
  image: hotel.images?.length
    ? (hotel.images[0].startsWith("http") ? hotel.images[0] : joinURL(baseURL, `data/images/${hotel.images[0]}`))
    : joinURL(baseURL, `data/images/${hotel.id}.jpg`),
  price: hotel.price || "",
  address: hotel.address,
  tags: hotel.tags || [],
}));

const { data: fetched } = await useAsyncData(key.value, async () => {
  if (liveCache.value?.key === key.value) return liveCache.value;
  const response = await $fetch<any>(`${config.public.backendApiUrl}/api/hotels`, {
    query: { page: currentPage.value, limit: 20, tag_id: tagId.value },
  });
  return { hotels: mapHotels(response.data || []), total: response.total || 0 };
}, {
  getCachedData(dataKey, nuxtApp) {
    if (liveCache.value?.key === dataKey) return liveCache.value;
    return nuxtApp.isHydrating ? nuxtApp.payload.data[dataKey] : nuxtApp.static.data[dataKey];
  },
});
const result = computed(() => fetched.value || { hotels: [], total: 0 });

useSeoMeta({
  title: computed(() => `${tag.value.name}飯店、商旅、汽車旅館住宿與休息推薦`),
  description: computed(() => `${tag.value.name}臨時需要假日休息、平日休息，還是規劃一趟輕旅行的假日住宿、平日住宿，這裡一次整理熱門的飯店、商旅與汽車旅館推薦清單，滿足不同族群與使用情境需求。從適合短暫放鬆的2小時、3小時休息方案，到高 CP 值的過夜住宿選擇，完整比較地點、價格與彈性時段，協助你快速找到最適合的住宿或休息空間，無論情侶約會、商務出差或臨時歇腳，都能安心入住、輕鬆選擇。`),
});

useHead({
  link: [{
    rel: "canonical",
    href: computed(() => `https://www.qk3houronline.com${route.path}`),
  }],
});
</script>
