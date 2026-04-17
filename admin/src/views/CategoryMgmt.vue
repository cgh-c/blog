<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCategories, createCategory, updateCategory, deleteCategory } from '../api'
import { ElMessage, ElMessageBox } from 'element-plus'

const categories = ref<any[]>([])
const loading = ref(true)
const dialogVisible = ref(false)
const editingId = ref<number | null>(null)
const form = ref({ name: '', description: '', sort_order: 0 })

async function loadCategories() {
  loading.value = true
  try {
    const res: any = await getCategories()
    categories.value = res.data || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  form.value = { name: '', description: '', sort_order: 0 }
  dialogVisible.value = true
}

function openEdit(cat: any) {
  editingId.value = cat.id
  form.value = { name: cat.name, description: cat.description || '', sort_order: cat.sort_order }
  dialogVisible.value = true
}

async function handleSave() {
  if (!form.value.name.trim()) {
    ElMessage.warning('Name is required')
    return
  }
  try {
    if (editingId.value) {
      await updateCategory(editingId.value, form.value)
      ElMessage.success('Updated')
    } else {
      await createCategory(form.value)
      ElMessage.success('Created')
    }
    dialogVisible.value = false
    loadCategories()
  } catch { /* handled */ }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('Delete this category?', 'Confirm')
    await deleteCategory(id)
    ElMessage.success('Deleted')
    loadCategories()
  } catch { /* cancelled */ }
}

onMounted(loadCategories)
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
      <h2>Categories</h2>
      <el-button type="primary" @click="openCreate">New Category</el-button>
    </div>

    <el-table :data="categories" v-loading="loading" stripe>
      <el-table-column prop="name" label="Name" />
      <el-table-column prop="slug" label="Slug" />
      <el-table-column prop="description" label="Description" />
      <el-table-column prop="sort_order" label="Sort" width="80" />
      <el-table-column label="Actions" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">Edit</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">Delete</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? 'Edit Category' : 'New Category'" width="500">
      <el-form label-position="top">
        <el-form-item label="Name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="Description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="Sort Order">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="handleSave">Save</el-button>
      </template>
    </el-dialog>
  </div>
</template>
