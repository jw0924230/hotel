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
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
</script>

<style scoped>
/* Reset & Base */
a { text-decoration: none; color: inherit; }
ul { list-style: none; padding: 0; margin: 0; }
.layout-container { font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; background: #f8f9fa; display: flex; flex-direction: column; min-height: 100vh; }
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

.content-wrapper { flex: 1; }

/* Footer */
.site-footer { background: #2C3E50; color: #ECF0F1; padding: 40px 0; margin-top: 60px; font-size: 14px; text-align: center; }
.footer-links { margin-bottom: 20px; }
.footer-links a { color: #BDC3C7; margin: 0 10px; transition: color 0.3s; }
.footer-links a:hover { color: white; text-decoration: none; }
</style>
