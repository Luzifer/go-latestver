<template>
  <table class="table table-sm table-striped">
    <thead>
      <tr>
        <th>Catalog Entry</th>
        <th>Version From</th>
        <th>Version To</th>
        <th>Updated At</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="log in logs"
        :key="`${log.name}:${log.tag}@${log.timestamp}`"
      >
        <td>
          <router-link
            :to="{name: 'entry', params: { name: log.catalog_name, tag: log.catalog_tag }}"
            class="text-decoration-none"
          >
            {{ log.catalog_name }}:{{ log.catalog_tag }}
          </router-link>
        </td>
        <td>{{ log.version_from }}</td>
        <td>{{ log.version_to }}</td>
        <td>{{ moment(log.timestamp).format('lll') }}</td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import moment from 'moment'

export default defineComponent({
  methods: {
    moment,
  },

  name: 'GoLatestVerLogTable',

  props: {
    logs: {
      required: true,
      type: Array<any>,
    },
  },
})
</script>
