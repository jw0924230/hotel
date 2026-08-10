<template>
  <nav class="township-filter" aria-label="鄉鎮市區篩選">
    <div class="township-filter-heading">
      <strong>鄉鎮市區</strong>
      <span>左右滑動選擇地區</span>
    </div>

    <div class="township-filter-scroller" role="radiogroup" aria-label="選擇鄉鎮市區">
      <NuxtLink
        :to="`/area/${city.id}/1`"
        class="township-pill"
        :class="{ active: selectedTownshipId === null }"
        role="radio"
        :aria-checked="selectedTownshipId === null"
        :aria-current="selectedTownshipId === null ? 'page' : undefined"
        @click.prevent="selectTownship(null)"
      >
        全部
      </NuxtLink>
      <NuxtLink
        v-for="township in city.townships || []"
        :key="township.id"
        :to="`/area/${city.id}/${township.id}/1`"
        class="township-pill"
        :class="{ active: Number(township.id) === selectedTownshipId }"
        role="radio"
        :aria-checked="Number(township.id) === selectedTownshipId"
        :aria-current="Number(township.id) === selectedTownshipId ? 'page' : undefined"
        @click.prevent="selectTownship(township)"
      >
        {{ township.name }}
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
defineProps<{
  city: any;
  selectedTownshipId: number | null;
}>();

const emit = defineEmits<{
  select: [township: any | null];
}>();

const selectTownship = (township: any | null) => {
  emit("select", township);
};
</script>

<style scoped>
.township-filter { margin-bottom: 24px; }
.township-filter-heading { display: flex; align-items: baseline; gap: 8px; margin-bottom: 10px; }
.township-filter-heading strong { color: #334155; font-size: 15px; }
.township-filter-heading span { color: #94a3b8; font-size: 12px; }
.township-filter-scroller {
  display: flex;
  gap: 9px;
  overflow-x: auto;
  padding: 2px 2px 10px;
  overscroll-behavior-inline: contain;
  scrollbar-width: thin;
  scroll-snap-type: x proximity;
  -webkit-overflow-scrolling: touch;
}
.township-pill {
  flex: 0 0 auto;
  scroll-snap-align: start;
  border: 1px solid #d7dde3;
  border-radius: 999px;
  background: #fff;
  padding: 8px 15px;
  color: #536273;
  cursor: pointer;
  font-size: 13px;
  font-weight: 700;
  line-height: 1.4;
  text-decoration: none;
  white-space: nowrap;
  transition: border-color .2s, background .2s, color .2s, transform .2s;
}
.township-pill:hover { border-color: #e74c3c; color: #c0392b; transform: translateY(-1px); }
.township-pill:focus-visible { outline: 3px solid rgba(231,76,60,.25); outline-offset: 2px; }
.township-pill.active { border-color: #e74c3c; background: #e74c3c; color: #fff; }
</style>
