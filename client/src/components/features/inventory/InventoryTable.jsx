import { Button } from '@/components/ui/button';
import React from 'react';
import { useDispatch } from 'react-redux';
import { deleteItem } from '@/redux/slices/inventorySlice';

export default function InventoryTable({ items }) {
  const dispatch = useDispatch();

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
              className="even:bg-gray-50 odd:bg-white dark:even:bg-gray-700 dark:odd:bg-gray-800 hover:scale-100 hover:shadow-sm transition-transform duration-400"
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
                {item.unitPrice.toFixed(2)}
              </td>
              <td className="py-3 px-4 border-b dark:border-gray-700 flex gap-2">
                <Button variant="outline" size="sm">
                  Edit
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  onClick={() => handleDelete(item.id)}
                >
                  Delete
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
