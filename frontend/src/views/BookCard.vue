<template>
    <div class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
      <div class="relative pb-[150%] bg-gray-100">
        <img
          :src="book.coverImage || 'https://via.placeholder.com/300x450?text=No+Cover'"
          :alt="book.title"
          class="absolute h-full w-full object-cover"
        >
        <div v-if="book.availableCopies === 0" class="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center">
<span class="text-white font-bold bg-red-500 px-2 py-1 rounded">Out of Stock</span>
</div>
</div>
<div class="p-4">
  <h3 class="text-lg font-semibold text-gray-900 mb-1 truncate">{{ book.title }}</h3>
  <p class="text-sm text-gray-600 mb-2">{{ book.author }}</p>
  
  <div class="flex items-center text-sm text-gray-500 mb-3">
    <span class="mr-2">{{ book.genre }}</span>
    <span class="mx-2">•</span>
    <span>{{ book.publicationYear }}</span>
  </div>
  
  <div class="flex justify-between items-center mb-3">
    <div class="text-sm">
      <span class="font-medium text-gray-900">{{ book.availableCopies }}</span>
      <span class="text-gray-500"> of {{ book.totalCopies }} available</span>
    </div>
    <span class="text-xs font-medium px-2 py-1 rounded-full" 
          :class="book.availableCopies > 0 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
      {{ book.availableCopies > 0 ? 'Available' : 'Unavailable' }}
    </span>
  </div>
  
  <div class="flex space-x-2">
    <button
      @click="$emit('view', book.id)"
      class="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-800 py-2 px-3 rounded text-sm font-medium transition"
    >
      View
    </button>
    
    <button
      v-if="book.availableCopies > 0"
      @click="$emit('borrow', book.id)"
      class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white py-2 px-3 rounded text-sm font-medium transition"
    >
      Borrow
    </button>
    
    <div v-else class="flex-1"></div>
  </div>
  
  <div v-if="isAdmin" class="flex space-x-2 mt-2">
    <button
      @click="$emit('edit', book.id)"
      class="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white py-2 px-3 rounded text-sm font-medium transition"
    >
      Edit
    </button>
    <button
      @click="$emit('delete', book.id)"
      class="flex-1 bg-red-500 hover:bg-red-600 text-white py-2 px-3 rounded text-sm font-medium transition"
    >
      Delete
    </button>
  </div>
</div>
</div> </template><script> import { computed } from 'vue' import { useStore } from 'vuex' export default { props: { book: { type: Object, required: true } }, setup() { const store = useStore() const isAdmin = computed(() => store.getters['auth/isAdmin']) return { isAdmin } } } </script>