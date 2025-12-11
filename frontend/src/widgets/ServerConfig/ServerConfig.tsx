import { useState, useEffect } from 'react';
import { getServerIPAddress, setServerIP } from '../../shared/config/apiConfig';

/**
 * Компонент для настройки IP адреса сервера в Tauri приложении
 * 
 * Отображается только в Tauri приложении
 * Позволяет пользователю указать IP адрес сервера в локальной сети
 */
export function ServerConfig() {
  const [ip, setIp] = useState('');
  const [isVisible, setIsVisible] = useState(false);
  
  useEffect(() => {
    // Проверяем Tauri несколькими способами
    const checkTauri = () => {
      if (typeof window === 'undefined') return false;
      
      // Способ 1: проверка через __TAURI__
      if ('__TAURI__' in window) return true;
      if ('__TAURI_INTERNALS__' in window) return true;
      if ((window as any).__TAURI_METADATA__) return true;
      
      // Способ 2: проверка через userAgent
      if (navigator.userAgent.includes('Tauri')) return true;
      
      // Способ 3: проверка через window.location (в Tauri обычно file:// или tauri://)
      if (window.location.protocol === 'tauri:') return true;
      
      // Способ 4: проверка через параметр URL (для тестирования)
      if (window.location.search.includes('tauri=true')) return true;
      
      return false;
    };
    
    const isTauri = checkTauri();
    setIsVisible(isTauri);
    
    if (isTauri) {
      const currentIP = getServerIPAddress();
      // Извлекаем IP из URL (убираем http:// и порт)
      if (currentIP) {
        // Поддерживаем и http:// и https://
        const ipMatch = currentIP.match(/https?:\/\/([^:]+)/);
        if (ipMatch) {
          setIp(ipMatch[1]);
        } else {
          // Если IP не настроен, пытаемся определить автоматически
          // Для локальной сети используем 192.168.1.64 (из скриншота видно)
          setIp('192.168.1.64');
        }
      } else {
        // По умолчанию используем IP из локальной сети
        setIp('192.168.1.64');
      }
    }
  }, []);
  
  // Показываем компонент если:
  // 1. Определили Tauri
  // 2. Или есть параметр в URL
  // 3. Или не в обычном браузере (file://, tauri:// протоколы)
  const isNotBrowser = window.location.protocol !== 'http:' && window.location.protocol !== 'https:';
  const shouldShow = isVisible || 
                     window.location.search.includes('show-config=true') || 
                     window.location.search.includes('tauri=true') ||
                     isNotBrowser;
  
  if (!shouldShow) {
    return null;
  }
  
  const handleSave = () => {
    const port = '8083';
    // Временно используем HTTP для разработки (HTTPS требует доверенный сертификат)
    const fullUrl = `http://${ip}:${port}`;
    setServerIP(fullUrl);
  };
  
  return (
    <div style={{
      position: 'fixed',
      top: '10px',
      right: '10px',
      background: '#fff',
      padding: '1rem',
      borderRadius: '8px',
      boxShadow: '0 2px 8px rgba(0,0,0,0.2)',
      zIndex: 1000,
      border: '2px solid #0d6efd',
      maxWidth: '300px'
    }}>
      <h4 style={{ marginBottom: '0.5rem', fontSize: '0.9rem', fontWeight: 'bold' }}>
        Настройка сервера (Tauri)
      </h4>
      <p style={{ fontSize: '0.8rem', color: '#6c757d', marginBottom: '0.5rem' }}>
        Текущий IP: <strong>{getServerIPAddress()}</strong>
      </p>
      <div style={{ display: 'flex', gap: '0.5rem', flexDirection: 'column' }}>
        <input
          type="text"
          placeholder="192.168.1.100"
          value={ip}
          onChange={(e) => setIp(e.target.value)}
          style={{
            padding: '0.5rem',
            border: '1px solid #dee2e6',
            borderRadius: '4px',
            fontSize: '0.9rem'
          }}
        />
        <button
          onClick={handleSave}
          style={{
            padding: '0.5rem',
            background: '#0d6efd',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '0.9rem',
            fontWeight: '500'
          }}
        >
          Сохранить IP
        </button>
      </div>
    </div>
  );
}

