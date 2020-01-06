import request from '@/utils/requestServer'

export function initialize(username, password) {
  return request(
    '/api/init',
    {
      method: 'POST',
      json: true,
      data: {
        username,
        password
      }
    }
  )
}

export function isLogined() {
  return request(
    '/api/ping'
  ).then(res => {
    console.log(res)
    return true
  }).catch(err => {
    console.log('err', err.response)
    if (err.response.status === 404) {
      console.log('noinit')
      return 'NoInit'
    } else {
      console.log('nologin')
      return 'NoLogin'
    }
  })
}

export function login(username, password) {
  return request(
    '/api/login',
    {
      method: 'POST',
      json: true,
      data: {
        username,
        password
      }
    }
  )
}

export function logout() {
  return request(
    '/api/logout',
    {
      method: 'POST'
    }
  )
}
