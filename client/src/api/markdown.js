import request from '@/utils/requestServer'

export function fetchMarkdown(path) {
  return request({
    url: `/docs/${path}/index.md`,
    method: 'get'
  })
}
