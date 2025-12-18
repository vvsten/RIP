import { useEffect } from 'react';
import { useAppDispatch } from '../../store/hooks';
import { restoreAuth } from '../../store/slices/authSlice';

/**
 * Компонент для инициализации состояния авторизации при загрузке приложения
 * Восстанавливает состояние из localStorage
 */
export function AuthInitializer() {
  const dispatch = useAppDispatch();

  useEffect(() => {
    // Восстанавливаем состояние авторизации из localStorage
    dispatch(restoreAuth());
  }, [dispatch]);

  return null;
}

