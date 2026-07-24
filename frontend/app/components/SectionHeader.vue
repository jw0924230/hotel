<template>
  <div class="sec-header">
    <slot />
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.sec-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  border-left: 5px solid #E74C3C;
  padding-left: 15px;
}
</style>
