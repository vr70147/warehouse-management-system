import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import dummyInventory from '@/data/dummyInventory';

const inventorySlice = createSlice({
    name: 'inventory',
    initialState: {
        items: dummyInventory,
        loading: false,
        error: null,
    },
    reducers: {
        addItem(state, action) {
            state.items.push(action.payload);
        },
        deleteItem(state, action) {
            state.items = state.items.filter((item) => item.id !== action.payload);
        },
        updateItem(state, action) {
            const index = state.items.findIndex((item) => item.id === action.payload.id);
            if (index !== -1) {
                state.items[index] = action.payload;
            }
        },
    },
});

export const { addItem, deleteItem, updateItem } = inventorySlice.actions;
export default inventorySlice.reducer;