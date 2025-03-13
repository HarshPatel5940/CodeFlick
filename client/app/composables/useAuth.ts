import { ref } from 'vue'

export function useAuth() {
  const isLoggedIn = ref(false)
  const username = ref('')
  const email = ref('')
  const profile = useProfile()
  const toast = useToast()
  const BE = import.meta.env.VITE_BE_URL
  const isLoading = ref(false)

  async function fetchSession() {
    isLoading.value = true

    try {
      const response = await fetch(`${BE}/api/auth/session`, {
        credentials: 'include',
      })

      if (!response.ok) {
        throw new Error('Failed to fetch session')
      }

      const body = await response.json()

      return {
        status: response.status,
        data: body.data,
      }
    }
    catch (e) {
      return {
        status: 500,
        error: e,
      }
    }
    finally {
      isLoading.value = false
    }
  }

  async function checkSession() {
    const { status, data, error } = await fetchSession()

    if (status !== 200) {
      console.error(error)
      return false
    }

    isLoggedIn.value = data.user_id !== ''
    username.value = data.name
    email.value = data.email

    if (isLoggedIn.value) {
      profile.set(
        data.user_id,
        data.name,
        data.email,
        data.isAdmin,
        data.isDeleted,
        data.isPremium,
      )

      toast.add({
        title: 'Welcome back!',
        description: `You are now logged in as ${data.name}`,
        color: 'green',
        timeout: 2000,
      })
    }

    return isLoggedIn.value
  }

  async function signOut() {
    toast.add({
      title: 'Signing out!',
      color: 'orange',
      timeout: 2000,
    })

    try {
      const response = await fetch(`${BE}/api/auth/logout`, {
        method: 'POST',
        credentials: 'include',
      })

      const body = await response.json()

      if (response.ok) {
        profile.$reset()
        isLoggedIn.value = false
        username.value = ''
        email.value = ''

        toast.add({
          title: 'Signed out!',
          color: 'green',
          timeout: 2000,
        })

        return { success: true }
      }

      throw new Error(body.message || 'Failed to sign out')
    }
    catch (error) {
      toast.add({
        title: 'Failed to sign out!',
        color: 'red',
        timeout: 2000,
      })

      return {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      }
    }
  }

  async function redirectToLogin() {
    if (isLoggedIn.value) {
      navigateTo('/dashboard')
      return
    }

    await navigateTo(`${BE}/api/auth/google/login?r=client`, { external: true })
  }

  return {
    isLoggedIn,
    username,
    email,
    isLoading,
    checkSession,
    signOut,
    redirectToLogin,
  }
}
