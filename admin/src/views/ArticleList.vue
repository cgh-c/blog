<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles, deleteArticle } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const articles = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = 15
const loading = ref(true)

async function loadArticles() {
  loading.value = true
  try {
    const res: any = await getArticles({ page: page.value, size })
    articles.value = res.data.list || []
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

function visibilityTag(v: string) {
  if (v === 'public') return 'success'
  if (v === 'draft') return 'warning'
  return 'info'
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('Are you sure to delete this article?', 'Confirm')
    await deleteArticle(id)
    ElMessage.success('Deleted')
    loadArticles()
  } catch {
    // cancelled
  }
}

onMounted(loadArticles)
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
      <h2>Articles</h2>
      <el-button type="primary" @click="router.push('/articles/new')">New Article</el-button>
    </div>

    <el-table :data="articles" v-loading="loading" stripe>
      <el-table-column prop="title" label="Title" min-width="200" />
      <el-table-column label="Category" width="120">
        <template #default="{ row }">
          {{ row.category?.name || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="Visibility" width="100">
        <template #default="{ row }">
          <el-tag :type="visibilityTag(row.visibility)" size="small">
            {{ row.visibility }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Views" prop="view_count" width="80" />
      <el-table-column label="Created" width="120">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="Actions" width="160" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="router.push(`/articles/${row.id}/edit`)">Edit</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-if="total > size"
      style="margin-top: 20px; justify-content: center;"
      :current-page="page"
      :page-size="size"
      :total="total"
      layout="prev, pager, next"
      @current-change="(p: number) => { page = p; loadArticles() }"
    />
  </div>
</template>
