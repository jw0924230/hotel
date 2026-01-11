<template>
  <div class="portal-page">
    <div class="container container-wide">
      <div class="sec-header">
        <h1 class="page-title">
           <span class="icon">📰</span> 精選專欄
        </h1>
        <div class="breadcrumbs">
           <NuxtLink to="/">首頁</NuxtLink> &gt; <span class="active">部落格</span>
        </div>
      </div>

      <div class="blog-grid">
        <article v-for="article in articles" :key="article.id" class="article-card">
           <NuxtLink :to="`/blog/${article.id}`" class="card-link">
             <div class="img-wrapper">
               <img :src="article.image" :alt="article.title" loading="lazy">
               <span class="category-tag">{{ article.category }}</span>
             </div>
             <div class="content">
               <div class="meta">{{ article.date }}</div>
               <h2 class="title">{{ article.title }}</h2>
               <p class="excerpt">{{ getExcerpt(article.content) }}...</p>
               <span class="read-more">閱讀更多 &rarr;</span>
             </div>
           </NuxtLink>
        </article>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import articles from '~/../data_json/blog/articles.json'

// Helper to strip HTML and get excerpt
const getExcerpt = (html: string) => {
    const tmp = html.replace(/<[^>]+>/g, '')
    return tmp.substring(0, 60)
}
</script>

<style scoped>
.portal-page { padding: 40px 0; background: #f8f9fa; min-height: 80vh; }
.container { padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

.sec-header { margin-bottom: 40px; border-bottom: 1px solid #ddd; padding-bottom: 15px; display: flex; justify-content: space-between; align-items: flex-end; }
.page-title { margin: 0; font-size: 28px; color: #2C3E50; font-weight: 700; }
.breadcrumbs { color: #7f8c8d; font-size: 14px; }
.breadcrumbs a { color: #2C3E50; text-decoration: none; }

.blog-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 30px; }

.article-card { background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 5px 15px rgba(0,0,0,0.05); transition: all 0.3s; }
.article-card:hover { transform: translateY(-5px); box-shadow: 0 10px 25px rgba(0,0,0,0.1); }
.card-link { text-decoration: none; color: inherit; display: flex; flex-direction: column; height: 100%; }

.img-wrapper { position: relative; height: 200px; overflow: hidden; }
.img-wrapper img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.5s; }
.article-card:hover .img-wrapper img { transform: scale(1.05); }

.category-tag { position: absolute; top: 15px; left: 15px; background: rgba(44, 62, 80, 0.8); color: white; padding: 5px 10px; border-radius: 20px; font-size: 12px; font-weight: 600; }

.content { padding: 25px; flex: 1; display: flex; flex-direction: column; }
.meta { font-size: 13px; color: #95A5A6; margin-bottom: 10px; }
.title { font-size: 20px; font-weight: 700; color: #2C3E50; margin: 0 0 15px 0; line-height: 1.4; }
.excerpt { font-size: 15px; color: #7F8C8D; line-height: 1.6; margin-bottom: 20px; flex: 1; }
.read-more { color: #E74C3C; font-weight: 600; font-size: 14px; }
</style>
