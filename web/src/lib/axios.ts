import axios from 'axios'

// Dev: gunakan proxy Vite (base relatif '/')
// Prod: gunakan VITE_API_BASE_URL, fallback ke 8080 jika tidak diset
const isDev = import.meta.env.DEV
const baseURL = isDev ? '/' : (import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080')

const api = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default api