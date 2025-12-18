import { useState } from 'react';
import { useAppDispatch, useAppSelector } from '../../shared/store/hooks';
import { addToCart } from '../../shared/store/slices/cartSlice';
import { useNavigate } from 'react-router-dom';
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
 * Использует Redux для добавления услуги в корзину (заявку)
 * 
 * @param props - содержит объект service с данными услуги
 */
export function ServiceCard({ service }: ServiceCardProps) {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const { isAuthenticated } = useAppSelector((state) => state.auth);
  const { isLoading } = useAppSelector((state) => state.cart);
  const [adding, setAdding] = useState(false);

  // Получаем base URL из Vite (для GitHub Pages это /RIP-2-mod-/)
  const baseUrl = import.meta.env.BASE_URL || '/';
  // URL изображения по умолчанию если поле пустое
  const defaultImageUrl = `${baseUrl}assets/default.svg`;
  const imageUrl = service.imageUrl || defaultImageUrl;
  
  // Обработчик добавления в корзину (заявку)
  const handleAddToCart = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    e.stopPropagation();
    
    // Если пользователь не авторизован, перенаправляем на страницу входа
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    
    setAdding(true);
    try {
      const result = await dispatch(addToCart(service.id));
      if (addToCart.fulfilled.match(result)) {
        // Показываем уведомление об успехе
        alert('Услуга добавлена в заявку');
      } else {
        alert('Ошибка: ' + (result.payload as string || 'Не удалось добавить в заявку'));
      }
    } catch (error) {
      console.error('Ошибка при добавлении в корзину:', error);
      alert('Произошла ошибка при добавлении в заявку');
    } finally {
      setAdding(false);
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
            disabled={adding || isLoading}
            style={{
              opacity: (adding || isLoading) ? 0.6 : 1,
              cursor: (adding || isLoading) ? 'not-allowed' : 'pointer',
            }}
          >
            {adding || isLoading ? 'Добавление...' : 'Добавить'}
          </button>
        </div>
      </div>
    </div>
  );
}