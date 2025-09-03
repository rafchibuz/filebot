import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || '';

// Создаем экземпляр axios с базовой конфигурацией
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 секунд
});

// Перехватчик для обработки ошибок
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error);
    
    if (error.response?.status === 500) {
      throw new Error('Ошибка сервера. Попробуйте позже.');
    } else if (error.response?.status === 404) {
      throw new Error('Данные не найдены.');
    } else if (error.code === 'ECONNABORTED') {
      throw new Error('Превышено время ожидания. Проверьте соединение.');
    } else if (!error.response) {
      throw new Error('Нет соединения с сервером.');
    }
    
    throw error;
  }
);

export const clipsAPI = {
  // Получение клипов с фильтрами
  async getClips(filters = {}) {
    try {
      const params = new URLSearchParams();
      
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== null && value !== undefined && value !== '') {
          params.append(key, value);
        }
      });

      const response = await apiClient.get(`/api/clips?${params.toString()}`);
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // Принудительное обновление клипов
  async refreshClips(daysBack = 7) {
    try {
      const response = await apiClient.get(`/api/clips/refresh?daysBack=${daysBack}`);
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // Получение списка стримеров
  async getStreamers() {
    try {
      const response = await apiClient.get('/api/streamers');
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // Проверка состояния сервера
  async healthCheck() {
    try {
      const response = await apiClient.get('/api/health');
      return response.data;
    } catch (error) {
      throw error;
    }
  }
};

export default apiClient;