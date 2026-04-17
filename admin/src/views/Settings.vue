<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSettings, updateSettings } from '../api'
import { ElMessage } from 'element-plus'

const settings = ref<any[]>([])
const loading = ref(true)
const saving = ref(false)

async function loadSettings() {
  loading.value = true
  try {
    const res: any = await getSettings()
    settings.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const map: Record<string, string> = {}
    for (const s of settings.value) {
      map[s.key] = s.value
    }
    await updateSettings(map)
    ElMessage.success('Settings saved')
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>

<template>
  <div v-loading="loading">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
      <h2>Site Settings</h2>
      <el-button type="primary" :loading="saving" @click="handleSave">Save All</el-button>
    </div>

    <el-form label-position="top">
      <el-form-item v-for="s in settings" :key="s.id" :label="s.description || s.key">
        <el-input
          v-if="s.type === 'text'"
          v-model="s.value"
          type="textarea"
          :rows="4"
        />
        <el-input v-else v-model="s.value" />
        <div style="font-size: 12px; color: #999; margin-top: 4px;">
          Key: {{ s.key }} | Public: {{ s.is_public ? 'Yes' : 'No' }}
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>
