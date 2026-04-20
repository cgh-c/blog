import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
})

api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    const message = err.response?.data?.message || '网络错误'
    const url = err.config?.url || ''
    // Login endpoint 401 means wrong credentials — show error, don't redirect
    if (err.response?.status === 401 && !url.includes('/auth/login')) {
      window.location.href = '/admin/login'
    } else {
      ElMessage.error(message)
    }
    return Promise.reject(new Error(message))
  }
)

// Auth
export const login = (data: { username: string; password: string }) =>
  api.post('/auth/login', data)

export const logout = () => api.post('/auth/logout')

export const getMe = () => api.get('/auth/me')

// Articles
export const getArticles = (params?: Record<string, any>) =>
  api.get('/admin/articles', { params })

export const getArticle = (id: number) =>
  api.get(`/admin/articles/${id}`)

export const createArticle = (data: any) =>
  api.post('/admin/articles', data)

export const updateArticle = (id: number, data: any) =>
  api.put(`/admin/articles/${id}`, data)

export const deleteArticle = (id: number) =>
  api.delete(`/admin/articles/${id}`)

// Categories
export const getCategories = () => api.get('/admin/categories')
export const createCategory = (data: any) => api.post('/admin/categories', data)
export const updateCategory = (id: number, data: any) => api.put(`/admin/categories/${id}`, data)
export const deleteCategory = (id: number) => api.delete(`/admin/categories/${id}`)

// Tags
export const getTags = () => api.get('/admin/tags')
export const createTag = (data: { name: string }) => api.post('/admin/tags', data)
export const deleteTag = (id: number) => api.delete(`/admin/tags/${id}`)

// Settings
export const getSettings = () => api.get('/admin/settings')
export const updateSettings = (settings: Record<string, string>) =>
  api.put('/admin/settings', { settings })

// Upload
export const uploadFile = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/admin/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// Dashboard
export const getDashboard = () => api.get('/admin/dashboard')

export default api
