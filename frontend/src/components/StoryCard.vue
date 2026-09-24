<template>
  <div class="card">
    <div style="display:flex;justify-content:space-between;align-items:flex-start;gap:12px;flex-wrap:wrap">
      <div>
        <div style="font-size:16px;font-weight:700">{{ story.title }}</div>
        <div style="font-size:13px;color:var(--muted);margin-top:2px">
          📍 {{ story.location }} · {{ story.user?.nickname }} · {{ formatTime(story.created_at) }}
        </div>
      </div>
      <div style="display:flex;align-items:center;gap:8px">
        <span class="coin-badge">🪙 悬赏 {{ story.bounty }}</span>
        <span class="tag" :class="'tag-' + story.status">{{ statusText }}</span>
      </div>
    </div>

    <div style="margin-top:12px;display:flex;gap:10px;flex-wrap:wrap">
      <button class="btn btn-ghost btn-sm" @click="expanded = !expanded">
        {{ expanded ? '收起明信片' : '展开双面明信片 🎴' }}
      </button>
      <template v-if="story.status === 'open'">
        <button v-if="me && me.id !== story.user_id" class="btn btn-primary btn-sm" @click="showRespond = true">
          📷 替他去拍
        </button>
        <template v-if="me && me.id === story.user_id">
          <button class="btn btn-gold btn-sm" @click="showAppend = true">追加悬赏</button>
          <button class="btn btn-ghost btn-sm" @click="cancel">取消需求</button>
        </template>
      </template>
    </div>

    <Postcard v-if="expanded" :story="story" />

    <!-- 发起人查看待确认的回应 -->
    <div v-if="isMine && pendingResponses.length" style="margin-top:14px">
      <div style="font-size:14px;font-weight:700;margin-bottom:8px">收到的回应（{{ pendingResponses.length }}）</div>
      <div v-for="r in pendingResponses" :key="r.id"
           style="display:flex;gap:12px;align-items:center;border:1px dashed var(--line);border-radius:10px;padding:10px;margin-bottom:8px">
        <img :src="r.new_photo" style="width:72px;height:72px;object-fit:cover;border-radius:8px" />
        <div style="flex:1;font-size:13px">
          <div><b>{{ r.user?.nickname }}</b>：{{ r.message }}</div>
          <div style="color:var(--muted)">{{ formatTime(r.created_at) }}</div>
        </div>
        <button v-if="story.status === 'open'" class="btn btn-primary btn-sm" :disabled="acting"
                @click="accept(r)">确认采纳并结算 🪙</button>
      </div>
    </div>

    <p v-if="error" class="error-text">{{ error }}</p>

    <!-- 替他去拍弹窗 -->
    <div v-if="showRespond" class="modal-mask" @click.self="showRespond = false">
      <div class="modal">
        <h3>📷 替 TA 去拍 · {{ story.location }}</h3>
        <div class="field">
          <label>现场新照片</label>
          <PhotoUpload v-model="respondForm.new_photo" />
        </div>
        <div class="field">
          <label>寄语</label>
          <textarea v-model.trim="respondForm.message" placeholder="把眼前的样子，说给回忆里的人听…"></textarea>
        </div>
        <p v-if="respondError" class="error-text">{{ respondError }}</p>
        <div class="modal-actions">
          <button class="btn btn-ghost" @click="showRespond = false">取消</button>
          <button class="btn btn-primary" :disabled="acting" @click="respond">提交回应</button>
        </div>
      </div>
    </div>

    <!-- 追加悬赏弹窗 -->
    <div v-if="showAppend" class="modal-mask" @click.self="showAppend = false">
      <div class="modal">
        <h3>追加悬赏</h3>
        <div class="field">
          <label>追加硬币数量（当前悬赏 {{ story.bounty }} 枚）</label>
          <input v-model.number="appendAmount" type="number" min="1" placeholder="例如 20" />
        </div>
        <p v-if="appendError" class="error-text">{{ appendError }}</p>
        <div class="modal-actions">
          <button class="btn btn-ghost" @click="showAppend = false">取消</button>
          <button class="btn btn-gold" :disabled="acting" @click="append">确认追加</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import api, { refreshMe } from '../api'
import Postcard from './Postcard.vue'
import PhotoUpload from './PhotoUpload.vue'

const props = defineProps({
  story: { type: Object, required: true },
  me: { type: Object, default: null }
})
const emit = defineEmits(['changed'])

const expanded = ref(false)
const showRespond = ref(false)
const showAppend = ref(false)
const acting = ref(false)
const error = ref('')
const respondError = ref('')
const appendError = ref('')
const appendAmount = ref(null)
const respondForm = reactive({ new_photo: '', message: '' })

const isMine = computed(() => props.me && props.me.id === props.story.user_id)
const pendingResponses = computed(() => (props.story.responses || []).filter(r => r.status === 'pending'))
const statusText = computed(() => ({ open: '等待拍摄', settled: '已连接', cancelled: '已取消' }[props.story.status] || props.story.status))

function formatTime(t) {
  return t ? new Date(t).toLocaleString('zh-CN', { month: 'numeric', day: 'numeric', hour: '2-digit', minute: '2-digit' }) : ''
}

async function run(fn, errRef = error) {
  acting.value = true
  errRef.value = ''
  try {
    await fn()
    await refreshMe().catch(() => {})
    emit('changed')
  } catch (e) {
    errRef.value = e.response?.data?.error || '操作失败，请稍后再试'
  } finally {
    acting.value = false
  }
}

async function respond() {
  if (!respondForm.new_photo || !respondForm.message) {
    respondError.value = '请上传现场照片并写下寄语'
    return
  }
  await run(async () => {
    await api.post(`/stories/${props.story.id}/respond`, respondForm)
    showRespond.value = false
    respondForm.new_photo = ''
    respondForm.message = ''
  }, respondError)
}

async function accept(r) {
  if (!confirm(`确认采纳 ${r.user?.nickname} 的回应？${props.story.bounty} 枚悬赏硬币将结算给对方。`)) return
  await run(() => api.post(`/responses/${r.id}/accept`))
}

async function append() {
  if (!appendAmount.value || appendAmount.value < 1) { appendError.value = '请输入追加数量'; return }
  await run(async () => {
    await api.post(`/stories/${props.story.id}/append`, { amount: appendAmount.value })
    showAppend.value = false
    appendAmount.value = null
  }, appendError)
}

async function cancel() {
  if (!confirm('确认取消该需求？冻结的悬赏硬币将退回。')) return
  await run(() => api.post(`/stories/${props.story.id}/cancel`))
}
</script>
