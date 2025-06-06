import { defineEventHandler } from 'h3'
import { useRuntimeConfig } from 'nuxt/app'
import type { Features, FeaturesResponse } from '~/types/features'

/**
 * Server endpoint that fetches enabled features from the API
 * This acts as a proxy to the backend /api/features endpoint
 */
export default defineEventHandler(async (event) => {
  try {
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:3000"
    const response = await fetch(`${apiUrl}/api/features`)

    if (!response.ok) {
      console.error('Failed to fetch features from API:', response.statusText)
      // Return default features if API call fails
      return {
        features: {
          authentication: {
            register: false // Default to disabled if API call fails
          }
        }
      }
    }

    // Parse and return the features
    const data = await response.json() as FeaturesResponse
    return data
  } catch (error) {
    console.error('Error fetching features:', error)

    // Return default features in case of error
    return {
      features: {
        authentication: {
          register: false // Default to disabled if there's an error
        }
      }
    }
  }
})
