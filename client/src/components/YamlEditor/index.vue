<template>
  <div class="yaml-editor">
    <textarea ref="textarea" />
  </div>
</template>

<script>
import CodeMirror from 'codemirror'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/rubyblue.css'
// require('script-loader!yamllint')
import 'codemirror/mode/yaml/yaml'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/yaml-lint'

const yaml = require('js-yaml')

export default {
  name: 'YamlEditor',
  /* eslint-disable vue/require-prop-types */
  props: ['value'],
  data() {
    return {
      yamlEditor: false
    }
  },
  watch: {
    value(value) {
      const editorValue = this.yamlEditor.getValue()
      if (value !== editorValue) {
        this.yamlEditor.setValue(value)
      }
    }
  },
  mounted() {
    this.yamlEditor = CodeMirror.fromTextArea(this.$refs.textarea, {
      lineNumbers: true,
      mode: 'text/x-yaml',
      gutters: ['CodeMirror-lint-markers'],
      theme: 'rubyblue',
      lint: true
    })

    this.yamlEditor.setValue(this.value)
    this.yamlEditor.on('change', cm => {
      this.$emit('changed', cm.getValue())
      this.$emit('input', cm.getValue())
    })
  },
  methods: {
    getValue() {
      return this.yamlEditor.getValue()
    }
  },
  created() {
    window.jsyaml = yaml
  }
}
</script>

<style scoped>
.yaml-editor{
  height: 100%;
  position: relative;
}
.yaml-editor >>> .CodeMirror {
  height: auto;
  min-height: 300px;
}
.yaml-editor >>> .CodeMirror-scroll{
  min-height: 300px;
}
.yaml-editor >>> .cm-s-rubyblue span.cm-string {
  color: #F08047;
}
</style>
