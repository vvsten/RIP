import { useState, useEffect } from 'react';
import { useAppSelector } from '../../shared/store/hooks';
import type { TransportServiceFilters } from '../../shared/types/TransportService';

/**
 * Props для компонента Filters
 */
interface FiltersProps {
  /** 
   * Callback функция, вызываемая при изменении фильтров
   * Принимает объект с параметрами фильтрации
   */
  onFilterChange: (filters: TransportServiceFilters) => void;
}

/**
 * Компонент фильтрации услуг
 * 
 * Использует Redux для получения сохраненных фильтров
 * Использует useState для управления состоянием полей формы
 * При изменении любого поля вызывает onFilterChange с обновленными фильтрами
 * 
 * @param props - содержит callback onFilterChange
 */
export function Filters({ onFilterChange }: FiltersProps) {
  // Получаем сохраненные фильтры из Redux
  const savedFilters = useAppSelector((state) => state.filters.filters);
  
  // Инициализируем поле поиска из Redux store
  const [search, setSearch] = useState(savedFilters.search || '');
  
  // Обновляем поле поиска при изменении фильтров в Redux (например, при навигации)
  useEffect(() => {
    setSearch(savedFilters.search || '');
  }, [savedFilters.search]);
  
  /**
   * Формирует объект фильтров из текущих полей
   */
  const buildFilters = (): TransportServiceFilters => {
    const filters: TransportServiceFilters = {};
    if (search) filters.search = search;
    return filters;
  };

  /**
   * Отправка формы поиска — как в шаблонах на бэкенде
   * Поиск инициируется по кнопке, а не при каждом вводе
   */
  const handleSubmit: React.FormEventHandler<HTMLFormElement> = (e) => {
    e.preventDefault();
    onFilterChange(buildFilters());
  };
  
  return (
    <div className="search-section">
      <form className="search-form" onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Поиск типа транспорта (фура, авиа, поезд...)"
          className="search-input"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <button type="submit" className="search-btn">🔍</button>
      </form>
    </div>
  );
}