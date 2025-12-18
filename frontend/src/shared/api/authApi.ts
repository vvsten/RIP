import { apiClient } from './axiosConfig';
import type { LoginRequest, RegisterRequest, AuthResponse, User, UpdateProfileRequest } from '../types/User';

/**
 * API функции для авторизации и работы с пользователями
 * Использует axios для HTTP запросов
 */

export const authApi = {
  /**
   * Вход пользователя
   */
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>('/api/users/login', credentials);
    return response.data;
  },

  /**
   * Регистрация нового пользователя
   */
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>('/api/users/register', data);
    return response.data;
  },

  /**
   * Выход пользователя
   */
  logout: async (): Promise<void> => {
    await apiClient.post('/api/users/logout');
  },

  /**
   * Получение профиля пользователя
   */
  getProfile: async (): Promise<User> => {
    const response = await apiClient.get<{ status: string; user: User }>('/api/users/profile');
    return response.data.user;
  },

  /**
   * Обновление профиля пользователя
   */
  updateProfile: async (data: UpdateProfileRequest): Promise<User> => {
    const response = await apiClient.put<{ status: string; user: User }>('/api/users/profile', data);
    return response.data.user;
  },
};

