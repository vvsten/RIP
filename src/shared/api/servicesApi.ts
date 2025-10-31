import type { Service, ServiceFilters, ServiceRaw } from '../types/Service';

/**
 * Конвертирует сырой формат Service из API (snake_case) в camelCase
 * Также преобразует полный URL MinIO в относительный путь для прокси
 */
function convertService(raw: ServiceRaw): Service {
  // Преобразуем полный URL в относительный путь для прокси
  let imageUrl = raw.image_url;
  if (imageUrl && imageUrl.startsWith('http://localhost:9003/')) {
    // Извлекаем путь после хоста
    imageUrl = imageUrl.replace('http://localhost:9003', '');
  }
  
  return {
    id: raw.id,
    name: raw.name,
    description: raw.description,
    price: raw.price,
    imageUrl: imageUrl,
    deliveryDays: raw.delivery_days,
    maxWeight: raw.max_weight,
    maxVolume: raw.max_volume,
    createdAt: raw.created_at,
    updatedAt: raw.updated_at
  };
}

/**
 * Mock данные для fallback когда бэкенд недоступен
 * Имитируют реальные данные услуг по транспортировке грузов
 */
const MOCK_SERVICES: Service[] = [
  {
    id: 1,
    name: 'Фура',
    description: 'Грузоперевозки на фурах для больших объемов. Идеально для габаритных грузов.',
    price: 50000,
    deliveryDays: 7,
    maxWeight: 20000,
    maxVolume: 80,
    imageUrl: undefined, // Используется default.svg
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-01-15T10:00:00Z'
  },
  {
    id: 2,
    name: 'Малотоннажный',
    description: 'Быстрые грузоперевозки на малотоннажных автомобилях. Подходит для малых партий.',
    price: 15000,
    deliveryDays: 3,
    maxWeight: 3000,
    maxVolume: 15,
    imageUrl: undefined,
    createdAt: '2024-01-10T10:00:00Z',
    updatedAt: '2024-01-10T10:00:00Z'
  },
  {
    id: 3,
    name: 'Авиа',
    description: 'Скоростная доставка грузов по воздуху. Самый быстрый способ доставки.',
    price: 150000,
    deliveryDays: 1,
    maxWeight: 5000,
    maxVolume: 25,
    imageUrl: undefined,
    createdAt: '2024-01-20T10:00:00Z',
    updatedAt: '2024-01-20T10:00:00Z'
  },
  {
    id: 4,
    name: 'Поезд',
    description: 'Надежные железнодорожные перевозки. Оптимально для больших объемов на дальние расстояния.',
    price: 80000,
    deliveryDays: 14,
    maxWeight: 40000,
    maxVolume: 100,
    imageUrl: undefined,
    createdAt: '2024-01-05T10:00:00Z',
    updatedAt: '2024-01-05T10:00:00Z'
  },
  {
    id: 5,
    name: 'Корабль',
    description: 'Морские грузоперевозки. Экономичный вариант для международной доставки.',
    price: 120000,
    deliveryDays: 30,
    maxWeight: 100000,
    maxVolume: 500,
    imageUrl: undefined,
    createdAt: '2024-01-12T10:00:00Z',
    updatedAt: '2024-01-12T10:00:00Z'
  },
  {
    id: 6,
    name: 'Мультимодальный',
    description: 'Комбинированная доставка разными видами транспорта. Максимальная гибкость.',
    price: 100000,
    deliveryDays: 10,
    maxWeight: 30000,
    maxVolume: 150,
    imageUrl: undefined,
    createdAt: '2024-01-18T10:00:00Z',
    updatedAt: '2024-01-18T10:00:00Z'
  }
];

/**
 * Конвертирует объект фильтров в query string для URL
 * @param filters - параметры фильтрации
 * @returns query string (например: "?search=truck&minPrice=1000")
 */
function buildQueryString(filters: ServiceFilters): string {
  const params = new URLSearchParams();
  
  if (filters.search) params.append('search', filters.search);
  if (filters.minPrice !== undefined) params.append('minPrice', filters.minPrice.toString());
  if (filters.maxPrice !== undefined) params.append('maxPrice', filters.maxPrice.toString());
  if (filters.dateFrom) params.append('dateFrom', filters.dateFrom);
  if (filters.dateTo) params.append('dateTo', filters.dateTo);
  
  const queryString = params.toString();
  return queryString ? `?${queryString}` : '';
}

/**
 * Имитирует серверную фильтрацию на mock данных
 * Используется когда бэкенд недоступен
 */
function filterMockServices(services: Service[], filters: ServiceFilters): Service[] {
  return services.filter(service => {
    // Поиск по названию и описанию
    if (filters.search) {
      const searchLower = filters.search.toLowerCase();
      const matchesSearch = 
        service.name.toLowerCase().includes(searchLower) ||
        service.description.toLowerCase().includes(searchLower);
      if (!matchesSearch) return false;
    }
    
    // Фильтр по цене
    if (filters.minPrice !== undefined && service.price < filters.minPrice) return false;
    if (filters.maxPrice !== undefined && service.price > filters.maxPrice) return false;
    
    // Фильтр по дате создания
    // dateFrom и dateTo приходят в формате YYYY-MM-DD, сравниваем как строки
    // createdAt в формате ISO (2024-01-15T10:00:00Z), берем только дату
    if (filters.dateFrom) {
      const serviceDate = service.createdAt.split('T')[0]; // YYYY-MM-DD
      if (serviceDate < filters.dateFrom) return false;
    }
    if (filters.dateTo) {
      const serviceDate = service.createdAt.split('T')[0]; // YYYY-MM-DD
      if (serviceDate > filters.dateTo) return false;
    }
    
    return true;
  });
}

/**
 * Получение списка услуг с сервера или mock данных
 * 
 * @param filters - параметры фильтрации (опционально)
 * @returns Promise с массивом услуг
 * 
 * Логика работы:
 * 1. Делает fetch запрос к бэкенду через proxy (/api/services)
 * 2. Если сервер недоступен → возвращает mock данные
 * 3. Если сервер доступен → серверная фильтрация, иначе локальная на mock
 */
export async function fetchServices(filters: ServiceFilters = {}): Promise<Service[]> {
  try {
    const queryString = buildQueryString(filters);
    const response = await fetch(`/api/services${queryString}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Если ответ содержит поле services, берем его
    let services: ServiceRaw[] = [];
    if (data.services && Array.isArray(data.services)) {
      services = data.services;
    } else if (Array.isArray(data)) {
      services = data;
    } else {
      throw new Error('Unexpected response format');
    }
    
    // Конвертируем из snake_case в camelCase
    return services.map(convertService);
  } catch (error) {
    console.warn('Failed to fetch services from backend, using mock data:', error);
    
    // Fallback на mock данные при недоступности сервера
    // Применяем локальную фильтрацию если сервер не ответил
    return filterMockServices(MOCK_SERVICES, filters);
  }
}

/**
 * Получение одной услуги по ID
 */
export async function fetchService(id: number): Promise<Service | null> {
  try {
    const response = await fetch(`/api/services/${id}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    
    // Если ответ содержит поле service
    if (data.service) {
      return convertService(data.service);
    }
    
    return null;
  } catch (error) {
    console.warn(`Failed to fetch service ${id} from backend:`, error);
    
    // Fallback на mock данные
    return MOCK_SERVICES.find(s => s.id === id) || null;
  }
}

