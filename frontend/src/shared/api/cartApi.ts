import { apiClient } from './axiosConfig';
import type { LogisticRequest } from '../types/Order';

/**
 * API функции для работы с корзиной
 * Использует axios для HTTP запросов
 */

export interface CartResponse {
  cart: LogisticRequest;
  services: any[];
  count: number;
}

export const cartApi = {
  /**
   * Добавление услуги в корзину
   */
  addToCart: async (serviceId: number): Promise<{ success: boolean; count: number; request_id: number; message: string }> => {
    const response = await apiClient.post(`/api/cart/add/${serviceId}`);
    return response.data;
  },

  /**
   * Получение корзины
   */
  getCart: async (): Promise<CartResponse> => {
    const response = await apiClient.get<CartResponse>('/api/cart');
    return response.data;
  },

  /**
   * Получение количества товаров в корзине
   */
  getCartCount: async (): Promise<number> => {
    const response = await apiClient.get<{ count: number }>('/api/cart/icon');
    return response.data.count || 0;
  },

  /**
   * Очистка корзины
   */
  clearCart: async (): Promise<void> => {
    await apiClient.delete('/api/cart');
  },
};

