import axios from 'axios'

const api = axios.create({ baseURL: '/api' })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.dispatchEvent(new Event('auth-changed'))
    }
    return Promise.reject(err)
  }
)

export function getUser() {
  try { return JSON.parse(localStorage.getItem('user')) } catch { return null }
}

export function setAuth(token, user) {
  localStorage.setItem('token', token)
  localStorage.setItem('user', JSON.stringify(user))
  window.dispatchEvent(new Event('auth-changed'))
}

export function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  window.dispatchEvent(new Event('auth-changed'))
}

export async function refreshMe() {
  const res = await api.get('/me')
  localStorage.setItem('user', JSON.stringify(res.data.user))
  window.dispatchEvent(new Event('auth-changed'))
  return res.data.user
}

export default api
