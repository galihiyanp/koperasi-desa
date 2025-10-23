<template>
  <div class="notifications" aria-live="polite" aria-atomic="true">
    <transition-group name="slide-down" tag="div" class="stack">
      <div
        v-for="n in store.list"
        :key="n.id"
        class="notice"
        :class="n.type"
        role="alert"
      >
        <div class="content">
          <strong v-if="n.title" class="title">{{ n.title }}</strong>
          <div class="message">{{ n.message }}</div>
          <div v-if="n.detail" class="detail">{{ n.detail }}</div>
        </div>
        <button class="close" @click="store.remove(n.id)" aria-label="Tutup">×</button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useNotificationsStore } from '@/stores/notifications'
const store = useNotificationsStore()
</script>

<style scoped>
.notifications {
  position: fixed;
  top: 12px;
  right: 12px;
  z-index: 2000;
  pointer-events: none;
}
.stack { display: flex; flex-direction: column; gap: 8px; }
.notice {
  pointer-events: auto;
  min-width: 280px;
  max-width: 360px;
  border-radius: 10px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-md);
  display: flex;
  align-items: start;
  justify-content: space-between;
  gap: 8px;
  padding: 10px 12px;
  background: var(--color-surface);
}
.notice.info { border-color: #3b82f6; background: #eff6ff; }
.notice.success { border-color: #10b981; background: #ecfdf5; }
.notice.warning { border-color: #f59e0b; background: #fffbeb; }
.notice.error { border-color: #ef4444; background: #fef2f2; }
.title { font-weight: 600; }
.message { font-size: 14px; }
.detail { font-size: 12px; color: var(--color-text-weak); margin-top: 2px; }
.close { border: none; background: transparent; cursor: pointer; font-size: 18px; line-height: 1; }

.slide-down-enter-active, .slide-down-leave-active { transition: all .25s ease; }
.slide-down-enter-from { opacity: 0; transform: translateY(-8px); }
.slide-down-leave-to { opacity: 0; transform: translateY(-8px); }
</style>