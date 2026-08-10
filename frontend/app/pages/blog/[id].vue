<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div v-if="article">
        <div class="breadcrumbs">
           <NuxtLink to="/">首頁</NuxtLink> &gt; 
           <NuxtLink to="/blog">部落格</NuxtLink> &gt; 
           <span class="active">{{ article.title }}</span>
        </div>

        <div :class="['article-layout', { 'without-toc': !parsedArticle.items.length }]">
        <aside v-if="parsedArticle.items.length" class="toc-rail" aria-label="文章目錄">
          <ArticleTableOfContents :items="parsedArticle.items" :content-root="articleBodyRoot" :enable-floating="false" />
        </aside>
        <article class="article-detail">
            <header class="article-header">
                <h1 class="main-title">{{ article.title }}</h1>
                <div class="meta-top">
                    <span class="date">{{ article.date }}</span>
                </div>
                <div v-if="article.featureTags.length" class="feature-tags">
                  <span v-for="tag in article.featureTags" :key="tag">{{ tag }}</span>
                </div>
            </header>

            <div class="mobile-inline-toc">
              <ArticleTableOfContents :items="parsedArticle.items" :content-root="articleBodyRoot" :enable-floating="false" />
            </div>

            <div
              ref="articleBodyRoot"
              class="article-body rich-html-content"
              v-html="parsedArticle.html"
            ></div>
            <div class="article-ad" v-html="parsedAdLink"></div>
            
            <div class="article-footer">
                <NuxtLink to="/blog" class="btn-back">← 返回文章列表</NuxtLink>
            </div>
        </article>
        <aside v-if="relatedCategories.length" class="category-sidebar" aria-label="文章分類快速連結">
          <h2>文章分類</h2>
          <div v-for="item in relatedCategories" :key="item.id" class="category-group">
            <div class="category-heading"><span class="category-name">{{ item.name }}</span><NuxtLink :to="`/blog/tag/${item.id}/1`" class="category-more">看更多 <span aria-hidden="true">→</span></NuxtLink></div>
            <NuxtLink v-if="item.post" :to="`/blog/${item.post.id}`" class="related-link">
              <strong>{{ item.post.title }}</strong>
              <time>{{ formatDate(item.post.created_at) }}</time>
            </NuxtLink>
            <span v-else class="no-related">目前沒有其他相關文章</span>
          </div>
        </aside>
        </div>

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
import { computed, ref } from 'vue'
import MarkdownIt from 'markdown-it'
import { joinURL } from 'ufo'
import { buildArticleToc } from '~/utils/articleToc'

const route = useRoute()
const config = useRuntimeConfig()
const baseURL = config.app.baseURL
const id = route.params.id as string
const md = new MarkdownIt({
    html: true,
    linkify: true,
    typographer: true
})

// Custom image renderer to handle baseURL
md.renderer.rules.image = (tokens, idx, options, env, self) => {
    const token = tokens[idx]
    if (!token || !token.attrs) return self.renderToken(tokens, idx, options)

    const srcIndex = token.attrIndex('src')
    if (srcIndex >= 0) {
        const attr = token.attrs[srcIndex]
        if (attr) {
             const src = attr[1]
             // If path is absolute (starts with /), prepend baseURL
             if (src.startsWith('/') && !src.startsWith('//')) {
                 attr[1] = joinURL(baseURL, src)
             }
        }
    }
    return self.renderToken(tokens, idx, options)
}

// Fetch single article exclusively from the Go API.
const { data: article } = await useAsyncData(`blog-post-${id}`, async () => {
  const post = await $fetch<any>(`${config.public.backendApiUrl}/api/posts/${id}`)
  return {
    id: String(post.id),
    title: post.title,
    date: formatDate(post.created_at),
    featureTags: post.tags || [],
    image: post.image,
    content: post.content,
    adLink: post.ad_link,
    seo_title: post.seo_title,
    seo_keywords: post.seo_keywords,
    seo_description: post.seo_description,
    categories: post.article_tags || []
  }
})

const { data: relatedCategoryData } = await useAsyncData(`blog-related-categories-${id}`, async () => {
  const categories = article.value?.categories || []
  return Promise.all(categories.map(async (category: any) => {
    const result = await $fetch<any>(`${config.public.backendApiUrl}/api/posts`, {
      query: { article_tag_id: category.id, exclude_id: id, page: 1, limit: 1 }
    })
    return { id: category.id, name: category.name, post: result.data?.[0] || null }
  }))
})
const relatedCategories = computed(() => relatedCategoryData.value || [])

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return dateStr.split('T')[0]
}

const articleBodyRoot = ref<HTMLElement | null>(null)

const parsedArticle = computed(() => {
    if (!article.value) return buildArticleToc('')
    const content = article.value.content || ''
    const isHtml = /<(?:p|h[1-6]|ul|ol|table)\b/i.test(content)
    
    let rendered = isHtml ? content : md.render(content)
    
    // Replace naked image links converted to anchors by linkify
    rendered = rendered.replace(/<a href="(https:\/\/(?:i\.imgur\.com|i\.meee\.com\.tw)\/[^"]+)">[^<]+<\/a>/g, (match, url) => {
        return `<img src="${url}" style="width: 100%; max-width: 100%; height: auto; display: block; margin: 20px 0;" />`
    })
    
    // Replace remaining naked image links not in attributes
    rendered = rendered.replace(/(?<!["'(\/])(https:\/\/(?:i\.imgur\.com|i\.meee\.com\.tw)\/[a-zA-Z0-9._-]+)(?!["'])/g, (match) => {
        return `<img src="${match}" style="width: 100%; max-width: 100%; height: auto; display: block; margin: 20px 0;" />`
    })
    
    return buildArticleToc(rendered)
})

const parsedAdLink = computed(() => {
    if (!article.value) return ''
    const ad = article.value.adLink || ''
    if (ad.includes('<a') || ad.includes('<p') || ad.includes('<span')) {
        return ad
    }
    return md.render(ad)
})

useSeoMeta({
  title: computed(() => article.value?.seo_title || article.value?.title || '精選專欄'),
  keywords: computed(() => article.value?.seo_keywords || ''),
  description: computed(() => article.value?.seo_description || '')
})
</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

.breadcrumbs { margin-bottom: 30px; color: #7f8c8d; font-size: 14px; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }
.breadcrumbs .active { color: #95A5A6; }

.article-layout { display:grid; grid-template-columns:180px minmax(0, 1fr) 220px; gap:14px; align-items:start; }
.article-layout.without-toc { grid-template-columns:minmax(0, 1fr) 220px; }
.article-detail { min-width:0; background: white; padding: 40px; border-radius: 12px; box-shadow: 0 5px 20px rgba(0,0,0,0.05); }
.toc-rail { position:sticky; top:24px; min-width:0; max-height:calc(100vh - 48px); overflow-y:auto; }
.mobile-inline-toc { display:none; }

.article-header { margin-bottom: 30px; text-align: left; }
.meta-top { margin-bottom: 15px; display: flex; justify-content: flex-end; align-items: center; }
.date { color: #95A5A6; font-size: 14px; }
.feature-tags { display:flex; flex-wrap:wrap; gap:7px; margin-bottom:15px; }
.feature-tags span { border-radius:999px; background:#f1f5f9; padding:5px 10px; color:#64748b; font-size:12px; font-weight:700; }
.main-title { font-size: 32px; color: #2C3E50; line-height: 1.4; margin: 0 0 15px 0; }

.featured-image { margin-bottom: 40px; border-radius: 8px; overflow: hidden; }
.featured-image img { width: 100%; height: auto; display: block; }

.article-body { font-size: 18px; line-height: 1.8; color: #2c3e50; }
.article-body :deep(h2), .article-body :deep(h3) { margin-top: 40px; margin-bottom: 20px; color: #2C3E50; font-weight: 700; }
.article-body :deep(p) { margin-bottom: 20px; }
.article-body :deep(img) { max-width: 100%; height: auto; border-radius: 8px; margin: 30px 0; display: block; }

.article-footer { margin-top: 50px; padding-top: 30px; border-top: 1px solid #eee; text-align: center; }
.btn-back { display: inline-block; padding: 10px 25px; border: 1px solid #ddd; border-radius: 30px; color: #7f8c8d; transition: all 0.2s; }
.btn-back:hover { border-color: #2C3E50; color: #2C3E50; background: transparent; }

.not-found { text-align: center; padding: 50px; }
.category-sidebar { position:sticky; top:24px; align-self:start; border-radius:12px; background:#fff; padding:18px; box-shadow:0 5px 20px rgba(0,0,0,.05); }
.category-sidebar h2 { margin:0 0 18px; color:#2c3e50; font-size:20px; }
.category-group { padding:16px 0; border-top:1px solid #eef2f6; }
.category-heading { display:flex; align-items:center; justify-content:space-between; gap:10px; margin-bottom:10px; }
.category-name { display:inline-block; border-radius:999px; background:#e74c3c; padding:5px 11px; color:#fff; font-size:12px; font-weight:700; }
.category-more { flex-shrink:0; color:#e74c3c; font-size:12px; font-weight:700; text-decoration:none; }
.category-more:hover { text-decoration:underline; text-underline-offset:3px; }
.related-link { display:flex; flex-direction:column; gap:7px; color:#2c3e50; line-height:1.5; text-decoration:none; }
.related-link:hover strong { color:#e74c3c; }
.related-link time,.no-related { color:#94a3b8; font-size:12px; }

@media(max-width: 1100px) {
    .article-layout { grid-template-columns:1fr; }
    .article-layout.without-toc { grid-template-columns:1fr; }
    .toc-rail { display:none; }
    .mobile-inline-toc { display:block; }
    .category-sidebar { position:static; max-height:none; overflow:visible; }
    .article-detail { padding: 25px; }
    .main-title { font-size: 24px; }
}
</style>
