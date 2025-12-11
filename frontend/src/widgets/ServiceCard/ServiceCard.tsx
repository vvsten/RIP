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
  // URL изображения по умолчанию если поле пустое
  const imageUrl = service.imageUrl || '/assets/default.svg';
  
  return (
    <div className="service-card">
      <img 
        src={imageUrl} 
        alt={service.name}
        className="service-image"
        onError={(e) => {
          // Fallback если изображение не загрузилось
          (e.target as HTMLImageElement).src = '/assets/default.svg';
        }}
      />
      <div className="service-content">
        <h3 className="service-title">{service.name}</h3>
        <p className="service-description">{service.description}</p>
        <div className="service-actions">
          <a href={`#service-${service.id}`} className="details-link">подробнее</a>
          <button className="consult-btn" data-service-id={service.id}>
            Получить консультацию
          </button>
        </div>
      </div>
    </div>
  );
}