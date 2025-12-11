import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import type { TransportServiceFilters } from '../../types/TransportService';

/**
 * Начальное состояние фильтров
 */
interface FiltersState {
  filters: TransportServiceFilters;
}

const initialState: FiltersState = {
  filters: {},
};

/**
 * Redux slice для управления фильтрами услуг
 * 
 * Позволяет сохранять состояние фильтров при навигации между страницами
 * Состояние сохраняется в Redux store и доступно через Redux DevTools
 */
const filtersSlice = createSlice({
  name: 'filters',
  initialState,
  reducers: {
    /**
     * Устанавливает новые фильтры
     */
    setFilters: (state, action: PayloadAction<TransportServiceFilters>) => {
      state.filters = action.payload;
    },
    /**
     * Очищает все фильтры
     */
    clearFilters: (state) => {
      state.filters = {};
    },
    /**
     * Обновляет конкретное поле фильтра
     */
    updateFilter: (state, action: PayloadAction<Partial<TransportServiceFilters>>) => {
      state.filters = { ...state.filters, ...action.payload };
    },
  },
});

export const { setFilters, clearFilters, updateFilter } = filtersSlice.actions;
export default filtersSlice.reducer;

