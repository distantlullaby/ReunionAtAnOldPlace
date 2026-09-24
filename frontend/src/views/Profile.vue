<template>
  <div>
    <div v-if="!me" class="empty">请先登录后查看个人中心</div>
    <template v-else>
      <!-- 硬币总览 -->
      <div class="card" style="display:flex;gap:32px;align-items:center;flex-wrap:wrap">
        <div>
          <div style="font-size:13px;color:var(--muted)">可用记忆硬币</div>
          <div style="font-size:34px;font-weight:700;color:var(--gold)">🪙 {{ me.balance }}</div>
        </div>
        <div>
          <div style="font-size:13px;color:var(--muted)">冻结中（悬赏）</div>
          <div style="font-size:24px;font-weight:600">{{ me.frozen }}</div>
        </div>
        <div style="margin-left:auto;font-size:13px;color:var(--muted)">
          {{ me.nickname }}（@{{ me.username }}）
        </div>
      </div>

      <!-- 标签页 -->
      <div style="display:flex;gap:8px;margin:20px 0 16px">
        <button v-for="t in tabs" :key="t.key" class="btn btn-sm"
                :class="tab === t.key ? 'btn-primary' : 'btn-ghost'"
                @click="tab = t.key">{{ t.label }}</button>
      </div>

      <!-- 我的发布 -->
      <div v-if="tab === 'stories'">
        <div v-if="!myStories.length" class="empty">你还没有发布过求看需求</div>
        <StoryCard v-for="s in myStories" :key="s.id" :story="s" :me="me" @changed="loadAll" />
      </div>

      <!-- 我的回应 -->
      <div v-else-if="tab === 'responses'">
        <div v-if="!myResponses.length" class="empty">你还没有替别人拍过照片</div>
        <div v-for="item in myResponses" :key="item.id" class="card">
          <div style="display:flex;gap:14px;align-items:flex-start;flex-wrap:wrap">
            <img :src="item.new_photo" style="width:110px;height:110px;object-fit:cover;border-radius:10px" />
            <div style="flex:1;min-width:220px">
              <div style="font-weight:700">{{ item.story?.title }}
                <span class="tag" :class="'tag-' + item.status">{{ respStatusText(item.status) }}</span>
              </div>
              <div style="font-size:13px;color:var(--muted);margin:4px 0">📍 {{ item.story?.location }} · 发起人 {{ item.story?.user?.nickname }}</div>
              <div style="font-size:14px">寄语：{{ item.message }}</div>
              <div v-if="item.status === 'accepted'" style="color:var(--gold);font-size:13px;margin-top:4px">
                🪙 已结算 +{{ item.story?.bounty }} 枚记忆硬币
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 硬币流水 -->
      <div v-else>
        <div v-if="!transactions.length" class="empty">暂无硬币流水</div>
        <div v-else class="card" style="padding:8px 20px">
          <div v-for="t in transactions" :key="t.id"
               style="display:flex;justify-content:space-between;gap:12px;padding:12px 0;border-bottom:1px dashed var(--line);font-size:14px">
            <div>
              <div>{{ t.remark }}</div>
              <div style="font-size:12px;color:var(--muted)">{{ new Date(t.created_at).toLocaleString('zh-CN') }}</div>
            </div>
            <div :style="{ color: t.amount > 0 ? '#2e7d32' : (t.amount < 0 ? '#b03a2e' : 'var(--muted)'), fontWeight: 700, whiteSpace: 'nowrap' }">
              {{ t.amount > 0 ? '+' + t.amount : (t.amount === 0 ? '结算' : t.amount) }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import api, { getUser, refreshMe } from '../api'
import StoryCard from '../components/StoryCard.vue'

const me = ref(getUser())
const tab = ref('stories')
const tabs = [
  { key: 'stories', label: '我的发布' },
  { key: 'responses', label: '我的回应' },
  { key: 'coins', label: '硬币流水' }
]
const myStories = ref([])
const myResponses = ref([])
const transactions = ref([])

function syncUser() {
  me.value = getUser()
  if (me.value) loadAll()
}

function respStatusText(s) {
  return { pending: '等待确认', accepted: '已被采纳', rejected: '未被采纳' }[s] || s
}

async function loadAll() {
  const [s, r, t] = await Promise.all([
    api.get('/my/stories'),
    api.get('/my/responses'),
    api.get('/my/transactions')
  ])
  myStories.value = s.data.stories || []
  myResponses.value = r.data.responses || []
  transactions.value = t.data.transactions || []
}

onMounted(async () => {
  if (me.value) {
    await refreshMe().catch(() => {})
    await loadAll()
  }
  window.addEventListener('auth-changed', syncUser)
})
onUnmounted(() => window.removeEventListener('auth-changed', syncUser))
</script>
