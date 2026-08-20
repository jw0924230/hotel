<template>
  <AreaResultsPage
    :city="currentCity"
    :township="currentTownship"
    :page="currentPage"
    :initial-hotels="result.hotels"
    :initial-total="result.total"
  />
</template>

<script setup lang="ts">
import { computed } from "vue";
import { joinURL } from "ufo";
import { buildAreaSeoDescription } from "~/utils/seo";

definePageMeta({ key: (route) => route.fullPath });

const route = useRoute();
const config = useRuntimeConfig();
const baseURL = config.app.baseURL;
const areaId = computed(() => String(route.params.id));
const townshipId = computed(() => String(route.params.townshipId));
const currentPage = computed(() => Number(route.params.page) || 1);
const liveCache = useState<any>("area-live-result", () => null);

const { data: locations } = await useAsyncData("locations", () =>
  $fetch<any>(`${config.public.backendApiUrl}/api/regions/combined`),
);
const currentCity = computed(() =>
  (locations.value?.cities || []).find((city: any) => String(city.id) === areaId.value) ||
  { id: Number(areaId.value), name: "未知地區", townships: [] },
);
const currentTownship = computed(() =>
  (currentCity.value.townships || []).find((township: any) => String(township.id) === townshipId.value) ||
  { id: Number(townshipId.value), name: "未知鄉鎮市區" },
);

const key = computed(() => `area-${areaId.value}-${townshipId.value}-${currentPage.value}`);
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
    query: {
      page: currentPage.value,
      limit: 20,
      area: currentCity.value.name,
      township_ids: String(currentTownship.value.id),
    },
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
  title: computed(() => `${currentCity.value.name}${currentTownship.value.name}飯店、商旅、汽車旅館住宿與休息推薦`),
  description: computed(() => buildAreaSeoDescription(`${currentCity.value.name}${currentTownship.value.name}`)),
});

</script>
