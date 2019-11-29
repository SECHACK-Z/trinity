import request from '@/utils/requestServer'

export function getWebhooks() {
  return request({
    url: '/api/webhooks',
    method: 'get'
  }).then(res => {
    return res.map(webhook => {
      const headerMap = JSON.parse(webhook.header || '{}')
      const header = Object.keys(headerMap).map(key => {
        return {
          key: key,
          value: headerMap[key]
        }
      })
      return { ...webhook, header }
    })
  })
}

export function postWebhook(webhook) {
  const header = {}
  webhook.header.forEach(h => {
    header[h.key] = h.value
  })

  const data = { ...webhook, header: JSON.stringify(header) }

  return request({
    url: `/api/webhooks`,
    method: 'post',
    data
  })
}

export function putWebhook(webhook) {
  const header = {}
  webhook.header.forEach(h => {
    header[h.key] = h.value
  })

  const data = { ...webhook, header: JSON.stringify(header) }

  return request({
    url: `/api/webhooks/${data.ID}`,
    method: 'put',
    data
  })
}

export function deleteWebhook(webhook) {
  return request({
    url: `/api/webhooks/${webhook.ID}`,
    method: 'delete'
  })
}
