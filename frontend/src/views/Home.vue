<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getArticles, getCategories, getTags } from '../api'
import { useSiteStore } from '../stores/site'
import ArticleCard from '../components/ArticleCard.vue'
import TagCloud from '../components/TagCloud.vue'
import Pagination from '../components/Pagination.vue'

const site = useSiteStore()
const articles = ref<any[]>([])
const categories = ref<any[]>([])
const tags = ref<any[]>([])
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

async function loadTags() {
  try {
    const res: any = await getTags()
    tags.value = res.data || []
  } catch {
    tags.value = []
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
  loadTags()
})
</script>

<template>
  <div>
  <!-- Hero Section -->
  <section class="hero">
    <div class="hero-content">
      <h1 class="hero-title">{{ site.settings.site_name || 'My Blog' }}</h1>
      <p class="hero-subtitle">{{ site.settings.site_description || 'A personal blog about coding and life' }}</p>
      <div class="hero-stats">
        <span class="hero-stat">{{ total }} articles</span>
        <span class="hero-stat-divider"></span>
        <span class="hero-stat">{{ categories.length }} categories</span>
        <span class="hero-stat-divider"></span>
        <span class="hero-stat">{{ tags.length }} tags</span>
      </div>
    </div>
  </section>

  <div class="container home-layout">
    <!-- Main content -->
    <div class="home-main">
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
        <div class="article-list">
          <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
        </div>
        <Pagination :page="page" :size="size" :total="total" @update:page="changePage" />
      </template>
    </div>

    <!-- Sidebar -->
    <aside class="home-sidebar">
      <TagCloud v-if="tags.length" :tags="tags" />
    </aside>
  </div>
  </div>
</template>
