import { useEffect, useState } from 'react';
import { getApiUrl } from '../../shared/config/apiConfig';

/**
 * Кнопка калькулятора/корзины под хедером
 * С бэйджем количества товаров в корзине
 */
export function CalculatorShortcut() {
  const [count, setCount] = useState<number>(0);
  const [logisticRequestId, setLogisticRequestId] = useState<number | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        // Совместимость: пробуем /api/cart (вернет id и count) и /api/cart/count
        const res = await fetch(getApiUrl('/api/cart'));
        if (res.ok) {
          const data = await res.json();
          const c = typeof data?.count === 'number' ? data.count : 0;
          const id = data?.cart?.id || data?.id || null;
          setCount(c);
          setLogisticRequestId(id);
          return;
        }
      } catch {}
      try {
        const res2 = await fetch(getApiUrl('/api/cart/count'));
        if (res2.ok) {
          const data2 = await res2.json();
          setCount(typeof data2?.count === 'number' ? data2.count : 0);
        }
      } catch {}
    };
    load();
    
    // Обновляем каждые 5 секунд
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, []);

  const href = logisticRequestId ? `/calculator?request_id=${logisticRequestId}` : '/calculator';
  const isDisabled = count <= 0;

  return (
    <div className="calculator-shortcut">
      {isDisabled ? (
        <a className="calculator-btn is-disabled" aria-disabled="true">
          🧮 Калькулятор
          <span className="cart-count" id="cartCount">{count || ''}</span>
        </a>
      ) : (
        <a href={href} className="calculator-btn" style={{ textDecoration: 'none' }}>
          🧮 Калькулятор
          <span className="cart-count" id="cartCount">{count}</span>
        </a>
      )}
    </div>
  );
}


