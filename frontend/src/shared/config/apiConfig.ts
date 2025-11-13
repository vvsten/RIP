/**
 * Конфигурация API для приложения
 * 
 * Поддерживает два режима:
 * - Web: использует относительные пути (/api) для работы через прокси
 * - Tauri: использует IP адрес локальной сети для прямого подключения
 */

// Определяем, запущено ли приложение в Tauri
// В Tauri v2 проверяем через window.__TAURI__ или userAgent
const isTauri = typeof window !== 'undefined' && (
  '__TAURI__' in window || 
  '__TAURI_INTERNALS__' in window ||
  (window as any).__TAURI_METADATA__ !== undefined ||
  navigator.userAgent.includes('Tauri')
);

// IP адрес сервера для Tauri (можно изменить через переменную окружения или конфиг)
// По умолчанию используется localhost, но в Tauri нужно указать IP локальной сети
const getServerIP = (): string => {
  // В Tauri приложении можно использовать переменную окружения или конфиг
  if (isTauri) {
    // Для Tauri используем IP из localStorage или переменной окружения
    const savedIP = localStorage.getItem('api_server_ip');
    if (savedIP) {
      return savedIP;
    }
    // По умолчанию для Tauri используем localhost (пользователь должен настроить)
    // Используем HTTPS если включен
    return 'https://localhost:8083';
  }
  // Для веб-версии используем относительные пути
  return '';
};

/**
 * Базовый URL для API запросов
 */
export const API_BASE_URL = getServerIP();

/**
 * Полный URL для API запроса
 * @param path - путь API (например, '/api/services')
 */
export const getApiUrl = (path: string): string => {
  if (isTauri && API_BASE_URL) {
    // Для Tauri используем полный URL с IP адресом
    return `${API_BASE_URL}${path}`;
  }
  // Для веб-версии используем относительный путь (работает через прокси)
  return path;
};

/**
 * Установка IP адреса сервера для Tauri
 * @param ip - IP адрес сервера (например, 'http://192.168.1.100:8083')
 */
export const setServerIP = (ip: string): void => {
  if (isTauri) {
    localStorage.setItem('api_server_ip', ip);
    // Перезагружаем страницу для применения изменений
    window.location.reload();
  }
};

/**
 * Получение текущего IP адреса сервера
 */
export const getServerIPAddress = (): string => {
  return API_BASE_URL;
};

/**
 * Проверка, запущено ли приложение в Tauri
 */
export const isTauriApp = (): boolean => {
  return isTauri;
};

