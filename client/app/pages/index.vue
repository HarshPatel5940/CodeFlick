<script setup>
const { isLoggedIn, username, checkSession, redirectToLogin } = useAuth()

onMounted(async () => {
  try {
    await checkSession()
  }
  catch (error) {
    console.error('Error checking session:', error)
  }
})
</script>

<template>
  <main class="min-h-screen dark:bg-mybg bg-myLightBg">
    <div class="fixed top-5 right-5 z-10">
      <Theme />
    </div>

    <LazyHomeHeroSection :is-logged-in="isLoggedIn" :username="username" @login="redirectToLogin" />
    <HomeFeaturesSection />
    <HomeCodeExamplesSection />
    <LazyHomeTestimonialsSection :is-logged-in="isLoggedIn" @login="redirectToLogin" />
    <HomeFooterSection />
  </main>
</template>
