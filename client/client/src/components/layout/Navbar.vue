<template>
  <header>
    <div class="flex justify-between items-center p-8 lg:px-12 relative z-20">
      <div class="text-xl lg:text-2xl font-bold dark:text-light">KBachStats</div>

      <!-- Mobile Toggle Button -->
      <div class="md:hidden z-30">
        <button @click="isMenuOpen = !isMenuOpen" class="block focus:outline-none">
          <span v-if="isMenuOpen" class="text-5xl md:text-primary text-primary dark:text-secondary">
            <Icon icon="material-symbols:close-rounded"/>
          </span>

          <span v-if="!isMenuOpen" class="text-5xl md:text-primary text-dark dark:text-light">
            <Icon icon="material-symbols-light:menu"/>
          </span>
        </button>
      </div>

      <!-- Nav Links -->
      <nav :class="[
        `fixed inset-0 z-20 flex flex-col items-center justify-center bg-primary md:relative
        md:bg-transparent md:flex md:justify-between md:flex-row ${isMenuOpen ? 'block' : 'hidden'}`
      ]">
        <ul class="flex flex-col items-center space-y-5 md:flex-row md:space-x-3 lg:space-x-6 md:space-y-0">
          <li v-for="item in menu" :key="item.name">
            <a :href="item.href" class="block transition ease-linear md:text-sm lg:text-lg font-bold text-white md:text-primary dark:text-light dark:hover:text-secondary">
              {{ item.name }}
            </a>
          </li>
        </ul>

        <div class="md:ml-4 flex flex-col md:flex-row">
          <PageLink dest="/signin" :primary="true" label="Signin"/>
          <PageLink dest="/register" :primary="false" label="Register"/>
        </div>

        <!-- Dark mode button -->
        <button @click="toggleDarkMode" class="text-primary dark:text-light dark:hover:text-secondary ml-8 lg:ml-16 z-10 hidden md:block">
          <Icon v-if="!isDarkmode" icon="tabler:moon" class="text-3xl"/>
          <Icon v-else icon="si:sun-line" class="text-3xl"/>
        </button>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import PageLink from '../UI/PageLink.vue'

const isMenuOpen = ref(false)

const menu = ref([
  {name:'Services', href: '#services'},
  {name:'Offer', href: '#offer'},
  {name:'Pricing', href: '#pricing'},
  {name:'Tetimonials', href:'#testimonials'}
])

const isDarkmode = ref(localStorage.getItem('theme') === 'dark')

const toggleDarkMode = () => {
  const html = document.documentElement
  if (isDarkmode.value) {
    html.removeAttribute('data-theme')
    localStorage.setItem('theme', 'light')
  } else {
    html.setAttribute('data-theme', 'dark')
    localStorage.setItem('theme', 'dark')
  }

  isDarkmode.value = !isDarkmode.value
}
</script>