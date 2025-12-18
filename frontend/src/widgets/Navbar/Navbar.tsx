import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../shared/store/hooks';
import { logoutUser } from '../../shared/store/slices/authSlice';
import { clearOrders } from '../../shared/store/slices/ordersSlice';
import { resetCart } from '../../shared/store/slices/cartSlice';
import { clearFilters } from '../../shared/store/slices/filtersSlice';

/**
 * Компонент навигационной панели
 * Использует существующие стили из style.css (header, logo, home-btn)
 * 
 * Отображает разные пункты меню для гостей и авторизованных пользователей
 * Показывает имя/логин пользователя после авторизации
 * Кнопка Вход/Выход переключает интерфейс
 */
export function Navbar() {
  const location = useLocation();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { isAuthenticated, user } = useAppSelector((state) => state.auth);

  const handleLogout = async () => {
    await dispatch(logoutUser());
    // Сбрасываем корзину и фильтры при выходе
    dispatch(resetCart());
    dispatch(clearFilters());
    dispatch(clearOrders());
    navigate('/');
  };

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
        {location.pathname !== '/transport-services' && (
          <Link to="/transport-services" className="home-btn">📦 Услуги</Link>
        )}
        {location.pathname !== '/about' && (
          <Link to="/about" className="home-btn">ℹ️ О компании</Link>
        )}
        
        {/* Пункты меню для авторизованных пользователей */}
        {isAuthenticated && (
          <>
            {location.pathname !== '/orders' && (
              <Link to="/orders" className="home-btn">📋 Мои заявки</Link>
            )}
            {location.pathname !== '/profile' && (
              <Link to="/profile" className="home-btn">👤 Профиль</Link>
            )}
            {/* Отображение имени/логина пользователя */}
            <span style={{ 
              color: 'white', 
              margin: '0 1rem',
              fontWeight: 'bold',
              fontSize: '0.9rem'
            }}>
              {user?.name || user?.login}
            </span>
            <button
              onClick={handleLogout}
              className="home-btn"
              style={{
                background: 'transparent',
                border: 'none',
                cursor: 'pointer',
                color: 'inherit',
                fontFamily: 'inherit',
                fontSize: 'inherit',
              }}
            >
              Выход
            </button>
          </>
        )}
        
        {/* Кнопка входа для гостей */}
        {!isAuthenticated && location.pathname !== '/login' && (
          <Link to="/login" className="home-btn">Вход</Link>
        )}
      </div>
    </header>
  );
}
