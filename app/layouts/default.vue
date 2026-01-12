<template>
  <div class="layout-container">
    <header class="site-header">
      <!-- Main Header -->
      <div class="main-header">
        <div class="container container-wide header-inner">
          <div class="logo-area">
             <NuxtLink to="/" class="logo-text">
               <span class="qk">QK</span><span class="text">休閒網</span>
               <div class="slogan">全國最大的汽車旅館入口網站</div>
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

    <footer class="site-footer">
      <div class="container">
        <div class="footer-links">
          <a href="#">關於我們</a> | <a href="#">隱私權聲明</a> | <a href="#">免責條款</a> | <a href="#">聯絡我們</a>
        </div>
        <p>Copyright © 2026 QK休閒網. All rights reserved.</p>
      </div>
    </footer>
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
.logo-text { display: block; line-height: 1; text-decoration: none; }
.qk { font-size: 48px; font-weight: 800; color: #E74C3C; letter-spacing: -2px; }
.text { font-size: 28px; color: #2C3E50; font-weight: 700; margin-left: 5px; }
.slogan { font-size: 13px; color: #7F8C8D; letter-spacing: 2px; margin-top: 5px; text-transform: uppercase; }

.main-nav { display: flex; align-items: center; }
.nav-link { font-size: 16px; font-weight: 700; color: #2C3E50; margin-left: 30px; padding: 5px 0; border-bottom: 2px solid transparent; transition: all 0.2s; }
.nav-link:hover { color: #E74C3C; border-bottom-color: #E74C3C; }
.nav-link.router-link-active { color: #E74C3C; }

.content-wrapper { flex: 1; }

/* Footer */
.site-footer { background: #2C3E50; color: #ECF0F1; padding: 40px 0; margin-top: 60px; font-size: 14px; text-align: center; }
.footer-links { margin-bottom: 20px; }
.footer-links a { color: #BDC3C7; margin: 0 10px; transition: color 0.3s; }
.footer-links a:hover { color: white; text-decoration: none; }

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
