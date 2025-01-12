<script setup lang="ts">
import type { Gist } from '../types/gists'

const myGists = ref<Gist[]>([])
const isResolved = ref(false)

async function fetchMyGists() {
  const BE = import.meta.env.VITE_BE_URL

  const response = await fetch(`${BE}/api/gists`, {
    method: 'GET',
    credentials: 'include',
  })

  const body = await response.json()

  return {
    status: response.status,
    bodyData: body,
  }
}

async function handleFetchMyGists() {
  const { status, bodyData } = await fetchMyGists()
  if (status !== 200) {
    isResolved.value = true
    console.error('Failed to fetch gists')
    return
  }

  isResolved.value = true
  myGists.value = bodyData.data as Gist[]
}

onMounted(async () => {
  await handleFetchMyGists()
})
</script>

<template>
  <div v-if="!isResolved">
    <div class="flex justify-center items-center h-96">
      <div class="flex flex-col items-center">
        <div class="animate-spin w-10 h-10 border-t-2 border-b-2 border-blue-500 rounded-full" />
        <div class="mt-4 text-lg font-semibold">
          Loading your gists...
        </div>
      </div>
    </div>
  </div>

  <div v-else>
    <div class="flex flex-col gap-y-4 p-5 md:px-10 w-full md:w-96 max-h-[85vh] overflow-y-scroll md:ml-28 my-8 md:my-14">
      <h2 class="text-md md:text-xl font-mono tracking-wider underline-offset-2 dark:no-underline">
        My Gists
      </h2>
      <div v-if="myGists.length === 0">
        <div class="flex justify-center items-center h-96">
          <div class="flex flex-col items-center">
            <div class="w-10 h-10">
              <Icon name="heroicons:document" class="w-10 h-10 text-gray-500" />
            </div>
            <div class="mt-4 text-lg font-semibold">
              No public gists found
            </div>
          </div>
        </div>
      </div>
      <div v-else class="flex flex-col space-y-4">
        <div v-for="gist in myGists" :key="gist.fieldId" class="bg-white dark:bg-gray-800/75 rounded-lg p-2">
          <a :href="gist.shortUrl ">
            <UCard>
              <template #header>
                {{ gist.gistTitle }}
              </template>
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
          </a>
        </div>
      </div>
      <UDivider class="pt-2 dark:pt-3 shadow-lg shadow-myLightBorder/30 dark:shadow-myborder/75" />
    </div>
  </div>
</template>
