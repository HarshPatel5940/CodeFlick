<script setup>
defineProps({
  title: {
    type: String,
    default: 'Something went wrong',
  },
  message: {
    type: String,
    required: true,
  },
  retryLabel: {
    type: String,
    default: 'Try again',
  },
  showRetry: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['retry'])

function handleRetry() {
  emit('retry')
}
</script>

<template>
  <div class="bg-red-50 dark:bg-red-900/20 p-4 rounded-lg">
    <div class="flex">
      <div class="flex-shrink-0">
        <Icon name="heroicons:exclamation-circle" class="h-5 w-5 text-red-400" aria-hidden="true" />
      </div>
      <div class="ml-3">
        <h3 class="text-sm font-medium text-red-800 dark:text-red-200">
          {{ title }}
        </h3>
        <div class="mt-2 text-sm text-red-700 dark:text-red-300">
          <p>{{ message }}</p>
        </div>
        <div v-if="showRetry" class="mt-4">
          <div class="flex">
            <UButton
              size="xs"
              color="red"
              variant="soft"
              @click="handleRetry"
            >
              {{ retryLabel }}
            </UButton>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
