<template><ArticleTagResultsPage :tag="tag" :page="page" :posts="result.data" :total="result.total" :all-tags="allTags || []" /></template>
<script setup lang="ts">
import { computed } from 'vue';
definePageMeta({ key: route => route.fullPath });
const route=useRoute(); const config=useRuntimeConfig();
const id=String(route.params.id); const page=Math.max(1,Number(route.params.page)||1);
const [{data:tagData},{data:allTags},{data:postsData}]=await Promise.all([
  useAsyncData(`article-tag-${id}`,()=> $fetch<any>(`${config.public.backendApiUrl}/api/article-tags/${id}`)),
  useAsyncData('article-tags-public',()=> $fetch<any[]>(`${config.public.backendApiUrl}/api/article-tags`)),
  useAsyncData(`article-tag-posts-${id}-${page}`,()=> $fetch<any>(`${config.public.backendApiUrl}/api/posts`,{query:{article_tag_id:id,page,limit:12}})),
]);
const tag=computed(()=>tagData.value||{id:Number(id),name:'文章標籤'}); const result=computed(()=>postsData.value||{data:[],total:0});
useSeoMeta({title:computed(()=>`${tag.value.name}文章推薦｜住宿、休息與旅遊攻略`),description:computed(()=>`整理「${tag.value.name}」相關的住宿、休息與旅遊文章，提供飯店選擇、行程規劃與實用攻略，快速找到符合需求的資訊。`)});
</script>
