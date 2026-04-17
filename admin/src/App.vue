<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { logout } from './api'
import { ElMessage } from 'element-plus'
import {
  DataAnalysis,
  Document,
  Folder,
  PriceTag,
  Setting,
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const isLoginPage = computed(() => route.name === 'Login')

async function handleLogout() {
  await logout()
  auth.clearAuth()
  ElMessage.success('Logged out')
  router.push('/login')
}
</script>

<template>
  <template v-if="isLoginPage">
    <RouterView />
  </template>

  <el-container v-else style="min-height: 100vh;">
    <el-aside width="200px" style="background: #304156;">
      <div style="padding: 20px; color: #fff; font-size: 18px; font-weight: bold; text-align: center;">
        Blog Admin
      </div>
      <el-menu
        :default-active="route.path"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/">
          <el-icon><DataAnalysis /></el-icon>
          <span>Dashboard</span>
        </el-menu-item>
        <el-menu-item index="/articles">
          <el-icon><Document /></el-icon>
          <span>Articles</span>
        </el-menu-item>
        <el-menu-item index="/categories">
          <el-icon><Folder /></el-icon>
          <span>Categories</span>
        </el-menu-item>
        <el-menu-item index="/tags">
          <el-icon><PriceTag /></el-icon>
          <span>Tags</span>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <span>Settings</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header style="display: flex; justify-content: flex-end; align-items: center; border-bottom: 1px solid #e6e6e6;">
        <span style="margin-right: 16px;">{{ auth.username }}</span>
        <el-button text @click="handleLogout">Logout</el-button>
      </el-header>
      <el-main>
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>
