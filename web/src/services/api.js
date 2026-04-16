import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const characterService = {
  async getCharacters() {
    const response = await api.get('/characters')
    return response.data
  },

  /**
   * Fetch one character by id via GET /characters/find?characterIds=:id
   */
  async getCharacter(id) {
    const response = await api.get('/characters/find', {
      params: { characterIds: id }
    })
    const data = response.data
    if (Array.isArray(data) && data.length > 0) return data[0]
    return null
  },

  /**
   * @param {{ edition?: string, name?: string, characterIds?: string|string[] }} params
   */
  async findCharacters(params = {}) {
    const response = await api.get('/characters/find', { params })
    return response.data
  },

  async createCharacter(character) {
    const response = await api.post('/characters/create', character)
    return response.data
  },

  async updateCharacter(id, character) {
    const response = await api.post(`/characters/update/${id}`, character)
    return response.data
  },

  async deleteCharacter(id) {
    const response = await api.delete(`/characters/delete/${id}`)
    return response.data
  },

  async getTemplates() {
    const response = await api.get('/templates')
    return response.data
  }
}

export const categoryService = {
  async list() {
    const response = await api.get('/categories')
    return response.data
  },
  async create(body) {
    const response = await api.post('/categories/create', body)
    return response.data
  },
  async update(id, body) {
    const response = await api.post(`/categories/update/${id}`, body)
    return response.data
  },
  async remove(id) {
    const response = await api.delete(`/categories/delete/${id}`)
    return response.data
  }
}

export const gameService = {
  async list() {
    const response = await api.get('/games')
    return response.data
  },
  async create(body) {
    const response = await api.post('/games/create', body)
    return response.data
  },
  async update(id, body) {
    const response = await api.post(`/games/update/${id}`, body)
    return response.data
  },
  async remove(id) {
    const response = await api.delete(`/games/delete/${id}`)
    return response.data
  }
}

export const stuntService = {
  async list() {
    const response = await api.get('/stunts')
    return response.data
  },
  async create(body) {
    const response = await api.post('/stunts/create', body)
    return response.data
  },
  async update(id, body) {
    const response = await api.post(`/stunts/update/${id}`, body)
    return response.data
  },
  async remove(id) {
    const response = await api.delete(`/stunts/delete/${id}`)
    return response.data
  }
}

export const userService = {
  /**
   * @param {{ username: string, password: string, role?: string }} body
   * @returns {Promise<{ user: object }>} Sets HttpOnly session cookie on success.
   */
  async register(body) {
    const response = await api.post('/users/register', body)
    return response.data
  },

  /**
   * Login with username/password. Session is the HttpOnly cookie; response body is `{ user }` only.
   * @param {{ username: string, password: string }} body
   */
  async auth(body) {
    const response = await api.post('/users/auth', body)
    return response.data
  },

  /**
   * Returns current authenticated user via middleware-protected route.
   */
  async me() {
    const response = await api.get('/users/me')
    return response.data
  },

  async logout() {
    await api.post('/users/logout')
  },

  /**
   * @param {string} id user _id
   * @param {object} body partial user fields
   */
  async update(id, body) {
    const response = await api.post(`/users/update/${id}`, body)
    return response.data
  }
}
