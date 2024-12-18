import { createSlice } from '@reduxjs/toolkit';
import dummyOrders from '@/data/dummyOrders';
import { toast } from 'sonner';

const ordersSlice = createSlice({
    name: "orders",
    initialState: {
        orders: [...dummyOrders],
        filteredOrders: [...dummyOrders].sort((a, b) => {
            const dateA = a.shipping_date ? new Date(a.shipping_date) : new Date(0);
            const dateB = b.shipping_date ? new Date(b.shipping_date) : new Date(0);
            return dateA - dateB;
        }),
        loading: false,
        error: null,
    },
    reducers: {
        setLoading(state, action) {
            state.loading = action.payload;
        },

        addOrder: (state, action) => {
            state.orders.push(action.payload);
            state.filteredOrders.push(action.payload);
            toast.success("Order added successfully!");
        },

        updateOrder: (state, action) => {
            const updateOrderInArray = (array, payload) => {
                const index = array.findIndex((item) => item.id === payload.id);
                if (index !== -1) {
                    array[index] = { ...array[index], ...payload };
                } else {
                    toast.error(`Order with ID ${payload.id} not found.`);
                }
            };

            updateOrderInArray(state.orders, action.payload);
            updateOrderInArray(state.filteredOrders, action.payload);

            toast.success("Order updated successfully!");
        },

        deleteOrder: (state, action) => {
            const orderIndex = state.orders.findIndex((order) => order.id === action.payload);
            if (orderIndex !== -1) {
                state.orders.splice(orderIndex, 1);
                state.filteredOrders = state.filteredOrders.filter(
                    (order) => order.id !== action.payload
                );
                toast.success("Order deleted successfully!");
            } else {
                toast.error(`Order with ID ${action.payload} not found.`);
            }
        },

        filterOrders: (state, action) => {
            const {
                status,
                customer_name,
                totalPriceMin,
                totalPriceMax,
                shippingDateFrom,
                shippingDateTo,
            } = action.payload;

            let filtered = [...state.orders];

            if (status) {
                filtered = filtered.filter((order) =>
                    order.status.toLowerCase().includes(status.toLowerCase())
                );
            }

            if (customer_name) {
                filtered = filtered.filter((order) =>
                    order.customer_name.toLowerCase().includes(customer_name.toLowerCase())
                );
            }

            if (totalPriceMin !== undefined && totalPriceMin !== null) {
                filtered = filtered.filter((order) => order.total_price >= totalPriceMin);
            }

            if (totalPriceMax !== undefined && totalPriceMax !== null) {
                filtered = filtered.filter((order) => order.total_price <= totalPriceMax);
            }

            if (shippingDateFrom) {
                filtered = filtered.filter(
                    (order) => new Date(order.shipping_date) >= new Date(shippingDateFrom)
                );
            }

            if (shippingDateTo) {
                filtered = filtered.filter(
                    (order) => new Date(order.shipping_date) <= new Date(shippingDateTo)
                );
            }

            state.filteredOrders = filtered;
        },
    },
});

export const { setLoading, addOrder, updateOrder, deleteOrder, filterOrders } = ordersSlice.actions;

export default ordersSlice.reducer;
