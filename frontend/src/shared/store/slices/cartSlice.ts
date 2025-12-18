import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { cartApi } from '../../api/cartApi';
import type { LogisticRequest } from '../../types/Order';

/**
 * Redux slice для управления корзиной (черновиком заявки)
 * Использует redux-thunk для асинхронных операций
 */

interface CartState {
  cart: LogisticRequest | null;
  count: number;
  isLoading: boolean;
  error: string | null;
}

const initialState: CartState = {
  cart: null,
  count: 0,
  isLoading: false,
  error: null,
};

/**
 * Thunk для добавления услуги в корзину
 */
export const addToCart = createAsyncThunk(
  'cart/addToCart',
  async (serviceId: number, { rejectWithValue }) => {
    try {
      const response = await cartApi.addToCart(serviceId);
      return response;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка добавления в корзину');
    }
  }
);

/**
 * Thunk для получения корзины
 */
export const fetchCart = createAsyncThunk(
  'cart/fetchCart',
  async (_, { rejectWithValue }) => {
    try {
      const response = await cartApi.getCart();
      return response;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка загрузки корзины');
    }
  }
);

/**
 * Thunk для получения количества товаров в корзине
 */
export const fetchCartCount = createAsyncThunk(
  'cart/fetchCartCount',
  async (_, { rejectWithValue }) => {
    try {
      const count = await cartApi.getCartCount();
      return count;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка загрузки количества');
    }
  }
);

/**
 * Thunk для очистки корзины
 */
export const clearCart = createAsyncThunk(
  'cart/clearCart',
  async (_, { rejectWithValue }) => {
    try {
      await cartApi.clearCart();
      return null;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка очистки корзины');
    }
  }
);

const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    /**
     * Очистка корзины (локально, без запроса к серверу)
     */
    resetCart: (state) => {
      state.cart = null;
      state.count = 0;
    },
    /**
     * Очистка ошибки
     */
    clearError: (state) => {
      state.error = null;
    },
    /**
     * Установка количества вручную
     */
    setCount: (state, action: PayloadAction<number>) => {
      state.count = action.payload;
    },
  },
  extraReducers: (builder) => {
    // Add to Cart
    builder
      .addCase(addToCart.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(addToCart.fulfilled, (state, action) => {
        state.isLoading = false;
        state.count = action.payload.count;
        state.error = null;
      })
      .addCase(addToCart.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });

    // Fetch Cart
    builder
      .addCase(fetchCart.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchCart.fulfilled, (state, action) => {
        state.isLoading = false;
        state.cart = action.payload.cart;
        state.count = action.payload.count;
        state.error = null;
      })
      .addCase(fetchCart.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });

    // Fetch Cart Count
    builder
      .addCase(fetchCartCount.fulfilled, (state, action) => {
        state.count = action.payload;
      });

    // Clear Cart
    builder
      .addCase(clearCart.fulfilled, (state) => {
        state.cart = null;
        state.count = 0;
      });
  },
});

export const { resetCart, clearError, setCount } = cartSlice.actions;
export default cartSlice.reducer;

