<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getTags, createTag, deleteTag } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const tags = ref<any[]>([])
const loading = ref(true)
const newTagName = ref('')

async function loadTags() {
  loading.value = true
  try {
    const res: any = await getTags()
    tags.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  if (!newTagName.value.trim()) return
  try {
    await createTag({ name: newTagName.value.trim() })
    ElMessage.success('Created')
    newTagName.value = ''
    loadTags()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('Delete this tag?', 'Confirm')
    await deleteTag(id)
    ElMessage.success('Deleted')
    loadTags()
  } catch { /* cancelled */ }
}

onMounted(loadTags)
</script>

<template>
  <div>
    <h2 style="margin-bottom: 20px;">Tags</h2>

    <div style="display: flex; gap: 12px; margin-bottom: 20px;">
      <el-input v-model="newTagName" placeholder="New tag name" @keyup.enter="handleCreate" style="width: 300px;" />
      <el-button type="primary" @click="handleCreate">Add</el-button>
    </div>

    <div style="display: flex; flex-wrap: wrap; gap: 10px;" v-loading="loading">
      <el-tag
        v-for="tag in tags"
        :key="tag.id"
        closable
        size="large"
        @close="handleDelete(tag.id)"
      >
        {{ tag.name }}
      </el-tag>
      <div v-if="!tags.length && !loading" style="color: #999;">No tags yet</div>
    </div>
  </div>
</template>
