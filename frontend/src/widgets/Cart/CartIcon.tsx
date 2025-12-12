import { useEffect, useState } from 'react';
import { getApiUrl } from '../../shared/config/apiConfig';

/**
 * Иконка корзины для страницы списка услуг
 * Использует метод /api/cart/icon, который всегда возвращает 0 или -1
 */
export function CartIcon() {
  const [count, setCount] = useState<number>(0);
  const [logisticRequestId, setLogisticRequestId] = useState<number | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        // Используем /api/cart/icon - всегда возвращает 200 OK
        // В Network tab будет видно count: 0 или -1
        // Но для отображения используем real_count (реальное количество)
        const res = await fetch(getApiUrl('/api/cart/icon'));
        if (res.ok) {
          const data = await res.json();
          // real_count - реальное количество для отображения на экране
          const realCount = typeof data?.real_count === 'number' ? data.real_count : 0;
          // count (0/-1) используется только для Network DevTools
          const statusValue = typeof data?.count === 'number' ? data.count : -1;
          const id = data?.request_id || null;
          
          console.log('CartIcon: загружены данные', { realCount, statusValue, id });
          // Используем реальное количество для отображения
          setCount(realCount);
          setLogisticRequestId(id);
        }
      } catch (err) {
        console.error('CartIcon: ошибка загрузки', err);
        // При ошибке считаем корзину пустой
        setCount(0);
      }
    };
    load();
    
    // Слушаем событие обновления корзины
    const handleCartUpdate = (event: Event) => {
      const customEvent = event as CustomEvent;
      console.log('CartIcon: получил событие cartUpdated', customEvent.detail);
      if (customEvent.detail?.count !== undefined) {
        const newCount = customEvent.detail.count;
        console.log('CartIcon: обновляю count на', newCount);
        // Используем count из события напрямую, не перезагружая данные
        setCount(newCount);
        // Обновляем request_id если он есть в событии
        if (customEvent.detail?.request_id !== undefined && customEvent.detail.request_id > 0) {
          console.log('CartIcon: обновляю request_id на', customEvent.detail.request_id);
          setLogisticRequestId(customEvent.detail.request_id);
        } else if (newCount > 0 && !logisticRequestId) {
          // Если request_id нет, но count > 0 и у нас еще нет ID, загружаем только ID
          // Делаем это с задержкой, чтобы БД успела обновиться
          console.log('CartIcon: загружаю request_id через 1 секунду');
          setTimeout(async () => {
            try {
              const res = await fetch(getApiUrl('/api/cart/icon'));
              if (res.ok) {
                const data = await res.json();
                if (data?.request_id && data.request_id > 0) {
                  console.log('CartIcon: получил request_id', data.request_id);
                  setLogisticRequestId(data.request_id);
                }
              }
            } catch (err) {
              console.error('CartIcon: ошибка загрузки request_id', err);
            }
          }, 1000);
        }
      }
    };
    
    window.addEventListener('cartUpdated', handleCartUpdate);
    
    // Обновляем каждые 5 секунд
    const interval = setInterval(load, 5000);
    return () => {
      clearInterval(interval);
      window.removeEventListener('cartUpdated', handleCartUpdate);
    };
  }, []);

  const href = logisticRequestId ? `/calculator?request_id=${logisticRequestId}` : '/calculator';
  const isDisabled = count <= 0; // кнопка неактивна если корзина пуста

  return (
    <div className="calculator-shortcut">
      {isDisabled ? (
        <a className="calculator-btn is-disabled" aria-disabled="true">
          🧮 Калькулятор
        </a>
      ) : (
        <a href={href} className="calculator-btn" style={{ textDecoration: 'none' }}>
          🧮 Калькулятор
          {count > 0 && <span className="cart-count" id="cartCount">{count}</span>}
        </a>
      )}
    </div>
  );
}

