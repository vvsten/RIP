import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { ordersApi, GetOrdersParams } from '../../api/ordersApi';
import type { LogisticRequest } from '../../types/Order';

/**
 * Redux slice для управления заявками
 * Использует redux-thunk для асинхронных операций
 */

interface OrdersState {
  orders: LogisticRequest[];
  currentOrder: LogisticRequest | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: OrdersState = {
  orders: [],
  currentOrder: null,
  isLoading: false,
  error: null,
};

/**
 * Thunk для получения списка заявок
 */
export const fetchOrders = createAsyncThunk(
  'orders/fetchOrders',
  async (params?: GetOrdersParams, { rejectWithValue }) => {
    try {
      const orders = await ordersApi.getOrders(params);
      return orders;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка загрузки заявок');
    }
  }
);

/**
 * Thunk для получения заявки по ID
 */
export const fetchOrder = createAsyncThunk(
  'orders/fetchOrder',
  async (id: number, { rejectWithValue }) => {
    try {
      const order = await ordersApi.getOrder(id);
      return order;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка загрузки заявки');
    }
  }
);

/**
 * Thunk для обновления заявки
 */
export const updateOrder = createAsyncThunk(
  'orders/updateOrder',
  async ({ id, data }: { id: number; data: any }, { rejectWithValue }) => {
    try {
      const order = await ordersApi.updateOrder(id, data);
      return order;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка обновления заявки');
    }
  }
);

/**
 * Thunk для удаления услуги из заявки
 */
export const removeServiceFromOrder = createAsyncThunk(
  'orders/removeService',
  async ({ orderId, serviceId }: { orderId: number; serviceId: number }, { rejectWithValue }) => {
    try {
      await ordersApi.removeServiceFromOrder(orderId, serviceId);
      // Перезагружаем заявку после удаления
      const order = await ordersApi.getOrder(orderId);
      return order;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка удаления услуги');
    }
  }
);

/**
 * Thunk для обновления количества услуги в заявке
 */
export const updateServiceInOrder = createAsyncThunk(
  'orders/updateService',
  async ({ orderId, serviceId, quantity }: { orderId: number; serviceId: number; quantity: number }, { rejectWithValue }) => {
    try {
      await ordersApi.updateServiceInOrder(orderId, serviceId, quantity);
      // Перезагружаем заявку после обновления
      const order = await ordersApi.getOrder(orderId);
      return order;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка обновления услуги');
    }
  }
);

/**
 * Thunk для подтверждения заявки
 */
export const formOrder = createAsyncThunk(
  'orders/formOrder',
  async ({ id, data }: { id: number; data: any }, { rejectWithValue }) => {
    try {
      const order = await ordersApi.formOrder(id, data);
      return order;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Ошибка подтверждения заявки');
    }
  }
);

const ordersSlice = createSlice({
  name: 'orders',
  initialState,
  reducers: {
    /**
     * Очистка текущей заявки
     */
    clearCurrentOrder: (state) => {
      state.currentOrder = null;
    },
    /**
     * Очистка ошибки
     */
    clearError: (state) => {
      state.error = null;
    },
    /**
     * Очистка всех заявок
     */
    clearOrders: (state) => {
      state.orders = [];
      state.currentOrder = null;
    },
  },
  extraReducers: (builder) => {
    // Fetch Orders
    builder
      .addCase(fetchOrders.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchOrders.fulfilled, (state, action) => {
        state.isLoading = false;
        state.orders = action.payload;
        state.error = null;
      })
      .addCase(fetchOrders.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });

    // Fetch Order
    builder
      .addCase(fetchOrder.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchOrder.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentOrder = action.payload;
        state.error = null;
      })
      .addCase(fetchOrder.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });

    // Update Order
    builder
      .addCase(updateOrder.fulfilled, (state, action) => {
        state.currentOrder = action.payload;
        // Обновляем заявку в списке
        const index = state.orders.findIndex(o => o.id === action.payload.id);
        if (index !== -1) {
          state.orders[index] = action.payload;
        }
      });

    // Remove Service
    builder
      .addCase(removeServiceFromOrder.fulfilled, (state, action) => {
        state.currentOrder = action.payload;
      });

    // Update Service
    builder
      .addCase(updateServiceInOrder.fulfilled, (state, action) => {
        state.currentOrder = action.payload;
      });

    // Form Order
    builder
      .addCase(formOrder.fulfilled, (state, action) => {
        state.currentOrder = action.payload;
        // Обновляем заявку в списке
        const index = state.orders.findIndex(o => o.id === action.payload.id);
        if (index !== -1) {
          state.orders[index] = action.payload;
        }
      });
  },
});

export const { clearCurrentOrder, clearError, clearOrders } = ordersSlice.actions;
export default ordersSlice.reducer;

