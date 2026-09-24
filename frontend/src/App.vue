<template>
  <div class="app">
    <header class="nav">
      <div class="nav-inner">
        <router-link to="/" class="brand">
          <span class="brand-mark">📮</span>
          <span class="brand-name">记忆连接</span>
        </router-link>
        <div class="slogan">你帮我再看一眼，我把记忆还给你</div>
        <nav class="nav-links">
          <template v-if="auth.isLogin">
            <router-link to="/" class="nav-link">求看广场</router-link>
            <router-link to="/profile" class="nav-link coin-badge">
              <span class="coin-dot">●</span>
              {{ auth.user?.coinBalance ?? 0 }}
            </router-link>
            <button class="btn sm ghost" @click="logout">退出</button>
          </template>
          <template v-else>
            <router-link to="/login" class="btn sm">登录 / 注册</router-link>
          </template>
        </nav>
      </div>
    </header>

    <main class="container">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" @toast="showToast" />
        </transition>
      </router-view>
    </main>

    <footer class="footer">记忆连接 · 让远方的人替你再看一眼故乡</footer>

    <transition name="fade">
      <div v-if="toast" class="toast" :class="{ err: toastType === 'err' }">{{ toast }}</div>
    </transition>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { auth } from './store/auth'

const router = useRouter()
const toast = ref('')
const toastType = ref('ok')
let timer = null

function showToast(msg, type = 'ok') {
  toast.value = msg
  toastType.value = type
  clearTimeout(timer)
  timer = setTimeout(() => (toast.value = ''), 2400)
}

function logout() {
  auth.logout()
  router.push('/login')
  showToast('已退出登录')
}
</script>

<style scoped>
.nav {
  position: sticky;
  top: 0;
  z-index: 50;
  background: rgba(246, 239, 227, 0.88);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--line);
}
.nav-inner {
  max-width: 860px;
  margin: 0 auto;
  padding: 12px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}
.brand {
  display: flex;
  align-items: center;
  gap: 8px;
}
.brand-mark {
  font-size: 22px;
}
.brand-name {
  font-size: 19px;
  font-weight: 800;
  letter-spacing: 1px;
}
.slogan {
  flex: 1;
  text-align: center;
  font-size: 13px;
  color: var(--ink-soft);
  font-style: italic;
}
.nav-links {
  display: flex;
  align-items: center;
  gap: 14px;
}
.nav-link {
  font-size: 14px;
  font-weight: 600;
  color: var(--ink);
}
.nav-link:hover {
  color: var(--amber-dark);
}
.coin-badge {
  background: rgba(196, 120, 58, 0.14);
  padding: 5px 13px;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}
.coin-dot {
  color: var(--amber);
  font-size: 10px;
}
.container {
  max-width: 860px;
  margin: 0 auto;
  padding: 26px 20px 60px;
  min-height: 70vh;
}
.footer {
  text-align: center;
  color: var(--ink-soft);
  font-size: 12px;
  padding: 24px;
  border-top: 1px solid var(--line);
}
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.18s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
@media (max-width: 640px) {
  .slogan {
    display: none;
  }
}
</style>
