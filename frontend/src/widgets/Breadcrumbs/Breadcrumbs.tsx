import { useLocation, Link } from 'react-router-dom';

/**
 * Конфигурация breadcrumbs для разных страниц
 * Ключ - путь страницы, значение - массив элементов навигационной цепочки
 */
const BREADCRUMBS_CONFIG: Record<string, Array<{ label: string; path: string }>> = {
  '/': [{ label: 'Главная', path: '/' }],
  '/transport-services': [
    { label: 'Главная', path: '/' },
    { label: 'Услуги', path: '/transport-services' }
  ],
  '/about': [
    { label: 'Главная', path: '/' },
    { label: 'О компании', path: '/about' }
  ]
};

/**
 * Компонент навигационной цепочки (breadcrumbs)
 * Самописная реализация без использования библиотек
 * 
 * Использует useLocation для определения текущего пути
 * Отрисовывает путь от главной страницы до текущей
 * 
 * Props: не требуются
 */
export function Breadcrumbs() {
  const location = useLocation();
  const items = BREADCRUMBS_CONFIG[location.pathname] || [];
  
  // Не показываем breadcrumbs на главной странице (только одна секция)
  if (items.length <= 1) {
    return null;
  }
  
  return (
    <div className="container" style={{ marginTop: '1rem', marginBottom: '1rem' }}>
      <nav style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', fontSize: '0.9rem' }}>
        {items.map((item, index) => (
          <span key={item.path} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
            {index > 0 && <span style={{ color: '#6c757d' }}>/</span>}
            {index === items.length - 1 ? (
              <span style={{ fontWeight: 'bold', color: '#212529' }}>{item.label}</span>
            ) : (
              <Link 
                to={item.path} 
                style={{ color: '#0d6efd', textDecoration: 'none' }}
              >
                {item.label}
              </Link>
            )}
          </span>
        ))}
      </nav>
    </div>
  );
}