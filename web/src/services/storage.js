import axios from 'axios'

// Storage service base URL - using proxy in development, or direct URL in production
const STORAGE_BASE_URL = import.meta.env.VITE_STORAGE_URL || '/storage'

const storageApi = axios.create({
  baseURL: STORAGE_BASE_URL,
  headers: {
    'Content-Type': 'multipart/form-data'
  }
})

export const storageService = {
  /**
   * Upload a file to storage
   * @param {File} file - The file to upload
   * @param {string} folder - Optional folder path (e.g., 'images/characters')
   * @returns {Promise} Response with filename and URL
   */
  async uploadFile(file, folder = '') {
    const formData = new FormData()
    formData.append('file', file)
    if (folder) {
      formData.append('folder', folder)
    }

    const response = await storageApi.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data
  },

  /**
   * Get the download URL for a file
   * @param {string} filename - The filename (can include folder path)
   * @returns {string} Full URL to download the file via storage service (which proxies to MinIO)
   */
  getFileUrl(filename) {
    // Remove leading slash if present
    let cleanFilename = filename.startsWith('/') ? filename.slice(1) : filename
    
    // Use storage service endpoint which handles MinIO authentication
    // Format: /storage/download/images/characters/...
    return `${STORAGE_BASE_URL}/download/${cleanFilename}`
  },

  /**
   * Delete a file from storage
   * @param {string} filename - The filename to delete
   * @returns {Promise} Response
   */
  async deleteFile(filename) {
    // Remove leading slash if present
    const cleanFilename = filename.startsWith('/') ? filename.slice(1) : filename
    const response = await storageApi.delete(`/delete/${cleanFilename}`)
    return response.data
  },

  /**
   * List files in storage
   * @param {string} folder - Optional folder to filter by
   * @returns {Promise} List of files
   */
  async listFiles(folder = '') {
    const params = folder ? { folder } : {}
    const response = await storageApi.get('/list', { params })
    return response.data
  },

  /**
   * Check if storage service is healthy
   * @returns {Promise} Health status
   */
  async healthCheck() {
    const response = await storageApi.get('/health')
    return response.data
  }
}

export default storageService

