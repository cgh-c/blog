<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticles } from '../api'
import ArticleCard from '../components/ArticleCard.vue'
import Pagination from '../components/Pagination.vue'

const route = useRoute()
const articles = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = 10
const loading = ref(true)

async function loadArticles() {
  loading.value = true
  try {
    const categoryId = route.params.id
    const res: any = await getArticles({ page: page.value, size, category_id: categoryId })
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

onMounted(loadArticles)

watch(() => route.params.id, () => {
  page.value = 1
  loadArticles()
})
</script>

<template>
  <div class="container">
    <div v-if="loading" class="loading">Loading...</div>
    <template v-else>
      <div v-if="articles.length === 0" class="loading">No articles in this category.</div>
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
      <Pagination :page="page" :size="size" :total="total" @update:page="changePage" />
    </template>
  </div>
</template>
