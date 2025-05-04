import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Auth/LoginPage.vue'
import Dashboard from '@/views/Dashboard.vue'
import BookList from '@/views/BookCard.vue'
import BookDetail from '@/views/BookDetail.vue'
import BorrowForm from '@/views/BorrowBook.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/books',
    name: 'BookList',
    component: BookList,
    meta: { requiresAuth: true }
  },
  {
    path: '/books/:id',
    name: 'BookDetail',
    component: BookDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/borrow/:id',
    name: 'BorrowForm',
    component: BorrowForm,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router