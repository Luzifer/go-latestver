<template>
  <b-row>
    <b-col>
      <log-table :logs="logs" />
    </b-col>
  </b-row>
</template>

<script>
import axios from 'axios'
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
      axios.get('/v1/log?num=50')
        .then(resp => {
          this.logs = resp.data
        })
    },
  },

  mounted() {
    this.fetchLog()
  },

  name: 'GoLatestVerLog',
}
</script>
