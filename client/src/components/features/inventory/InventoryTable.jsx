import { Button } from '@/components/ui/button';
import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { deleteItem, updateItem } from '@/redux/slices/inventorySlice';
import UnifiedItemModal from '@/components/features/inventory/UnifiedItemModal';

export default function InventoryTable({ items }) {
  const dispatch = useDispatch();
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [currentItem, setCurrentItem] = useState(null);

  const handleEdit = (item) => {
    setCurrentItem(item);
    setIsEditModalOpen(true);
  };

  const handleUpdateItem = (updatedItem) => {
    dispatch(updateItem(updatedItem));
  };

  const handleDelete = (id) => {
    dispatch(deleteItem(id));
  };
  return (
    <div className="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4">
      <table className="w-full text-left text-sm text-gray-600 dark:text-gray-400">
        <thead>
          <tr>
            <th className="py-2 px-4">Name</th>
            <th className="py-2 px-4">Category</th>
            <th className="py-2 px-4">Quantity</th>
            <th className="py-2 px-4">Unit Price</th>
            <th className="py-2 px-4">Actions</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr
              key={item.id}
              className="even:bg-gray-50 odd:bg-white dark:even:bg-gray-700 dark:odd:bg-gray-800 hover:scale-102 hover:shadow transition-transform duration-300 ease-out"
            >
              <td className="py-3 px-4 border-b dark:border-gray-700">
                {item.name}
              </td>
              <td className="py-3 px-4 border-b dark:border-gray-700">
                {item.category}
              </td>
              <td className="py-3 px-4 border-b dark:border-gray-700">
                {item.quantity}
              </td>
              <td className="py-3 px-4 border-b dark:border-gray-700">
                $
                {isNaN(item.unitPrice)
                  ? '0.00'
                  : parseFloat(item.unitPrice).toFixed(2)}
              </td>
              <td className="py-3 px-4 border-b dark:border-gray-700 flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  className="hover:scale-105 transition-transform duration-200"
                  onClick={() => handleEdit(item)}
                >
                  Edit
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  className="hover:scale-105 transition-transform duration-200"
                  onClick={() => handleDelete(item.id)}
                >
                  Delete
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <UnifiedItemModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        mode="edit"
        item={currentItem}
        onSubmit={handleUpdateItem}
      />
    </div>
  );
}
