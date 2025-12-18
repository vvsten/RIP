import { configureStore } from '@reduxjs/toolkit';
import filtersReducer from './slices/filtersSlice';
import authReducer from './slices/authSlice';
import ordersReducer from './slices/ordersSlice';
import cartReducer from './slices/cartSlice';

/**
 * Redux store для приложения
 * 
 * Использует Redux Toolkit для управления состоянием
 * Включает redux-thunk middleware для асинхронных операций
 * Включает:
 * - filters: состояние фильтров услуг
 * - auth: состояние авторизации пользователя
 * - orders: состояние заявок
 * - cart: состояние корзины (черновика заявки)
 */
export const store = configureStore({
  reducer: {
    filters: filtersReducer,
    auth: authReducer,
    orders: ordersReducer,
    cart: cartReducer,
  },
  // Включаем Redux DevTools для отладки
  devTools: import.meta.env.DEV,
  // redux-thunk уже включен по умолчанию в Redux Toolkit
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

