<template>
  <article class="card">
    <!-- 折叠态：明信片信封式摘要 -->
    <header class="card-head" @click="toggle">
      <div class="head-main">
        <div class="title-row">
          <h3>{{ data.story.title }}</h3>
          <span class="tag" :class="data.story.status">
            {{ data.story.status === 'open' ? '求看中' : '已重逢' }}
          </span>
        </div>
        <div class="meta">
          <span class="loc">📍 {{ data.story.location || '某处旧时光' }}</span>
          <span class="author">· {{ data.authorName || '匿名' }}</span>
        </div>
      </div>
      <div class="head-right">
        <div class="reward"><span class="coin">●</span>{{ data.story.reward }}</div>
        <div class="resp-count">{{ data.responseCount }} 个现场</div>
        <span class="toggle-icon" :class="{ open: expanded }">⌄</span>
      </div>
    </header>

    <!-- 展开态：双面明信片 -->
    <transition name="expand">
      <div v-if="expanded" class="card-body">
        <div v-if="loading" class="loading">明信片展开中…</div>

        <template v-else-if="detail">
          <div class="face-switch">
            <button :class="{ active: face === 'past' }" @click="face = 'past'">
              过去回忆
            </button>
            <button :class="{ active: face === 'now' }" @click="face = 'now'">
              当下现场 <i>{{ detail.responses.length }}</i>
            </button>
          </div>

          <div class="postcard-stage">
            <transition name="flip" mode="out-in">
              <!-- 过去面 -->
              <section v-if="face === 'past'" key="past" class="face past-face">
                <div class="face-label">— 过 去 · PAST —</div>
                <p class="memory">{{ detail.memoryText || '发起人没有写下更多回忆。' }}</p>
                <div v-if="detail.oldPhotos?.length" class="photos old">
                  <img
                    v-for="(p, i) in detail.oldPhotos"
                    :key="i"
                    :src="p"
                    :alt="`老照片${i + 1}`"
                    @click="$emit('preview', p)"
                    @error="onImgError"
                  />
                </div>
                <div v-else class="no-photo">没有附上老照片</div>
                <div class="stamp old-stamp">老照片</div>
              </section>

              <!-- 当下面 -->
              <section v-else key="now" class="face now-face">
                <div class="face-label">— 当 下 · NOW —</div>

                <div v-if="!detail.responses.length" class="no-response">
                  <div class="big">📷</div>
                  <p>还没有人替 TA 去看一眼</p>
                </div>

                <div v-else class="resp-list">
                  <div
                    v-for="r in detail.responses"
                    :key="r.id"
                    class="resp"
                    :class="{ adopted: detail.acceptedResponseId === r.id }"
                  >
                    <img
                      :src="r.newPhoto"
                      alt="现场新照"
                      class="resp-img"
                      @click="$emit('preview', r.newPhoto)"
                      @error="onImgError"
                    />
                    <div class="resp-info">
                      <div class="resp-msg">{{ r.message || '（没有寄语）' }}</div>
                      <div class="resp-by">
                        {{ r.nickname || '路人' }} · {{ formatTime(r.createdAt) }}
                        <span v-if="detail.acceptedResponseId === r.id" class="adopted-tag">
                          已采纳 ✓
                        </span>
                      </div>
                    </div>
                    <button
                      v-if="isOwner && detail.status === 'open'"
                      class="btn sm teal"
                      @click="$emit('accept', detail, r)"
                    >
                      确认采纳
                    </button>
                  </div>
                </div>
                <div class="stamp now-stamp">新现场</div>
              </section>
            </transition>
          </div>

          <!-- 底部操作 -->
          <div class="card-actions">
            <template v-if="detail.status === 'open'">
              <button v-if="isOwner" class="btn ghost sm" @click="$emit('append', detail)">
                ＋ 追加悬赏
              </button>
              <button
                v-else-if="auth.isLogin"
                class="btn sm"
                @click="$emit('respond', detail)"
              >
                替他去拍 📸
              </button>
              <router-link v-else to="/login" class="btn sm">登录后替他去拍</router-link>
            </template>
            <div v-else class="done-tip">
              这张明信片已通过他人的镜头重逢，悬赏已结算
            </div>
          </div>
        </template>
      </div>
    </transition>
  </article>
</template>

<script setup>
import { computed, ref } from 'vue'
import { storyApi } from '../api'
import { auth } from '../store/auth'

const props = defineProps({
  data: { type: Object, required: true }
})
defineEmits(['respond', 'append', 'accept', 'preview'])

const expanded = ref(false)
const loading = ref(false)
const detail = ref(null)
const face = ref('past')

const isOwner = computed(() => auth.user?.id === props.data.story.userId)

async function toggle() {
  expanded.value = !expanded.value
  if (expanded.value && !detail.value) {
    await load()
  }
}

async function load() {
  loading.value = true
  try {
    detail.value = await storyApi.detail(props.data.story.id)
  } finally {
    loading.value = false
  }
}

function formatTime(t) {
  return (t || '').replace('T', ' ').slice(0, 16)
}

// 图片缺失/失效时用内联占位图兜底，避免裂图。
const fallbackImg =
  'data:image/svg+xml,' +
  encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="240" height="160">
      <rect width="100%" height="100%" fill="#e3d4bd"/>
      <text x="50%" y="46%" font-size="14" fill="#8a7563" text-anchor="middle" font-family="sans-serif">照片待补</text>
      <text x="50%" y="68%" font-size="22" text-anchor="middle">🌫️</text>
    </svg>`
  )
function onImgError(e) {
  if (e.target.src !== fallbackImg) e.target.src = fallbackImg
}

defineExpose({ reload: load })
</script>

<style scoped>
.card {
  background: var(--paper-2);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  margin-bottom: 18px;
  overflow: hidden;
}
.card-head {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px;
  cursor: pointer;
}
.card-head:hover {
  background: rgba(196, 120, 58, 0.04);
}
.head-main {
  flex: 1;
  min-width: 0;
}
.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.title-row h3 {
  font-size: 17px;
  font-weight: 700;
}
.meta {
  margin-top: 3px;
  font-size: 13px;
  color: var(--ink-soft);
}
.head-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  position: relative;
}
.reward {
  font-size: 18px;
  font-weight: 800;
}
.reward .coin {
  font-size: 12px;
  margin-right: 3px;
}
.resp-count {
  font-size: 12px;
  color: var(--ink-soft);
}
.toggle-icon {
  position: absolute;
  bottom: -6px;
  right: -4px;
  color: var(--ink-soft);
  font-size: 18px;
  transition: transform 0.2s;
}
.toggle-icon.open {
  transform: rotate(180deg);
}

.card-body {
  padding: 0 18px 18px;
}
.loading,
.no-response {
  text-align: center;
  color: var(--ink-soft);
  padding: 30px 0;
}

.face-switch {
  display: flex;
  background: var(--paper);
  border-radius: 999px;
  padding: 4px;
  margin-bottom: 14px;
}
.face-switch button {
  flex: 1;
  padding: 8px;
  border-radius: 999px;
  background: transparent;
  font-size: 14px;
  font-weight: 600;
  color: var(--ink-soft);
}
.face-switch button.active {
  background: #fff;
  color: var(--amber-dark);
  box-shadow: 0 2px 8px rgba(98, 70, 40, 0.12);
}
.face-switch i {
  font-style: normal;
  color: var(--teal);
  font-weight: 700;
}

.postcard-stage {
  position: relative;
}
.face {
  position: relative;
  border-radius: 12px;
  padding: 20px 20px 24px;
  min-height: 200px;
}
.past-face {
  background: linear-gradient(135deg, #f3e6d0, #efe0c6);
  border: 1px dashed #cdb48e;
}
.now-face {
  background: linear-gradient(135deg, #e4efee, #d9e9e7);
  border: 1px dashed #9fc4c2;
}
.face-label {
  text-align: center;
  font-size: 12px;
  letter-spacing: 3px;
  color: var(--ink-soft);
  margin-bottom: 14px;
}
.memory {
  font-size: 15px;
  white-space: pre-wrap;
  margin-bottom: 14px;
}
.photos {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 10px;
}
.photos img {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border-radius: 8px;
  cursor: zoom-in;
  border: 3px solid #fff;
  box-shadow: 0 3px 10px rgba(98, 70, 40, 0.18);
  filter: sepia(0.25);
}
.no-photo {
  color: var(--ink-soft);
  font-size: 13px;
  font-style: italic;
}
.stamp {
  position: absolute;
  top: 14px;
  right: 16px;
  font-size: 11px;
  padding: 3px 8px;
  border: 1.5px solid;
  border-radius: 4px;
  transform: rotate(8deg);
  opacity: 0.65;
}
.old-stamp {
  border-color: var(--amber-dark);
  color: var(--amber-dark);
}
.now-stamp {
  border-color: var(--teal);
  color: var(--teal);
}

.resp-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.resp {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 10px;
  padding: 10px;
}
.resp.adopted {
  outline: 2px solid var(--teal);
}
.resp-img {
  width: 72px;
  height: 72px;
  object-fit: cover;
  border-radius: 8px;
  cursor: zoom-in;
  flex-shrink: 0;
}
.resp-info {
  flex: 1;
  min-width: 0;
}
.resp-msg {
  font-size: 14px;
  word-break: break-word;
}
.resp-by {
  margin-top: 4px;
  font-size: 12px;
  color: var(--ink-soft);
}
.adopted-tag {
  color: var(--teal);
  font-weight: 700;
}
.no-response .big {
  font-size: 34px;
  margin-bottom: 6px;
}

.card-actions {
  margin-top: 14px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
.done-tip {
  width: 100%;
  text-align: center;
  font-size: 13px;
  color: var(--ink-soft);
  font-style: italic;
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s ease;
}
.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
.flip-enter-active,
.flip-leave-active {
  transition: all 0.28s ease;
}
.flip-enter-from {
  opacity: 0;
  transform: rotateY(-35deg) translateX(10px);
}
.flip-leave-to {
  opacity: 0;
  transform: rotateY(35deg) translateX(-10px);
}
</style>
