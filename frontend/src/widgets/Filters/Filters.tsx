import { useState } from 'react';
import type { ServiceFilters } from '../../shared/types/Service';

/**
 * Props для компонента Filters
 */
interface FiltersProps {
  /** 
   * Callback функция, вызываемая при изменении фильтров
   * Принимает объект с параметрами фильтрации
   */
  onFilterChange: (filters: ServiceFilters) => void;
}

/**
 * Компонент фильтрации услуг
 * 
 * Использует useState для управления состоянием полей формы
 * При изменении любого поля вызывает onFilterChange с обновленными фильтрами
 * 
 * @param props - содержит callback onFilterChange
 */
export function Filters({ onFilterChange }: FiltersProps) {
  // Только строка поиска — как в шаблонах бэкенда
  const [search, setSearch] = useState('');
  
  /**
   * Формирует объект фильтров из текущих полей
   */
  const buildFilters = (): ServiceFilters => {
    const filters: ServiceFilters = {};
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
  
  /**
   * Обработчик очистки фильтров
   * Сбрасывает все поля и вызывает onFilterChange с пустым объектом
   */
  // Очистка не требуется — поведение как в шаблонах
  
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