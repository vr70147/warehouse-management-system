import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    pendingOrders: 120,
    shippedOrders: 98,
    lowInventory: 15,
    customers: 250,
};

const dashboardSlice = createSlice({
    name: 'dashboard',
    initialState,
    reducers: {
        updatePendingOrders(state, action) {
            state.pendingOrders = action.payload;
        },
        updateShippedOrders(state, action) {
            state.shippedOrders = action.payload;
        },
        updateLoweInventory(state, action) {
            state.lowInventory = action.payload;
        },
        updateCustomers(state, action) {
            state.customers = action.payload;
        },
    },
});

export const {
    updatePendingOrders,
    updateShippedOrders,
    updateLowInventory,
    updateCustomers,
} = dashboardSlice.actions;

export default dashboardSlice.reducer;