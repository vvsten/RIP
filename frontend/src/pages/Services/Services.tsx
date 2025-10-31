import { useState, useEffect } from 'react';
import { Filters } from '../../widgets/Filters/Filters';
import { ServiceCard } from '../../widgets/ServiceCard/ServiceCard';
import { fetchServices } from '../../shared/api/servicesApi';
import type { Service, ServiceFilters } from '../../shared/types/Service';

/**
 * Страница списка услуг с фильтрацией
 * 
 * Использует useState для управления состоянием списка услуг и фильтров
 * Использует useEffect для загрузки данных при монтировании компонента
 * 
 * Особенности:
 * - Фильтрация на бэкенде через API
 * - Fallback на mock данные при недоступности сервера
 * - Обработка состояний загрузки и ошибок
 */
export function Services() {
  // useState для списка услуг
  const [services, setServices] = useState<Service[]>([]);
  
  // useState для состояния загрузки
  const [loading, setLoading] = useState(true);
  
  // useState для обработки ошибок
  const [error, setError] = useState<string | null>(null);
  
  // useEffect вызывается при монтировании компонента
  // Загружает все услуги без фильтров
  useEffect(() => {
    loadServices();
  }, []);
  
  /**
   * Загрузка услуг с сервера
   * @param filters - опциональные параметры фильтрации
   */
  const loadServices = async (filters: ServiceFilters = {}) => {
    setLoading(true);
    setError(null);
    
    try {
      // Вызываем fetchServices из servicesApi
      // Функция автоматически обрабатывает fallback на mock
      const data = await fetchServices(filters);
      setServices(data);
    } catch (err) {
      setError('Не удалось загрузить услуги');
      console.error('Error loading services:', err);
    } finally {
      setLoading(false);
    }
  };
  
  /**
   * Обработчик изменения фильтров
   * Вызывается из компонента Filters
   * @param filters - новые параметры фильтрации
   */
  const handleFilterChange = (filters: ServiceFilters) => {
    loadServices(filters);
  };
  
  return (
    <div className="container">
      <h2 style={{ marginBottom: '2rem', fontSize: '2rem', fontWeight: 'bold' }}>
        Наши услуги
      </h2>
      
      {/* Компонент фильтров */}
      <Filters onFilterChange={handleFilterChange} />
      
      {/* Индикатор загрузки */}
      {loading && (
        <div style={{ textAlign: 'center', padding: '2rem' }}>
          <div style={{ 
            border: '3px solid #f3f3f3',
            borderTop: '3px solid #0d6efd',
            borderRadius: '50%',
            width: '40px',
            height: '40px',
            animation: 'spin 1s linear infinite',
            margin: '0 auto'
          }} />
          <p style={{ marginTop: '1rem', color: '#6c757d' }}>Загрузка...</p>
        </div>
      )}
      
      {/* Обработка ошибок */}
      {error && (
        <div style={{ 
          background: '#f8d7da',
          color: '#721c24',
          padding: '1rem',
          borderRadius: '6px',
          marginBottom: '2rem'
        }}>
          {error}
        </div>
      )}
      
      {/* Список услуг */}
      {!loading && !error && (
        services.length > 0 ? (
          <div className="services-grid">
            {services.map((service) => (
              <ServiceCard key={service.id} service={service} />
            ))}
          </div>
        ) : (
          <div className="no-services">
            <h2>Услуги не найдены</h2>
            <p>Попробуйте изменить параметры фильтрации</p>
          </div>
        )
      )}
      
      <style>{`
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}