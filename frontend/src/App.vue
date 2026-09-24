<template>
  <div>
    <header class="topbar">
      <div class="brand" @click="$router.push('/')">
        <span class="logo">🎞️</span>
        <div>
          <div class="brand-name">记忆连接</div>
          <div class="slogan">你帮我再看一眼，我把记忆还给你</div>
        </div>
      </div>
      <nav>
        <router-link to="/" class="nav-link">首页</router-link>
        <router-link to="/profile" class="nav-link">个人中心</router-link>
      </nav>
      <div class="auth-area">
        <template v-if="user">
          <span class="coin-badge">🪙 {{ user.balance }}<em v-if="user.frozen">（冻结 {{ user.frozen }}）</em></span>
          <span class="nickname">{{ user.nickname }}</span>
          <button class="btn btn-ghost" @click="logout">退出</button>
        </template>
        <button v-else class="btn btn-primary" @click="showAuth = true">登录 / 注册</button>
      </div>
    </header>

    <main class="container">
      <router-view :key="route.fullPath" />
    </main>

    <AuthModal v-if="showAuth" @close="showAuth = false" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { getUser, clearAuth } from './api'
import AuthModal from './components/AuthModal.vue'

const route = useRoute()
const user = ref(getUser())
const showAuth = ref(false)

function sync() { user.value = getUser() }
function logout() { clearAuth() }

onMounted(() => window.addEventListener('auth-changed', sync))
onUnmounted(() => window.removeEventListener('auth-changed', sync))
</script>
