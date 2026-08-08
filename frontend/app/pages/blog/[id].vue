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
                <h1 class="main-title">{{ article.title }}</h1>
                <div class="meta-top">
                    <span class="category">{{ article.category }}</span>
                    <span class="date">{{ article.date }}</span>
                </div>
            </header>

            <ArticleTableOfContents
              :items="parsedArticle.items"
              :content-root="articleBodyRoot"
            />
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
    category: post.tags && post.tags.length > 0 ? post.tags[0] : '精選專欄',
    image: post.image,
    content: post.content,
    adLink: post.ad_link,
    seo_title: post.seo_title,
    seo_keywords: post.seo_keywords,
    seo_description: post.seo_description
  }
})

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
.container-narrow { max-width: 800px; }

.breadcrumbs { margin-bottom: 30px; color: #7f8c8d; font-size: 14px; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }
.breadcrumbs .active { color: #95A5A6; }

.article-detail { background: white; padding: 40px; border-radius: 12px; box-shadow: 0 5px 20px rgba(0,0,0,0.05); }

.article-header { margin-bottom: 30px; text-align: left; }
.meta-top { margin-bottom: 15px; display: flex; justify-content: space-between; align-items: center; }
.category { background: #E74C3C; color: white; padding: 4px 12px; border-radius: 20px; font-size: 13px; font-weight: 600; }
.date { color: #95A5A6; font-size: 14px; }
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

@media(max-width: 768px) {
    .article-detail { padding: 25px; }
    .main-title { font-size: 24px; }
}
</style>
