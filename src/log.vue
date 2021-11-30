<template>
  <b-row>
    <b-col>
      <b-table
        :fields="fields"
        :items="logs"
        small
        striped
      >
        <template #cell(_key)="data">
          <router-link :to="{name: 'entry', params: { name:data.item.catalog_name, tag: data.item.catalog_tag }}">
            {{ data.item.catalog_name }}:{{ data.item.catalog_tag }}
          </router-link>
        </template>

        <template #cell(timestamp)="data">
          {{ moment(data.item.timestamp).format('lll') }}
        </template>
      </b-table>
    </b-col>
  </b-row>
</template>

<script>
import axios from 'axios'
import moment from 'moment'

export default {
  data() {
    return {
      fields: [
        { key: '_key', label: 'Catalog Entry' },
        { key: 'version_from', label: 'Version From' },
        { key: 'version_to', label: 'Version To' },
        { key: 'timestamp', label: 'Updated At' },
      ],

      logs: [],
    }
  },

  methods: {
    fetchLog() {
      axios.get('/v1/log')
        .then(resp => {
          this.logs = resp.data
        })
    },

    moment,
  },

  mounted() {
    this.fetchLog()
  },

  name: 'GoLatestVerLog',
}
</script>
