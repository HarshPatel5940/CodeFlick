<script setup>
const isLoggedIn = ref(false)
const username = ref('')
const email = ref('')
const profile = useProfile()

async function checkSession() {
  const { data, error } = await useFetch('/api/session')
  if (error.value || !data.value) {
    console.error(error.value)
    return
  }

  const details = data.value.data.data
  isLoggedIn.value = details.user_id !== ''
  username.value = details.name
  email.value = details.email

  profile.set(
    details.user_id,
    details.name,
    details.email,
    details.isAdmin,
    details.isDeleted,
    details.isPremium,
  )
}
checkSession()

async function redirectToLogin() {
  const { data } = await useFetch('/api/login')
  if (data.value && data.value.status === 307) {
    navigateTo(data.value.data.redirectURI, { external: true })
  }
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
