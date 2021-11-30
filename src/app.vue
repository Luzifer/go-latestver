<template>
  <div>
    <b-navbar
      class="mb-2"
      toggleable="lg"
      type="dark"
      variant="primary"
    >
      <b-navbar-brand :to="{name: 'index'}">
        <font-awesome-icon
          fixed-width
          :icon="['fas', 'cloud-download-alt']"
        />
        Go-LatestVer
      </b-navbar-brand>

      <b-navbar-toggle target="nav-collapse" />

      <b-collapse
        id="nav-collapse"
        is-nav
      >
        <b-navbar-nav>
          <b-nav-item :to="{name: 'log'}">
            Log
          </b-nav-item>
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-item href="/log.rss">
            <font-awesome-icon
              fixed-width
              :icon="['fas', 'rss']"
            />
            All Updates
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-container>
      <router-view />
    </b-container>
  </div>
</template>

<script>
export default {
  name: 'GoLatestVerApp',
}
</script>
