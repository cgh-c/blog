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

// Generate a gradient based on article id for a visual accent
function cardAccent(id: number): string {
  const gradients = [
    'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    'linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)',
    'linear-gradient(135deg, #fccb90 0%, #d57eeb 100%)',
    'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
  ]
  return gradients[id % gradients.length]
}
</script>

<template>
  <article class="article-card">
    <div class="card-accent" :style="{ background: cardAccent(article.id) }"></div>
    <div class="card-body">
      <h2>
        <RouterLink :to="`/article/${article.slug}`">{{ article.title }}</RouterLink>
      </h2>
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
      <p v-if="article.summary" class="article-summary">{{ article.summary }}</p>
      <div class="card-footer">
        <div v-if="article.tags?.length" class="article-tags">
          <span v-for="tag in article.tags" :key="tag.id" class="tag">{{ tag.name }}</span>
        </div>
        <RouterLink :to="`/article/${article.slug}`" class="read-more">Read more &rarr;</RouterLink>
      </div>
    </div>
  </article>
</template>
