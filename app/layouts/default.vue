<template>
  <div class="layout-container">
    <header class="site-header">
      <!-- Top Bar (Desktop Only) -->
      <div class="top-bar d-none d-md-block">
        <div class="container container-wide">
          <div class="top-links">
             <a href="#">設為首頁</a> | <a href="#">加入最愛</a> | <a href="#">聯絡我們</a>
          </div>
        </div>
      </div>

      <!-- Main Header -->
      <div class="main-header">
        <div class="container container-wide header-inner">
          <div class="logo-area">
             <NuxtLink to="/" class="logo-text">
               <span class="qk">QK</span><span class="text">休閒網</span>
               <div class="slogan">全國最大的汽車旅館入口網站</div>
             </NuxtLink>
          </div>
          <!-- Header Banner (Desktop) -->
          <div class="header-banner d-none d-md-block">
             <img src="https://placehold.co/468x60/dcdcdc/333333?text=全台毛小孩住宿懶人包" alt="Banner">
          </div>
          <!-- Mobile Menu Toggle -->
          <button class="menu-toggle d-md-none" @click="isMenuOpen = !isMenuOpen">
            ☰
          </button>
        </div>
      </div>

      <!-- Navbar (Blue Strip) -->
      <nav class="main-nav d-none d-md-block">
        <div class="container container-wide">
          <ul class="nav-links">
            <li><a href="#" class="active">首頁</a></li>
            <li><a href="#">優惠訂房</a></li>
            <li><a href="#">空房查詢</a></li>
            <li><a href="#">地圖搜尋</a></li>
            <li><a href="#">討論版</a></li>
            <li><a href="#">評鑑版</a></li>
            <li><a href="#">最新消息</a></li>
            <li><a href="#">業者登入</a></li>
          </ul>
        </div>
      </nav>

      <!-- Mobile Nav Drawer -->
       <div class="mobile-nav" :class="{ open: isMenuOpen }">
          <a href="#" @click="isMenuOpen = false">首頁</a>
          <a href="#" @click="isMenuOpen = false">優惠訂房</a>
          <a href="#" @click="isMenuOpen = false">找旅館</a>
          <a href="#" @click="isMenuOpen = false">最新消息</a>
          <a href="#" @click="isMenuOpen = false" class="close-btn">關閉選單</a>
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
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const isMenuOpen = ref(false)
</script>

<style scoped>
/* Reset & Base */
a { text-decoration: none; color: inherit; }
ul { list-style: none; padding: 0; margin: 0; }
.layout-container { font-family: Arial, "Helvetica Neue", Helvetica, sans-serif; background: #fff; }
.container { width: 100%; padding: 0 10px; margin: 0 auto; }
.container-wide { max-width: 1000px; } /* Typical old web width */

/* Utilities */
.d-none { display: none !important; }
.d-md-block { display: block !important; }
@media (min-width: 768px) {
  .d-none { display: block !important; } /* Reset for desktop override if needed, but using specific classes preferably */
  .d-md-block { display: block !important; }
  .d-md-none { display: none !important; }
}

/* Top Bar */
.top-bar { background: #f5f5f5; border-bottom: 1px solid #ddd; padding: 5px 0; font-size: 12px; color: #666; }
.top-links { text-align: right; }
.top-links a:hover { text-decoration: underline; color: #CC0000; }

/* Main Header */
.main-header { padding: 15px 0; background: white; }
.header-inner { display: flex; justify-content: space-between; align-items: center; }
.logo-text { display: block; line-height: 1; }
.qk { font-size: 40px; font-weight: bold; color: #000080; }
.text { font-size: 24px; color: #000080; font-weight: bold; }
.slogan { font-size: 12px; color: #666; letter-spacing: 1px; margin-top: 5px; }

/* Main Nav */
.main-nav { background: #E6E6FA; border-top: 1px solid #000080; border-bottom: 1px solid #000080; }
.nav-links { display: flex; }
.nav-links li { flex: 1; text-align: center; }
.nav-links a {
  display: block; padding: 10px 0; font-weight: bold; color: #000080; font-size: 14px;
  border-right: 1px solid #ccc;
}
.nav-links a:hover, .nav-links a.active { background: #FFFACD; color: #CC0000; }

/* Mobile Nav */
.menu-toggle { font-size: 24px; background: none; border: none; cursor: pointer; }
.mobile-nav {
  position: fixed; top: 0; right: -250px; width: 250px; height: 100vh;
  background: white; z-index: 1001; box-shadow: -2px 0 5px rgba(0,0,0,0.1);
  transition: right 0.3s ease; display: flex; flex-direction: column;
}
.mobile-nav.open { right: 0; }
.mobile-nav a { padding: 15px; border-bottom: 1px solid #eee; font-weight: bold; }
.close-btn { color: red; text-align: center; margin-top: auto; }

/* Footer */
.site-footer { background: #336699; color: white; padding: 20px 0; margin-top: 30px; font-size: 12px; text-align: center; }
.footer-links { margin-bottom: 10px; }
.footer-links a { color: #fff; margin: 0 5px; }
.footer-links a:hover { text-decoration: underline; }
</style>
