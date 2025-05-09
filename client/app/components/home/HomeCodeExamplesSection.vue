<script setup lang="ts">
const reactComponentCode = `import React from 'react';

const Button = ({
  children,
  variant = 'primary',
  size = 'medium',
  onClick
}) => {
  const baseClasses = 'rounded font-medium focus:outline-none';

  const variantClasses = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700',
    secondary: 'bg-gray-200 text-gray-800 hover:bg-gray-300',
    danger: 'bg-red-600 text-white hover:bg-red-700'
  };

  const sizeClasses = {
    small: 'text-sm py-1 px-3',
    medium: 'py-2 px-4',
    large: 'text-lg py-3 px-6'
  };

  const classes = \`\${baseClasses} \${variantClasses[variant]} \${sizeClasses[size]}\`;

  return (
    <button
      className={classes}
      onClick={onClick}
    >
      {children}
    </button>
  );
};

export default Button;`

const pythonCode = `def quicksort(arr):
    if len(arr) <= 1:
        return arr

    pivot = arr[len(arr) 
    left = [x for x in arr if x < pivot]
    middle = [x for x in arr if x == pivot]
    right = [x for x in arr if x > pivot]

    return quicksort(left) + middle + quicksort(right)

# Example usage
if __name__ == "__main__":
    unsorted_list = [3, 6, 8, 10, 1, 2, 1]
    sorted_list = quicksort(unsorted_list)
    print(sorted_list)  # Output: [1, 1, 2, 3, 6, 8, 10]`

const steps = [
  {
    title: 'Sign Up',
    description: 'Create your account with Google authentication in seconds',
  },
  {
    title: 'Create Snippet',
    description: 'Write or paste your code with full syntax highlighting',
  },
  {
    title: 'Share',
    description: 'Get a unique link to share with your team or the community',
  },
  {
    title: 'Collaborate',
    description: 'Receive feedback, suggestions, and improvements',
  },
]
</script>

<template>
  <section class="py-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-4xl font-bold dark:text-white text-gray-900 mb-4">
          Code For Every Need
        </h2>
        <p class="text-lg dark:text-gray-300 text-gray-600 max-w-2xl mx-auto">
          Share any code in any language - from simple snippets to complex algorithms
        </p>
      </div>

      <div class="grid md:grid-cols-2 gap-8 mb-8">
        <div>
          <h3 class="text-xl font-semibold dark:text-white text-gray-900 mb-4">
            React Component
          </h3>
          <CodePreview
            filename="Button.jsx"
            lang="jsx"
            :code="reactComponentCode"
          />
        </div>
        <div>
          <h3 class="text-xl font-semibold dark:text-white text-gray-900 mb-4">
            Python Algorithm
          </h3>
          <CodePreview
            filename="quicksort.py"
            lang="python"
            :code="pythonCode"
          />
        </div>
      </div>
    </div>
  </section>

  <section class="py-16 bg-gray-50 dark:bg-mybg/50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-4xl font-bold dark:text-white text-gray-900 mb-4">
          How It Works
        </h2>
        <p class="text-lg dark:text-gray-300 text-gray-600 max-w-2xl mx-auto">
          Get started with CodeFlick in just a few simple steps
        </p>
      </div>

      <div class="grid md:grid-cols-4 gap-8">
        <div v-for="(step, index) in steps" :key="index" class="text-center">
          <div class="flex justify-center mb-4">
            <div class="w-12 h-12 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
              <span class="text-green-600 dark:text-green-400 font-bold text-lg">{{ index + 1 }}</span>
            </div>
          </div>
          <h3 class="text-xl font-semibold dark:text-white text-gray-900 mb-2">
            {{ step.title }}
          </h3>
          <p class="dark:text-gray-400 text-gray-600">
            {{ step.description }}
          </p>
        </div>
      </div>
    </div>
  </section>
</template>
