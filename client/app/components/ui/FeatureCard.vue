<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
  icon: {
    type: String,
    default: 'tabler:code',
  },
})
</script>

<template>
  <div
    v-motion-slide-visible-bottom
    class="bg-white dark:bg-gray-800 p-6 rounded-xl shadow-md border border-gray-100 dark:border-gray-700 transition-transform hover:scale-105"
  >
    <div class="mb-4 text-green-500">
      <Icon :name="icon" class="w-10 h-10" />
    </div>
    <h3 class="text-xl font-semibold dark:text-white text-gray-900 mb-3">
      {{ title }}
    </h3>
    <p class="dark:text-gray-400 text-gray-600">
      {{ description }}
    </p>
  </div>
</template>
