import axios from 'axios';
import { getApiUrl } from '../config/apiConfig';

/**
 * Базовый экземпляр axios для API запросов
 * 
 * Используется для всех запросов к бэкенду
 * Автоматически добавляет токен авторизации из localStorage
 * Обрабатывает ошибки авторизации
 */
export const apiClient = axios.create({
  baseURL: '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Интерцептор для добавления токена авторизации
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // Преобразуем относительный URL в полный через getApiUrl
    if (config.url && !config.url.startsWith('http')) {
      config.url = getApiUrl(config.url);
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Интерцептор для обработки ошибок
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Токен истек или невалиден
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('user');
      // Перенаправляем на страницу входа
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

