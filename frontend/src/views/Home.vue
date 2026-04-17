<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getArticles, getCategories } from '../api'
import ArticleCard from '../components/ArticleCard.vue'
import Pagination from '../components/Pagination.vue'

const articles = ref<any[]>([])
const categories = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const size = 10
const loading = ref(true)
const selectedCategory = ref<number | undefined>(undefined)

async function loadArticles() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, size }
    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }
    const res: any = await getArticles(params)
    articles.value = res.data.list || []
    total.value = res.data.total
  } catch {
    articles.value = []
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    const res: any = await getCategories()
    categories.value = res.data || []
  } catch {
    categories.value = []
  }
}

function selectCategory(id?: number) {
  selectedCategory.value = id
  page.value = 1
  loadArticles()
}

function changePage(p: number) {
  page.value = p
  loadArticles()
}

onMounted(() => {
  loadCategories()
  loadArticles()
})
</script>

<template>
  <div class="container">
    <div v-if="categories.length" class="category-list">
      <a
        class="category-item"
        :class="{ active: !selectedCategory }"
        @click.prevent="selectCategory(undefined)"
        href="#"
      >
        All
      </a>
      <a
        v-for="cat in categories"
        :key="cat.id"
        class="category-item"
        :class="{ active: selectedCategory === cat.id }"
        @click.prevent="selectCategory(cat.id)"
        href="#"
      >
        {{ cat.name }}
        <span class="category-count">({{ cat.article_count }})</span>
      </a>
    </div>

    <div v-if="loading" class="loading">Loading...</div>

    <template v-else>
      <div v-if="articles.length === 0" class="loading">No articles yet.</div>
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
      <Pagination :page="page" :size="size" :total="total" @update:page="changePage" />
    </template>
  </div>
</template>
