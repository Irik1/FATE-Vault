import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

export const characterService = {
  async getCharacters() {
    const response = await api.get('/charactersList')
    return response.data
  },
  
  async getCharacter(id) {
    const response = await api.get(`/charactersList/${id}`)
    return response.data
  },
  
  async updateCharacter(id, character) {
    const response = await api.put(`/charactersList/${id}`, character)
    return response.data
  }
}

export default api


