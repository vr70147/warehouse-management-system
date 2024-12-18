import { configureStore } from '@reduxjs/toolkit';
import dashboardReducer from './slices/dashboardSlice';
import inventoryReducer from './slices/inventorySlice';
import ordersReducer from './slices/ordersSlice';

export const store = configureStore({
    reducer: {
        dashboard: dashboardReducer,
        inventory: inventoryReducer,
        orders: ordersReducer,
    },
});