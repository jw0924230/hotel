<template>
  <section
    class="klook-search-widget"
    :class="{ 'is-failed': loadFailed, 'is-loaded': loadSucceeded }"
    aria-label="Klook 行程搜尋"
  >
    <ins
      ref="widgetRoot"
      class="klk-aff-widget"
      data-wid="122345"
      data-height="340px"
      data-adid="1369782"
      data-lang="zh-TW"
      data-prod="search_vertical"
      data-currency="TWD"
    >
      <a
        href="https://www.klook.com/zh-TW/"
        target="_blank"
        rel="sponsored noopener noreferrer"
        >Klook.com</a
      >
    </ins>

    <div v-if="loadFailed" class="klook-blocked-overlay" role="alert">
      <div class="klook-blocked-icon" aria-hidden="true">
        <svg viewBox="0 0 24 24" fill="none">
          <path
            d="M12 8v4m0 4h.01M10.3 4.4 3.2 17a2 2 0 0 0 1.75 3h14.1a2 2 0 0 0 1.75-3L13.7 4.4a2 2 0 0 0-3.4 0Z"
            stroke="currentColor"
            stroke-width="1.8"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
      <div class="klook-blocked-message">
        <strong>Klook 搜尋功能被阻擋</strong>
        <p>請暫停此網站的 AdBlock 或內容阻擋功能後，再重新整理頁面。</p>
      </div>
      <div class="klook-blocked-actions">
        <button type="button" @click="reloadPage">暫停後重新整理</button>
        <a
          href="https://www.klook.com/zh-TW/"
          target="_blank"
          rel="sponsored noopener noreferrer"
          >前往 Klook</a
        >
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from "vue";

const KLOOK_LOADER_URL =
  "https://affiliate.klook.com/widget/fetch-iframe-init.js";
const widgetRoot = ref<HTMLElement | null>(null);
const loadFailed = ref(false);
const loadSucceeded = ref(false);
const adBlockDetected = ref(false);
const loaderScripts: HTMLScriptElement[] = [];
let verificationTimer: ReturnType<typeof window.setTimeout> | null = null;
let verificationAttempts = 0;
let adBlockBait: HTMLDivElement | null = null;

const detectAdBlock = async () => {
  adBlockBait = document.createElement("div");
  adBlockBait.className = "adsbox ad-banner ad-placement advertisement";
  adBlockBait.setAttribute("aria-hidden", "true");
  adBlockBait.style.cssText =
    "position:absolute;left:-10000px;top:-10000px;width:1px;height:1px;pointer-events:none;";
  document.body.appendChild(adBlockBait);

  await new Promise<void>((resolve) => {
    window.requestAnimationFrame(() =>
      window.requestAnimationFrame(() => resolve()),
    );
  });

  const styles = window.getComputedStyle(adBlockBait);
  adBlockDetected.value =
    styles.display === "none" ||
    styles.visibility === "hidden" ||
    adBlockBait.offsetHeight === 0;
  adBlockBait.remove();
  adBlockBait = null;
};

const showBlockedOverlay = () => {
  adBlockDetected.value = true;
  loadSucceeded.value = false;
  loadFailed.value = true;
};

const reloadPage = () => window.location.reload();

const renderedSuccessfully = () => {
  const root = widgetRoot.value;
  const iframe = root?.querySelector("iframe");
  if (!root || !iframe) return false;

  const rootStyles = window.getComputedStyle(root);
  const iframeStyles = window.getComputedStyle(iframe);
  return (
    rootStyles.display !== "none" &&
    rootStyles.visibility !== "hidden" &&
    iframeStyles.display !== "none" &&
    iframeStyles.visibility !== "hidden" &&
    iframe.getBoundingClientRect().height > 0
  );
};

const appendLoaderScript = (source: string) => {
  if (!widgetRoot.value || renderedSuccessfully()) return;

  const script = document.createElement("script");
  script.src = source;
  script.async = true;
  script.type = "text/javascript";
  script.dataset.klookWidgetLoader = "true";
  script.addEventListener("load", verifyRendered, { once: true });
  script.addEventListener("error", showBlockedOverlay, { once: true });
  loaderScripts.push(script);
  document.head.appendChild(script);
};

const verifyRendered = () => {
  const root = widgetRoot.value;
  if (!root) return;

  if (renderedSuccessfully()) {
    loadFailed.value = false;
    loadSucceeded.value = true;
    return;
  }

  verificationAttempts += 1;
  if (verificationAttempts >= 40) {
    showBlockedOverlay();
    return;
  }

  verificationTimer = window.setTimeout(verifyRendered, 200);
};

const initializeWidget = async () => {
  await nextTick();
  await detectAdBlock();
  if (adBlockDetected.value) {
    showBlockedOverlay();
    return;
  }
  const root = widgetRoot.value;
  if (!root || root.querySelector("iframe") || root.dataset.rendered === "true")
    return;

  loadFailed.value = false;
  loadSucceeded.value = false;
  verificationAttempts = 0;
  appendLoaderScript(KLOOK_LOADER_URL);
};

onMounted(initializeWidget);

onBeforeUnmount(() => {
  if (verificationTimer) window.clearTimeout(verificationTimer);
  adBlockBait?.remove();
  loaderScripts.forEach((script) => script.remove());
});
</script>

<style scoped>
.klook-search-widget {
  position: relative;
  width: 100%;
  min-width: 0;
  height: 120px;
  min-height: 120px;
  margin: 0 0 32px;
  overflow: hidden;
  border-radius: 10px;
  background: #f4f6f8;
}

.klook-search-widget.is-loaded {
  height: 340px;
  min-height: 340px;
}

.klk-aff-widget {
  display: block;
  width: 100%;
  max-width: 100%;
  min-height: 340px;
  overflow: hidden;
}

.klk-aff-widget > a {
  display: grid;
  min-height: 340px;
  place-items: center;
  color: #536273;
  font-weight: 600;
}

.is-failed .klk-aff-widget {
  visibility: hidden;
}

.klook-search-widget.is-failed {
  height: 120px;
  min-height: 120px;
}

.klook-blocked-overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 16px 20px;
  background: rgba(244, 246, 248, 0.97);
  color: #334155;
  text-align: center;
}

.klook-blocked-icon {
  display: grid;
  width: 46px;
  height: 46px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 50%;
  background: #fff;
  color: #e74c3c;
  box-shadow: 0 4px 14px rgba(44, 62, 80, 0.1);
}

.klook-blocked-icon svg {
  width: 25px;
  height: 25px;
}

.klook-blocked-message {
  min-width: 0;
  flex: 1 1 auto;
  text-align: left;
}

.klook-blocked-message strong {
  font-size: 18px;
}

.klook-blocked-message p {
  max-width: 520px;
  margin: 4px 0 0;
  color: #64748b;
  font-size: 14px;
  line-height: 1.6;
}

.klook-blocked-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex: 0 0 auto;
}

.klook-blocked-actions button {
  border: 0;
  border-radius: 7px;
  background: #e74c3c;
  padding: 9px 16px;
  color: #fff;
  font-weight: 700;
  cursor: pointer;
}

.klook-blocked-actions a {
  color: #536273;
  font-size: 14px;
  font-weight: 600;
  text-decoration: underline;
  text-underline-offset: 3px;
}

@media (max-width: 768px) {
  .klook-search-widget {
    margin-bottom: 24px;
    border-radius: 8px;
  }

  .klook-blocked-overlay {
    gap: 10px;
    padding: 12px;
  }

  .klook-blocked-icon {
    display: none;
  }

  .klook-blocked-message strong {
    font-size: 15px;
  }

  .klook-blocked-message p {
    font-size: 12px;
    line-height: 1.4;
  }

  .klook-blocked-actions {
    flex-direction: column;
    gap: 6px;
  }

  .klook-blocked-actions button {
    padding: 7px 10px;
    font-size: 12px;
  }

  .klook-blocked-actions a {
    font-size: 12px;
  }
}
</style>
