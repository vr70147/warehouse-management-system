import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button } from '@/components/ui/button';
import InventoryTable from '@/components/features/inventory/InventoryTable';
import InventoryFilters from '@/components/features/inventory/InventoryFilters';
import AddItemModal from '@/components/features/inventory/AddItemModal';
import {
  filterItems,
  searchItems,
  sortItem,
} from '@/redux/slices/inventorySlice';

export default function InventoryPage() {
  const dispatch = useDispatch();
  const allItems = useSelector((state) => state.inventory.items);
  const visibleItems = useSelector((state) => state.inventory.filteredItems);
  const [isModalOpen, setIsModalOpen] = useState(false);

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
        <Button variant="outline" onClick={() => setIsModalOpen(true)}>
          Add New Item
        </Button>
      </div>
      <div className="flex justify-end">
        <InventoryFilters
          onFilter={handleFilter}
          onSearch={handleSearch}
          onSort={handleSort}
        />
      </div>
      <InventoryTable items={visibleItems} />
      <AddItemModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </div>
  );
}
