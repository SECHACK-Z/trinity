<template>
  <div class="app-container">
    <div v-html="html" />
  </div>
</template>

<script>
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import { fetchMarkdown } from '@/api/markdown'

export default {
  name: 'Documentation',
  data() {
    return {
      md: new MarkdownIt({
        highlight: function(code, lang) {
          return hljs.highlightAuto(code, [lang]).value
        },
        html: true,
        linkify: true,
        breaks: true,
        typographer: true
      }),
      markdownText: '# hello'
    }
  },
  computed: {
    html() {
      return this.md.render(this.markdownText)
    }
  },
  async mounted() {
    const path = this.$route.meta.path
    this.markdownText = await fetchMarkdown(path)
  }
}
</script>

<style lang="scss" scoped>
</style>
