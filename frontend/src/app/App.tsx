import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from '../shared/store/store';
import { Navbar } from '../widgets/Navbar/Navbar';
import { Breadcrumbs } from '../widgets/Breadcrumbs/Breadcrumbs';
import { ServerConfig } from '../widgets/ServerConfig/ServerConfig';
import { AuthInitializer } from '../shared/components/AuthInitializer/AuthInitializer';
import { Home } from '../pages/Home/Home';
import { Services } from '../pages/Services/Services';
import { About } from '../pages/About/About';
import { Login } from '../pages/Login/Login';
import { Register } from '../pages/Register/Register';
import { OrdersList } from '../pages/OrdersList/OrdersList';
import { OrderDetails } from '../pages/OrderDetails/OrderDetails';
import { Profile } from '../pages/Profile/Profile';
import { ProtectedRoute } from '../shared/components/ProtectedRoute/ProtectedRoute';
import '../css/style.css';

/**
 * Главный компонент приложения
 * 
 * Настраивает роутинг, подключает глобальные компоненты (Navbar, Breadcrumbs)
 * Использует BrowserRouter для SPA навигации
 * Подключает Redux Provider для управления состоянием
 * Подключает существующие стили из style.css
 * 
 * Использует ProtectedRoute для защиты страниц, требующих авторизации
 */
export function App() {
  return (
    <Provider store={store}>
      <BrowserRouter>
        {/* Компонент настройки сервера для Tauri (отображается только в Tauri) */}
        <ServerConfig />
        
        {/* Инициализация состояния авторизации из localStorage */}
        <AuthInitializer />
        
        {/* Навигационная панель - всегда вверху */}
        <Navbar />
        
        {/* Навигационная цепочка - отображается на нужных страницах */}
        <Breadcrumbs />
        
        {/* Маршруты для страниц */}
        <Routes>
          {/* Публичные маршруты */}
          <Route path="/" element={<Home />} />
          <Route path="/transport-services" element={<Services />} />
          <Route path="/about" element={<About />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          
          {/* Защищенные маршруты (требуют авторизации) */}
          <Route
            path="/orders"
            element={
              <ProtectedRoute>
                <OrdersList />
              </ProtectedRoute>
            }
          />
          <Route
            path="/orders/:id"
            element={
              <ProtectedRoute>
                <OrderDetails />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />
          
          {/* Fallback для несуществующих маршрутов */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </Provider>
  );
}