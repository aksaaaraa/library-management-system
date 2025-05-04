<template>
    <div class="min-h-screen bg-gray-50 py-8">
      <div class="max-w-md mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white shadow rounded-lg overflow-hidden">
          <div class="bg-indigo-600 py-4 px-6">
            <h2 class="text-xl font-bold text-white">Borrow Book</h2>
          </div>
          
          <div class="p-6">
            <div class="flex items-center mb-6">
              <img
                :src="book.coverImage || 'https://via.placeholder.com/300x450?text=No+Cover'"
                :alt="book.title"
                class="h-24 w-16 object-cover rounded-md"
              >
              <div class="ml-4">
                <h3 class="text-lg font-medium text-gray-900">{{ book.title }}</h3>
                <p class="text-sm text-gray-600">{{ book.author }}</p>
                <p class="text-sm text-gray-600">Available: {{ book.availableCopies }} of {{ book.totalCopies }}</p>
              </div>
            </div>
            
            <form @submit.prevent="handleSubmit" class="space-y-6">
              <div>
                <label for="member" class="block text-sm font-medium text-gray-700 mb-1">Select Member</label>
                <select
                  id="member"
                  v-model="form.memberId"
                  required
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                >
                  <option value="" disabled>Select a member</option>
                  <option v-for="member in members" :key="member.id" :value="member.id">
                    {{ member.name }} ({{ member.email }})
                  </option>
                </select>
              </div>
              
              <div>
                <label for="dueDate" class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
                <input
                  id="dueDate"
                  v-model="form.dueDate"
                  type="date"
                  required
                  :min="minDueDate"
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                >
                <p class="mt-1 text-xs text-gray-500">Books are typically due in 14 days</p>
              </div>
              
              <div>
                <label for="notes" class="block text-sm font-medium text-gray-700 mb-1">Notes (Optional)</label>
                <textarea
                  id="notes"
                  v-model="form.notes"
                  rows="3"
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                ></textarea>
              </div>
              
              <div class="flex items-center">
                <input
                  id="terms"
                  v-model="form.agreeTerms"
                  type="checkbox"
                  required
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                >
                <label for="terms" class="ml-2 block text-sm text-gray-700">
                  I agree to the <a href="#" class="text-indigo-600 hover:text-indigo-500">borrowing terms</a>
                </label>
              </div>
              
              <div class="flex justify-end space-x-3">
                <button
                  type="button"
                  @click="$router.go(-1)"
                  class="inline-flex justify-center py-2 px-4 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="loading"
                  class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-70"
                >
                  <span v-if="!loading">Confirm Borrow</span>
                  <svg v-else class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  
  export default {
    setup() {
      const route = useRoute()
      const router = useRouter()
      const loading = ref(false)
      const book = ref({})
      const members = ref([])
      
      const minDueDate = computed(() => {
        const today = new Date()
        today.setDate(today.getDate() + 1)
        return today.toISOString().split('T')[0]
      })
      
      const form = ref({
        memberId: '',
        dueDate: '',
        notes: '',
        agreeTerms: false
      })
      
      // Set default due date to 14 days from now
      const setDefaultDueDate = () => {
        const today = new Date()
        today.setDate(today.getDate() + 14)
        form.value.dueDate = today.toISOString().split('T')[0]
      }
      
      const fetchData = async () => {
        try {
          // In a real app, you would fetch book and members from API
          await new Promise(resolve => setTimeout(resolve, 300))
          
          // Mock book data
          book.value = {
            id: route.params.id,
            title: 'The Great Gatsby',
            author: 'F. Scott Fitzgerald',
            availableCopies: 4,
            totalCopies: 10,
            coverImage: 'https://images.unsplash.com/photo-1544947950-fa07a98d237f?ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=80'
          }
          
          // Mock members data
          members.value = [
            { id: 1, name: 'John Doe', email: 'john@example.com' },
            { id: 2, name: 'Jane Smith', email: 'jane@example.com' },
            { id: 3, name: 'Michael Johnson', email: 'michael@example.com' }
          ]
          
          setDefaultDueDate()
        } catch (error) {
          console.error('Failed to fetch data:', error)
        }
      }
      
      const handleSubmit = async () => {
        try {
          loading.value = true
          // Simulate API call
          await new Promise(resolve => setTimeout(resolve, 800))
          
          // In a real app, you would submit the borrow request
          console.log('Borrow request submitted:', {
            bookId: book.value.id,
            ...form.value
          })
          
          // Show success message and redirect
          router.push({
            path: '/borrow/success',
            query: { bookId: book.value.id }
          })
        } catch (error) {
          console.error('Failed to submit borrow request:', error)
        } finally {
          loading.value = false
        }
      }
      
      onMounted(fetchData)
      
      return {
        book,
        members,
        form,
        loading,
        minDueDate,
        handleSubmit
      }
    }
  }
  </script>