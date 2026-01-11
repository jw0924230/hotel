<template>
  <div class="portal-page">
    <div class="container container-narrow">
      <div v-if="article">
        <div class="breadcrumbs">
           <NuxtLink to="/">首頁</NuxtLink> &gt; 
           <NuxtLink to="/blog">部落格</NuxtLink> &gt; 
           <span class="active">{{ article.title }}</span>
        </div>

        <article class="article-detail">
            <header class="article-header">
                <div class="meta-top">
                    <span class="category">{{ article.category }}</span>
                    <span class="date">{{ article.date }}</span>
                </div>
                <h1 class="main-title">{{ article.title }}</h1>
            </header>

            <div class="featured-image">
                <img :src="article.image" :alt="article.title">
            </div>

            <div class="article-body" v-html="article.content"></div>
            
            <div class="article-footer">
                <NuxtLink to="/blog" class="btn-back">← 返回文章列表</NuxtLink>
            </div>
        </article>

      </div>
      <div v-else class="not-found">
        <p>找不到此文章</p>
        <NuxtLink to="/blog" class="btn-back">返回列表</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed } from 'vue'
import articles from '~/../data_json/blog/articles.json'

const route = useRoute()
const id = route.params.id as string

const article = computed(() => {
    return articles.find(a => a.id === id)
})
</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-narrow { max-width: 800px; }

.breadcrumbs { margin-bottom: 30px; color: #7f8c8d; font-size: 14px; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }
.breadcrumbs .active { color: #95A5A6; }

.article-detail { background: white; padding: 40px; border-radius: 12px; box-shadow: 0 5px 20px rgba(0,0,0,0.05); }

.article-header { margin-bottom: 30px; text-align: center; }
.meta-top { margin-bottom: 15px; display: flex; justify-content: center; gap: 15px; align-items: center; }
.category { background: #E74C3C; color: white; padding: 4px 12px; border-radius: 20px; font-size: 13px; font-weight: 600; }
.date { color: #95A5A6; font-size: 14px; }
.main-title { font-size: 32px; color: #2C3E50; line-height: 1.4; margin: 0; }

.featured-image { margin-bottom: 40px; border-radius: 8px; overflow: hidden; }
.featured-image img { width: 100%; height: auto; display: block; }

.article-body { font-size: 18px; line-height: 1.8; color: #2c3e50; }
.article-body :deep(h2), .article-body :deep(h3) { margin-top: 40px; margin-bottom: 20px; color: #2C3E50; font-weight: 700; }
.article-body :deep(p) { margin-bottom: 20px; }
.article-body :deep(img) { max-width: 100%; height: auto; border-radius: 8px; margin: 20px 0; }

.article-footer { margin-top: 50px; padding-top: 30px; border-top: 1px solid #eee; text-align: center; }
.btn-back { display: inline-block; padding: 10px 25px; border: 1px solid #ddd; border-radius: 30px; color: #7f8c8d; transition: all 0.2s; }
.btn-back:hover { border-color: #2C3E50; color: #2C3E50; background: transparent; }

.not-found { text-align: center; padding: 50px; }

@media(max-width: 768px) {
    .article-detail { padding: 25px; }
    .main-title { font-size: 24px; }
}
</style>
