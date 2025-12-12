import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from '../shared/store/store';
import { Navbar } from '../widgets/Navbar/Navbar';
import { Breadcrumbs } from '../widgets/Breadcrumbs/Breadcrumbs';
import { ServerConfig } from '../widgets/ServerConfig/ServerConfig';
import { Home } from '../pages/Home/Home';
import { Services } from '../pages/Services/Services';
import { About } from '../pages/About/About';
import '../css/style.css';

/**
 * Главный компонент приложения
 * 
 * Настраивает роутинг, подключает глобальные компоненты (Navbar, Breadcrumbs)
 * Использует BrowserRouter для SPA навигации
 * Подключает Redux Provider для управления состоянием
 * Подключает существующие стили из style.css
 */
export function App() {
  return (
    <Provider store={store}>
    <BrowserRouter>
        {/* Компонент настройки сервера для Tauri (отображается только в Tauri) */}
        <ServerConfig />
        
      {/* Навигационная панель - всегда вверху */}
      <Navbar />
      
      {/* Навигационная цепочка - отображается на нужных страницах */}
      <Breadcrumbs />
      
      {/* Маршруты для страниц */}
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/transport-services" element={<Services />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </BrowserRouter>
    </Provider>
  );
}