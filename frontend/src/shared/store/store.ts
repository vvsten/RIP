import { configureStore } from '@reduxjs/toolkit';
import filtersReducer from './slices/filtersSlice';

/**
 * Redux store для приложения
 * 
 * Использует Redux Toolkit для управления состоянием
 * Включает:
 * - filters: состояние фильтров услуг
 */
export const store = configureStore({
  reducer: {
    filters: filtersReducer,
  },
  // Включаем Redux DevTools для отладки
  devTools: import.meta.env.DEV,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

