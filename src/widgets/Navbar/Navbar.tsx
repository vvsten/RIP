import { Link, useLocation } from 'react-router-dom';

/**
 * Компонент навигационной панели
 * Использует существующие стили из style.css (header, logo, home-btn)
 * 
 * Props: не требуются (использует useLocation из react-router-dom для определения активной страницы)
 */
export function Navbar() {
  const location = useLocation();
  
  return (
    <header className="header">
      <Link to="/" className="logo">
        <div className="logo-icon">🚚</div>
        GruzDelivery
      </Link>
      <div className="header-actions">
        {location.pathname !== '/' && (
          <Link to="/" className="home-btn">🏠 Главная</Link>
        )}
        {location.pathname !== '/TransportService' && (
          <Link to="/TransportService" className="home-btn">📦 Услуги</Link>
        )}
        {location.pathname !== '/about' && (
          <Link to="/about" className="home-btn">ℹ️ О компании</Link>
        )}
      </div>
    </header>
  );
}