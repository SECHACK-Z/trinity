<template>
  <div class="app-container">
    <div class="editor-container">
      <yaml-editor ref="yamlEditor" v-model="configText"></yaml-editor>
    </div>
    <el-row>
      <el-button type="primary" @click="setConfig">反映</el-button>
      <el-button type="success" @click="saveConfig">保存</el-button>
    </el-row>
  </div>
</template>

<script>
import {getConfig, setConfig, saveConfig} from '@/api/config'
import YamlEditor from '@/components/YamlEditor'

export default {
  name: 'ConfigEditor',
  components: {
    YamlEditor
  },
  data() {
    return {
      configText: ''
    }
  },
  methods: {
    async setConfig() {
      await setConfig({Yaml: this.configText})
      const res = await getConfig()
      this.configText = res.Yaml
      this.$notify({
        title: '反映成功',
        message: '設定を反映しました',
        type: 'success'
      })
    },
    async saveConfig() {
      await saveConfig({Yaml: this.configText})
      const res = await getConfig()
      this.configText = res.Yaml
      this.$notify({
        title: '保存成功',
        message: '設定を保存しました',
        type: 'success'
      })
    }
  },
  async mounted() {
    const res = await getConfig()
    console.log(res)
    this.configText = res.Yaml
  }
}
</script>

<style>
  .editor-container {
    position: relative;
    height: 100%;
  }

  .CodeMirror {
    line-height: 18px;
  }
</style>
