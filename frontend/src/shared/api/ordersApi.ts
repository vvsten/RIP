import { apiClient } from './axiosConfig';
import type { OrdersListResponse, OrderResponse, LogisticRequest } from '../types/Order';

/**
 * API функции для работы с заявками
 * Использует axios для HTTP запросов
 */

export interface GetOrdersParams {
  status?: string;
  date_from?: string;
  date_to?: string;
}

export const ordersApi = {
  /**
   * Получение списка заявок с фильтрацией
   */
  getOrders: async (params?: GetOrdersParams): Promise<LogisticRequest[]> => {
    const response = await apiClient.get<OrdersListResponse>('/api/orders', { params });
    return response.data.orders;
  },

  /**
   * Получение заявки по ID
   */
  getOrder: async (id: number): Promise<LogisticRequest> => {
    const response = await apiClient.get<OrderResponse>(`/api/orders/${id}`);
    return response.data.order;
  },

  /**
   * Обновление заявки (только для черновиков)
   */
  updateOrder: async (id: number, data: {
    from_city?: string;
    to_city?: string;
    weight?: number;
    length?: number;
    width?: number;
    height?: number;
  }): Promise<LogisticRequest> => {
    const response = await apiClient.put<OrderResponse>(`/api/orders/${id}`, data);
    return response.data.order;
  },

  /**
   * Удаление услуги из заявки
   */
  removeServiceFromOrder: async (orderId: number, serviceId: number): Promise<void> => {
    await apiClient.delete(`/api/orders/${orderId}/services/${serviceId}`);
  },

  /**
   * Обновление количества услуги в заявке
   */
  updateServiceInOrder: async (
    orderId: number,
    serviceId: number,
    quantity: number
  ): Promise<void> => {
    await apiClient.put(`/api/orders/${orderId}/services/${serviceId}`, { quantity });
  },

  /**
   * Подтверждение заявки (формирование)
   */
  formOrder: async (id: number, data: {
    from_city: string;
    to_city: string;
    weight: number;
    length: number;
    width: number;
    height: number;
  }): Promise<LogisticRequest> => {
    const response = await apiClient.put<OrderResponse>(`/api/orders/${id}/form`, data);
    // После формирования нужно получить обновленную заявку
    const orderResponse = await apiClient.get<OrderResponse>(`/api/orders/${id}`);
    return orderResponse.data.order;
  },

  /**
   * Отправка заявки на грузоперевозку
   */
  submitOrder: async (data: {
    from_city: string;
    to_city: string;
    weight: number;
    length: number;
    width: number;
    height: number;
  }): Promise<LogisticRequest> => {
    const response = await apiClient.post<OrderResponse>('/api/submitcargoorder', data);
    return response.data.order;
  },
};

