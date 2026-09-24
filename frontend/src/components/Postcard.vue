<template>
  <div class="postcard-wrap">
    <div class="postcard" :class="{ flipped }" @click="flipped = !flipped">
      <!-- 正面：过去回忆 + 老照片 -->
      <div class="postcard-face front">
        <div class="face-label past">过去 · 回忆那一面</div>
        <div class="face-body">
          <div class="face-photo">
            <img v-if="story.old_photo" :src="story.old_photo" alt="老照片" />
            <div v-else class="photo-placeholder">没有留下老照片<br />只有一段回忆</div>
          </div>
          <div class="face-text">{{ story.memory_text }}</div>
        </div>
      </div>
      <!-- 背面：当下现场 + 新照片 -->
      <div class="postcard-face back">
        <div class="face-label now">当下 · 现场这一面</div>
        <div class="face-body">
          <div class="face-photo">
            <img v-if="accepted" :src="accepted.new_photo" alt="现场新照片" />
            <div v-else class="photo-placeholder">还没有人替 TA 去拍<br />或许就是你？</div>
          </div>
          <div class="face-text">
            <template v-if="accepted">
              <p>{{ accepted.message }}</p>
              <p style="margin-top:10px;color:var(--muted);font-size:13px">—— {{ accepted.user?.nickname }} 摄于当下</p>
            </template>
            <template v-else>
              <span style="color:var(--muted)">这里还是空白，等待一位有心人路过 {{ story.location }}，替 TA 再看一眼。</span>
            </template>
          </div>
        </div>
      </div>
    </div>
    <div class="postcard-flip-hint">点击明信片翻面 · {{ flipped ? '当下现场' : '过去回忆' }}</div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({ story: { type: Object, required: true } })
const flipped = ref(false)
const accepted = computed(() => (props.story.responses || []).find(r => r.status === 'accepted'))
</script>
