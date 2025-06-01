<script setup>
import { ref, onMounted, watch, computed, onUnmounted } from 'vue'
import router from 'src/router/index'
import algorithms from 'configs/algorithms'

const links = router.links.concat(Object.keys(algorithms).map(
  name => ({
    name: name.toUpperCase(),
    path: "/algo/" + name,
  })
));

const currentTheme = ref('')

const updateTheme = (event) => {
  if (event.key === "theme" && currentTheme.value != event.newValue) {
    currentTheme.value = event.newValue;
    document.body.setAttribute('data-theme', currentTheme.value);
  }
}

const toggleTheme = () => {
  currentTheme.value = currentTheme.value === 'light' ? 'dark' : 'light';

  document.body.setAttribute('data-theme', currentTheme.value);
  localStorage.setItem('theme', currentTheme.value)
}

onMounted(() => {
  // сохраненная тема в localStorage
  const savedTheme = localStorage.getItem('theme')

  currentTheme.value = savedTheme ? savedTheme :
    // системные предпочтения
    window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';

  document.body.setAttribute('data-theme', currentTheme.value);

  window.addEventListener('storage', updateTheme);
})

onUnmounted(() => {
  window.removeEventListener("storage", updateTheme);
})

</script>

<template>
  <header class="menu">
    <RouterLink class="router-link" v-for="link in links" :to="link.path">{{ link.name }}</RouterLink>

    <label class="theme-toggle">
      <input type="checkbox" class="toggle-input" readonly @click="toggleTheme" />
      <span class="toggle-slider">
        <span class="icon" :class="{ moveRight: currentTheme === 'dark' }">
          <span v-if="currentTheme === 'light'">☀</span>
          <span v-else>🌙</span>
        </span>
      </span>
    </label>

  </header>
  <main class="main-layout">
    <RouterView />
  </main>
</template>

<style scoped lang="scss">
.menu {
  position: fixed;
  top: 0;
  width: 100%;
  display: flex;
  justify-content: centers;
  background-color: var(--color-menu-bg);
  gap: 1rem;
  height: 2rem;
  z-index: 9999;


  .router-link {
    flex-grow: 1;
    text-align: center;
  }

  .router-link-active {
    font-weight: bold;
  }

  .router-link-exact-active {
    color: var(--color-menu-active);
  }

  .theme-toggle {
    position: relative;
    display: inline-block;
    width: 3rem;
    height: 1.6rem;
    margin-top: 0.2rem;
    margin-right: 1rem;
    justify-self: center;

    .toggle-input {
      opacity: 0;
      width: 0;
      height: 0;
    }

    .toggle-slider:hover {
      border-color: var(--color-border-hover);
    }


    .toggle-slider {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: var(--color-card-bg);
      border: 1px solid var(--color-border);
      border-radius: 0.8rem;
      cursor: pointer;
      transition: background-color 0.3s, border-color 0.3s;

      .icon {
        position: absolute;
        top: 0.05rem;
        left: 0.05rem;
        width: 1.4rem;
        height: 1.4rem;
        border-radius: 50%;
        background-color: var(--color-menu-bg);
        color: var(--color-text);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 0.7rem;
        box-shadow: 0 0 2px rgba(0, 0, 0, 0.2);
        transition: left 0.3s ease;
      }

      .moveRight {
        left: 1.5rem;
      }
    }
  }


}

.main-layout {
  margin: 0 auto;
  height: 100%;
  min-height: min-content;
  padding: 3rem 0;
}
</style>
