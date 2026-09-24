<template>
  <div>
    <!-- 顶部横幅 -->
    <section class="hero">
      <h1>那些回不去的角落，<br />总有人能替你再看一眼</h1>
      <p>贴上老照片、留下回忆与悬赏硬币，让此刻路过那里的人，拍下它现在的模样。</p>
      <button v-if="auth.isLogin" class="btn" @click="openCreate">📮 发布求看</button>
      <router-link v-else to="/login" class="btn">登录后发布求看</router-link>
    </section>

    <!-- 状态过滤 -->
    <div class="filter-bar">
      <button
        v-for="f in filters"
        :key="f.value"
        class="filter"
        :class="{ active: status === f.value }"
        @click="changeFilter(f.value)"
      >
        {{ f.label }}
      </button>
    </div>

    <!-- Feed 列表 -->
    <div v-if="loading && !list.length" class="empty">加载中…</div>
    <div v-else-if="!list.length" class="empty">
      <div class="big">🗺️</div>
      <p>这里还没有求看，发布第一张双面明信片吧</p>
    </div>
    <template v-else>
      <PostcardCard
        v-for="item in list"
        :key="item.story.id"
        :data="item"
        @respond="openRespond"
        @append="openAppend"
        @accept="onAccept"
        @preview="previewImg"
      />

      <div class="pager">
        <button class="btn ghost sm" :disabled="page <= 1" @click="page--; load()">
          上一页
        </button>
        <span class="page-info">第 {{ page }} 页</span>
        <button
          class="btn ghost sm"
          :disabled="page * size >= total"
          @click="page++; load()"
        >
          下一页
        </button>
      </div>
    </template>

    <!-- 发布求看弹窗 -->
    <Modal v-if="showCreate" @close="showCreate = false" title="发布一张「过去」明信片">
      <div class="field">
        <label>标题 *</label>
        <input v-model.trim="createForm.title" class="input" placeholder="如：再看一眼三中老校门" />
      </div>
      <div class="field">
        <label>地点</label>
        <input v-model.trim="createForm.location" class="input" placeholder="如：某市某区·城南中学" />
      </div>
      <div class="field">
        <label>过去回忆</label>
        <textarea v-model="createForm.memoryText" class="textarea" placeholder="讲讲那里曾经的故事…"></textarea>
      </div>
      <div class="field">
        <label>老照片</label>
        <ImageUploader v-model="createForm.oldPhotos" multiple :max="6" add-text="老照片" />
      </div>
      <div class="field">
        <label>悬赏记忆硬币（发布即从余额冻结）</label>
        <input v-model.number="createForm.reward" type="number" min="0" class="input" />
        <small>当前余额 <span class="coin">{{ auth.user?.coinBalance ?? 0 }}</span> 枚</small>
      </div>
      <div class="modal-actions">
        <button class="btn ghost" @click="showCreate = false">取消</button>
        <button class="btn" :disabled="submitting" @click="submitCreate">
          {{ submitting ? '发布中…' : '冻结并发布' }}
        </button>
      </div>
    </Modal>

    <!-- 替他去拍弹窗 -->
    <Modal v-if="respondTarget" @close="respondTarget = null" title="替 TA 去看一眼">
      <p class="modal-tip">拍下「{{ respondTarget.title }}」此刻的样子，并附上一句寄语。</p>
      <div class="field">
        <label>当下现场新照 *</label>
        <ImageUploader v-model="respondPhoto" :max="1" add-text="现场新照" />
      </div>
      <div class="field">
        <label>给 TA 的寄语</label>
        <textarea v-model="respondMessage" class="textarea" placeholder="如：树还在，墙翻新了…"></textarea>
      </div>
      <div class="modal-actions">
        <button class="btn ghost" @click="respondTarget = null">取消</button>
        <button class="btn teal" :disabled="submitting" @click="submitRespond">
          {{ submitting ? '提交中…' : '送出当下' }}
        </button>
      </div>
    </Modal>

    <!-- 追加悬赏弹窗 -->
    <Modal v-if="appendTarget" @close="appendTarget = null" title="追加悬赏硬币">
      <p class="modal-tip">
        当前托管悬赏 <span class="coin">{{ appendTarget.reward }}</span> 枚，
        余额 <span class="coin">{{ auth.user?.coinBalance ?? 0 }}</span> 枚。
      </p>
      <div class="field">
        <label>追加数量</label>
        <input v-model.number="appendAmount" type="number" min="1" class="input" />
      </div>
      <div class="modal-actions">
        <button class="btn ghost" @click="appendTarget = null">取消</button>
        <button class="btn" :disabled="submitting" @click="submitAppend">确认追加冻结</button>
      </div>
    </Modal>

    <!-- 大图预览 -->
    <div v-if="previewUrl" class="img-mask" @click="previewUrl = ''">
      <img :src="previewUrl" alt="预览" />
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import PostcardCard from '../components/PostcardCard.vue'
import ImageUploader from '../components/ImageUploader.vue'
import Modal from '../components/Modal.vue'
import { storyApi } from '../api'
import { auth } from '../store/auth'

const emit = defineEmits(['toast'])

const filters = [
  { value: '', label: '全部' },
  { value: 'open', label: '求看中' },
  { value: 'completed', label: '已重逢' }
]
const status = ref('')
const page = ref(1)
const size = ref(5)
const total = ref(0)
const list = ref([])
const loading = ref(false)

const showCreate = ref(false)
const submitting = ref(false)
const createForm = reactive({
  title: '',
  location: '',
  memoryText: '',
  oldPhotos: [],
  reward: 10
})

const respondTarget = ref(null)
const respondPhoto = ref([])
const respondMessage = ref('')

const appendTarget = ref(null)
const appendAmount = ref(5)

const previewUrl = ref('')

async function load() {
  loading.value = true
  try {
    const res = await storyApi.list({ page: page.value, size: size.value, status: status.value })
    list.value = res.list
    total.value = res.total
  } catch (e) {
    emit('toast', e.message, 'err')
  } finally {
    loading.value = false
  }
}
load()

function changeFilter(v) {
  status.value = v
  page.value = 1
  load()
}

function previewImg(url) {
  previewUrl.value = url
}

function openCreate() {
  Object.assign(createForm, {
    title: '',
    location: '',
    memoryText: '',
    oldPhotos: [],
    reward: 10
  })
  showCreate.value = true
}

async function submitCreate() {
  if (!createForm.title) {
    emit('toast', '请填写标题', 'err')
    return
  }
  submitting.value = true
  try {
    await storyApi.create({ ...createForm })
    showCreate.value = false
    emit('toast', '已发布，悬赏硬币已冻结')
    page.value = 1
    status.value = ''
    await load()
    await refreshBalance()
  } catch (e) {
    emit('toast', e.message, 'err')
  } finally {
    submitting.value = false
  }
}

function openRespond(story) {
  respondTarget.value = story
  respondPhoto.value = []
  respondMessage.value = ''
}

async function submitRespond() {
  if (!respondPhoto.value.length) {
    emit('toast', '请上传现场照片', 'err')
    return
  }
  submitting.value = true
  try {
    await storyApi.respond(respondTarget.value.id, {
      newPhoto: respondPhoto.value[0],
      message: respondMessage.value
    })
    respondTarget.value = null
    emit('toast', '已送出当下现场，等待发起人确认')
    await load()
  } catch (e) {
    emit('toast', e.message, 'err')
  } finally {
    submitting.value = false
  }
}

function openAppend(story) {
  appendTarget.value = story
  appendAmount.value = 5
}

async function submitAppend() {
  if (!appendAmount.value || appendAmount.value <= 0) {
    emit('toast', '请输入正确的追加数量', 'err')
    return
  }
  submitting.value = true
  try {
    const updated = await storyApi.append(appendTarget.value.id, appendAmount.value)
    appendTarget.value = null
    emit('toast', `已追加，当前托管悬赏 ${updated.reward} 枚`)
    await refreshBalance()
    await load()
  } catch (e) {
    emit('toast', e.message, 'err')
  } finally {
    submitting.value = false
  }
}

async function onAccept({ id: storyId }, response) {
  if (!window.confirm(`确认采纳「${response.nickname || '路人'}」的现场，并结算悬赏硬币吗？结算后不可撤销。`)) {
    return
  }
  try {
    await storyApi.accept(storyId, response.id)
    emit('toast', '已确认采纳，悬赏硬币已结算给对方 🎉')
    await load()
    await refreshBalance()
  } catch (e) {
    emit('toast', e.message, 'err')
  }
}

// 发帖/追加/采纳后同步顶部硬币余额。
async function refreshBalance() {
  try {
    const { userApi } = await import('../api')
    const data = await userApi.profile()
    auth.user = data.user
  } catch (e) {
    /* 忽略 */
  }
}
</script>

<style scoped>
.hero {
  text-align: center;
  padding: 26px 16px 30px;
}
.hero h1 {
  font-size: 26px;
  line-height: 1.5;
  font-weight: 800;
}
.hero p {
  color: var(--ink-soft);
  font-size: 14px;
  margin: 12px 0 20px;
}
.filter-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 18px;
}
.filter {
  padding: 6px 18px;
  border-radius: 999px;
  background: transparent;
  border: 1.5px solid var(--line);
  color: var(--ink-soft);
  font-size: 13px;
  font-weight: 600;
}
.filter.active {
  background: var(--amber);
  border-color: var(--amber);
  color: #fff;
}
.pager {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 10px 0 20px;
}
.page-info {
  font-size: 13px;
  color: var(--ink-soft);
}
.modal-tip {
  font-size: 13px;
  color: var(--ink-soft);
  margin-bottom: 16px;
}
.img-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.82);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 300;
  padding: 30px;
  cursor: zoom-out;
}
.img-mask img {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
}
</style>
