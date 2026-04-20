<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getArticles, getTags } from '../api'
import ArticleCard from '../components/ArticleCard.vue'
import Pagination from '../components/Pagination.vue'

const route = useRoute()
const router = useRouter()
const articles = ref<any[]>([])
const tagName = ref('')
const total = ref(0)
const page = ref(1)
const size = 10
const loading = ref(true)

async function loadArticles() {
  loading.value = true
  try {
    const slug = route.params.slug as string
    // Find tag id by slug
    const tagsRes: any = await getTags()
    const allTags: any[] = tagsRes.data || []
    const tag = allTags.find((t: any) => t.slug === slug)
    if (!tag) {
      articles.value = []
      total.value = 0
      tagName.value = slug
      return
    }
    tagName.value = tag.name
    const res: any = await getArticles({ page: page.value, size, tag_id: tag.id })
    articles.value = res.data.list || []
    total.value = res.data.total
  } catch {
    articles.value = []
  } finally {
    loading.value = false
  }
}

function changePage(p: number) {
  page.value = p
  loadArticles()
}

function goBack() {
  router.back()
}

onMounted(loadArticles)

watch(() => route.params.slug, () => {
  page.value = 1
  loadArticles()
})
</script>

<template>
  <div class="container">
    <div class="page-header">
      <button class="back-btn" @click="goBack">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
        Back
      </button>
      <h1 class="page-title">Tag: {{ tagName }}</h1>
      <p class="page-info">{{ total }} articles</p>
    </div>

    <div v-if="loading" class="loading">Loading...</div>
    <template v-else>
      <div v-if="articles.length === 0" class="loading">No articles with this tag.</div>
      <div class="article-list">
        <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
      </div>
      <Pagination :page="page" :size="size" :total="total" @update:page="changePage" />
    </template>
  </div>
</template>
