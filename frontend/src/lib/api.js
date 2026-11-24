const API_BASE_URL = ''

export class ApiError extends Error {
  constructor(message, status, data = null) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

export async function apiRequest(path, options = {}) {
  const token = localStorage.getItem('token')

  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  let response

  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...options,
      headers,
    })
  } catch {
    throw new ApiError(
      'Tidak dapat terhubung ke server. Pastikan backend Go sedang berjalan.',
      0
    )
  }

  const text = await response.text()
  let result = {}

  try {
    result = text ? JSON.parse(text) : {}
  } catch {
    result = { message: text }
  }

  if (!response.ok) {
    const message =
      result.message ||
      result.error ||
      result.data?.message ||
      'Terjadi kesalahan pada server.'

    if (response.status === 401 && token) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')

      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }

    if (response.status === 403) {
      throw new ApiError(
        'Akses ditolak. Kamu tidak memiliki izin untuk melakukan tindakan ini.',
        403,
        result
      )
    }

    throw new ApiError(message, response.status, result)
  }

  return result
}

export const api = {
  get: (path) => apiRequest(path),

  post: (path, body) =>
    apiRequest(path, {
      method: 'POST',
      body: JSON.stringify(body),
    }),

  put: (path, body) =>
    apiRequest(path, {
      method: 'PUT',
      body: JSON.stringify(body),
    }),

  delete: (path) =>
    apiRequest(path, {
      method: 'DELETE',
    }),
}
