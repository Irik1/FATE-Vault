import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

export const characterService = {
  async getCharacters() {
    const response = await api.get('/characters/list')
    return response.data
  },
  
  async getCharacter(id) {
    const response = await api.get(`/charactersList/${id}`)
    return response.data
  },
  
  async updateCharacter(id, character) {
    const response = await api.post(`/characters/update/${id}`, character)
    return response.data
  },
  
  async createCharacter(character) {
    const response = await api.post('/characters/create', character)
    return response.data
  },
  
  async getTemplates() {
    const response = await api.get('/templates')
    return response.data
  }
}

export default api


