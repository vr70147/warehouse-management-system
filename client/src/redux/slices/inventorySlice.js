import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import dummyInventory from '@/data/dummyInventory';

const inventorySlice = createSlice({
    name: 'inventory',
    initialState: {
        items: dummyInventory,
        filteredItems: dummyInventory,
        loading: false,
        error: null,
    },
    reducers: {
        addItem(state, action) {
            state.items.push(action.payload);
            state.filteredItems.push(action.payload)
        },
        deleteItem(state, action) {
            state.items = state.items.filter((item) => item.id !== action.payload);
        },
        updateItem(state, action) {
            const index = state.items.findIndex((item) => item.id === action.payload.id);
            if (index !== -1) {
                state.items[index] = action.payload;
            }

            const filteredIndex = state.filteredItems.findIndex((item) => item.id === action.payload.id);
            if (filteredIndex !== -1) {
                state.filteredItems[filteredIndex] = action.payload;
            }
        },

        filterItems(state, action) {
            const { category, supplier, priceMin, priceMax } = action.payload;
            let filtered = [...state.items];

            if (category) {
                filtered = filtered.filter((item) => item.category === category);
            }
            if (supplier) {
                filtered = filtered.filter((item) => item.supplier.toLowerCase().includes(supplier.toLowerCase()));
            }
            if (priceMin) {
                filtered = filtered.filter((item) => item.unitPrice >= priceMin);
            }
            if (priceMax) {
                filtered = filtered.filter((item) => item.unitPrice <= priceMax);
            }

            state.filteredItems = filtered;
        },

        searchItems(state, action) {
            const searchTerm = action.payload.toLowerCase();
            state.filteredItems = state.items.filter((item) => item.name.toLocaleLowerCase().includes(searchTerm));
        },
        sortItem(state, action) {
            const order = action.payload;
            state.filteredItems = [...state.filteredItems].sort((a, b) => order === 'asc' ? a.quantity - b.quantity : b.quantity - a.quantity);
        },
    },
});

export const { addItem, deleteItem, updateItem, filterItems, searchItems, sortItem } = inventorySlice.actions;
export default inventorySlice.reducer;