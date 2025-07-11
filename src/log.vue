<template>
  <div class="row">
    <div class="col">
      <log-table :logs="logs" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import LogTable from './logtable.vue'

export default defineComponent({
  components: { LogTable },

  data() {
    return {
      logs: [] as any[],
    }
  },

  methods: {
    fetchLog(): Promise<void> {
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
})
</script>
