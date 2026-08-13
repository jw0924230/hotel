<template>
  <div class="klook-deals-widget">
    <ins
      class="klk-aff-widget"
      data-aid="122345"
      :data-city_id="cityId"
      data-country_id="1014"
      data-tag_id="0"
      data-currency=""
      data-lang=""
      data-label1=""
      data-label2=""
      data-label3=""
      data-prod="deals_widget"
      data-total="4"
      data-width="1176"
    >
      <a href="//www.klook.com/">Klook.com</a>
    </ins>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted } from "vue";

defineProps<{
  cityId: string;
}>();

const KLOOK_WIDGET_URL =
  "https://affiliate.klook.com/widget/fetch-iframe-init.js";
let loaderScript: HTMLScriptElement | null = null;

onMounted(() => {
  loaderScript = document.createElement("script");
  loaderScript.type = "text/javascript";
  loaderScript.async = true;
  loaderScript.src = KLOOK_WIDGET_URL;
  document.head.appendChild(loaderScript);
});

onBeforeUnmount(() => {
  loaderScript?.remove();
  loaderScript = null;
});
</script>

<style scoped>
.klook-deals-widget {
  width: 100%;
  min-width: 0;
  margin: 0 0 20px;
  overflow-x: auto;
  overflow-y: hidden;
  -webkit-overflow-scrolling: touch;
}

.klk-aff-widget {
  display: block;
  width: 1176px;
  min-width: 1176px;
  margin: 0 auto;
}
</style>
