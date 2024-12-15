import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import InventorySearch from '@/components/features/inventory/InventorySearch';
import InventoryTable from '@/components/features/inventory/InventoryTable';
import AddItemButton from '@/components/features/inventory/AddItemButton';
import {
  filteredItems,
  searchItems,
  sortItem,
} from '@/redux/slices/inventorySlice';
import InventoryFilters from '@/components/features/inventory/InventoryFilters';

export default function InventoryPage() {
  const { items } = useSelector((state) => state.inventory);
  const dispatch = useDispatch();
  const filteredItems = useSelector((state) => state.inventory.filteredItems);

  const handleFilter = (filter) => {
    dispatch(filterItems(filter));
  };

  const handleSearch = (searchTerm) => {
    dispatch(searchItems(searchTerm));
  };

  const handleSort = (order) => {
    dispatch(sortItem(order));
  };

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          Inventory Management
        </h1>
        <AddItemButton />
      </div>
      <div className="flex justify-end">
        <InventoryFilters
          onFilter={handleFilter}
          onSearch={handleSearch}
          onSort={handleSort}
        />
      </div>
      <InventoryTable items={filteredItems} />
    </div>
  );
}
