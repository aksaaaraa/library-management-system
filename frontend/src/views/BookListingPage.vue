<template>
    <div class="min-h-screen bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="sm:flex sm:items-center sm:justify-between mb-8">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Book Catalog</h1>
            <p class="mt-2 text-sm text-gray-700">Browse our collection of {{ totalBooks }} books</p>
          </div>
          <div class="mt-4 sm:mt-0">
            <router-link
              to="/books/add"
              class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <PlusIcon class="-ml-1 mr-2 h-5 w-5" />
              Add Book
            </router-link>
          </div>
        </div>
        
        <!-- Search and Filters -->
        <div class="mb-8 bg-white shadow rounded-lg p-4">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label for="search" class="block text-sm font-medium text-gray-700 mb-1">Search</label>
              <div class="relative rounded-md shadow-sm">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <SearchIcon class="h-5 w-5 text-gray-400" />
                </div>
                <input
                  id="search"
                  v-model="searchQuery"
                  type="text"
                  class="focus:ring-indigo-500 focus:border-indigo-500 block w-full pl-10 sm:text-sm border-gray-300 rounded-md py-2 px-3 border"
                  placeholder="Title, author, ISBN..."
                  @input="handleSearch"
                >
              </div>
            </div>
            
            <div>
              <label for="genre" class="block text-sm font-medium text-gray-700 mb-1">Genre</label>
              <select
                id="genre"
                v-model="selectedGenre"
                class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                @change="fetchBooks"
              >
                <option value="">All Genres</option>
                <option v-for="genre in genres" :key="genre" :value="genre">{{ genre }}</option>
              </select>
            </div>
            
            <div>
              <label for="availability" class="block text-sm font-medium text-gray-700 mb-1">Availability</label>
              <select
                id="availability"
                v-model="availabilityFilter"
                class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                @change="fetchBooks"
              >
                <option value="all">All Books</option>
                <option value="available">Available Only</option>
                <option value="unavailable">Unavailable Only</option>
              </select>
            </div>
          </div>
        </div>
        
        <!-- Book Grid -->
        <div v-if="loading" class="flex justify-center py-12">
          <Spinner class="h-12 w-12 text-indigo-500" />
        </div>
        
        <div v-else>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <BookCard
              v-for="book in books"
              :key="book.id"
              :book="book"
              @borrow="handleBorrow"
              @view="handleView"
              @edit="handleEdit"
              @delete="handleDelete"
            />
          </div>
          
          <!-- Pagination -->
          <div class="mt-8 flex items-center justify-between">
            <div class="text-sm text-gray-700">
              Showing <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span> to 
              <span class="font-medium">{{ Math.min(currentPage * pageSize, totalBooks) }}</span> of 
              <span class="font-medium">{{ totalBooks }}</span> books
            </div>
            <nav class="flex space-x-2">
              <button
                @click="prevPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-3 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronLeftIcon class="h-5 w-5" />
              </button>
              
              <template v-for="page in visiblePages" :key="page">
                <button
                  @click="goToPage(page)"
                  :class="[page === currentPage ? 'bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50', 'relative inline-flex items-center px-4 py-2 border text-sm font-medium']"
                >
                  {{ page }}
                </button>
              </template>
              
              <button
                @click="nextPage"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-3 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <ChevronRightIcon class="h-5 w-5" />
              </button>
            </nav>
          </div>
        </div>
      </div>
      
      <!-- Delete Confirmation Modal -->
      <ConfirmationModal
        :open="deleteModalOpen"
        title="Delete Book"
        message="Are you sure you want to delete this book? This action cannot be undone."
        confirmText="Delete"
        cancelText="Cancel"
        confirmColor="red"
        @confirm="confirmDelete"
        @cancel="deleteModalOpen = false"
      />
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import {
    PlusIcon,
    SearchIcon,
    ChevronLeftIcon,
    ChevronRightIcon
  } from '@heroicons/vue/outline'
  import BookCard from '@/components/BookCard.vue'
  import Spinner from '@/components/Spinner.vue'
  import ConfirmationModal from '@/components/ConfirmationModal.vue'
  
  export default {
    components: {
      PlusIcon,
      SearchIcon,
      ChevronLeftIcon,
      ChevronRightIcon,
      BookCard,
      Spinner,
      ConfirmationModal
    },
    setup() {
      const router = useRouter()
      const books = ref([])
      const loading = ref(true)
      const totalBooks = ref(0)
      const currentPage = ref(1)
      const pageSize = 12
      const searchQuery = ref('')
      const selectedGenre = ref('')
      const availabilityFilter = ref('all')
      const deleteModalOpen = ref(false)
      const bookToDelete = ref(null)
      
      const genres = ref([
        'Fiction',
        'Non-Fiction',
        'Science Fiction',
        'Fantasy',
        'Mystery',
        'Romance',
        'Biography',
        'History',
        'Self-Help',
        'Science'
      ])
      
      const totalPages = computed(() => Math.ceil(totalBooks.value / pageSize))
      
      const visiblePages = computed(() => {
        const pages = []
        const maxVisible = 5
        let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
        let end = Math.min(totalPages.value, start + maxVisible - 1)
        
        if (end - start + 1 < maxVisible) {
          start = Math.max(1, end - maxVisible + 1)
        }
        
        for (let i = start; i <= end; i++) {
          pages.push(i)
        }
        
        return pages
      })
      
      const fetchBooks = async () => {
        try {
          loading.value = true
          // Simulate API call
          await new Promise(resolve => setTimeout(resolve, 800))
          
          // Mock data - in real app you would fetch from API
          const mockBooks = Array.from({ length: 36 }, (_, i) => ({
            id: i + 1,
            title: `Book Title ${i + 1}`,
            author: `Author ${String.fromCharCode(65 + (i % 26))}`,
            isbn: `978-${Math.floor(100000 + Math.random() * 900000)}-${Math.floor(10 + Math.random() * 90)}-${Math.floor(1 + Math.random() * 9)}`,
            genre: genres.value[i % genres.value.length],
            availableCopies: Math.floor(Math.random() * 5),
            totalCopies: Math.floor(3 + Math.random() * 5),
            coverImage: `https://source.unsplash.com/random/300x450/?book,cover&${i}`
          }))
          
          // Apply filters
          let filteredBooks = mockBooks
          
          if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            filteredBooks = filteredBooks.filter(book => 
              book.title.toLowerCase().includes(query) ||
              book.author.toLowerCase().includes(query) ||
              book.isbn.toLowerCase().includes(query)
          }
          
          if (selectedGenre.value) {
            filteredBooks = filteredBooks.filter(book => book.genre === selectedGenre.value)
          }
          
          if (availabilityFilter.value === 'available') {
            filteredBooks = filteredBooks.filter(book => book.availableCopies > 0)
          } else if (availabilityFilter.value === 'unavailable') {
            filteredBooks = filteredBooks.filter(book => book.availableCopies === 0)
          }
          
          totalBooks.value = filteredBooks.length
          
          // Apply pagination
          const start = (currentPage.value - 1) * pageSize
          books.value = filteredBooks.slice(start, start + pageSize)
        } catch (error) {
          console.error('Failed to fetch books:', error)
        } finally {
          loading.value = false
        }
      }
      
      const handleSearch = () => {
        currentPage.value = 1
        fetchBooks()
      }
      
      const prevPage = () => {
        if (currentPage.value > 1) {
          currentPage.value--
          fetchBooks()
        }
      }
      
      const nextPage = () => {
        if (currentPage.value < totalPages.value) {
          currentPage.value++
          fetchBooks()
        }
      }
      
      const goToPage = (page) => {
        currentPage.value = page
        fetchBooks()
      }
      
      const handleBorrow = (bookId) => {
        router.push(`/borrow/${bookId}`)
      }
      
      const handleView = (bookId) => {
        router.push(`/books/${bookId}`)
      }
      
      const handleEdit = (bookId) => {
        router.push(`/books/edit/${bookId}`)
      }
      
      const handleDelete = (bookId) => {
        bookToDelete.value = bookId
        deleteModalOpen.value = true
      }
      
      const confirmDelete = async () => {
        try {
          // Simulate API call to delete
          await new Promise(resolve => setTimeout(resolve, 500))
          console.log('Book deleted:', bookToDelete.value)
          fetchBooks()
        } catch (error) {
          console.error('Failed to delete book:', error)
        } finally {
          deleteModalOpen.value = false
          bookToDelete.value = null
        }
      }
      
      onMounted(fetchBooks)
      
      return {
        books,
        loading,
        totalBooks,
        currentPage,
        pageSize,
        totalPages,
        visiblePages,
        searchQuery,
        selectedGenre,
        availabilityFilter,
        genres,
        deleteModalOpen,
        bookToDelete,
        fetchBooks,
        handleSearch,
        prevPage,
        nextPage,
        goToPage,
        handleBorrow,
        handleView,
        handleEdit,
        handleDelete,
        confirmDelete
      }
    }
  }
  </script>