<script setup lang="ts">
import SidebarMenu from '@/components/SidebarMenu.vue'
import { RouterView } from 'vue-router'

const appName = import.meta.env.VITE_APP_NAME || 'Sistem Koperasi Kantor Desa'
</script>

<template>
  <div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2 class="app-name">{{ appName }}</h2>
        <p class="subtitle">Admin</p>
      </div>
      <div class="sidebar-content">
        <SidebarMenu />
      </div>
    </aside>

    <!-- Main -->
    <div class="main">
      <!-- Topbar -->
      <header class="topbar">
        <div class="topbar-inner">
          <div class="topbar-title">{{ appName }}</div>
          <div class="flex items-center gap-3">
            <span class="muted">Logged in</span>
            <button class="btn-light">Logout</button>
          </div>
        </div>
      </header>

      <!-- Content -->
      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>