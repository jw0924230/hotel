<template>
  <div class="layout-container">
    <header class="site-header">
      <!-- Main Header -->
      <div class="main-header">
        <div class="container container-wide header-inner">
          <div class="logo-area">
             <NuxtLink to="/" class="logo-container">
               <img src="/logo.png" alt="休息3小時" class="logo-img" />
               <div class="logo-text-group">
                 <span class="qk">休息</span><span class="text">3小時</span>
               </div>
             </NuxtLink>
          </div>
          <nav class="main-nav">
             <NuxtLink to="/blog" class="nav-link">部落格</NuxtLink>
          </nav>
        </div>
      </div>
    </header>

    <div class="content-wrapper">
       <slot />
    </div>

    <AppFooter />
    <button 
      v-show="showGoTop" 
      class="go-top-btn" 
      @click="scrollToTop"
      aria-label="Go to Top"
    >
      ↑
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const showGoTop = ref(false)

const handleScroll = () => {
    showGoTop.value = window.scrollY > 100
}

const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
    window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
/* Reset & Base */
*, *::before, *::after { box-sizing: border-box; }
html, body { margin: 0; padding: 0; }
a { text-decoration: none; color: inherit; }
ul { list-style: none; padding: 0; margin: 0; }
.layout-container { font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; background: #f8f9fa; display: flex; flex-direction: column; min-height: 100vh; overflow-x: hidden; }
.container { width: 100%; padding: 0 15px; margin: 0 auto; }
.container-wide { max-width: 1200px; }

/* Utilities */
.d-none { display: none !important; }
.d-md-block { display: block !important; }
@media (min-width: 768px) {
  .d-none { display: block !important; }
  .d-md-block { display: block !important; }
  .d-md-none { display: none !important; }
}

/* Main Header */
.main-header { padding: 20px 0; background: white; box-shadow: 0 2px 10px rgba(0,0,0,0.05); }
.header-inner { display: flex; justify-content: space-between; align-items: center; }
.logo-container { display: flex; align-items: center; text-decoration: none; }
.logo-img { height: 50px; width: 50px; object-fit: contain; margin-right: 10px; }
.logo-text-group { display: flex; align-items: flex-end; line-height: 1; }
.qk { font-size: 40px; font-weight: 800; color: #E74C3C; letter-spacing: -2px; line-height: 1; }
.text { font-size: 28px; color: #2C3E50; font-weight: 700; margin-left: 5px; line-height: 1; margin-bottom: 3px; }
.slogan { font-size: 13px; color: #7F8C8D; letter-spacing: 2px; margin-top: 5px; text-transform: uppercase; }

.main-nav { display: flex; align-items: center; }
.nav-link { font-size: 16px; font-weight: 700; color: #2C3E50; margin-left: 30px; padding: 5px 0; border-bottom: 2px solid transparent; transition: all 0.2s; }
.nav-link:hover { color: #E74C3C; border-bottom-color: #E74C3C; }
.nav-link.router-link-active { color: #E74C3C; }

.content-wrapper { flex: 1; }



.go-top-btn {
    position: fixed;
    bottom: 30px;
    right: 30px;
    width: 50px;
    height: 50px;
    background-color: #E74C3C;
    color: white;
    border: none;
    border-radius: 50%;
    font-size: 24px;
    cursor: pointer;
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
    display: flex;
    justify-content: center;
    align-items: center;
    transition: all 0.3s;
    z-index: 1000;
}

.go-top-btn:hover {
    background-color: #c0392b;
    transform: translateY(-3px);
    box-shadow: 0 6px 15px rgba(0,0,0,0.3);
}
</style>
