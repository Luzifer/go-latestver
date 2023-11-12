<template>
  <b-row>
    <b-col>
      <b-table
        :fields="fields"
        :items="catalog"
        small
        striped
      >
        <template #cell(_key)="data">
          <router-link :to="{name: 'entry', params: { name:data.item.name, tag: data.item.tag }}">
            {{ data.item.name }}:{{ data.item.tag }}
          </router-link>
        </template>

        <template #cell(version_time)="data">
          {{ moment(data.item.version_time).format('lll') }}
        </template>

        <template #cell(_links)="data">
          <a
            v-for="link in data.item.links"
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
          </a>
        </template>
      </b-table>
    </b-col>
  </b-row>
</template>

<script>
import moment from 'moment'

export default {
  data() {
    return {
      catalog: [],
      fields: [
        { key: '_key', label: 'Catalog Entry', sortable: true },
        { key: 'current_version', label: 'Version' },
        { key: 'version_time', label: 'Updated At', sortable: true },
        { key: '_links', label: 'External Links' },
      ],
    }
  },

  methods: {
    fetchCatalogIndex() {
      return fetch('/v1/catalog')
        .then(resp => resp.json())
        .then(data => {
          this.catalog = data
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
    this.fetchCatalogIndex()
  },

  name: 'GoLatestVerCatalogIndex',
}
</script>
