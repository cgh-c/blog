<script setup lang="ts">
defineProps<{
  article: {
    id: number
    slug: string
    title: string
    summary: string
    category?: { id: number; name: string }
    tags?: { id: number; name: string }[]
    created_at: string
    view_count: number
  }
}>()

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <article class="article-card">
    <h2>
      <RouterLink :to="`/article/${article.slug}`">{{ article.title }}</RouterLink>
    </h2>
    <div class="article-meta">
      <span>{{ formatDate(article.created_at) }}</span>
      <span v-if="article.category">
        <RouterLink :to="`/category/${article.category.id}`">
          {{ article.category.name }}
        </RouterLink>
      </span>
      <span>{{ article.view_count }} views</span>
    </div>
    <p v-if="article.summary" class="article-summary">{{ article.summary }}</p>
    <div v-if="article.tags?.length" class="article-tags">
      <span v-for="tag in article.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
    </div>
  </article>
</template>
