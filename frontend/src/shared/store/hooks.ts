import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
import type { RootState, AppDispatch } from './store';

/**
 * Типизированные хуки для работы с Redux
 * 
 * Используются вместо стандартных useDispatch и useSelector
 * для обеспечения типобезопасности
 */
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

