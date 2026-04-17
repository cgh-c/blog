<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticle, createArticle, updateArticle, getCategories, getTags, uploadFile } from '../api'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const articleId = computed(() => route.params.id ? Number(route.params.id) : null)
const isEdit = computed(() => !!articleId.value)

const form = ref({
  title: '',
  content: '',
  summary: '',
  cover_image: '',
  category_id: null as number | null,
  tag_ids: [] as number[],
  visibility: 'draft' as string,
})

const categories = ref<any[]>([])
const tags = ref<any[]>([])
const loading = ref(false)
const saving = ref(false)

async function loadData() {
  loading.value = true
  try {
    const [catRes, tagRes]: any[] = await Promise.all([getCategories(), getTags()])
    categories.value = catRes.data || []
    tags.value = tagRes.data || []

    if (isEdit.value) {
      const res: any = await getArticle(articleId.value!)
      const a = res.data
      form.value = {
        title: a.title,
        content: a.content,
        summary: a.summary || '',
        cover_image: a.cover_image || '',
        category_id: a.category_id || null,
        tag_ids: a.tags?.map((t: any) => t.id) || [],
        visibility: a.visibility,
      }
    }
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!form.value.title.trim()) {
    ElMessage.warning('Title is required')
    return
  }
  if (!form.value.content.trim()) {
    ElMessage.warning('Content is required')
    return
  }

  saving.value = true
  try {
    const data = { ...form.value }
    if (isEdit.value) {
      await updateArticle(articleId.value!, data)
      ElMessage.success('Updated')
    } else {
      await createArticle(data)
      ElMessage.success('Created')
    }
    router.push('/articles')
  } finally {
    saving.value = false
  }
}

async function handleUploadImg(files: File[], callback: (urls: string[]) => void) {
  const urls: string[] = []
  for (const file of files) {
    try {
      const res: any = await uploadFile(file)
      urls.push(res.data.url)
    } catch {
      ElMessage.error('Upload failed')
    }
  }
  callback(urls)
}

onMounted(loadData)
</script>

<template>
  <div v-loading="loading">
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
      <h2>{{ isEdit ? 'Edit Article' : 'New Article' }}</h2>
      <div>
        <el-button @click="router.back()">Cancel</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">Save</el-button>
      </div>
    </div>

    <el-form label-position="top">
      <el-form-item label="Title">
        <el-input v-model="form.title" placeholder="Article title" size="large" />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="Category">
            <el-select v-model="form.category_id" placeholder="Select category" clearable style="width: 100%;">
              <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Tags">
            <el-select v-model="form.tag_ids" multiple placeholder="Select tags" style="width: 100%;">
              <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.id" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Visibility">
            <el-select v-model="form.visibility" style="width: 100%;">
              <el-option label="Draft" value="draft" />
              <el-option label="Public" value="public" />
              <el-option label="Private" value="private" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="Summary">
        <el-input v-model="form.summary" type="textarea" :rows="2" placeholder="Brief summary" />
      </el-form-item>

      <el-form-item label="Cover Image URL">
        <el-input v-model="form.cover_image" placeholder="URL or upload via editor" />
      </el-form-item>

      <el-form-item label="Content">
        <MdEditor
          v-model="form.content"
          style="height: 500px;"
          @onUploadImg="handleUploadImg"
        />
      </el-form-item>
    </el-form>
  </div>
</template>
