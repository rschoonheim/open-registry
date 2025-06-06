import { ref } from 'vue'
import type { User, LoginCredentials, LoginResponse } from '~/types/auth'

export const useAuth = () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isAuthenticated = ref(false)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  /**
   * Login user with credentials
   */
  const login = async (credentials: LoginCredentials) => {
    isLoading.value = true
    error.value = null


    try {
      // Call the server endpoint created in /server/api/v1/auth/login.post.ts
      const response = await $fetch<LoginResponse>('/api/v1/auth/login', {
        method: 'POST',
        body: credentials
      })

      // Store user information and token
      user.value = response.user
      token.value = response.token
      isAuthenticated.value = true

      // Store token in localStorage for persistence
      if (process.client) {
        localStorage.setItem('auth-token', response.token)
      }

      return { success: true, user: response.user }
    } catch (err: any) {
      error.value = err.statusMessage || 'Authentication failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Logout the current user
   */
  const logout = () => {
    user.value = null
    token.value = null
    isAuthenticated.value = false

    // Remove token from localStorage
    if (process.client) {
      localStorage.removeItem('auth-token')
    }
  }

  /**
   * Initialize auth state from localStorage on client side
   */
  const initAuth = async () => {
    if (process.client) {
      const savedToken = localStorage.getItem('auth-token')

      if (savedToken) {
        token.value = savedToken
        // Here you could add a function to validate the token with your API
        // and fetch the current user information

        // For now, just set as authenticated
        isAuthenticated.value = true
      }
    }
  }

  // Initialize auth state when composable is used
  if (process.client) {
    initAuth()
  }

  return {
    user,
    token,
    isAuthenticated,
    isLoading,
    error,
    login,
    logout,
    initAuth
  }
}
