import { ref } from 'vue'
import type { Features } from '~/types/features'

/**
 * Composable to fetch and manage feature flags from the API
 * Usage: const { features, isLoading, refresh } = useFeatures()
 */
export const useFeatures = () => {
  const features = useState<Features>('features', () => ({
    authentication: {
      register: false
    }
  }))

  const isLoading = ref(false)
  const error = ref<Error | null>(null)

  /**
   * Fetch features from the server endpoint
   */
  const fetchFeatures = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch('/api/v1/features')

      if (!response.ok) {
        throw new Error(`Failed to fetch features: ${response.statusText}`)
      }

      const data = await response.json()

      // Update the features state with the fetched data
      features.value = { ...features.value, ...data.features }
    } catch (err) {
      console.error('Error fetching features:', err)
      error.value = err instanceof Error ? err : new Error(String(err))
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Refresh the features by fetching them again
   */
  const refresh = () => fetchFeatures()

  // Fetch features immediately when the composable is used
  fetchFeatures()

  return {
    features,
    isLoading,
    error,
    refresh
  }
}
