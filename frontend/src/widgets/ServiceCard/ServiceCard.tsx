import { getApiUrl } from '../../shared/config/apiConfig';
import type { TransportService } from '../../shared/types/TransportService';

/**
 * Props для компонента ServiceCard
 */
interface ServiceCardProps {
  /** Данные услуги для отображения */
  service: TransportService;
}

/**
 * Компонент карточки услуги
 * 
 * Отображает информацию об услуге в виде карточки с существующими стилями
 * Если изображения нет, подставляет изображение по-умолчанию
 * 
 * @param props - содержит объект service с данными услуги
 */
export function ServiceCard({ service }: ServiceCardProps) {
  // Получаем base URL из Vite (для GitHub Pages это /RIP-2-mod-/)
  const baseUrl = import.meta.env.BASE_URL || '/';
  // URL изображения по умолчанию если поле пустое
  const defaultImageUrl = `${baseUrl}assets/default.svg`;
  const imageUrl = service.imageUrl || defaultImageUrl;
  
  // Обработчик добавления в корзину
  const handleAddToCart = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    e.stopPropagation();
    
    try {
      const response = await fetch(getApiUrl(`/api/cart/add/${service.id}`), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        }
      });
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      
      if (data.success) {
        console.log('ServiceCard: товар добавлен, count:', data.count, 'request_id:', data.request_id);
        
        // Обновляем счетчик корзины через событие
        const event = new CustomEvent('cartUpdated', { 
          detail: { 
            count: data.count,
            request_id: data.request_id 
          } 
        });
        console.log('ServiceCard: отправляю событие cartUpdated', event.detail);
        window.dispatchEvent(event);
        
        // Показываем уведомление об успехе
        alert(data.message || 'Услуга добавлена в корзину');
      } else {
        alert('Ошибка: ' + (data.error || 'Не удалось добавить в корзину'));
      }
    } catch (error) {
      console.error('Ошибка при добавлении в корзину:', error);
      alert('Произошла ошибка при добавлении в корзину: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    }
  };
  
  return (
    <div className="service-card">
      <img 
        src={imageUrl} 
        alt={service.name}
        className="service-image"
        onError={(e) => {
          // Fallback если изображение не загрузилось
          (e.target as HTMLImageElement).src = defaultImageUrl;
        }}
      />
      <div className="service-content">
        <h3 className="service-title">{service.name}</h3>
        <p className="service-description">{service.description}</p>
        <div className="service-actions">
          <a href={`#service-${service.id}`} className="details-link">подробнее</a>
          <button 
            type="button"
            className="consult-btn" 
            onClick={handleAddToCart}
          >
            Получить консультацию
          </button>
        </div>
      </div>
    </div>
  );
}