<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getArticleBySlug } from '../api'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

const route = useRoute()
const article = ref<any>(null)
const loading = ref(true)
const error = ref('')

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(async () => {
  try {
    const slug = route.params.slug as string
    const res: any = await getArticleBySlug(slug)
    article.value = res.data
  } catch (e: any) {
    error.value = 'Article not found'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="container">
    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="error" class="loading">{{ error }}</div>
    <article v-else class="article-detail">
      <h1>{{ article.title }}</h1>
      <div class="article-meta">
        <span>{{ formatDate(article.created_at) }}</span>
        <span v-if="article.category">
          <RouterLink :to="`/category/${article.category.id}`">
            {{ article.category.name }}
          </RouterLink>
        </span>
        <span>{{ article.view_count }} views</span>
      </div>
      <div v-if="article.tags?.length" class="article-tags">
        <span v-for="tag in article.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
      </div>
      <div class="article-body">
        <MdPreview :modelValue="article.content" />
      </div>
    </article>
  </div>
</template>
