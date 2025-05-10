<script setup lang="ts">
import type { Gist } from '~/types/gists'
import { nextTick, onMounted, ref, watch } from 'vue'

const route = useRoute()
const BE = import.meta.env.VITE_BE_URL
const gistId = route.params.id
const queryGid = route.query.gid as string

const profile = useProfile()
const toast = useToast()

const gist = ref<Gist>({
  fileId: '',
  userId: '',
  fileName: '',

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

const isEditUrlModalOpen = ref(false)
const newCustomUrl = ref('')
const isEditTitleModalOpen = ref(false)
const newTitle = ref('')

function openEditUrlModal() {
  newCustomUrl.value = gist.value.shortUrl
  isEditUrlModalOpen.value = true
}

function openEditTitleModal() {
  newTitle.value = gist.value.gistTitle
  isEditTitleModalOpen.value = true
}

async function updateGist(options: {
  updateType: string
  validationCheck?: () => { isValid: boolean, message?: string }
  formData?: Record<string, string | File>
  successCallback?: () => void
  errorPrefix?: string
  successMessage?: string
}) {
  const {
    updateType,
    validationCheck,
    formData = {},
    successCallback,
    errorPrefix = `Failed to update ${updateType}`,
    successMessage,
  } = options

  if (validationCheck) {
    const validation = validationCheck()
    if (!validation.isValid) {
      notificationMessage.value = validation.message || `Invalid ${updateType}`
      notificationColor.value = 'red'
      showNotification.value = true
      setTimeout(() => {
        showNotification.value = false
      }, 3000)
      return
    }
  }

  isLoading.value = true

  try {
    const endpoint = `${BE}/api/gists/${gist.value.fileId}`

    const requestFormData = new FormData()

    if (!('gist_title' in formData)) {
      requestFormData.append('gist_title', gist.value.gistTitle)
    }

    if (!('is_public' in formData)) {
      requestFormData.append('is_public', String(gist.value.isPublic))
    }

    for (const [key, value] of Object.entries(formData)) {
      if (value instanceof File) {
        requestFormData.append(key, value, value.name)
      }
      else {
        requestFormData.append(key, value)
      }
    }

    if (updateType !== 'visibility' && !('file' in formData)) {
      const file = new File([gistContent.value], gist.value.fileName, { type: 'text/plain' })
      requestFormData.append('file', file)
    }

    const response = await fetch(endpoint, {
      method: 'PUT',
      credentials: 'include',
      body: requestFormData,
    })

    if (!response.ok) {
      throw new Error(`${errorPrefix}: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || errorPrefix)
    }

    if (successCallback) {
      successCallback()
    }

    if (successMessage) {
      notificationMessage.value = successMessage
      notificationColor.value = 'green'
      showNotification.value = true
    }
  }
  catch (err) {
    notificationMessage.value = err instanceof Error ? err.message : errorPrefix
    notificationColor.value = 'red'
    showNotification.value = true
    console.error(`Error updating ${updateType}:`, err)
  }
  finally {
    isLoading.value = false
    setTimeout(() => {
      showNotification.value = false
    }, 3000)
  }
}

async function updateCustomUrl() {
  updateGist({
    updateType: 'url',
    validationCheck: () => ({
      isValid: !!newCustomUrl.value && newCustomUrl.value.trim() !== '',
      message: 'Custom URL cannot be empty',
    }),
    formData: {
      custom_url: newCustomUrl.value,
    },
    successCallback: () => {
      const oldShortUrl = gist.value.shortUrl
      gist.value.shortUrl = newCustomUrl.value
      isEditUrlModalOpen.value = false
      notificationMessage.value = `Custom URL updated from ${oldShortUrl} to ${newCustomUrl.value}`
    },
  })
}

async function updateTitle() {
  updateGist({
    updateType: 'title',
    validationCheck: () => ({
      isValid: !!newTitle.value && newTitle.value.trim() !== '',
      message: 'Title cannot be empty',
    }),
    formData: {
      gist_title: newTitle.value,
    },
    successCallback: () => {
      const oldTitle = gist.value.gistTitle
      gist.value.gistTitle = newTitle.value
      isEditTitleModalOpen.value = false

      useHead({
        title: `Codeflick | ${gist.value.gistTitle}`,
      })

      notificationMessage.value = `Title updated from "${oldTitle}" to "${newTitle.value}"`
    },
  })
}

async function toggleVisibility() {
  updateGist({
    updateType: 'visibility',
    formData: {
      is_public: String(!gist.value.isPublic),
    },
    successCallback: () => {
      gist.value.isPublic = !gist.value.isPublic
    },
    successMessage: `Gist is now ${!gist.value.isPublic ? 'public' : 'private'}`,
  })
}

async function handleOnClickCopy() {
  try {
    await navigator.clipboard.writeText(window.location.href)
    toast.add({
      title: 'Link Copied!',
      description: 'Link copied to clipboard',
      color: 'green',
      icon: 'heroicons:check-circle',
      timeout: 2000,
    })
  }
  catch (err) {
    console.error('Failed to copy link:', err)
    toast.add({
      title: 'Copy Failed',
      description: 'Failed to copy link to clipboard',
      color: 'red',
      icon: 'heroicons:exclamation-circle',
      timeout: 2000,
    })
  }
}

async function handleCodeUpdate(newCode: string) {
  if (newCode !== gistContent.value) {
    hasChanges.value = true
    gistContent.value = newCode

    await nextTick()
  }
}

async function saveChanges() {
  if (!hasChanges.value)
    return

  const file = new File([gistContent.value], gist.value.fileName, { type: 'text/plain' })

  updateGist({
    updateType: 'content',
    formData: {
      file,
    },
    successCallback: () => {
      hasChanges.value = false
    },
    successMessage: 'Changes saved successfully!',
    errorPrefix: 'Failed to save changes',
  })
}

async function fetchGistDetails(retry = false) {
  isLoading.value = true
  error.value = null

  try {
    const endpoint = queryGid
      ? `${BE}/api/gists/${gistId}?gid=${queryGid}`
      : `${BE}/api/gists/${gistId}`

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
      fileName: data.metadata.fileName,
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
        title: `Codeflick | ${gist.value.gistTitle}`,
      })
    }
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

watch(() => route.params.id, (newId, oldId) => {
  if (newId !== oldId) {
    fetchGistDetails()
  }
})

watch(() => route.query.gid, (newGid, oldGid) => {
  if (newGid !== oldGid) {
    fetchGistDetails()
  }
})

onMounted(() => {
  fetchGistDetails()
})
</script>

<template>
  <main class="w-full min-h-screen dark:bg-mybg bg-myLightBg overflow-x-hidden">
    <UModal v-model="isEditUrlModalOpen" :ui="{ width: 'sm:max-w-md' }">
      <div class="p-4">
        <h2 class="text-lg font-bold mb-4">
          Edit Custom URL
        </h2>
        <div class="mb-4">
          <UInput
            v-model="newCustomUrl"
            label="Custom URL"
            placeholder="Enter custom URL"
            class="w-full"
          />
        </div>
        <div class="flex justify-end gap-2">
          <UButton
            variant="ghost"
            color="gray"
            @click="isEditUrlModalOpen = false"
          >
            Cancel
          </UButton>
          <UButton
            variant="solid"
            color="primary"
            :loading="isLoading"
            @click="updateCustomUrl"
          >
            Update URL
          </UButton>
        </div>
      </div>
    </UModal>

    <UModal v-model="isEditTitleModalOpen" :ui="{ width: 'sm:max-w-md' }">
      <div class="p-4">
        <h2 class="text-lg font-bold mb-4">
          Edit Title
        </h2>
        <div class="mb-4">
          <UInput
            v-model="newTitle"
            label="Title"
            placeholder="Enter new title"
            class="w-full"
          />
        </div>
        <div class="flex justify-end gap-2">
          <UButton
            variant="ghost"
            color="gray"
            @click="isEditTitleModalOpen = false"
          >
            Cancel
          </UButton>
          <UButton
            variant="solid"
            color="primary"
            :loading="isLoading"
            @click="updateTitle"
          >
            Update Title
          </UButton>
        </div>
      </div>
    </UModal>

    <div class="flex flex-row pt-4 md:pt-5 justify-center">
      <Navbar />
    </div>

    <div class="max-w-7xl mx-auto px-4 pt-8 pb-4 sm:px-6 lg:px-8">
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
        <div class="flex flex-col md:flex-row justify-between items-start md:items-center space-y-4 md:space-y-0 w-full pt-5 md:pt-10">
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
          <div class="flex flex-row gap-4 flex-wrap">
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

            <UButton class="rounded-full opacity-75" icon="i-heroicons-share" variant="outline" @click="handleOnClickCopy">
              Share
            </UButton>
            <UDropdown
              v-if="gist.userId === profile.data.UserID"
              :items="[
                [
                  {
                    label: 'Edit Title',
                    icon: 'heroicons:link',
                    onClick: openEditTitleModal,
                  },
                  {
                    label: 'Edit URL',
                    icon: 'heroicons:link',
                    onClick: openEditUrlModal,
                  },
                  {
                    label: gist.isPublic ? 'Make Private' : 'Make Public',
                    icon: gist.isPublic ? 'heroicons:lock-closed' : 'heroicons:lock-open',
                    onClick: toggleVisibility,
                  },
                ],
              ]"
              :popper="{ placement: 'bottom-end' }"
            >
              <UButton
                class="rounded-full opacity-75"
                variant="outline"
                icon="i-heroicons-adjustments-horizontal"
              >
                Edit
              </UButton>

              <template #item="{ item }">
                <div :onclick="item.onClick" class="w-full flex flex-row justify-between">
                  <span class="truncate">{{ item.label }}</span>
                  <UIcon :name="item.icon" class="flex-shrink-0 h-4 w-4 text-gray-400 dark:text-gray-500 ms-auto ml-2" />
                </div>
              </template>
            </UDropdown>
          </div>
        </div>
        <LazyCodePreview
          :title="gist.gistTitle"
          :filename="gist.fileName"
          :code="gistContent"
          :lang="gist.fileName.split('.').pop() || 'go'"
          :show-edit-button="gist.userId === profile.data.UserID"
          :show-copy-button="true"
          @update:code="handleCodeUpdate"
        />
      </div>
    </div>

    <div v-if="hasChanges" class="max-w-7xl mx-auto px-4 mt-4 sm:px-6 lg:px-8 flex justify-end">
      <UButton
        color="primary"
        :loading="isLoading"
        icon="heroicons:check"
        class="mb-4"
        @click="saveChanges"
      >
        Save Changes
      </UButton>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <LazyGistReply v-if="gist && gist.fileId && !isLoading" :gist-id="gist.fileId" />
    </div>
  </main>
</template>
