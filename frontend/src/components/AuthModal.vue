<template>
  <div class="modal-mask" @click.self="$emit('close')">
    <div class="modal">
      <h3>{{ isLogin ? '登录' : '注册' }}</h3>
      <div class="field">
        <label>用户名</label>
        <input v-model.trim="form.username" placeholder="请输入用户名" @keyup.enter="submit" />
      </div>
      <div class="field" v-if="!isLogin">
        <label>昵称（可选）</label>
        <input v-model.trim="form.nickname" placeholder="怎么称呼你" />
      </div>
      <div class="field">
        <label>密码</label>
        <input v-model="form.password" type="password" placeholder="请输入密码" @keyup.enter="submit" />
      </div>
      <p v-if="!isLogin" style="font-size:13px;color:var(--muted)">🪙 注册即赠送 100 枚记忆硬币</p>
      <p v-if="error" class="error-text">{{ error }}</p>
      <div class="modal-actions">
        <button class="btn btn-ghost" @click="isLogin = !isLogin; error = ''">
          {{ isLogin ? '没有账号？去注册' : '已有账号？去登录' }}
        </button>
        <button class="btn btn-primary" :disabled="loading" @click="submit">
          {{ loading ? '请稍候…' : (isLogin ? '登录' : '注册') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import api, { setAuth } from '../api'

const emit = defineEmits(['close'])
const isLogin = ref(true)
const loading = ref(false)
const error = ref('')
const form = reactive({ username: '', password: '', nickname: '' })

async function submit() {
  error.value = ''
  if (!form.username || !form.password) { error.value = '请填写用户名和密码'; return }
  loading.value = true
  try {
    const url = isLogin.value ? '/login' : '/register'
    const res = await api.post(url, form)
    setAuth(res.data.token, res.data.user)
    emit('close')
  } catch (e) {
    error.value = e.response?.data?.error || '网络错误，请稍后再试'
  } finally {
    loading.value = false
  }
}
</script>
