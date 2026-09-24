import { reactive, watch } from 'vue'

// 轻量全局登录态：token + 用户信息，持久化到 localStorage。
const stored = JSON.parse(localStorage.getItem('ml_auth') || 'null')

export const auth = reactive({
  token: stored?.token || '',
  user: stored?.user || null,
  get isLogin() {
    return !!this.token
  },
  setSession(token, user) {
    this.token = token
    this.user = user
  },
  logout() {
    this.token = ''
    this.user = null
  }
})

watch(
  () => [auth.token, auth.user],
  () => {
    if (auth.token) {
      localStorage.setItem('ml_auth', JSON.stringify({ token: auth.token, user: auth.user }))
    } else {
      localStorage.removeItem('ml_auth')
    }
  },
  { deep: true }
)
