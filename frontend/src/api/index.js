import request from './request'

export const authApi = {
  register: (data) => request.post('/auth/register', data),
  login: (data) => request.post('/auth/login', data)
}

export const storyApi = {
  list: (params) => request.get('/stories', { params }),
  detail: (id) => request.get(`/stories/${id}`),
  create: (data) => request.post('/stories', data),
  append: (id, add) => request.post(`/stories/${id}/append`, { add }),
  respond: (id, data) => request.post(`/stories/${id}/respond`, data),
  accept: (id, rid) => request.post(`/stories/${id}/accept/${rid}`)
}

export const userApi = {
  profile: () => request.get('/user/profile'),
  update: (data) => request.put('/user/profile', data)
}

export const uploadUrl = '/api/upload'
