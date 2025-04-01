import type { Gist } from '~/types/gists'
import { ref } from 'vue'

export function useGists() {
  const BE = import.meta.env.VITE_BE_URL
  const gists = ref<Gist[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function fetchGists(params = {}) {
    isLoading.value = true
    error.value = null

    try {
      const queryParams = new URLSearchParams()
      for (const [key, value] of Object.entries(params)) {
        queryParams.append(key, String(value))
      }

      const queryString = queryParams.toString()
      const url = `${BE}/api/gists${queryString ? `?${queryString}` : ''}`

      const response = await fetch(url, {
        method: 'GET',
        credentials: 'include',
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch gists: ${response.statusText}`)
      }

      const body = await response.json()

      if (!body.success) {
        throw new Error(body.message || 'Failed to fetch gists')
      }

      gists.value = body.data as Gist[]
      return gists.value
    }
    catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      return []
    }
    finally {
      isLoading.value = false
    }
  }

  async function fetchSingleGist(id: string) {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${BE}/api/gists/${id}`, {
        method: 'GET',
        credentials: 'include',
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch gist: ${response.statusText}`)
      }

      const body = await response.json()

      if (!body.success) {
        throw new Error(body.message || 'Failed to fetch gist')
      }

      return body.data as Gist
    }
    catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      throw err
    }
    finally {
      isLoading.value = false
    }
  }

  return {
    gists,
    isLoading,
    error,
    fetchGists,
    fetchSingleGist,
  }
}
