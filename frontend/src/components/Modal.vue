<template>
  <teleport to="body">
    <div class="modal-mask" @click.self="$emit('close')">
      <div class="modal">
        <h3>{{ title }}</h3>
        <slot />
      </div>
    </div>
  </teleport>
</template>

<script setup>
defineProps({ title: { type: String, default: '' } })
defineEmits(['close'])
</script>
