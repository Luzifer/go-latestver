<template>
  <div>
    <div class="row">
      <div class="col">
        <div class="card-group">
          <div class="card">
            <div class="card-header">
              Current Version
            </div>
            <div class="card-body">
              {{ entry.current_version }}<br>
              <small>{{ moment(entry.version_time).format('lll') }}</small>
            </div>
          </div>

          <div class="card">
            <div class="card-header">
              Last Checked
            </div>
            <div class="card-body">
              {{ moment(entry.last_checked).format('lll') }}
            </div>
          </div>

          <div class="card">
            <div class="card-header">
              Badges
            </div>
            <div class="card-body">
              <p class="text-center">
                Current Version:<br>
                <img
                  class="clickable"
                  :src="badgeURL"
                  @click="copyURL"
                >
              </p>
              <p class="text-center">
                Compare to Version:<br>
                <img
                  class="clickable"
                  :src="`${badgeURL}?compare=otherversion`"
                  @click="copyURL"
                >
              </p>
              <p class="text-center">
                <small>(Click badge to copy URL)</small>
              </p>
            </div>
          </div>

          <div class="card">
            <div class="card-header">
              External Links
            </div>
            <div class="list-group list-group-flush">
              <a
                v-for="link in entry.links"
                :key="link.name"
                class="list-group-item list-group-item-action"
                :href="link.url"
                rel="noopener noreferrer"
                target="_blank"
              >
                <i :class="iconClassesToIcon(link.icon_class)" />
                {{ link.name }}
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row mt-3">
      <div class="col">
        <log-table :logs="logs" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { iconClassesToIcon } from './helpers'
import LogTable from './logtable.vue'
import moment from 'moment'

export default defineComponent({
  components: { LogTable },

  computed: {
    badgeURL(): string {
      return `${window.location.href.split('?')[0]}.svg`
    },
  },

  data() {
    return {
      entry: {} as any,
      logs: [] as any[],
    }
  },

  methods: {
    copyURL(evt): void {
      navigator.clipboard.writeText(evt.target.attributes.src.value)
        .then(() => this.$bvToast.toast('URL copied to clipboard', {
          autoHideDelay: 2000,
          solid: true,
          title: 'Copy Badge URL',
          variant: 'success',
        }))
        .catch(() => this.$bvToast.toast('Something went wrong', {
          autoHideDelay: 2000,
          solid: true,
          title: 'Copy Badge URL',
          variant: 'danger',
        }))
    },

    fetchEntry(): Promise<void> {
      return fetch(`/v1/catalog/${this.$route.params.name}/${this.$route.params.tag}`)
        .then(resp => resp.json())
        .then(data => {
          this.entry = data
        })
    },

    fetchLog(): Promise<void> {
      return fetch(`/v1/catalog/${this.$route.params.name}/${this.$route.params.tag}/log`)
        .then(resp => resp.json())
        .then(data => {
          this.logs = data
        })
    },

    iconClassesToIcon,
    moment,
  },

  mounted() {
    this.fetchEntry()
    this.fetchLog()
  },

  name: 'GoLatestVerCatalogEntry',
})
</script>

<style>
img.clickable {
  cursor:pointer;
}
</style>
