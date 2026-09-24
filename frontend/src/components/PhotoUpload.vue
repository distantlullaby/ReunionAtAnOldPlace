<template>
  <div class="upload-box">
    <img v-if="modelValue" :src="modelValue" alt="预览" />
    <div>
      <input ref="fileInput" type="file" accept="image/*" style="display:none" @change="onPick" />
      <button type="button" class="btn btn-ghost btn-sm" :disabled="uploading" @click="fileInput.click()">
        {{ uploading ? '上传中…' : (modelValue ? '重新选择' : '选择照片') }}
      </button>
      <p v-if="error" class="error-text">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import api from '../api'

defineProps({ modelValue: String })
const emit = defineEmits(['update:modelValue'])

const fileInput = ref(null)
const uploading = ref(false)
const error = ref('')

async function onPick(e) {
  const file = e.target.files[0]
  if (!file) return
  error.value = ''
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    const res = await api.post('/upload', fd)
    emit('update:modelValue', res.data.url)
  } catch (err) {
    error.value = err.response?.data?.error || '上传失败'
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}
</script>
