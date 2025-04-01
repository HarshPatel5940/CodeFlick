<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useGists } from '~/composables/useGists'
import type { Gist } from '~/types/gists'

const props = defineProps({
  fetchUrl: {
    type: String,
    required: true,
  },
  emptyMessage: {
    type: String,
    default: 'No gists found',
  },
  listHeight: {
    type: String,
    default: '65vh',
  },
})

const { gists, isLoading, error, fetchGists } = useGists()
const gistsList = ref<Gist[]>([])
const originalGists = ref<Gist[]>([])

async function loadGists() {
  const url = new URL(props.fetchUrl, window.location.origin)
  await fetchGists(Object.fromEntries(url.searchParams.entries()))
  gistsList.value = [...gists.value]
  originalGists.value = [...gists.value]
}

function addGist(newGist: Gist) {
  if (!props.fetchUrl.includes('fetchPublic') || (props.fetchUrl.includes('fetchPublic') && newGist.isPublic)) {
    gistsList.value = [newGist, ...gistsList.value]
    originalGists.value = [newGist, ...originalGists.value]
  }
}

function updateGist(updatedGist: Gist) {
  if (!props.fetchUrl.includes('fetchPublic')) {
    const index = gistsList.value.findIndex(g => g.fileId === updatedGist.fileId)
    if (index !== -1) {
      gistsList.value[index] = updatedGist

      const originalIndex = originalGists.value.findIndex(g => g.fileId === updatedGist.fileId)
      if (originalIndex !== -1) {
        originalGists.value[originalIndex] = updatedGist
      }
    }
  }

  else if (props.fetchUrl.includes('fetchPublic')) {
    const index = gistsList.value.findIndex(g => g.fileId === updatedGist.fileId)

    if (index !== -1 && !updatedGist.isPublic) {
      gistsList.value = gistsList.value.filter(g => g.fileId !== updatedGist.fileId)
      originalGists.value = originalGists.value.filter(g => g.fileId !== updatedGist.fileId)
    }

    else if (index !== -1 && updatedGist.isPublic) {
      gistsList.value[index] = updatedGist

      const originalIndex = originalGists.value.findIndex(g => g.fileId === updatedGist.fileId)
      if (originalIndex !== -1) {
        originalGists.value[originalIndex] = updatedGist
      }
    }

    else if (index === -1 && updatedGist.isPublic) {
      loadGists()
    }
  }
}

function searchGists(query: string) {
  if (!query) {
    gistsList.value = [...originalGists.value]
    return
  }

  const lowercaseQuery = query.toLowerCase()
  gistsList.value = originalGists.value.filter(gist =>
    gist.gistTitle.toLowerCase().includes(lowercaseQuery),
  )
}

defineExpose({ loadGists, addGist, updateGist, searchGists })

onMounted(loadGists)

watch(() => props.fetchUrl, loadGists, { immediate: false })
</script>

<template>
  <div
    class="flex flex-col gap-y-4 p-5 w-full overflow-y-auto bg-white/10 dark:bg-gray-800/20 rounded-lg shadow-sm"
    :style="{ maxHeight: listHeight }"
  >
    <div v-if="isLoading" class="flex justify-center items-center h-64">
      <div class="animate-spin w-10 h-10 border-t-2 border-b-2 border-blue-500 rounded-full" />
      <div class="ml-4 text-lg">
        Loading...
      </div>
    </div>

    <div v-else-if="error" class="p-4 bg-red-50 dark:bg-red-900/20 text-red-500 rounded-lg">
      {{ error }}
    </div>

    <div v-else-if="gistsList.length === 0" class="flex justify-center items-center h-64">
      <div class="flex flex-col items-center">
        <Icon name="heroicons:document" class="w-10 h-10 text-gray-500" />
        <div class="mt-4 text-lg font-semibold">
          {{ emptyMessage }}
        </div>
      </div>
    </div>

    <div v-else class="flex flex-col space-y-4">
      <div v-for="gist in gistsList" :key="gist.fileId" class="bg-white dark:bg-gray-800/75 rounded-lg p-2 transition-all hover:shadow-md hover:scale-[1.01]">
        <a :href="gist.shortUrl" class="w-full block">
          <UTooltip :text="gist.shortUrl" :popper="{ placement: 'top-start' }" class="w-full">
            <UCard class="w-full">
              <div class="font-medium">
                {{ gist.gistTitle }}
              </div>
              <template #footer>
                <div class="flex justify-between items-center">
                  <div class="flex items-center">
                    <Icon name="heroicons:clock" class="w-5 h-5 text-yellow-500" />
                    <div class="ml-1">
                      <NuxtTime :datetime="gist.updatedAt" relative />
                    </div>
                  </div>
                  <div class="flex items-center">
                    <Icon name="heroicons:eye" class="w-5 h-5 text-blue-500" />
                    <span class="ml-1">{{ gist.viewCount }}</span>
                  </div>
                </div>
              </template>
            </UCard>
          </UTooltip>
        </a>
      </div>
    </div>
  </div>
</template>
