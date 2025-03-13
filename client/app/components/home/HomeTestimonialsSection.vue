<script setup>
const { isLoggedIn } = defineProps({
  isLoggedIn: Boolean,
})

const emit = defineEmits(['login'])
</script>

<template>
  <section class="py-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-4xl font-bold dark:text-white text-gray-900 mb-4">
          Join Our Community
        </h2>
        <p class="text-lg dark:text-gray-300 text-gray-600 max-w-2xl mx-auto">
          Connect with developers from around the world
        </p>
      </div>

      <div class="text-center">
        <p class="text-4xl font-bold mb-6 dark:text-white text-gray-900">
          <HomeCountUp :end-val="1200" suffix="+" />
        </p>
        <p class="text-lg dark:text-gray-300 text-gray-600 mb-8">
          Developers sharing code every day
        </p>

        <UButton
          icon="tabler:login"
          class="text-md md:text-lg font-mono tracking-wide p-3 px-6 rounded-xl border-2 border-neutral-700/50 hover:border-green-500 dark:bg-mybg dark:text-white"
          @click="emit('login')"
        >
          {{ isLoggedIn ? "Go to Dashboard" : "Join Now" }}
        </UButton>
      </div>
    </div>
  </section>
</template>
