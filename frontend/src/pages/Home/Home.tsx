/**
 * Главная страница приложения
 * 
 * Статическая информационная страница с описанием сервиса
 * Использует существующие стили для красивого отображения
 */
export function Home() {
  return (
    <div className="container">
      <div style={{
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '4rem 2rem',
        borderRadius: '15px',
        color: 'white',
        textAlign: 'center',
        marginBottom: '3rem'
      }}>
        <h1 style={{ fontSize: '3rem', fontWeight: 'bold', marginBottom: '1rem' }}>
          Добро пожаловать в GruzDelivery
        </h1>
        <p style={{ fontSize: '1.3rem', opacity: 0.95 }}>
          Ваш надежный партнер в грузоперевозках
        </p>
      </div>
      
      <div style={{ 
        background: '#ffffff',
        padding: '2rem',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        marginBottom: '2rem'
      }}>
        <p style={{ fontSize: '1.1rem', color: '#6c757d', marginBottom: '1rem' }}>
          Мы предлагаем полный спектр услуг по транспортировке грузов различными видами транспорта.
          От быстрой доставки малых партий до перевозки крупногабаритных грузов на дальние расстояния.
        </p>
        <p style={{ fontSize: '1.1rem', color: '#6c757d' }}>
          Наша компания работает на рынке грузоперевозок уже более 10 лет, обеспечивая 
          качественный сервис и надежную доставку для наших клиентов.
        </p>
      </div>
      
      <h3 style={{ fontSize: '1.8rem', fontWeight: 'bold', marginBottom: '1.5rem', color: '#212529' }}>
        Почему выбирают нас:
      </h3>
      
      <div className="services-grid" style={{ 
        gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))',
        marginBottom: '3rem'
      }}>
        <div style={{
          background: '#ffffff',
          padding: '2rem',
          borderRadius: '8px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>📦</div>
          <h4 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
            Более 10 лет опыта
          </h4>
          <p style={{ color: '#6c757d' }}>
            Мы на рынке логистики с 2014 года
          </p>
        </div>
        
        <div style={{
          background: '#ffffff',
          padding: '2rem',
          borderRadius: '8px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>🚚</div>
          <h4 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
            Разные виды транспорта
          </h4>
          <p style={{ color: '#6c757d' }}>
            Фуры, малотоннажные, авиа, поезд, корабль
          </p>
        </div>
        
        <div style={{
          background: '#ffffff',
          padding: '2rem',
          borderRadius: '8px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>💰</div>
          <h4 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
            Конкурентные цены
          </h4>
          <p style={{ color: '#6c757d' }}>
            Лучшие тарифы на рынке
          </p>
        </div>
        
        <div style={{
          background: '#ffffff',
          padding: '2rem',
          borderRadius: '8px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
          textAlign: 'center'
        }}>
          <div style={{ fontSize: '3rem', marginBottom: '1rem' }}>⚡</div>
          <h4 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
            Быстрая доставка
          </h4>
          <p style={{ color: '#6c757d' }}>
            Гарантируем сроки доставки
          </p>
        </div>
      </div>
    </div>
  );
}