import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button } from '@/components/ui/button';
import InventoryTable from '@/components/features/inventory/InventoryTable';
import UnifiedItemModal from '@/components/features/inventory/UnifiedItemModal';
import { addItem } from '@/redux/slices/inventorySlice';

export default function InventoryPage() {
  const dispatch = useDispatch();
  const filteredItems = useSelector((state) => state.inventory.filteredItems);

  const [isAddModalOpen, setIsAddModalOpen] = useState(false);

  const handleAddItem = (newItem) => {
    dispatch(addItem({ ...newItem, id: Date.now() })); // מוסיף פריט חדש עם ID ייחודי
  };

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
