import { configureStore } from '@reduxjs/toolkit';
import dashboardReducer from './slices/dashboardSlice';
import inventoryReducer from './slices/inventorySlice';

export const store = configureStore({
    reducer: {
        dashboard: dashboardReducer,
        inventory: inventoryReducer,
    },
});