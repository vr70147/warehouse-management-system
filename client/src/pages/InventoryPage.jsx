import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button } from '@/components/ui/button';
import InventoryTable from '@/components/features/inventory/InventoryTable';
import UnifiedItemModal from '@/components/features/inventory/UnifiedItemModal';
import {
  addItem,
  filterItems,
  searchItems,
  sortItem,
} from '@/redux/slices/inventorySlice';
import Filter from '@/components/shared/Filters';

export default function InventoryPage() {
  const dispatch = useDispatch();
  const filteredItems = useSelector((state) => state.inventory.filteredItems);

  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [sortOrder, setSortOrder] = useState('asc');

  const handleAddItem = (newItem) => {
    dispatch(addItem({ ...newItem, id: Date.now() }));
  };

  const handleSortToggle = () => {
    const newOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    setSortOrder(newOrder);
    dispatch(sortItem(newOrder));
  };

  const handleFilter = (filters) => {
    dispatch(filterItems(filters));
  };

  const filterFields = [
    {
      name: 'category',
      type: 'select',
      placeholder: 'Select Category',
      options: ['Electronics', 'Furniture', 'Clothing', 'Accessories'],
    },
    {
      name: 'supplier',
      type: 'text',
      placeholder: 'Search by Supplier',
    },
    {
      name: 'priceMin',
      type: 'number',
      placeholder: 'Min Price',
    },
    {
      name: 'priceMax',
      type: 'number',
      placeholder: 'Max Price',
    },
  ];

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          Inventory Management
        </h1>
        <Button variant="blue" onClick={() => setIsAddModalOpen(true)}>
          Add New Item
        </Button>
      </div>
      <Filter onFilter={handleFilter} fields={filterFields} />
      <InventoryTable items={filteredItems} />
      <UnifiedItemModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        mode="add"
        onSubmit={handleAddItem}
      />
    </div>
  );
}
