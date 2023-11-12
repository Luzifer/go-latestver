<template>
  <b-row>
    <b-col>
      <log-table :logs="logs" />
    </b-col>
  </b-row>
</template>

<script>
import LogTable from './logtable.vue'

export default {
  components: { LogTable },

  data() {
    return {
      logs: [],
    }
  },

  methods: {
    fetchLog() {
      return fetch('/v1/log?num=50')
        .then(resp => resp.json())
        .then(data => {
          this.logs = data
        })
    },
  },

  mounted() {
    this.fetchLog()
  },

  name: 'GoLatestVerLog',
}
</script>
