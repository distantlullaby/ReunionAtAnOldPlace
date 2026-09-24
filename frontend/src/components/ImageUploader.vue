<template>
  <div class="uploader">
    <div class="thumbs">
      <div v-for="(url, i) in modelValue" :key="url" class="thumb">
        <img :src="url" alt="预览" />
        <button v-if="!disabled" type="button" class="thumb-del" @click="remove(i)">×</button>
      </div>
      <button
        v-if="!disabled && (!max || modelValue.length < max)"
        type="button"
        class="thumb-add"
        @click="pick"
        :disabled="uploading"
      >
        <span v-if="uploading">上传中…</span>
        <span v-else>＋<br />{{ addText }}</span>
      </button>
    </div>
    <input ref="fileRef" type="file" accept="image/*" :multiple="multiple" hidden @change="onChange" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import request from '../api/request'

const props = defineProps({
  modelValue: { type: Array, default: () => [] },
  multiple: { type: Boolean, default: false },
  max: { type: Number, default: 0 },
  disabled: { type: Boolean, default: false },
  addText: { type: String, default: '上传照片' }
})
const emit = defineEmits(['update:modelValue'])

const fileRef = ref(null)
const uploading = ref(false)

function pick() {
  fileRef.value?.click()
}

function remove(i) {
  const next = [...props.modelValue]
  next.splice(i, 1)
  emit('update:modelValue', next)
}

async function onChange(e) {
  const files = Array.from(e.target.files || [])
  if (!files.length) return
  uploading.value = true
  try {
    const urls = [...props.modelValue]
    for (const file of files) {
      if (props.max && urls.length >= props.max) break
      const fd = new FormData()
      fd.append('file', file)
      const res = await request.post('/upload', fd, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      urls.push(res.url)
    }
    emit('update:modelValue', urls)
  } catch (err) {
    alert(err.message || '图片上传失败')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}
</script>

<style scoped>
.thumbs {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.thumb {
  position: relative;
  width: 84px;
  height: 84px;
  border-radius: 10px;
  overflow: hidden;
  border: 1.5px solid var(--line);
}
.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.thumb-del {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(67, 52, 42, 0.7);
  color: #fff;
  font-size: 14px;
  line-height: 1;
}
.thumb-add {
  width: 84px;
  height: 84px;
  border-radius: 10px;
  border: 1.5px dashed var(--amber);
  background: rgba(196, 120, 58, 0.06);
  color: var(--amber-dark);
  font-size: 12px;
  line-height: 1.5;
}
.thumb-add:hover {
  background: rgba(196, 120, 58, 0.13);
}
</style>
