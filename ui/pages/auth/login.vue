<script setup lang="ts">
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import type { LoginCredentials } from '~/types/auth'
import { useFeatures } from "~/composables/useFeatures"

// Set the layout for this page
definePageMeta({
  layout: 'auth'
})

// Use our authentication composable
const { login, isLoading, error: authError } = useAuth()
const { features } = useFeatures()
const router = useRouter()

// Define validation schema with Yup
const validationSchema = yup.object({
  username: yup.string().required('Username is required'),
  password: yup.string().required('Password is required').min(6, 'Password must be at least 6 characters')
})

// Initialize form with VeeValidate
const { handleSubmit, errors, resetForm, defineField } = useForm({
  validationSchema
})

// Define fields with defineField
const [username, usernameAttrs] = defineField('username', '')
const [password, passwordAttrs] = defineField('password', '')

// Form submission handler
const onSubmit = handleSubmit(async (values) => {
  const credentials: LoginCredentials = {
    username: values.username,
    password: values.password
  }

  const result = await login(credentials)
  if (result.success) {
    // Redirect to dashboard after successful login
    router.push('/')
  }
})
</script>

<template>
  <div>
    <form @submit.prevent="onSubmit" class="mt-8 space-y-6">
      <div class="rounded-md shadow-sm space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <input
            id="username"
            v-model="username"
            v-bind="usernameAttrs"
            type="text"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Username"
          />
          <p v-if="errors.username" class="mt-1 text-sm text-red-600">{{ errors.username }}</p>
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            id="password"
            v-model="password"
            v-bind="passwordAttrs"
            type="password"
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Password"
          />
          <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
        </div>
      </div>

      <!-- Display error message if login fails -->
      <div v-if="authError" class="text-sm text-red-600 mt-2">
        {{ authError }}
      </div>

      <div>
        <button
          type="submit"
          :disabled="isLoading"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-75"
        >
          <span v-if="isLoading">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          </span>
          <span v-else>Sign in</span>
        </button>

        <button
          v-if="features.authentication.register"
          @click.prevent="router.push('/auth/register')"
          class="mt-2 w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-indigo-600 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Create an account
        </button>
      </div>
    </form>
  </div>
</template>

