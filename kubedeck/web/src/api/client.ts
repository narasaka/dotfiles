import ky from 'ky'

export const api = ky.create({
  prefixUrl: '/api/v1',
  credentials: 'include',
  hooks: {
    afterResponse: [
      async (_request, _options, response) => {
        if (response.status === 401 && !response.url.includes('/auth/')) {
          window.location.href = '/login'
        }
      },
    ],
  },
})
