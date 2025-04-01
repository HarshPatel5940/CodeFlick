<script setup lang="ts">
import type { Gist } from '~/types/gists'
import { onMounted, ref } from 'vue'

const route = useRoute()
const BE = import.meta.env.VITE_BE_URL
const gistId = route.params.id
// const queryGid = route.query.gid

const profile = useProfile()

const gist = ref<Gist>({
  fileId: '',
  userId: '',

  gistTitle: '',
  forkedFrom: '',
  shortUrl: '',

  viewCount: 0,
  isPublic: false,
  isDeleted: false,

  auditLog: [],

  createdAt: '',
  updatedAt: '',
})

const gistContent = ref('')
const isLoading = ref(true)
const error = ref<string | null>(null)

const hasChanges = ref(false)
const showNotification = ref(false)
const notificationMessage = ref('')
const notificationColor = ref('green')

async function handleOnClickCopy() {
  try {
    await navigator.clipboard.writeText(window.location.href)
    notificationMessage.value = 'Link copied to clipboard!'
    showNotification.value = true
    setTimeout(() => {
      showNotification.value = false
    }, 3000)
  }
  catch (err) {
    console.error('Failed to copy link:', err)
  }
}

async function handleCodeUpdate(newCode: string) {
  if (newCode !== gistContent.value) {
    hasChanges.value = true
    gistContent.value = newCode
  }
}

async function saveChanges() {
  if (!hasChanges.value)
    return
  isLoading.value = true

  try {
    const endpoint = `${BE}/api/gists/${gist.value.fileId}`

    const formData = new FormData()
    formData.append('file', gistContent.value)

    const response = await fetch(endpoint, {
      method: 'PUT',
      credentials: 'include',
      body: formData,
    })

    if (!response.ok) {
      throw new Error(`Failed to update gist: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || 'Failed to save changes')
    }

    hasChanges.value = false
    notificationMessage.value = 'Changes saved successfully!'
    notificationColor.value = 'green'
    showNotification.value = true
  }
  catch (err) {
    notificationMessage.value = err instanceof Error ? err.message : 'Failed to save changes'
    notificationColor.value = 'red'
    showNotification.value = true
    console.error('Error saving gist:', err)
  }
  finally {
    isLoading.value = false
    setTimeout(() => {
      showNotification.value = false
    }, 3000)
  }
}

async function fetchGistDetails(retry = false) {
  isLoading.value = true
  error.value = null

  try {
    const endpoint = `${BE}/api/gists/${gistId}`

    const response = await fetch(endpoint, {
      method: 'GET',
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Failed to fetch gist: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || 'Failed to load gist')
    }

    gist.value = {
      fileId: data.metadata.fileID,
      userId: data.metadata.userID,
      gistTitle: data.metadata.gistTitle,
      forkedFrom: data.metadata.forkedFrom,
      shortUrl: data.metadata.shortUrl,
      viewCount: data.metadata.viewCount,
      isPublic: data.metadata.isPublic,
      isDeleted: data.metadata.isDeleted,
      auditLog: data.metadata.auditLog,
      createdAt: data.metadata.createdAt,
      updatedAt: data.metadata.updatedAt,
    }
    if (!gist.value) {
      throw new Error('Failed to load gist')
    }
    gistContent.value = data.content

    if (gist.value?.gistTitle) {
      useHead({
        title: `${gist.value.gistTitle}`,
      })
    }

    console.log('Profile:', profile.data.UserID, gist.value)
  }
  catch (err) {
    if (!retry) {
      fetchGistDetails(true)
    }
    else {
      error.value = err instanceof Error ? err.message : 'An error occurred loading this gist'
      console.error('Error fetching gist:', error.value)
    }
  }
  finally {
    isLoading.value = false
  }
}
onMounted(() => {
  fetchGistDetails()
})
</script>

<template>
  <main class="w-screen min-h-screen dark:bg-mybg bg-myLightBg">
    <UNotification
      v-if="showNotification"
      :timeout="3000"
      :description="notificationMessage"
      class="dark:text-white"
      :color="notificationColor"
      :pause-timeout-on-hover="false"
    />
    <div class="flex flex-row pt-4 md:pt-5 justify-center">
      <Navbar />
    </div>

    <div class="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
      <div v-if="isLoading" class="flex justify-center py-12">
        <div class="flex flex-col items-center">
          <div class="animate-spin w-10 h-10 border-t-2 border-b-2 border-blue-500 rounded-full" />
          <div class="mt-4 text-lg font-semibold dark:text-white">
            Loading gist...
          </div>
        </div>
      </div>

      <ErrorDisplay
        v-else-if="error"
        :message="error"
        @retry="fetchGistDetails"
      />

      <EmptyState
        v-else-if="!gist"
        title="Gist Not Found"
        description="The gist you're looking for doesn't exist or has been removed."
        action-label="Go to Dashboard"
        @action="navigateTo('/dashboard')"
      />

      <div v-else>
        <div class="pt-10 mb-6 flex flex-col md:flex-row justify-between items-start md:items-center space-y-4 md:space-y-0">
          <div class="flex items-center gap-4">
            <UButton
              icon="i-heroicons-arrow-left"
              variant="ghost"
              class="font-medium rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
              @click="navigateTo('/dashboard')"
            />
            <h1 class="text-2xl md:text-3xl font-bold dark:text-white text-gray-900 transition-colors">
              {{ gist.gistTitle }}
            </h1>
            <div class="flex items-center space-x-6 text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800/40 px-4 py-2 rounded-full shadow-sm">
              /{{ gist.shortUrl }}
            </div>
          </div>
          <div class="flex flex-row gap-4">
            <div class="flex items-center space-x-6 text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800/40 px-4 py-2 rounded-full shadow-sm">
              <div>
                <span v-if="gist.isPublic" class="flex items-center">
                  <Icon name="heroicons:lock-open" class="w-4 h-4 mr-2" />
                  Public
                </span>
                <span v-else class="flex items-center">
                  <Icon name="heroicons:lock-closed" class="w-4 h-4 mr-2" />
                  Private
                </span>
              </div>
              <div class="h-4 w-px bg-gray-300 dark:bg-gray-600" />
              <span class="flex items-center">
                <Icon name="heroicons:eye" class="w-4 h-4 mr-2" />
                {{ gist.viewCount }} views
              </span>
              <div class="h-4 w-px bg-gray-300 dark:bg-gray-600" />
              <UTooltip :text="`${gist.updatedAt}`" :close-delay="750">
                <span class="flex flex-row gap-2 items-center">
                  <Icon name="heroicons:clock" class="w-4 h-4" />
                  Last Updated
                  <NuxtTime :datetime="gist.updatedAt" relative />
                </span>
              </UTooltip>
            </div>
            <UButton class="rounded-full opacity-75" variant="outline" :onclick="handleOnClickCopy">
              Share
            </UButton>
            <!-- TODO: Implement this edit button and also add switch visibility button. -->
            <UButton class="rounded-full opacity-75" variant="outline">
              Edit URL
            </UButton>
          </div>
        </div>

        <LazyCodePreview
          :title="gist.gistTitle"
          :filename="gist.fileId"
          :code="gistContent"
          lang="go"
          :show-edit-button="gist.userId === profile.data.UserID"
          :show-copy-button="true"
          @update:code="handleCodeUpdate"
        />
      </div>

      <div v-if="hasChanges" class="mt-4 flex justify-end">
        <UButton
          color="primary"
          :loading="isLoading"
          icon="heroicons:check"
          @click="saveChanges"
        >
          Save Changes
        </UButton>
      </div>
    </div>
  </main>
</template>
