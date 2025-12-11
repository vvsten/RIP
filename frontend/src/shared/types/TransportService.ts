/**
 * Интерфейс услуги для транспортировки грузов
 * Соответствует модели Service из бэкенда
 * 
 * Примечание: Бэкенд возвращает JSON с полями в snake_case (image_url, delivery_days и т.д.)
 * Этот интерфейс использует camelCase для совместимости с TypeScript
 */
export interface TransportService {
  id: number;
  name: string;
  description: string;
  price: number;
  imageUrl?: string;
  deliveryDays: number;
  maxWeight: number;
  maxVolume: number;
  createdAt: string;
  updatedAt: string;
}

/**
 * Сырой формат сервиса из API (snake_case)
 * Используется для десериализации ответа от бэкенда
 */
export interface TransportServiceRaw {
  id: number;
  name: string;
  description: string;
  price: number;
  image_url?: string;
  delivery_days: number;
  max_weight: number;
  max_volume: number;
  created_at: string;
  updated_at: string;
}

/**
 * Параметры фильтрации услуг
 * Используется в компоненте Filters для передачи query-параметров в API
 */
export interface TransportServiceFilters {
  search?: string;      // Поиск по названию/описанию
  minPrice?: number;    // Минимальная цена
  maxPrice?: number;    // Максимальная цена
  dateFrom?: string;    // Дата с (в формате YYYY-MM-DD)
  dateTo?: string;      // Дата по (в формате YYYY-MM-DD)
}

