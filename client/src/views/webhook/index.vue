<template>
  <div class="app-container">
    <el-table
      :data="webhooks"
    >
      <el-table-column
        prop="ID"
        label="ID"
        width="80px"
      />
      <el-table-column
        prop="url"
        label="URL"
      />
      <el-table-column
        prop="body"
        label="Body"
      />
      <el-table-column
        label="Action"
      >
        <template slot-scope="scope">
          <router-link
            :to="'/webhook/'+scope.row.ID"
          >
            詳細
          </router-link>
        </template>
      </el-table-column>
    </el-table>
    <el-button @click="$router.push('/webhook/new')">New</el-button>
  </div>
</template>

<script>
import { getWebhooks } from '@/api/webhook'

export default {
  name: 'WebhookList',
  data() {
    return {
      webhooks: []
    }
  },
  async mounted() {
    const res = await getWebhooks()
    console.log(res)
    this.webhooks = res
  },
  methods: {
  }
}
</script>

<style>
</style>
