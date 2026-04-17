<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  size: number
  total: number
}>()

const emit = defineEmits<{
  (e: 'update:page', page: number): void
}>()

const totalPages = computed(() => Math.ceil(props.total / props.size))

const pages = computed(() => {
  const arr: number[] = []
  const start = Math.max(1, props.page - 2)
  const end = Math.min(totalPages.value, props.page + 2)
  for (let i = start; i <= end; i++) {
    arr.push(i)
  }
  return arr
})
</script>

<template>
  <div v-if="totalPages > 1" class="pagination">
    <button :disabled="page <= 1" @click="emit('update:page', page - 1)">Prev</button>
    <button
      v-for="p in pages"
      :key="p"
      :class="{ active: p === page }"
      @click="emit('update:page', p)"
    >
      {{ p }}
    </button>
    <button :disabled="page >= totalPages" @click="emit('update:page', page + 1)">Next</button>
  </div>
</template>
