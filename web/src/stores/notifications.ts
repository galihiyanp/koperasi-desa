import { defineStore } from 'pinia'

export type NotificationType = 'info' | 'success' | 'warning' | 'error'
export type NotificationItem = {
  id: number
  type: NotificationType
  title?: string
  message: string
  detail?: string
  timestamp: number
  timeout?: number
}

let counter = 1

export const useNotificationsStore = defineStore('notifications', {
  state: () => ({
    list: [] as NotificationItem[],
  }),
  actions: {
    notify(n: Omit<NotificationItem, 'id' | 'timestamp'>) {
      const id = counter++
      const item: NotificationItem = {
        id,
        timestamp: Date.now(),
        ...n,
      }
      this.list.unshift(item)
      const to = item.timeout ?? 4000
      if (to > 0) {
        setTimeout(() => { this.remove(id) }, to)
      }
      return id
    },
    remove(id: number) {
      this.list = this.list.filter(x => x.id !== id)
    },
    clear() {
      this.list = []
    },
  },
})