<script setup lang="ts">
const colorMode = useColorMode()
const theme = useTheme()

function toggleColorMode() {
  colorMode.value = theme.toggleTheme()
}

onMounted(() => {
  colorMode.value = theme.theme
})
</script>

<template>
  <button class="p-2 bg-inherit border-none" @click="toggleColorMode">
    <Icon v-if="theme.themeBool" name="solar:moon-bold" class="p-1" />
    <Icon v-if="!theme.themeBool" name="solar:sun-bold" class="p-1 " style="color: #000;" />
  </button>
</template>
