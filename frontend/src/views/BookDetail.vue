
## 5. Book Detail Page

```vue
<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="bg-white shadow rounded-lg overflow-hidden">
        <!-- Book Header -->
        <div class="bg-indigo-700 px-6 py-6 sm:px-8 sm:py-8 lg:px-10 lg:py-10">
          <div class="max-w-3xl mx-auto flex flex-col md:flex-row items-start">
            <div class="flex-shrink-0 mb-6 md:mb-0 md:mr-8">
              <img
                :src="book.coverImage || 'https://via.placeholder.com/300x450?text=No+Cover'"
                :alt="book.title"
                class="h-48 w-32 object-cover rounded-md shadow-lg sm:h-64 sm:w-48"
              >
            </div>
            
            <div class="flex-1 text-white">
              <h1 class="text-2xl font-bold sm:text-3xl">{{ book.title }}</h1>
              <p class="mt-2 text-lg sm:text-xl">{{ book.author }}</p>
              
              <div class="mt-4 flex flex-wrap items-center gap-2">
                <span class="px-3 py-1 rounded-full bg-indigo-600 text-sm font-medium">
                  {{ book.genre }}
                </span>
                <span class="px-3 py-1 rounded-full bg-indigo-600 text-sm font-medium">
                  {{ book.publicationYear }}
                </span>
                <span class="px-3 py-1 rounded-full text-sm font-medium"
                      :class="book.availableCopies > 0 ? 'bg-green-500' : 'bg-red-500'">
                  {{ book.availableCopies > 0 ? `${book.availableCopies} Available` : 'Out of Stock' }}
                </span>
              </div>
              
              <div class="mt-6 flex space-x-4">
                <button
                  v-if="book.availableCopies > 0"
                  @click="handleBorrow"
                  class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-indigo-700 bg-white hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  <BookmarkIcon class="-ml-1 mr-2 h-5 w-5" />
                  Borrow Book
                </button>
                
                <button
                  v-if="isAdmin"
                  @click="handleEdit"
                  class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-yellow-500 hover:bg-yellow-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
                >
                  <PencilIcon class="-ml-1 mr-2 h-5 w-5" />
                  Edit Book
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Book Content -->
        <div class="px-6 py-8 sm:px-8 sm:py-10 lg:px-10 lg:py-12">
          <div class="max-w-3xl mx-auto">
            <div class="prose prose-indigo max-w-none">
              <h2 class="text-lg font-medium text-gray-900">About this book</h2>
              <p class="mt-4 text-gray-600">{{ book.description || 'No description available.' }}</p>
              
              <div class="mt-8 border-t border-gray-200 pt-8">
                <h2 class="text-lg font-medium text-gray-900">Details</h2>
                <div class="mt-4 grid grid-cols-1 gap-y-4 sm:grid-cols-2 sm:gap-x-6">
                  <div>
                    <p class="text-sm text-gray-500">ISBN</p>
                    <p class="mt-1 text-sm text-gray-900">{{ book.isbn }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500">Publisher</p>
                    <p class="mt-1 text-sm text-gray-900">{{ book.publisher }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500">Total Copies</p>
                    <p class="mt-1 text-sm text-gray-900">{{ book.totalCopies }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-500">Available Copies</p>
                    <p class="mt-1 text-sm text-gray-900">{{ book.availableCopies }}</p>
                  </div>
                </div>
              </div>
              
              <div v-if="similarBooks.length > 0" class="mt-12">
                <h2 class="text-lg font-medium text-gray-900">Similar Books</h2>
                <div class="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
                  <BookCard
                    v-for="similarBook in similarBooks"
                    :key="similarBook.id"
                    :book="similarBook"
                    class="cursor-pointer"
                    @click="$router.push(`/books/${similarBook.id}`)"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { BookmarkIcon, PencilIcon } from '@heroicons/vue/outline'
import BookCard from '@/components/BookCard.vue'

export default {
  components: {
    BookmarkIcon,
    PencilIcon,
    BookCard
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const store = useStore()
    const book = ref({})
    const similarBooks = ref([])
    const loading = ref(true)
    
    const isAdmin = computed(() => store.getters['auth/isAdmin'])
    
    const fetchBook = async () => {
      try {
        loading.value = true
        // In a real app, you would fetch from API using route.params.id
        await new Promise(resolve => setTimeout(resolve, 500))
        
        // Mock data
        book.value = {
          id: route.params.id,
          title: 'The Great Gatsby',
          author: 'F. Scott Fitzgerald',
          isbn: '978-0743273565',
          genre: 'Classic',
          publisher: 'Scribner',
          publicationYear: 1925,
          description: 'The Great Gatsby is a 1925 novel by American writer F. Scott Fitzgerald. Set in the Jazz Age on Long Island, the novel depicts narrator Nick Carraway\'s interactions with mysterious millionaire Jay Gatsby and Gatsby\'s obsession to reunite with his former lover, Daisy Buchanan.',
          totalCopies: 10,
          availableCopies: 4,
          coverImage: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80'
        }
        
        // Mock similar books
        similarBooks.value = [
          {
            id: 2,
            title: 'To Kill a Mockingbird',
            author: 'Harper Lee',
            genre: 'Classic',
            publicationYear: 1960,
            totalCopies: 8,
            availableCopies: 3,
            coverImage: 'https://images.unsplash.com/photo-1589998059171-988d887df646?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80'
          },
          {
            id: 3,
            title: '1984',
            author: 'George Orwell',
            genre: 'Dystopian',
            publicationYear: 1949,
            totalCopies: 7,
            availableCopies: 0,
            coverImage: 'https://images.unsplash.com/photo-1531346878377-a5be20888e57?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80'
          },
          {
            id: 4,
            title: 'Pride and Prejudice',
            author: 'Jane Austen',
            genre: 'Romance',
            publicationYear: 1813,
            totalCopies: 6,
            availableCopies: 2,
            coverImage: 'https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80'
          }
        ]
      } catch (error) {
        console.error('Failed to fetch book:', error)
      } finally {
        loading.value = false
      }
    }
    
    const handleBorrow = () => {
      router.push(`/borrow/${book.value.id}`)
    }
    
    const handleEdit = () => {
      router.push(`/books/edit/${book.value.id}`)
    }
    
    onMounted(fetchBook)
    
    return {
      book,
      similarBooks,
      loading,
      isAdmin,
      handleBorrow,
      handleEdit
    }
  }
}
</script>