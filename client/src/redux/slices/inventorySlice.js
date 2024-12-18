import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import dummyInventory from '@/data/dummyInventory';
import { toast } from 'sonner';

const inventorySlice = createSlice({
    name: 'inventory',
    initialState: {
        items: dummyInventory,
        filteredItems: dummyInventory.sort((a, b) => a.quantity - b.quantity),
        loading: false,
        error: null,
    },
    reducers: {
        setLoading(state, action) {
            state.loading = action.payload;
        },
        addItem(state, action) {
            state.items = [...state.items, action.payload];
            state.filteredItems = [...state.filteredItems, action.payload];
            state.filteredItems.sort((a, b) => a.quantity - b.quantity);
            toast.success("Item added successfully!");
        },
        deleteItem(state, action) {
            state.items = state.items.filter((item) => item.id !== action.payload);
            state.filteredItems = state.filteredItems.filter((item) => item.id !== action.payload);

            toast.error("Item deleted successfully!");
        },

        updateItem(state, action) {
            const updateItemInArray = (array, payload) => {
                const index = array.findIndex((item) => item.id === payload.id);
                if (index !== -1) {
                    array[index] = { ...array[index], ...payload };
                }
            };

            updateItemInArray(state.items, action.payload);
            updateItemInArray(state.filteredItems, action.payload);

            state.filteredItems.sort((a, b) => a.quantity - b.quantity);
            toast.success("Item updated successfully!");
        },

        filterItems(state, action) {
            const { name, category, supplier, priceMin, priceMax } = action.payload;
            let filtered = [...state.items];
            if (name) {
                filtered = filtered.filter((item) => item.name.toLowerCase().includes(name.toLowerCase()));
            }
            if (category) {
                filtered = filtered.filter((item) => item.category.toLowerCase().includes(category.toLowerCase()));
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
            const { field, order } = action.payload;
            state.filteredItems = [...state.filteredItems].sort((a, b) => order === 'asc' ? a[field] - b[field] : b[field] - a[field]);
        },
    },
});

export const { addItem, deleteItem, updateItem, filterItems, searchItems, sortItem, setLoading } = inventorySlice.actions;
export default inventorySlice.reducer;