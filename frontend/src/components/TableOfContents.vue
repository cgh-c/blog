<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'

const props = defineProps<{ content: string }>()

interface TocItem {
  id: string
  text: string
  level: number
}

const tocItems = ref<TocItem[]>([])
const activeId = ref('')

function extractHeadings() {
  // Wait for md-editor to render, then extract headings from the article-body
  nextTick(() => {
    setTimeout(() => {
      const container = document.querySelector('.article-body')
      if (!container) return
      const headings = container.querySelectorAll('h1, h2, h3')
      const items: TocItem[] = []
      headings.forEach((el, index) => {
        const id = el.id || `heading-${index}`
        if (!el.id) el.id = id
        items.push({
          id,
          text: el.textContent || '',
          level: parseInt(el.tagName[1]),
        })
      })
      tocItems.value = items
      if (items.length > 0) {
        activeId.value = items[0].id
      }
    }, 500)
  })
}

function handleScroll() {
  const headings = tocItems.value
  if (!headings.length) return
  for (let i = headings.length - 1; i >= 0; i--) {
    const el = document.getElementById(headings[i].id)
    if (el && el.getBoundingClientRect().top <= 100) {
      activeId.value = headings[i].id
      return
    }
  }
  activeId.value = headings[0].id
}

function scrollTo(id: string) {
  const el = document.getElementById(id)
  if (el) {
    window.scrollTo({ top: el.offsetTop - 80, behavior: 'smooth' })
  }
}

watch(() => props.content, extractHeadings)
onMounted(() => {
  extractHeadings()
  window.addEventListener('scroll', handleScroll)
})
onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
  <nav v-if="tocItems.length" class="toc-widget">
    <h3 class="widget-title">Table of Contents</h3>
    <ul class="toc-list">
      <li
        v-for="item in tocItems"
        :key="item.id"
        :class="['toc-item', `toc-level-${item.level}`, { active: activeId === item.id }]"
      >
        <a @click.prevent="scrollTo(item.id)" href="#">{{ item.text }}</a>
      </li>
    </ul>
  </nav>
</template>
