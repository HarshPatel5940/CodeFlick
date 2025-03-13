<template>
  <section class="py-16 bg-gray-50 dark:bg-mybg/50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-4xl font-bold dark:text-white text-gray-900 mb-4">
          Why Choose CodeFlick?
        </h2>
        <p class="text-lg dark:text-gray-300 text-gray-600 max-w-2xl mx-auto">
          The modern platform for developers to share, discover, and collaborate on code snippets
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-8 mb-16">
        <FeatureCard
          v-motion-slide-visible-once-bottom
          title="Easy Code Sharing"
          description="Create and share code snippets with syntax highlighting for over 100 programming languages and frameworks."
          icon="tabler:code"
        />
        <FeatureCard
          v-motion-slide-visible-once-bottom
          :delay="100"
          title="Collaborate"
          description="Work together with other developers on your code snippets. Add comments, suggestions, and improvements."
          icon="tabler:users"
        />
        <FeatureCard
          v-motion-slide-visible-once-bottom
          :delay="200"
          title="Discover"
          description="Find code snippets that solve your problems. Learn from the community and improve your skills."
          icon="tabler:globe"
        />
      </div>
    </div>
  </section>
</template>
