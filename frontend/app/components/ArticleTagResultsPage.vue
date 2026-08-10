<template>
  <div class="portal-page"><div class="container">
    <header class="header"><h1>{{ tag.name }}文章推薦</h1><p>住宿、休息與旅遊攻略</p></header>
    <nav class="tag-nav" aria-label="文章特色標籤">
      <NuxtLink v-for="item in allTags" :key="item.id" :to="`/blog/tag/${item.id}/1`" :class="{ active: Number(item.id) === Number(tag.id) }">{{ item.name }}</NuxtLink>
    </nav>
    <div v-if="posts.length" class="grid">
      <article v-for="post in posts" :key="post.id" class="card">
        <NuxtLink :to="`/blog/${post.id}`" class="main-link">
          <div class="image"><img v-if="post.image" :src="imageUrl(post.image)" :alt="post.title" loading="lazy"><span>{{ post.tags?.[0] || '精選專欄' }}</span></div>
          <div class="content"><time>{{ formatDate(post.created_at) }}</time><h2>{{ post.title }}</h2><p>{{ excerpt(post.content) }}</p></div>
        </NuxtLink>
        <div v-if="post.article_tags?.length" class="tags">
          <NuxtLink v-for="item in post.article_tags" :key="item.id" :to="`/blog/tag/${item.id}/1`">{{ item.name }}</NuxtLink>
        </div>
      </article>
    </div>
    <p v-else class="empty">此標籤目前沒有文章</p>
    <nav v-if="totalPages > 1" class="pagination"><NuxtLink v-if="page > 1" :to="`/blog/tag/${tag.id}/${page - 1}`">上一頁</NuxtLink><span>第 {{ page }} / {{ totalPages }} 頁</span><NuxtLink v-if="page < totalPages" :to="`/blog/tag/${tag.id}/${page + 1}`">下一頁</NuxtLink></nav>
  </div></div>
</template>
<script setup lang="ts">
import { computed } from 'vue'; import { joinURL } from 'ufo';
const props = defineProps<{ tag:any; page:number; posts:any[]; total:number; allTags:any[] }>();
const config=useRuntimeConfig(); const totalPages=computed(()=>Math.max(1,Math.ceil(props.total/12)));
const imageUrl=(value:string)=>value?.startsWith('/')?joinURL(config.app.baseURL,value):value;
const formatDate=(value:string)=>value ? value.split('T')[0] : '';
const excerpt=(html:string)=>(html||'').replace(/<[^>]+>/g,'').slice(0,90);
</script>
<style scoped>
.portal-page{padding:40px 0;background:#f8f9fa;min-height:80vh}.container{max-width:1200px;margin:auto;padding:0 15px}.header{border-bottom:1px solid #ddd;margin-bottom:20px}.header h1{color:#2c3e50;margin:0 0 6px}.header p{color:#7f8c8d}.tag-nav,.tags{display:flex;flex-wrap:wrap;gap:8px}.tag-nav{background:#eef1f4;padding:18px;border-radius:10px;margin-bottom:28px}.tag-nav a,.tags a{border-radius:999px;background:#fff;padding:7px 13px;color:#536273;text-decoration:none;font-weight:700;font-size:13px}.tag-nav a.active,.tag-nav a:hover,.tags a:hover{background:#e74c3c;color:#fff}.grid{display:grid;grid-template-columns:repeat(3,1fr);gap:24px}.card{background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 4px 14px #0000000d}.main-link{text-decoration:none;color:inherit}.image{height:190px;position:relative;background:#e2e8f0}.image img{width:100%;height:100%;object-fit:cover}.image span{position:absolute;left:14px;top:14px;background:#2c3e50dd;color:#fff;border-radius:999px;padding:5px 10px;font-size:12px}.content{padding:20px}.content time{color:#95a5a6;font-size:13px}.content h2{font-size:19px;color:#2c3e50;line-height:1.4}.content p{color:#7f8c8d;line-height:1.6}.tags{padding:0 20px 20px}.tags a{background:#f1f5f9;padding:4px 9px;font-size:11px}.pagination{display:flex;justify-content:center;gap:16px;margin-top:36px}.pagination a{color:#e74c3c}.empty{text-align:center;padding:50px;color:#7f8c8d}@media(max-width:800px){.grid{grid-template-columns:1fr}}
</style>
