<template>
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
</template>

<script>
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
    }
  },

  methods: {
    moment,
  },

  name: 'GoLatestVerLogTable',

  props: {
    logs: {
      required: true,
      type: Array,
    },
  },
}
</script>
