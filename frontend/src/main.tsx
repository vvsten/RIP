import React from 'react';
import ReactDOM from 'react-dom/client';
import { App } from './app/App';

/**
 * Точка входа в приложение
 * 
 * Монтирует корневой компонент App в DOM элемент #root
 */
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);


