/**
 * Страница "О компании"
 * 
 * Информационная страница с описанием компании, истории и контактов
 */
export function About() {
  return (
    <div className="container">
      <h1 style={{ marginBottom: '2rem', fontSize: '2.5rem', fontWeight: 'bold' }}>
        О компании GruzDelivery
      </h1>
      
      <div style={{ 
        background: '#ffffff',
        padding: '2rem',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        marginBottom: '2rem'
      }}>
        <h2 style={{ 
          fontSize: '1.8rem', 
          fontWeight: 'bold', 
          marginBottom: '1rem',
          color: '#212529'
        }}>
          Наша история
        </h2>
        <p style={{ 
          fontSize: '1.1rem', 
          color: '#6c757d', 
          lineHeight: '1.8',
          marginBottom: '1rem'
        }}>
          Компания GruzDelivery была основана в 2014 году с целью предоставления 
          качественных услуг в сфере грузоперевозок. За годы работы мы зарекомендовали 
          себя как надежный партнер для тысяч клиентов по всей России.
        </p>
        <p style={{ 
          fontSize: '1.1rem', 
          color: '#6c757d', 
          lineHeight: '1.8'
        }}>
          Мы начинали с малого — предоставляли услуги малотоннажных перевозок 
          в пределах одного региона. Сегодня мы работаем по всей стране и предлагаем 
          широкий спектр транспортных услуг.
        </p>
      </div>
      
      <div style={{ 
        background: '#ffffff',
        padding: '2rem',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        marginBottom: '2rem'
      }}>
        <h2 style={{ 
          fontSize: '1.8rem', 
          fontWeight: 'bold', 
          marginBottom: '1rem',
          color: '#212529'
        }}>
          Наши ценности
        </h2>
        <ul style={{ fontSize: '1.1rem', color: '#6c757d', lineHeight: '2' }}>
          <li><strong style={{ color: '#212529' }}>Надежность</strong> — мы всегда выполняем свои обязательства</li>
          <li><strong style={{ color: '#212529' }}>Качество</strong> — высокий уровень обслуживания клиентов</li>
          <li><strong style={{ color: '#212529' }}>Инновации</strong> — использование современных технологий</li>
          <li><strong style={{ color: '#212529' }}>Ответственность</strong> — забота о грузах наших клиентов</li>
        </ul>
      </div>
      
      <div style={{ 
        background: '#ffffff',
        padding: '2rem',
        borderRadius: '8px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        marginBottom: '2rem'
      }}>
        <h2 style={{ 
          fontSize: '1.8rem', 
          fontWeight: 'bold', 
          marginBottom: '1.5rem',
          color: '#212529'
        }}>
          Наши преимущества
        </h2>
        
        <div className="services-grid" style={{ 
          gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))'
        }}>
          <div style={{
            background: '#f8f9fa',
            padding: '1.5rem',
            borderRadius: '8px',
            textAlign: 'center'
          }}>
            <h5 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
              Опыт
            </h5>
            <p style={{ color: '#6c757d' }}>
              Более 10 лет на рынке грузоперевозок
            </p>
          </div>
          
          <div style={{
            background: '#f8f9fa',
            padding: '1.5rem',
            borderRadius: '8px',
            textAlign: 'center'
          }}>
            <h5 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
              Разнообразие
            </h5>
            <p style={{ color: '#6c757d' }}>
              Широкий выбор видов транспорта
            </p>
          </div>
          
          <div style={{
            background: '#f8f9fa',
            padding: '1.5rem',
            borderRadius: '8px',
            textAlign: 'center'
          }}>
            <h5 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
              География
            </h5>
            <p style={{ color: '#6c757d' }}>
              Доставка по всей России
            </p>
          </div>
          
          <div style={{
            background: '#f8f9fa',
            padding: '1.5rem',
            borderRadius: '8px',
            textAlign: 'center'
          }}>
            <h5 style={{ fontSize: '1.3rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
              Цены
            </h5>
            <p style={{ color: '#6c757d' }}>
              Конкурентные тарифы на рынке
            </p>
          </div>
        </div>
      </div>
      
      <div style={{ 
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        padding: '2rem',
        borderRadius: '15px',
        color: 'white'
      }}>
        <h2 style={{ 
          fontSize: '1.8rem', 
          fontWeight: 'bold', 
          marginBottom: '1.5rem'
        }}>
          Контакты
        </h2>
        <p style={{ fontSize: '1.1rem', lineHeight: '2' }}>
          <strong>Адрес:</strong> г. Москва, ул. Транспортная, д. 10<br />
          <strong>Телефон:</strong> +7 (495) 123-45-67<br />
          <strong>Email:</strong> info@gruzdelivery.ru<br />
          <strong>Время работы:</strong> Пн-Пт: 9:00 - 18:00
        </p>
      </div>
    </div>
  );
}