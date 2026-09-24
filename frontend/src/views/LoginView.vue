<template>
  <div class="auth-wrap">
    <div class="auth-card">
      <div class="auth-emoji">📮</div>
      <h1 class="auth-title">{{ mode === 'login' ? '欢迎回来' : '加入记忆连接' }}</h1>
      <p class="auth-sub">
        {{ mode === 'login' ? '登录后，让远方的人替你再看一眼' : '注册即赠送 100 枚记忆硬币，开启你的求看之旅' }}
      </p>

      <form @submit.prevent="submit">
        <div class="field">
          <label>用户名</label>
          <input v-model.trim="form.username" class="input" placeholder="至少 3 个字符" />
        </div>
        <div class="field">
          <label>密码</label>
          <input v-model="form.password" type="password" class="input" placeholder="至少 6 位" />
        </div>
        <div v-if="mode === 'register'" class="field">
          <label>昵称（选填）</label>
          <input v-model.trim="form.nickname" class="input" placeholder="别人怎么称呼你" />
        </div>

        <p v-if="error" class="error">{{ error }}</p>

        <button class="btn block" :disabled="loading">
          {{ loading ? '请稍候…' : mode === 'login' ? '登 录' : '注册并领取硬币' }}
        </button>
      </form>

      <div class="switch-mode">
        {{ mode === 'login' ? '还没有账号？' : '已经有账号了？' }}
        <a @click.prevent="toggle" href="#">
          {{ mode === 'login' ? '立即注册' : '去登录' }}
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { authApi } from '../api'
import { auth } from '../store/auth'

const mode = ref('login')
const loading = ref(false)
const error = ref('')
const form = reactive({ username: '', password: '', nickname: '' })
const router = useRouter()
const route = useRoute()

function toggle() {
  mode.value = mode.value === 'login' ? 'register' : 'login'
  error.value = ''
}

async function submit() {
  error.value = ''
  if (form.username.length < 3 || form.password.length < 6) {
    error.value = '用户名至少3位、密码至少6位'
    return
  }
  loading.value = true
  try {
    const payload =
      mode.value === 'login'
        ? { username: form.username, password: form.password }
        : { username: form.username, password: form.password, nickname: form.nickname }
    const res = await (mode.value === 'login' ? authApi.login(payload) : authApi.register(payload))
    auth.setSession(res.token, res.user)
    router.push(route.query.redirect || '/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-wrap {
  display: flex;
  justify-content: center;
  padding: 30px 0;
}
.auth-card {
  width: 100%;
  max-width: 400px;
  background: var(--paper-2);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 34px 30px;
  text-align: center;
}
.auth-emoji {
  font-size: 44px;
}
.auth-title {
  margin: 10px 0 6px;
  font-size: 23px;
}
.auth-sub {
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 24px;
}
form {
  text-align: left;
}
.error {
  color: #b4452f;
  font-size: 13px;
  margin-bottom: 12px;
}
.switch-mode {
  margin-top: 18px;
  font-size: 13px;
  color: var(--ink-soft);
}
.switch-mode a {
  color: var(--amber-dark);
  font-weight: 700;
  margin-left: 4px;
}
</style>
