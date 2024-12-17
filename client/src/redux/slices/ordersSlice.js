import { createSlice } from '@reduxjs/toolkit';
import dummyOrders from '@/data/dummyOrders';

const ordersSlice = createSlice({
    name: "orders",
    initialState: {
        orders: dummyOrders,
        sortedOrders: dummyOrders.sort((a, b) => a.deliveryDate - b.deliveryDate),
        loading: false,
        error: null
    },
    reducers: {
    }
})

export default ordersSlice.reducer;