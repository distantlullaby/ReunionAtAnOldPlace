<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:18px">
      <div>
        <h2 style="font-size:20px">求看广场</h2>
        <p style="font-size:13px;color:var(--muted)">每一条需求，都是一段想被再看一眼的回忆</p>
      </div>
      <button class="btn btn-primary" @click="openCreate">✏️ 发布求看需求</button>
    </div>

    <div v-if="loading" class="empty">加载中…</div>
    <div v-else-if="!stories.length" class="empty">还没有人发布求看需求，来发第一条吧</div>
    <StoryCard v-for="s in stories" :key="s.id" :story="s" :me="me" @changed="load" />

    <!-- 发布需求弹窗 -->
    <div v-if="showCreate" class="modal-mask" @click.self="showCreate = false">
      <div class="modal">
        <h3>✏️ 发布求看需求</h3>
        <div class="field">
          <label>标题</label>
          <input v-model.trim="form.title" placeholder="例如：想看看老家巷口的那棵槐树" />
        </div>
        <div class="field">
          <label>回忆地点</label>
          <input v-model.trim="form.location" placeholder="越具体越好，例如：XX市XX区XX巷口" />
        </div>
        <div class="field">
          <label>过去回忆</label>
          <textarea v-model.trim="form.memory_text" placeholder="写下你和这个地方的故事…"></textarea>
        </div>
        <div class="field">
          <label>老照片（可选）</label>
          <PhotoUpload v-model="form.old_photo" />
        </div>
        <div class="field">
          <label>悬赏记忆硬币（发布即冻结，确认后结算给拍摄者）</label>
          <input v-model.number="form.bounty" type="number" min="1" placeholder="例如 30" />
        </div>
        <p v-if="error" class="error-text">{{ error }}</p>
        <div class="modal-actions">
          <button class="btn btn-ghost" @click="showCreate = false">取消</button>
          <button class="btn btn-primary" :disabled="submitting" @click="create">
            {{ submitting ? '发布中…' : '冻结硬币并发布' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import api, { getUser, refreshMe } from '../api'
import StoryCard from '../components/StoryCard.vue'
import PhotoUpload from '../components/PhotoUpload.vue'

const stories = ref([])
const me = ref(getUser())
const loading = ref(true)
const showCreate = ref(false)
const submitting = ref(false)
const error = ref('')
const form = reactive({ title: '', location: '', memory_text: '', old_photo: '', bounty: null })

function syncUser() { me.value = getUser() }

async function load() {
  loading.value = true
  try {
    const res = await api.get('/stories')
    stories.value = res.data.stories || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  if (!me.value) { alert('请先登录 / 注册'); return }
  showCreate.value = true
}

async function create() {
  error.value = ''
  if (!form.title || !form.location || !form.memory_text) { error.value = '请填写标题、地点和回忆'; return }
  if (!form.bounty || form.bounty < 1) { error.value = '悬赏至少 1 枚硬币'; return }
  submitting.value = true
  try {
    await api.post('/stories', form)
    await refreshMe().catch(() => {})
    showCreate.value = false
    Object.assign(form, { title: '', location: '', memory_text: '', old_photo: '', bounty: null })
    await load()
  } catch (e) {
    error.value = e.response?.data?.error || '发布失败，请稍后再试'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  load()
  if (me.value) refreshMe().catch(() => {})
  window.addEventListener('auth-changed', syncUser)
})
onUnmounted(() => window.removeEventListener('auth-changed', syncUser))
</script>
