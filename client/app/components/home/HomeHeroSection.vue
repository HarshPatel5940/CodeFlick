<script setup>
const { isLoggedIn, username } = defineProps({
  isLoggedIn: Boolean,
  username: String,
})

const emit = defineEmits(['login'])

const backgroundCodePattern = `
import { useState, useEffect } from 'react';
function CodeSnippet({ id, language }) {
  const [snippet, setSnippet] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchSnippet() {
      try {
        setLoading(true);
        const response = await fetch(\`/api/snippets/\${id}\`);
        if (!response.ok) {
          throw new Error('Failed to fetch snippet');
        }
        const data = await response.json();
        setSnippet(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchSnippet();
  }, [id]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!snippet) return <div>No snippet found</div>;

  return (
    <div className="code-snippet">
      <div className="snippet-header">
        <h2>{snippet.title}</h2>
        <div className="language-badge">{language}</div>
      </div>
      <pre>{snippet.code}</pre>
    </div>
  );
}

function Spinner() {
  return (
    <div className="spinner">
      <div className="bounce1"></div>
      <div className="bounce2"></div>
      <div className="bounce3"></div>
    </div>
  );
}
`.repeat(5)

const sampleCode = ref(`package main

import "fmt"

// fib returns a function that returns
// successive Fibonacci numbers.
func fib() func() int {
    a, b := 0, 1
    return func() int {
        a, b = b, a+b
        return a
    }
}

func main() {
    f := fib()
    // Function calls are evaluated left-to-right.
    fmt.Println(f(), f(), f(), f(), f())
}
`)
</script>

<template>
  <section class="relative overflow-hidden pt-20 pb-16 md:pt-32 md:pb-24">
    <div class="absolute inset-0 opacity-5 dark:opacity-10 overflow-hidden z-0">
      <pre class="text-xs sm:text-sm text-mybg dark:text-white font-mono leading-tight opacity-50">
        <code>{{ backgroundCodePattern }}</code>
      </pre>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
      <div class="text-center mb-16">
        <div class="flex justify-center mb-4">
          <Icon name="tabler:code" class="text-green-500 w-16 h-16 md:w-20 md:h-20" />
        </div>
        <h1 class="dark:text-gray-100 text-gray-900 text-5xl md:text-7xl font-bold tracking-wider mb-6">
          CodeFlick
        </h1>
        <p class="text-xl md:text-2xl dark:text-gray-300 text-gray-700 max-w-3xl mx-auto mb-8">
          Flick it, share it! Open-source Gists for the dev community.
        </p>

        <div v-if="!isLoggedIn" class="dark:text-white/80 text-md mb-6">
          Join our growing community of developers sharing code snippets
        </div>
        <div v-else class="dark:text-white/80 text-lg md:text-xl mb-6">
          Welcome back, <span class="font-semibold text-green-500">{{ username }}</span>!
        </div>

        <UButton
          v-motion-slide-visible-bottom
          icon="tabler:login"
          class="text-md md:text-lg font-mono tracking-wide p-3 px-6 rounded-xl border-2 border-neutral-700/50 hover:border-green-500 dark:bg-mybg dark:text-white"
          @click="emit('login')"
        >
          {{ isLoggedIn ? "Go to Dashboard" : "Get Started" }}
        </UButton>
      </div>

      <div class="max-w-4xl mx-auto">
        <CodePreview
          filename="fibonacci.go"
          lang="go"
          :code="sampleCode"
          :show-copy-button="true"
          :show-edit-button="false"
        />
      </div>
    </div>

    <div class="absolute top-0 left-0 w-full h-full overflow-hidden -z-5">
      <div class="absolute top-1/4 left-1/4 w-64 h-64 bg-green-500/10 rounded-full filter blur-3xl" />
      <div class="absolute bottom-1/3 right-1/3 w-96 h-96 bg-blue-500/10 rounded-full filter blur-3xl" />
    </div>
  </section>
</template>

<style>
.code-pattern {
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
  background-size: 20px 20px;
  background-position: center center;
}
</style>
