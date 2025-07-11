<template>
  <div>
    <nav class="navbar navbar-expand-lg bg-body-tertiary mb-2">
      <div class="container-fluid">
        <RouterLink
          class="navbar-brand"
          :to="{name: 'index'}"
        >
          <i class="fas fa-cloud-download-alt fa-fw me-1" />
          Go-LatestVer
        </RouterLink>

        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon" />
        </button>

        <div
          id="navbarSupportedContent"
          class="collapse navbar-collapse"
        >
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">
            <li class="nav-item">
              <RouterLink
                :to="{name: 'log'}"
                class="nav-link"
              >
                Log
              </RouterLink>
            </li>
          </ul>
          <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
            <li class="nav-item">
              <a
                href="/log.rss"
                class="nav-link"
              >
                <i class="fas fa-rss fa-fw" />
                All Updates
              </a>
            </li>
            <li class="nav-item">
              <a
                class="nav-link"
                href="https://github.com/Luzifer/go-latestver"
                rel="noopener noreferrer"
                target="_blank"
                title="Go-LatestVer on Github"
              >
                <i class="fab fa-github fa-fw" />
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { defineComponent } from 'vue'

export default defineComponent({
  components: { RouterLink, RouterView },
  name: 'GoLatestVerApp',
})
</script>
