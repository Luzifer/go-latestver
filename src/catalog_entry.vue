<template>
  <div>
    <b-row>
      <b-col>
        <b-card-group>
          <b-card header="Current Version">
            {{ entry.current_version }}<br>
            <small>{{ moment(entry.version_time).format('lll') }}</small>
          </b-card>
          <b-card header="Last Checked">
            {{ moment(entry.last_checked).format('lll') }}
          </b-card>
          <b-card header="Badges">
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
          </b-card>
          <b-card
            header="External Links"
            no-body
          >
            <b-list-group flush>
              <b-list-group-item
                v-for="link in entry.links"
                :key="link.name"
                :href="link.url"
                rel="noopener noreferrer"
                target="_blank"
              >
                <font-awesome-icon
                  fixed-width
                  :icon="iconClassesToIcon(link.icon_class)"
                />
                {{ link.name }}
              </b-list-group-item>
            </b-list-group>
          </b-card>
        </b-card-group>
      </b-col>
    </b-row>

    <b-row class="mt-3">
      <b-col>
        <log-table :logs="logs" />
      </b-col>
    </b-row>
  </div>
</template>

<script>
import axios from 'axios'
import LogTable from './logtable.vue'
import moment from 'moment'

export default {
  components: { LogTable },

  computed: {
    badgeURL() {
      return `${window.location.href.split('?')[0]}.svg`
    },
  },

  data() {
    return {
      entry: {},
      logs: [],
    }
  },

  methods: {
    copyURL(evt) {
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

    fetchEntry() {
      axios.get(`/v1/catalog/${this.$route.params.name}/${this.$route.params.tag}`)
        .then(resp => {
          this.entry = resp.data
        })
    },

    fetchLog() {
      axios.get(`/v1/catalog/${this.$route.params.name}/${this.$route.params.tag}/log`)
        .then(resp => {
          this.logs = resp.data
        })
    },

    iconClassesToIcon(ic) {
      let namespace = 'fas'
      let icon = ''

      for (const c of ic.split(' ')) {
        if (c === 'fa-fw') {
          continue
        }

        if (['fab', 'fas'].includes(c)) {
          namespace = c
        }

        if (c.startsWith('fa-')) {
          icon = c.replace('fa-', '')
        }
      }

      return [namespace, icon]
    },

    moment,
  },

  mounted() {
    this.fetchEntry()
    this.fetchLog()
  },

  name: 'GoLatestVerCatalogEntry',
}
</script>

<style>
img.clickable {
  cursor:pointer;
}
</style>
