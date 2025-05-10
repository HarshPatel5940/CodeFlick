<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const props = defineProps({
  gistId: {
    type: String,
    required: true,
  },
})

const BE = import.meta.env.VITE_BE_URL
const replies = ref<any[]>([])
const newReply = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)
const profile = useProfile()
const toast = useToast()
const editReplyId = ref<string | null>(null)
const editReplyContent = ref('')
const currentUserName = computed(() => profile.data.name || 'You')

onMounted(() => {
  fetchReplies()
})

async function fetchReplies() {
  isLoading.value = true
  error.value = null

  try {
    const response = await fetch(`${BE}/api/gists/${props.gistId}/reply`, {
      method: 'GET',
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Failed to fetch replies: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || 'Failed to fetch replies')
    }

    replies.value = (data.data || []).map((reply: { replyId: any, userId: string, message: any, createdAt: any, updatedAt: any, userName: any, name: any }) => {
      return {
        ...reply,
        replyId: reply.replyId || `temp-${Date.now()}`,
        userId: reply.userId || '',
        message: reply.message || '',
        createdAt: reply.createdAt || new Date().toISOString(),
        updatedAt: reply.updatedAt || reply.createdAt || new Date().toISOString(),
        userName: reply.name || (reply.userId === profile.data.UserID ? currentUserName.value : 'User'),
      }
    })
  }
  catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load replies'
    console.error('Error fetching replies:', err)
  }
  finally {
    isLoading.value = false
  }
}

async function addReply() {
  if (!newReply.value.trim()) {
    toast.add({
      title: 'Empty Reply',
      description: 'Reply cannot be empty',
      color: 'red',
      timeout: 3000,
    })
    return
  }

  isLoading.value = true

  try {
    const response = await fetch(`${BE}/api/gists/${props.gistId}/reply`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        message: newReply.value,
      }),
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Failed to add reply: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || 'Failed to add reply')
    }

    newReply.value = ''
    await fetchReplies()

    toast.add({
      title: 'Reply Added',
      description: 'Your reply has been added successfully',
      color: 'green',
      timeout: 3000,
    })
  }
  catch (err) {
    toast.add({
      title: 'Error',
      description: err instanceof Error ? err.message : 'Failed to add reply',
      color: 'red',
      timeout: 3000,
    })
    console.error('Error adding reply:', err)
  }
  finally {
    isLoading.value = false
  }
}

function startEditReply(replyId: string, content: string) {
  editReplyId.value = replyId
  editReplyContent.value = content
}

async function updateReply() {
  if (!editReplyId.value || !editReplyContent.value.trim()) {
    toast.add({
      title: 'Empty Reply',
      description: 'Reply cannot be empty',
      color: 'red',
      timeout: 3000,
    })
    return
  }

  isLoading.value = true

  try {
    const response = await fetch(`${BE}/api/gists/${props.gistId}/reply/${editReplyId.value}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        message: editReplyContent.value,
      }),
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Failed to update reply: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || 'Failed to update reply')
    }

    editReplyId.value = null
    editReplyContent.value = ''
    await fetchReplies()

    toast.add({
      title: 'Reply Updated',
      description: 'Your reply has been updated successfully',
      color: 'green',
      timeout: 3000,
    })
  }
  catch (err) {
    toast.add({
      title: 'Error',
      description: err instanceof Error ? err.message : 'Failed to update reply',
      color: 'red',
      timeout: 3000,
    })
    console.error('Error updating reply:', err)
  }
  finally {
    isLoading.value = false
  }
}

async function deleteReply(replyId: string) {
  isLoading.value = true

  try {
    const response = await fetch(`${BE}/api/gists/${props.gistId}/reply/${replyId}`, {
      method: 'DELETE',
      credentials: 'include',
    })

    if (!response.ok) {
      throw new Error(`Failed to delete reply: ${response.statusText}`)
    }

    const data = await response.json()

    if (!data.success) {
      throw new Error(data.message || 'Failed to delete reply')
    }

    await fetchReplies()

    toast.add({
      title: 'Reply Deleted',
      description: 'Your reply has been deleted successfully',
      color: 'green',
      timeout: 3000,
    })
  }
  catch (err) {
    toast.add({
      title: 'Error',
      description: err instanceof Error ? err.message : 'Failed to delete reply',
      color: 'red',
      timeout: 3000,
    })
    console.error('Error deleting reply:', err)
  }
  finally {
    isLoading.value = false
  }
}

function formatDate(dateString: string) {
  if (!dateString)
    return 'Unknown date'

  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime()))
      return 'Unknown date'

    const options: Intl.DateTimeFormatOptions = {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }

    return new Intl.DateTimeFormat('en-US', options).format(date)
  }
  catch (err) {
    console.error('Date formatting error:', err, dateString)
    return 'Unknown date'
  }
}

function cancelEdit() {
  editReplyId.value = null
  editReplyContent.value = ''
}
</script>

<template>
  <div class="self-center mt-8 border-t border-gray-200 dark:border-gray-700 pt-8">
    <h3 class="text-xl font-bold mb-6 dark:text-white">
      Replies
    </h3>

    <div v-if="isLoading && replies.length === 0" class="flex justify-center py-4">
      <div class="animate-spin w-6 h-6 border-t-2 border-b-2 border-blue-500 rounded-full" />
      <div class="ml-2">
        Loading replies...
      </div>
    </div>

    <div v-else-if="error" class="p-4 bg-red-50 dark:bg-red-900/20 text-red-500 rounded-lg">
      {{ error }}
      <UButton class="mt-2" size="sm" @click="fetchReplies">
        Try Again
      </UButton>
    </div>

    <div v-else>
      <div v-if="replies.length === 0" class="text-center py-6 text-gray-500 dark:text-gray-400">
        No replies yet. Be the first to reply!
      </div>

      <div v-else class="space-y-4 mb-6">
        <div v-for="reply in replies" :key="reply.replyId" class="p-4 bg-white dark:bg-gray-800/30 rounded-lg shadow-sm">
          <div class="flex justify-between">
            <div class="font-semibold dark:text-white">
              {{ reply.userName }}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              {{ formatDate(reply.createdAt) }}
            </div>
          </div>

          <div v-if="editReplyId === reply.replyId" class="mt-2">
            <UTextarea
              v-model="editReplyContent"
              :rows="3"
              class="w-full"
              placeholder="Edit your reply..."
              :disabled="isLoading"
            />
            <div class="mt-2 flex justify-end space-x-2">
              <UButton size="sm" variant="ghost" @click="cancelEdit">
                Cancel
              </UButton>
              <UButton size="sm" color="primary" :loading="isLoading" @click="updateReply">
                Save
              </UButton>
            </div>
          </div>

          <div v-else class="mt-2 dark:text-gray-200">
            {{ reply.message }}
          </div>

          <div v-if="reply.userId === profile.data.UserID && editReplyId !== reply.replyId" class="mt-2 flex justify-end space-x-2">
            <UButton size="xs" variant="ghost" icon="i-heroicons-pencil" @click="startEditReply(reply.replyId, reply.message)">
              Edit
            </UButton>
            <UButton size="xs" variant="ghost" color="red" icon="i-heroicons-trash" @click="deleteReply(reply.replyId)">
              Delete
            </UButton>
          </div>
        </div>
      </div>

      <div class="mt-6">
        <h4 class="text-lg font-semibold mb-2 dark:text-white">
          Add a Reply
        </h4>
        <UTextarea
          v-model="newReply"
          :rows="3"
          class="w-full"
          placeholder="Write your reply..."
          :disabled="isLoading"
        />
        <div class="mt-2 flex justify-between items-center">
          <div class="text-sm text-gray-500 dark:text-gray-400">
            Posting as <span class="font-medium">{{ currentUserName }}</span>
          </div>
          <UButton color="primary" :loading="isLoading" @click="addReply">
            Post Reply
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>
