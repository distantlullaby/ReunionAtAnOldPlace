<template>
  <div v-if="data" class="profile">
    <!-- 用户信息与硬币 -->
    <section class="user-card">
      <div class="avatar">{{ (data.user.nickname || data.user.username).slice(0, 1) }}</div>
      <div class="user-info">
        <h2>{{ data.user.nickname || data.user.username }}</h2>
        <div class="username">@{{ data.user.username }}</div>
      </div>
      <div class="coin-box">
        <div class="coin-num"><span class="coin">●</span>{{ data.user.coinBalance }}</div>
        <div class="coin-label">记忆硬币</div>
      </div>
    </section>

    <div class="stat-row">
      <div class="stat">
        <b>{{ data.stats.storyCount }}</b><span>发布求看</span>
      </div>
      <div class="stat">
        <b>{{ data.stats.responseCount }}</b><span>替人去拍</span>
      </div>
      <div class="stat tip">硬币在「发布冻结 → 被采纳结算」中循环流转</div>
    </div>

    <!-- 标签页 -->
    <div class="tabs">
      <button v-for="t in tabs" :key="t.key" :class="{ active: tab === t.key }" @click="tab = t.key">
        {{ t.label }}
      </button>
    </div>

    <!-- 我发布的 -->
    <section v-if="tab === 'stories'" class="panel">
      <div v-if="!data.myStories.length" class="empty">
        <div class="big">📭</div>还没有发布过求看
      </div>
      <div v-for="s in data.myStories" :key="s.id" class="rec">
        <div class="rec-main">
          <div class="rec-title">{{ s.title }}</div>
          <div class="rec-sub">
            <span class="tag" :class="s.status">{{ s.status === 'open' ? '求看中' : '已重逢' }}</span>
            <span>托管悬赏 <span class="coin">{{ s.reward }}</span></span>
            <span>{{ fmt(s.createdAt) }}</span>
          </div>
        </div>
        <router-link class="btn sm ghost" to="/">去广场查看</router-link>
      </div>
    </section>

    <!-- 我去拍的 -->
    <section v-if="tab === 'responses'" class="panel">
      <div v-if="!data.myResponses.length" class="empty">
        <div class="big">📸</div>还没有替别人去拍过
      </div>
      <div v-for="r in data.myResponses" :key="r.id" class="rec">
        <img :src="r.newPhoto" class="rec-img" @error="onImgError" />
        <div class="rec-main">
          <div class="rec-title">{{ r.storyTitle }}</div>
          <div class="rec-sub">
            <span class="tag" :class="r.status">{{ r.status === 'open' ? '等待采纳' : '已重逢' }}</span>
            <span>{{ r.message || '（无寄语）' }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 硬币流水 -->
    <section v-if="tab === 'coins'" class="panel">
      <div v-if="!data.transactions.length" class="empty">
        <div class="big">🪙</div>暂无硬币流水
      </div>
      <div class="ledger">
        <div v-for="t in data.transactions" :key="t.id" class="ledger-item">
          <div class="ledger-ic">{{ iconOf(t.type) }}</div>
          <div class="ledger-main">
            <div class="ledger-title">{{ labelOf(t.type) }}</div>
            <div class="ledger-remark">{{ t.remark }}</div>
            <div class="ledger-time">{{ fmt(t.createdAt) }}</div>
          </div>
          <div class="ledger-amount" :class="t.amount >= 0 ? 'plus' : 'minus'">
            {{ t.amount >= 0 ? '+' : '' }}{{ t.amount }}
            <div class="ledger-balance">余额 {{ t.balanceAfter }}</div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { userApi } from '../api'
import { auth } from '../store/auth'

const data = ref(null)
const tab = ref('stories')
const tabs = [
  { key: 'stories', label: '我发布的' },
  { key: 'responses', label: '我去拍的' },
  { key: 'coins', label: '硬币流水' }
]

const typeMap = {
  register_gift: { label: '注册赠送', icon: '🎁' },
  freeze: { label: '发布冻结', icon: '🧊' },
  append_freeze: { label: '追加冻结', icon: '🧊' },
  reward_income: { label: '悬赏收入', icon: '🪙' },
  refund: { label: '解冻退回', icon: '↩️' }
}
const labelOf = (t) => typeMap[t]?.label || t
const iconOf = (t) => typeMap[t]?.icon || '·'

function fmt(t) {
  return (t || '').replace('T', ' ').slice(0, 16)
}

const fallbackImg =
  'data:image/svg+xml,' +
  encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120"><rect width="100%" height="100%" fill="#e3d4bd"/><text x="50%" y="54%" font-size="26" text-anchor="middle">🌫️</text></svg>`
  )
function onImgError(e) {
  if (e.target.src !== fallbackImg) e.target.src = fallbackImg
}

onMounted(async () => {
  try {
    const res = await userApi.profile()
    data.value = res
    auth.user = res.user // 同步顶部余额
  } catch (e) {
    data.value = null
  }
})
</script>

<style scoped>
.user-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: linear-gradient(135deg, #e8d3b6, #f3e6d0);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 22px;
}
.avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--amber);
  color: #fff;
  font-size: 28px;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.user-info {
  flex: 1;
}
.user-info h2 {
  font-size: 21px;
}
.username {
  color: var(--ink-soft);
  font-size: 13px;
}
.coin-box {
  text-align: center;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 14px;
  padding: 12px 20px;
}
.coin-num {
  font-size: 28px;
  font-weight: 800;
}
.coin-num .coin {
  font-size: 16px;
  margin-right: 3px;
}
.coin-label {
  font-size: 12px;
  color: var(--ink-soft);
}
.stat-row {
  display: flex;
  gap: 12px;
  margin: 16px 0;
}
.stat {
  flex: 1;
  background: var(--paper-2);
  border: 1px solid var(--line);
  border-radius: 12px;
  padding: 14px;
  text-align: center;
}
.stat b {
  display: block;
  font-size: 22px;
}
.stat span {
  font-size: 12px;
  color: var(--ink-soft);
}
.stat.tip {
  flex: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--ink-soft);
  font-style: italic;
}
.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
.tabs button {
  flex: 1;
  padding: 10px;
  border-radius: 10px;
  background: transparent;
  border: 1.5px solid var(--line);
  font-weight: 600;
  color: var(--ink-soft);
}
.tabs button.active {
  background: var(--teal);
  border-color: var(--teal);
  color: #fff;
}
.panel {
  background: var(--paper-2);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  padding: 8px 16px;
}
.rec {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 0;
  border-bottom: 1px solid var(--line);
}
.rec:last-child {
  border-bottom: none;
}
.rec-main {
  flex: 1;
  min-width: 0;
}
.rec-title {
  font-weight: 700;
  font-size: 15px;
}
.rec-sub {
  margin-top: 5px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: var(--ink-soft);
  flex-wrap: wrap;
}
.rec-img {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 8px;
}
.ledger-item {
  display: flex;
  gap: 12px;
  padding: 14px 0;
  border-bottom: 1px solid var(--line);
}
.ledger-item:last-child {
  border-bottom: none;
}
.ledger-ic {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--paper);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}
.ledger-main {
  flex: 1;
}
.ledger-title {
  font-weight: 700;
  font-size: 14px;
}
.ledger-remark {
  font-size: 13px;
  color: var(--ink);
}
.ledger-time {
  font-size: 12px;
  color: var(--ink-soft);
}
.ledger-amount {
  text-align: right;
  font-weight: 800;
  font-size: 17px;
}
.ledger-amount.plus {
  color: var(--teal);
}
.ledger-amount.minus {
  color: var(--amber-dark);
}
.ledger-balance {
  font-size: 11px;
  font-weight: 400;
  color: var(--ink-soft);
  margin-top: 2px;
}
</style>
