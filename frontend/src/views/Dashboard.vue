<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
        <div class="flex items-center space-x-4">
          <button class="p-2 rounded-full text-gray-500 hover:text-gray-600 hover:bg-gray-100">
            <BellIcon class="h-6 w-6" />
          </button>
          <div class="relative">
            <button @click="profileMenuOpen = !profileMenuOpen" class="flex items-center space-x-2">
              <img class="h-8 w-8 rounded-full" :src="user.avatar" alt="User profile">
              <span class="hidden md:inline text-sm font-medium text-gray-700">{{ user.name }}</span>
              <ChevronDownIcon class="h-5 w-5 text-gray-500" />
            </button>
            
            <transition
              enter-active-class="transition ease-out duration-100"
              enter-from-class="transform opacity-0 scale-95"
              enter-to-class="transform opacity-100 scale-100"
              leave-active-class="transition ease-in duration-75"
              leave-from-class="transform opacity-100 scale-100"
              leave-to-class="transform opacity-0 scale-95"
            >
              <div v-show="profileMenuOpen" class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white ring-1 ring-black ring-opacity-5 z-10">
                <router-link to="/profile" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Your Profile</router-link>
                <router-link to="/settings" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Settings</router-link>
                <button @click="logout" class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Sign out</button>
              </div>
            </transition>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-8">
        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-indigo-500 rounded-md p-3">
                <BookOpenIcon class="h-6 w-6 text-white" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 truncate">Total Books</dt>
                <dd class="flex items-baseline">
                  <div class="text-2xl font-semibold text-gray-900">{{ stats.totalBooks }}</div>
                </dd>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-green-500 rounded-md p-3">
                <UsersIcon class="h-6 w-6 text-white" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 truncate">Total Members</dt>
                <dd class="flex items-baseline">
                  <div class="text-2xl font-semibold text-gray-900">{{ stats.totalMembers }}</div>
                </dd>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-yellow-500 rounded-md p-3">
                <ClockIcon class="h-6 w-6 text-white" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 truncate">Borrowed Books</dt>
                <dd class="flex items-baseline">
                  <div class="text-2xl font-semibold text-gray-900">{{ stats.borrowedBooks }}</div>
                </dd>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-red-500 rounded-md p-3">
                <ExclamationIcon class="h-6 w-6 text-white" />
              </div>
              <div class="ml-5 w-0 flex-1">
                <dt class="text-sm font-medium text-gray-500 truncate">Overdue Books</dt>
                <dd class="flex items-baseline">
                  <div class="text-2xl font-semibold text-gray-900">{{ stats.overdueBooks }}</div>
                </dd>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity and Quick Actions -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Recent Activity -->
        <div class="lg:col-span-2">
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <div class="px-4 py-5 sm:px-6 border-b border-gray-200">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Recent Activity</h3>
            </div>
            <div class="divide-y divide-gray-200">
              <div v-for="activity in recentActivities" :key="activity.id" class="px-4 py-4 sm:px-6">
                <div class="flex items-center">
                  <div class="flex-shrink-0">
                    <img class="h-10 w-10 rounded-full" :src="activity.user.avatar" :alt="activity.user.name">
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ activity.user.name }}</div>
                    <div class="text-sm text-gray-500">{{ activity.description }}</div>
                  </div>
                  <div class="ml-auto text-sm text-gray-500">
                    <time :datetime="activity.datetime">{{ activity.time }}</time>
                  </div>
                </div>
              </div>
            </div>
            <div class="px-4 py-4 sm:px-6 border-t border-gray-200">
              <router-link to="/activity" class="text-sm font-medium text-indigo-600 hover:text-indigo-500">
                View all activity
              </router-link>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div>
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <div class="px-4 py-5 sm:px-6 border-b border-gray-200">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Quick Actions</h3>
            </div>
            <div class="px-4 py-5 sm:p-6">
              <div class="space-y-4">
                <router-link
                  to="/books/add"
                  class="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  <PlusIcon class="-ml-1 mr-2 h-5 w-5" />
                  Add New Book
                </router-link>
                
                <router-link
                  to="/members/add"
                  class="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                >
                  <UserAddIcon class="-ml-1 mr-2 h-5 w-5" />
                  Add New Member
                </router-link>
                
                <router-link
                  to="/borrow"
                  class="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
                >
                  <BookmarkIcon class="-ml-1 mr-2 h-5 w-5" />
                  Borrow Book
                </router-link>
                
                <router-link
                  to="/return"
                  class="w-full flex items-center justify-center px-4 py-3 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  <ReplyIcon class="-ml-1 mr-2 h-5 w-5" />
                  Return Book
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import {
  BellIcon,
  ChevronDownIcon,
  BookOpenIcon,
  UsersIcon,
  ClockIcon,
  ExclamationIcon,
  PlusIcon,
  UserAddIcon,
  BookmarkIcon,
  ReplyIcon
} from '@heroicons/vue/outline'

export default {
  components: {
    BellIcon,
    ChevronDownIcon,
    BookOpenIcon,
    UsersIcon,
    ClockIcon,
    ExclamationIcon,
    PlusIcon,
    UserAddIcon,
    BookmarkIcon,
    ReplyIcon
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    const profileMenuOpen = ref(false)
    
    const user = {
      name: 'John Doe',
      avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
      role: 'Librarian'
    }
    
    const stats = {
      totalBooks: 1243,
      totalMembers: 342,
      borrowedBooks: 87,
      overdueBooks: 12
    }
    
    const recentActivities = [
      {
        id: 1,
        user: {
          name: 'Jane Smith',
          avatar: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80'
        },
        description: 'Borrowed "The Great Gatsby"',
        datetime: '2023-05-01T15:23',
        time: '2 hours ago'
      },
      {
        id: 2,
        user: {
          name: 'Michael Johnson',
          avatar: 'https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80'
        },
        description: 'Returned "To Kill a Mockingbird"',
        datetime: '2023-05-01T12:45',
        time: '5 hours ago'
      },
      {
        id: 3,
        user: {
          name: 'Sarah Williams',
          avatar: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80'
        },
        description: 'Registered as new member',
        datetime: '2023-05-01T09:30',
        time: '8 hours ago'
      }
    ]
    
    const logout = async () => {
      await store.dispatch('auth/logout')
      router.push('/login')
    }
    
    return {
      user,
      stats,
      recentActivities,
      profileMenuOpen,
      logout
    }
  }
}
</script>