<script setup lang="ts">
import { useSiteStore } from '../stores/site'
import { MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/preview.css'

const site = useSiteStore()
</script>

<template>
  <div class="container about-page">
    <img
      v-if="site.settings.avatar_url"
      :src="site.settings.avatar_url"
      alt="avatar"
      class="avatar"
    />
    <div class="article-body">
      <MdPreview :modelValue="site.settings.about_content || '# About\n\nWelcome to my blog!'" />
    </div>
  </div>
</template>
