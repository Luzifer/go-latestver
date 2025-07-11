<template>
  <div class="row">
    <div class="col">
      <table class="table table-sm table-striped">
        <thead>
          <tr>
            <th>Catalog Entry</th>
            <th>Version</th>
            <th>Updated At</th>
            <th>External Links</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="entry in catalog"
            :key="entry.key"
          >
            <td>
              <router-link
                :to="{name: 'entry', params: {name: entry.name, tag: entry.tag}}"
                class="text-decoration-none"
              >
                {{ entry.key }}
              </router-link>
            </td>
            <td>{{ entry.current_version }}</td>
            <td>{{ moment(entry.version_time).format('lll') }}</td>
            <td>
              <a
                v-for="link in entry.links"
                :key="link.name"
                :href="link.url"
                rel="noopener noreferrer"
                target="_blank"
                class="text-decoration-none"
              >
                <i :class="iconClassesToIcon(link.icon_class)" />
                {{ link.name }}
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { iconClassesToIcon } from './helpers'
import moment from 'moment'

export default defineComponent({
  data() {
    return {
      catalog: [] as any[],
    }
  },

  methods: {
    fetchCatalogIndex(): Promise<void> {
      return fetch('/v1/catalog')
        .then(resp => resp.json())
        .then(data => {
          this.catalog = data.map((e: any) => ({ ...e, key: `${e.name}:${e.tag}` }))
        })
    },

    iconClassesToIcon,
    moment,
  },

  mounted() {
    this.fetchCatalogIndex()
  },

  name: 'GoLatestVerCatalogIndex',
})
</script>
