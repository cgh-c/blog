<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getArticleBySlug } from '../api'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'
import TableOfContents from '../components/TableOfContents.vue'

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
    <div v-else class="article-layout">
      <article class="article-detail">
        <h1>{{ article.title }}</h1>
        <div class="article-meta">
          <span class="meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
            {{ formatDate(article.created_at) }}
          </span>
          <span v-if="article.category" class="meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
            <RouterLink :to="`/category/${article.category.id}`">
              {{ article.category.name }}
            </RouterLink>
          </span>
          <span class="meta-item">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
            {{ article.view_count }} views
          </span>
        </div>
        <div v-if="article.tags?.length" class="article-tags">
          <span v-for="tag in article.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
        </div>
        <div class="article-body">
          <MdPreview :modelValue="article.content" previewTheme="github" :codeTheme="'github'" />
        </div>
      </article>
      <aside class="article-sidebar">
        <TableOfContents :content="article.content" />
      </aside>
    </div>
  </div>
</template>
