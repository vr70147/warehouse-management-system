import React from 'react';
import { useSelector } from 'react-redux';
import InventorySearch from '@/components/features/inventory/InventorySearch';
import InventoryTable from '@/components/features/inventory/InventoryTable';
import AddItemButton from '@/components/features/inventory/AddItemButton';

export default function InventoryPage() {
  const { items } = useSelector((state) => state.inventory);

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          Inventory Management
        </h1>
        <AddItemButton />
      </div>
      <InventorySearch />
      <InventoryTable items={items} />
    </div>
  );
}
