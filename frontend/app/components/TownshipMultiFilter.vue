<template>
  <section ref="root" class="township-multi-filter" aria-label="鄉鎮市區篩選">
    <div class="filter-heading">
      <div>
        <strong>鄉鎮市區</strong>
        <span>可複選多個地區</span>
      </div>
      <button v-if="selectedIds.length || (areaScope && !scopeCleared)" type="button" class="clear-all" @click="clearAll">
        清除條件
      </button>
    </div>

    <div class="filter-controls">
      <div class="dropdown-wrap">
        <button type="button" class="dropdown-trigger" :aria-expanded="open" @click="open = !open">
          <span>{{ triggerLabel }}</span>
          <svg viewBox="0 0 20 20" fill="none" aria-hidden="true"><path d="m5 7.5 5 5 5-5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>
        </button>
        <div v-if="open" class="dropdown-panel">
          <div class="dropdown-search">
            <input v-model="search" type="search" placeholder="搜尋縣市或鄉鎮市區..." />
          </div>
          <div class="dropdown-options">
            <section v-for="group in filteredGroups" :key="group.id" class="option-group">
              <div class="group-name">{{ group.name }}</div>
              <label v-for="township in group.townships" :key="township.id" class="township-option">
                <input v-model="selectedIds" type="checkbox" :value="Number(township.id)" />
                <span>{{ township.name }}</span>
              </label>
            </section>
            <div v-if="filteredGroups.length === 0" class="no-options">找不到符合的鄉鎮市區</div>
          </div>
          <div class="dropdown-footer">
            <span>已選 {{ selectedIds.length }} 個</span>
            <button type="button" @click="open = false">完成</button>
          </div>
        </div>
      </div>
      <button type="button" class="apply-filter" @click="apply">套用篩選</button>
    </div>

    <div v-if="visibleChips.length || (areaScope && !scopeCleared)" class="selected-chips">
      <button v-if="areaScope && !scopeCleared" type="button" class="scope-chip" @click="scopeCleared = true">
        城市範圍：{{ areaScope }} <span>×</span>
      </button>
      <button v-for="township in visibleChips" :key="township.id" type="button" @click="removeTownship(township.id)">
        {{ township.name }} <span>×</span>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";

const props = withDefaults(defineProps<{
  cities: any[];
  initialTownshipIds?: number[];
  areaScope?: string;
}>(), {
  initialTownshipIds: () => [],
  areaScope: "",
});
const emit = defineEmits<{
  apply: [payload: { townshipIds: number[]; preserveArea: boolean }];
}>();

const root = ref<HTMLElement | null>(null);
const open = ref(false);
const search = ref("");
const selectedIds = ref<number[]>([...props.initialTownshipIds].map(Number));
const scopeCleared = ref(false);
const allTownships = computed(() => props.cities.flatMap((city) => city.townships || []));
const townshipById = computed(() => new Map(allTownships.value.map((item) => [Number(item.id), item])));
const visibleChips = computed(() => selectedIds.value.map((id) => townshipById.value.get(Number(id))).filter(Boolean));
const triggerLabel = computed(() => selectedIds.value.length ? `已選 ${selectedIds.value.length} 個鄉鎮市區` : "選擇鄉鎮市區");
const filteredGroups = computed(() => {
  const keyword = search.value.trim().toLowerCase();
  return props.cities.map((city) => ({
    ...city,
    townships: (city.townships || []).filter((township: any) =>
      !keyword || String(city.name).toLowerCase().includes(keyword) || String(township.name).toLowerCase().includes(keyword),
    ),
  })).filter((city) => city.townships.length > 0);
});

const removeTownship = (id: number) => {
  selectedIds.value = selectedIds.value.filter((value) => Number(value) !== Number(id));
};
const clearAll = () => {
  selectedIds.value = [];
  scopeCleared.value = Boolean(props.areaScope);
};
const apply = () => {
  open.value = false;
  emit("apply", {
    townshipIds: [...new Set(selectedIds.value.map(Number))],
    preserveArea: Boolean(props.areaScope) && !scopeCleared.value,
  });
};
const closeOnOutsideClick = (event: PointerEvent) => {
  if (root.value && !root.value.contains(event.target as Node)) open.value = false;
};
const closeOnEscape = (event: KeyboardEvent) => {
  if (event.key === "Escape") open.value = false;
};
onMounted(() => {
  document.addEventListener("pointerdown", closeOnOutsideClick);
  document.addEventListener("keydown", closeOnEscape);
});
onBeforeUnmount(() => {
  document.removeEventListener("pointerdown", closeOnOutsideClick);
  document.removeEventListener("keydown", closeOnEscape);
});
</script>

<style scoped>
.township-multi-filter { position: relative; z-index: 5; margin-bottom: 24px; border: 1px solid #e2e8f0; border-radius: 12px; background: #fff; padding: 16px; box-shadow: 0 2px 8px rgba(15,23,42,.04); }
.filter-heading, .filter-controls { display: flex; align-items: center; justify-content: space-between; gap: 14px; }
.filter-heading { margin-bottom: 12px; }
.filter-heading > div { display: flex; align-items: baseline; gap: 8px; }
.filter-heading strong { color: #334155; font-size: 15px; }
.filter-heading span { color: #94a3b8; font-size: 12px; }
.clear-all { border: 0; background: none; color: #64748b; cursor: pointer; font-size: 13px; }
.filter-controls { justify-content: flex-start; }
.dropdown-wrap { position: relative; width: min(440px, 100%); }
.dropdown-trigger { display: flex; width: 100%; min-height: 44px; align-items: center; justify-content: space-between; border: 1px solid #cbd5e1; border-radius: 8px; background: #fff; padding: 10px 13px; color: #475569; cursor: pointer; text-align: left; }
.dropdown-trigger svg { width: 18px; height: 18px; }
.dropdown-panel { position: absolute; top: calc(100% + 7px); left: 0; width: 100%; overflow: hidden; border: 1px solid #cbd5e1; border-radius: 10px; background: #fff; box-shadow: 0 16px 35px rgba(15,23,42,.16); }
.dropdown-search { padding: 12px; border-bottom: 1px solid #e2e8f0; }
.dropdown-search input { width: 100%; box-sizing: border-box; border: 1px solid #cbd5e1; border-radius: 7px; padding: 9px 11px; font-size: 14px; outline: none; }
.dropdown-search input:focus { border-color: #e74c3c; box-shadow: 0 0 0 3px rgba(231,76,60,.1); }
.dropdown-options { max-height: 330px; overflow-y: auto; padding: 6px 0; }
.option-group + .option-group { border-top: 1px solid #f1f5f9; }
.group-name { position: sticky; top: 0; padding: 8px 13px 5px; background: #f8fafc; color: #64748b; font-size: 12px; font-weight: 800; }
.township-option { display: flex; align-items: center; gap: 9px; padding: 8px 14px; color: #334155; cursor: pointer; font-size: 14px; }
.township-option:hover { background: #fff7f6; }
.township-option input { accent-color: #e74c3c; }
.dropdown-footer { display: flex; align-items: center; justify-content: space-between; border-top: 1px solid #e2e8f0; padding: 10px 12px; color: #64748b; font-size: 12px; }
.dropdown-footer button, .apply-filter { border: 0; border-radius: 7px; background: #e74c3c; padding: 9px 16px; color: #fff; font-weight: 700; cursor: pointer; }
.apply-filter { min-height: 44px; }
.selected-chips { display: flex; flex-wrap: wrap; gap: 7px; margin-top: 12px; }
.selected-chips button { border: 0; border-radius: 999px; background: #f1f5f9; padding: 6px 10px; color: #475569; cursor: pointer; font-size: 12px; }
.selected-chips button.scope-chip { background: #e0f2fe; color: #075985; }
.selected-chips span { margin-left: 4px; font-weight: 800; }
.no-options { padding: 30px 15px; color: #94a3b8; text-align: center; }
@media (max-width: 600px) {
  .filter-heading > div { align-items: flex-start; flex-direction: column; gap: 2px; }
  .filter-controls { align-items: stretch; flex-direction: column; }
  .dropdown-wrap { width: 100%; }
  .dropdown-panel { position: fixed; z-index: 20; top: 10vh; right: 14px; bottom: 10vh; left: 14px; width: auto; }
  .dropdown-options { max-height: calc(80vh - 126px); }
}
</style>
