import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Navbar } from '../widgets/Navbar/Navbar';
import { Breadcrumbs } from '../widgets/Breadcrumbs/Breadcrumbs';
import { Home } from '../pages/Home/Home';
import { Services } from '../pages/Services/Services';
import { About } from '../pages/About/About';
import '../css/style.css';

/**
 * Главный компонент приложения
 * 
 * Настраивает роутинг, подключает глобальные компоненты (Navbar, Breadcrumbs)
 * Использует BrowserRouter для SPA навигации
 * Подключает существующие стили из style.css
 */
export function App() {
  return (
    <BrowserRouter>
      {/* Навигационная панель - всегда вверху */}
      <Navbar />
      
      {/* Навигационная цепочка - отображается на нужных страницах */}
      <Breadcrumbs />
      
      {/* Маршруты для трех страниц */}
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/services" element={<Services />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </BrowserRouter>
  );
}