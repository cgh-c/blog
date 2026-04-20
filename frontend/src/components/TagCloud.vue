<script setup lang="ts">
defineProps<{
  tags: { id: number; name: string; slug: string; article_count?: number }[]
}>()

function tagSize(count: number): string {
  if (count >= 5) return '1.15rem'
  if (count >= 3) return '1rem'
  return '0.85rem'
}
</script>

<template>
  <div class="tag-cloud-widget">
    <h3 class="widget-title">Tags</h3>
    <div class="tag-cloud">
      <RouterLink
        v-for="tag in tags"
        :key="tag.id"
        :to="`/tag/${tag.slug}`"
        class="tag-cloud-item"
        :style="{ fontSize: tagSize(tag.article_count || 0) }"
      >
        {{ tag.name }}
        <span v-if="tag.article_count" class="tag-cloud-count">{{ tag.article_count }}</span>
      </RouterLink>
    </div>
  </div>
</template>
