import request from '@/utils/requestServer'

export function getConfig() {
  return request({
    url: '/api/config',
    method: 'get'
  })
}

export function setConfig(data) {
  return request({
    url: '/api/config',
    method: 'post',
    data
  })
}

export function saveConfig(data) {
  return request({
    url: '/api/config/save',
    method: 'post',
    data
  })
}
