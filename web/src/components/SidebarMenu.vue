<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { menuItems } from '@/config/menu'

const route = useRoute()
const isActive = (path?: string) => computed(() => path && route.path.startsWith(path))
</script>

<template>
  <nav class="menu">
    <RouterLink
      v-for="item in menuItems"
      :key="item.label"
      :to="item.path || '/dashboard'"
      class="menu-item"
      :class="{ active: isActive(item.path).value }"
    >
      <span>{{ item.label }}</span>
    </RouterLink>
  </nav>
</template>