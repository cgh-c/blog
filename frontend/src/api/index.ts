import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true,
})

// Response interceptor: unwrap the data envelope
api.interceptors.response.use(
  (res) => res.data,
  (err) => {
    const message = err.response?.data?.message || 'Network error'
    return Promise.reject(new Error(message))
  }
)

export interface PaginatedResponse<T> {
  code: number
  message: string
  data: {
    list: T[]
    total: number
    page: number
    size: number
  }
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// Public API
export const getArticles = (params?: Record<string, any>) =>
  api.get('/articles', { params })

export const getArticleBySlug = (slug: string) =>
  api.get(`/articles/${slug}`)

export const getCategories = () =>
  api.get('/categories')

export const getTags = () =>
  api.get('/tags')

export const getArchives = () =>
  api.get('/archives')

export const getSettings = () =>
  api.get('/settings')

export default api
