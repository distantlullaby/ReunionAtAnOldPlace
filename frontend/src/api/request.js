import axios from 'axios'
import { auth } from '../store/auth'
import router from '../router'

const request = axios.create({ baseURL: '/api', timeout: 10000 })

// 请求自动带上 token。
request.interceptors.request.use((config) => {
  if (auth.token) config.headers.Authorization = `Bearer ${auth.token}`
  return config
})

// 统一解包 { code, msg, data }，401 时清登录态并跳转。
request.interceptors.response.use(
  (resp) => {
    const body = resp.data
    if (body && typeof body === 'object' && 'code' in body) {
      if (body.code === 0) return body.data
      return Promise.reject(new Error(body.msg || '请求失败'))
    }
    return body
  },
  (error) => {
    const status = error.response?.status
    const msg = error.response?.data?.msg || error.message || '网络错误'
    if (status === 401) {
      auth.logout()
      if (router.currentRoute.value.name !== 'login') {
        router.push({ name: 'login' })
      }
    }
    return Promise.reject(new Error(msg))
  }
)

export default request
