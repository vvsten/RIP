import { Link, useLocation } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { getApiUrl } from '../../shared/config/apiConfig';

/**
 * Компонент навигационной панели
 * Использует существующие стили из style.css (header, logo, home-btn)
 * 
 * Props: не требуются (использует useLocation из react-router-dom для определения активной страницы)
 */
export function Navbar() {
  const location = useLocation();
  const [cartCount, setCartCount] = useState<number>(0);
  const [orderId, setOrderId] = useState<number | null>(null);

  // Загружаем количество в корзине при монтировании и при изменении страницы
  useEffect(() => {
    const loadCart = async () => {
      try {
        const res = await fetch(getApiUrl('/api/cart'));
        if (res.ok) {
          const data = await res.json();
          const count = typeof data?.count === 'number' ? data.count : 0;
          const id = data?.cart?.id || data?.id || null;
          setCartCount(count);
          setOrderId(id);
          return;
        }
      } catch {}
      try {
        const res2 = await fetch(getApiUrl('/api/cart/count'));
        if (res2.ok) {
          const data2 = await res2.json();
          setCartCount(typeof data2?.count === 'number' ? data2.count : 0);
        }
      } catch {}
    };
    loadCart();
  }, [location.pathname]);

  const calculatorHref = orderId ? `/calculator?order_id=${orderId}` : '/calculator';
  const isCalculatorDisabled = cartCount <= 0;

  return (
    <header className="header">
      <Link to="/" className="logo">
        <div className="logo-icon">🚚</div>
        GruzDelivery
      </Link>
      <div className="header-actions">
        {/* Кнопки навигации */}
        {location.pathname !== '/' && (
          <Link to="/" className="home-btn">🏠 Главная</Link>
        )}
        {location.pathname !== '/services' && (
          <Link to="/services" className="home-btn">📦 Услуги</Link>
        )}
        {location.pathname !== '/about' && (
          <Link to="/about" className="home-btn">ℹ️ О компании</Link>
        )}
        {/* Кнопка калькулятора с бейджем количества */}
        {isCalculatorDisabled ? (
          <span className="home-btn" style={{ opacity: 0.5, cursor: 'not-allowed' }}>
            🧮 Калькулятор
            {cartCount > 0 && <span style={{ 
              marginLeft: '0.5rem',
              background: '#ff4444',
              color: 'white',
              borderRadius: '50%',
              padding: '0.1rem 0.4rem',
              fontSize: '0.8rem'
            }}>{cartCount}</span>}
          </span>
        ) : (
          <Link to={calculatorHref} className="home-btn" style={{ position: 'relative' }}>
            🧮 Калькулятор
            {cartCount > 0 && <span style={{ 
              position: 'absolute',
              top: '-8px',
              right: '-8px',
              background: '#ff4444',
              color: 'white',
              borderRadius: '50%',
              width: '20px',
              height: '20px',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '0.7rem',
              fontWeight: 'bold',
              border: '2px solid white',
              boxShadow: '0 2px 6px rgba(255, 68, 68, 0.4)'
            }}>{cartCount}</span>}
          </Link>
        )}
      </div>
    </header>
  );
}