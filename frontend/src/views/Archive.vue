<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getArchives } from '../api'

const archives = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res: any = await getArchives()
    archives.value = res.data || []
  } catch {
    archives.value = []
  } finally {
    loading.value = false
  }
})

function monthName(month: number): string {
  const names = ['', 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  return names[month] || ''
}
</script>

<template>
  <div class="container">
    <h1 style="margin-bottom: 32px;">Archive</h1>
    <div v-if="loading" class="loading">Loading...</div>
    <template v-else>
      <div v-if="archives.length === 0" class="loading">No articles yet.</div>
      <div v-for="item in archives" :key="`${item.year}-${item.month}`" class="archive-item">
        <span>{{ item.year }} {{ monthName(item.month) }}</span>
        <span class="date">{{ item.count }} articles</span>
      </div>
    </template>
  </div>
</template>
