<template>
  <footer class="py-8 border-t dark:border-myborder border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center">
        <div class="flex items-center mb-4 md:mb-0">
          <Icon name="tabler:code" class="text-green-500 w-6 h-6 mr-2" />
          <span class="dark:text-white text-gray-900 font-bold text-xl">CodeFlick</span>
        </div>
        <div class="flex space-x-6">
          <a href="#" class="text-gray-500 hover:text-green-500 dark:text-gray-400 dark:hover:text-green-400">
            <Icon name="tabler:brand-github" class="w-6 h-6" />
          </a>
          <a href="#" class="text-gray-500 hover:text-green-500 dark:text-gray-400 dark:hover:text-green-400">
            <Icon name="tabler:brand-twitter" class="w-6 h-6" />
          </a>
          <a href="#" class="text-gray-500 hover:text-green-500 dark:text-gray-400 dark:hover:text-green-400">
            <Icon name="tabler:brand-discord" class="w-6 h-6" />
          </a>
        </div>
      </div>
      <div class="text-center mt-6 text-sm text-gray-500 dark:text-gray-400">
        &copy; {{ new Date().getFullYear() }} CodeFlick. All rights reserved.
      </div>
    </div>
  </footer>
</template>
