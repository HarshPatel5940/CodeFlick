<script setup>
defineProps({
  type: {
    type: String,
    default: 'card',
    validator: val => ['card', 'list', 'text'].includes(val),
  },
  count: {
    type: Number,
    default: 1,
  },
})
</script>

<template>
  <div>
    <!-- Card skeleton -->
    <template v-if="type === 'card'">
      <div
        v-for="i in count"
        :key="i"
        class="bg-white dark:bg-gray-800 rounded-lg p-4 mb-4 animate-pulse"
      >
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-4" />
        <div class="h-20 bg-gray-200 dark:bg-gray-700 rounded mb-4" />
        <div class="flex justify-between">
          <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4" />
          <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4" />
        </div>
      </div>
    </template>

    <!-- List skeleton -->
    <template v-else-if="type === 'list'">
      <div
        v-for="i in count"
        :key="i"
        class="flex items-center py-3 animate-pulse"
      >
        <div class="h-10 w-10 bg-gray-200 dark:bg-gray-700 rounded-full mr-3" />
        <div class="flex-1">
          <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4 mb-2" />
          <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded w-1/2" />
        </div>
      </div>
    </template>

    <!-- Text skeleton -->
    <template v-else>
      <div class="animate-pulse">
        <div
          v-for="i in count"
          :key="i"
          class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full mb-2"
        />
      </div>
    </template>
  </div>
</template>
