<script setup>
const isLoggedIn = ref(false)
const username = ref('')
const email = ref('')
const profile = useProfile()
const toast = useToast()
const BE = import.meta.env.VITE_BE_URL

async function fetchSession() {
  try {
    const response = await fetch(`${BE}/api/auth/session`, {
      credentials: 'include',
    })
    if (!response.ok) {
      console.error('Failed to fetch session')
      return
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
}

async function checkSession() {
  const { status, data, error } = await fetchSession()
  if (status !== 200) {
    console.error(error.value)
    return
  }

  isLoggedIn.value = data.user_id !== ''
  username.value = data.name
  email.value = data.email

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
    type: 'success',
  })
}

onMounted(async () => {
  await checkSession()
})

async function redirectToLogin() {
  if (isLoggedIn.value) {
    navigateTo('/dashboard')
    return
  }
  await navigateTo(`${BE}/api/auth/google/login?r=client`, { external: true })
}
</script>

<template>
  <section
    class="bg-inherit w-full h-full justify-center overflow-auto items-center flex flex-col"
  >
    <div class="max-w-6xl mx-auto flex flex-col items-center gap-3">
      <Icon name="tabler:code" class="text-green-500 w-16 h-16" />
      <h1
        class="dark:text-gray-100 text-gray-900 text-7xl font-bold tracking-wider"
      >
        CodeFlick
      </h1>
      <div
        v-if="!isLoggedIn"
        class="dark:text-white/80 text-sm md:text-md italic"
      >
        To get started / login, click the button below
      </div>
      <div v-else class="dark:text-white/80 text-sm md:text-md italic">
        Welcome, {{ username }}
      </div>
      <div>
        <UButton
          v-motion-slide-visible-bottom
          icon="tabler:login"
          class="text-md md:text-lg font-mono tracking-wide p-2 rounded-xl border-2 border-neutral-700/50 hover:border-green-500 dark:bg-neutral-800/50 dark:text-white mt-4"
          @click="redirectToLogin()"
        >
          {{ isLoggedIn ? "Continue" : "Login" }}
        </UButton>
      </div>
    </div>
  </section>
</template>
