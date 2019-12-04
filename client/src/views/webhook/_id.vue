<template>
  <div class="app-container">
    <el-form :model="webhook" label-width="120px">
      <el-form-item label="URL">
        <el-input v-model="webhook.url" />
      </el-form-item>
      <el-form-item label="Body">
        <el-input v-model="webhook.body" />
      </el-form-item>
      <el-form-item v-for="(header, idx) in webhook.header" :key="header.key" label="Header">
        <el-input v-model="header.key" />
        <el-input v-model="header.value" />
        <el-button @click="deleteHeader(idx)">Delete Header</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="addHeader">New Header</el-button>
      </el-form-item>
      <el-form-item v-for="(secret, idx) in webhook.secrets" :key="secret.place_holder" label="Secrets">
        <el-input v-model="secret.place_holder" />
        <el-input v-model="secret.secret" />
        <el-button @click="deleteSecret(idx)">Delete Secret</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="addSecret">New Secret</el-button>
      </el-form-item>
      <el-form-item>
        <el-tag
          v-for="e in webhook.event"
          :key="e.event"
          closable
          @close="delSubEvent(e)"
        >
          {{ e.event }}
        </el-tag>
        <el-select
          size="mini"
          value=""
          placeholder="購読するイベント"
          @change="addSubEvent"
        >
          <el-option
            v-for="option in options"
            :key="option"
            :value="option"
            :label="option"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <template v-if="webhookID === 'new'">
          <el-button @click="createWebhook">Create</el-button>
        </template>
        <template v-else>
          <el-button @click="saveWebhook">Save</el-button>
          <el-button @click="deleteWebhook">Delete</el-button>
        </template>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { getWebhooks, postWebhook, putWebhook, deleteWebhook } from '@/api/webhook'

export default {
  name: 'WebhookList',
  data() {
    return {
      webhookID: '',
      webhook: {},
      allOptions: ['poi', 'po']
    }
  },
  computed: {
    options() {
      return this.allOptions.filter(opt => {
        return !(this.webhook.event || []).find(e => e.event === opt)
      })
    }
  },
  async mounted() {
    await this.updateWebhook()
  },
  methods: {
    async updateWebhook() {
      this.webhookID = this.$route.params.id

      if (!this.webhookID) {
        this.webhookID = 'new'
        this.webhook = {
          url: '',
          body: '',
          header: [],
          event: [],
          secrets: []
        }
      } else {
        const res = await getWebhooks()
        console.log(res)
        this.webhook = res.find(w => w.ID === Number(this.webhookID))
      }
    },
    deleteHeader(idx) {
      this.webhook.header.splice(idx, 1)
    },
    addHeader() {
      this.webhook.header.push({ key: '', value: '' })
    },
    deleteSecret(idx) {
      this.webhook.secrets.splice(idx, 1)
    },
    addSecret() {
      this.webhook.secrets.push({ place_holder: '', secret: '' })
    },
    async createWebhook() {
      const res = await postWebhook(this.webhook)
      this.$router.push(`/webhook/${res.ID}`)
      await this.updateWebhook()
    },
    async saveWebhook() {
      await putWebhook(this.webhook)
      await this.updateWebhook()
    },
    async deleteWebhook() {
      await deleteWebhook(this.webhook)
      this.$router.push('/webhook')
    },
    addSubEvent(event) {
      this.webhook.event.push({ event })
    },
    delSubEvent(event) {
      console.log(event)
      this.webhook.event = this.webhook.event.filter(e => e.event !== event.event)
    }
  }
}
</script>

<style>
</style>
