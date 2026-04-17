<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboard } from '../api'

const stats = ref<any>({})
const loading = ref(true)

onMounted(async () => {
  try {
    const res: any = await getDashboard()
    stats.value = res.data
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <h2>Dashboard</h2>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="8">
        <el-card shadow="hover">
          <el-statistic title="Total Articles" :value="stats.total_articles || 0" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <el-statistic title="Published" :value="stats.public_articles || 0" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <el-statistic title="Drafts" :value="stats.draft_articles || 0" />
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="8">
        <el-card shadow="hover">
          <el-statistic title="Categories" :value="stats.total_categories || 0" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <el-statistic title="Tags" :value="stats.total_tags || 0" />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <el-statistic title="Total Views" :value="stats.total_views || 0" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>
