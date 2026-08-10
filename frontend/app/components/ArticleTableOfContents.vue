<template>
  <div v-if="items.length" class="article-toc-host">
    <section ref="inlineToc" class="article-toc article-toc-inline">
      <button
        type="button"
        class="article-toc-toggle"
        :aria-expanded="isOpen"
        aria-controls="article-toc-list"
        @click="isOpen = !isOpen"
      >
        <span>文章目錄</span>
        <span class="article-toc-icon" :class="{ open: isOpen }" aria-hidden="true">
          <svg viewBox="0 0 24 24" fill="none">
            <path d="m7 10 5 5 5-5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
      </button>
      <nav v-show="isOpen" id="article-toc-list" aria-label="文章目錄">
        <ol>
          <li v-for="item in items" :key="item.id">
            <button type="button" @click="goToHeading(item.id)">{{ item.title }}</button>
          </li>
        </ol>
      </nav>
    </section>

    <aside
      v-show="enableFloating && showFloating"
      class="article-toc article-toc-floating"
      :class="`article-toc-floating-${context}`"
      aria-label="浮動文章目錄"
    >
      <div class="article-toc-floating-title">文章目錄</div>
      <nav>
        <ol>
          <li v-for="item in items" :key="`floating-${item.id}`">
            <button type="button" @click="goToHeading(item.id)">{{ item.title }}</button>
          </li>
        </ol>
      </nav>
    </aside>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch, type PropType } from 'vue'
import type { ArticleTocItem } from '~/utils/articleToc'

const props = defineProps({
  items: {
    type: Array as PropType<ArticleTocItem[]>,
    required: true
  },
  contentRoot: {
    type: Object as PropType<HTMLElement | null>,
    default: null
  },
  scrollRoot: {
    type: Object as PropType<HTMLElement | null>,
    default: null
  },
  context: {
    type: String as PropType<'page' | 'modal'>,
    default: 'page'
  },
  enableFloating: {
    type: Boolean,
    default: true
  }
})

const inlineToc = ref<HTMLElement | null>(null)
const isOpen = ref(true)
const showFloating = ref(false)
let listeningTarget: HTMLElement | Window | null = null
let frameId = 0

const updateFloatingState = () => {
  frameId = 0
  if (!inlineToc.value) {
    showFloating.value = false
    return
  }

  const rootTop = props.scrollRoot?.getBoundingClientRect().top || 0
  showFloating.value = inlineToc.value.getBoundingClientRect().bottom < rootTop + 16
}

const handleScroll = () => {
  if (frameId) return
  frameId = window.requestAnimationFrame(updateFloatingState)
}

const unbindScroll = () => {
  listeningTarget?.removeEventListener('scroll', handleScroll)
  window.removeEventListener('resize', handleScroll)
  listeningTarget = null
}

const bindScroll = () => {
  unbindScroll()
  listeningTarget = props.scrollRoot || window
  listeningTarget.addEventListener('scroll', handleScroll, { passive: true })
  window.addEventListener('resize', handleScroll, { passive: true })
  handleScroll()
}

const goToHeading = (id: string) => {
  const target = props.contentRoot?.querySelector<HTMLElement>(
    `[data-article-toc-id="${id}"]`
  ) || document.getElementById(id)

  target?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

onMounted(() => nextTick(bindScroll))
watch(() => props.scrollRoot, () => nextTick(bindScroll))
watch(() => props.items, () => nextTick(handleScroll), { deep: true })

onBeforeUnmount(() => {
  unbindScroll()
  if (frameId) window.cancelAnimationFrame(frameId)
})
</script>

<style scoped>
.article-toc-host {
  position: relative;
  margin: 0 0 30px;
}

.article-toc {
  border: 1px solid #dfe5eb;
  border-top: 0;
  border-radius: 10px;
  background: #f3f5f7;
  box-shadow: 0 3px 12px rgba(44, 62, 80, 0.05);
}

.article-toc-toggle {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  border: 0;
  background: transparent;
  padding: 15px 18px;
  color: #2c3e50;
  font-size: 17px;
  font-weight: 700;
  cursor: pointer;
}

.article-toc-icon {
  display: grid;
  width: 30px;
  height: 30px;
  place-items: center;
  border-radius: 50%;
  background: #fff;
  color: #657382;
  box-shadow: 0 1px 4px rgba(44, 62, 80, 0.12);
  transform: rotate(-90deg);
  transition: transform 0.2s ease;
}

.article-toc-icon svg {
  width: 18px;
  height: 18px;
}

.article-toc-icon.open {
  transform: rotate(0deg);
}

.article-toc nav {
  border-top: 1px solid #edf0f3;
  padding: 12px 18px 16px;
}

.article-toc ol {
  display: grid;
  gap: 5px;
  list-style: none;
  margin: 0;
  padding: 0;
}

.article-toc li {
  position: relative;
  padding-left: 15px;
}

.article-toc li::before {
  position: absolute;
  top: 0.7em;
  left: 0;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #e74c3c;
  content: '';
}

.article-toc li button {
  border: 0;
  background: transparent;
  padding: 5px 0;
  color: #536273;
  font-size: 15px;
  font-weight: 600;
  line-height: 1.5;
  text-align: left;
  cursor: pointer;
}

.article-toc li button:hover,
.article-toc li button:focus-visible {
  color: #e74c3c;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.article-toc-floating {
  display: none;
}

@media (min-width: 1101px) {
  .article-toc-floating {
    position: fixed;
    top: 24px;
    z-index: 1010;
    display: block;
    width: 210px;
    max-height: calc(100vh - 48px);
    overflow-y: auto;
  }

  .article-toc-floating-page {
    right: auto;
    left: max(15px, calc(50% - 600px));
  }

  .article-toc-floating-modal {
    right: calc(50% + 445px);
  }

  .article-toc-floating-title {
    padding: 13px 15px;
    color: #2c3e50;
    font-size: 15px;
    font-weight: 700;
  }

  .article-toc-floating nav {
    padding: 10px 14px 13px;
  }

  .article-toc-floating li button {
    font-size: 13px;
  }

}

@media (max-width: 768px) {
  .article-toc-host {
    margin-bottom: 24px;
  }

  .article-toc-toggle {
    padding: 13px 15px;
    font-size: 16px;
  }

  .article-toc nav {
    padding: 10px 15px 13px;
  }
}
</style>
