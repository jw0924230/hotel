<template>
  <noscript><iframe src="https://www.googletagmanager.com/ns.html?id=GTM-NTR8MG73"
height="0" width="0" style="display:none;visibility:hidden"></iframe></noscript>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
const route = useRoute()
const config = useRuntimeConfig()

useHead({
  link: [
    {
      key: 'canonical',
      rel: 'canonical',
      href: computed(() => `${config.public.siteUrl}${route.path}`),
    },
  ],
})
</script>
