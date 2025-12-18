import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../../shared/store/hooks';
import { fetchCartCount, fetchCart } from '../../shared/store/slices/cartSlice';

/**
 * Иконка корзины для страницы списка услуг
 * Использует Redux для управления состоянием корзины
 * Отображает кнопку перехода на страницу заявки (черновика)
 * Кнопка меняет состояние: если черновик есть - доступна, если нет - недоступна
 */
export function CartIcon() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { count, cart, isLoading } = useAppSelector((state) => state.cart);
  const { isAuthenticated } = useAppSelector((state) => state.auth);

  useEffect(() => {
    if (isAuthenticated) {
      // Загружаем количество товаров в корзине
      dispatch(fetchCartCount());
      // Загружаем корзину для получения ID черновика
      dispatch(fetchCart());
    }
  }, [dispatch, isAuthenticated]);

  const handleClick = () => {
    if (cart && cart.id && count > 0) {
      navigate(`/orders/${cart.id}`);
    } else if (count > 0) {
      // Если есть товары, но нет ID черновика, загружаем корзину
      dispatch(fetchCart()).then((result) => {
        if (fetchCart.fulfilled.match(result) && result.payload.cart?.id) {
          navigate(`/orders/${result.payload.cart.id}`);
        }
      });
    }
  };

  const isDisabled = count <= 0 || isLoading;

  return (
    <div className="calculator-shortcut">
      {isDisabled ? (
        <a 
          className="calculator-btn is-disabled" 
          aria-disabled="true"
          style={{ cursor: 'not-allowed', opacity: 0.6 }}
        >
          🧮 Заявка
        </a>
      ) : (
        <a 
          onClick={handleClick}
          className="calculator-btn" 
          style={{ textDecoration: 'none', cursor: 'pointer' }}
        >
          🧮 Заявка
          {count > 0 && <span className="cart-count" id="cartCount">{count}</span>}
        </a>
      )}
    </div>
  );
}

