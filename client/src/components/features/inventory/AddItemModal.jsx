import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { addItem } from '@/redux/slices/inventorySlice';
import { Button } from '@/components/ui/button';

export default function AddItemModal({ isOpen, onClose }) {
  const dispatch = useDispatch();
  const [formData, setFormData] = useState({
    name: '',
    category: '',
    quantity: 0,
    unitPrice: 0.0,
    supplier: '',
    lastUpdated: new Date().toISOString(),
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const newItem = {
      id: Date.now(),
      ...formData,
      quantity: parseInt(formData.quantity, 10),
      unitPrice: parseFloat(formData.unitPrice),
    };
    dispatch(addItem(newItem));
    onClose();
    setFormData({
      name: '',
      category: '',
      quantity: 0,
      unitPrice: 0.0,
      supplier: '',
      lastUpdated: new Date().toISOString(),
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50">
      {/* Overlay */}
      <div className="absolute inset-0 bg-black bg-opacity-50"></div>

      {/* Modal Content */}
      <div className="relative bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg w-1/3 z-10">
        <h2 className="text-xl font-bold text-gray-800 dark:text-white mb-4">
          Add New Item
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Name Field */}
          <div>
            <label className="block font-medium text-gray-700 dark:text-gray-300">
              Name
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className="w-full bg-gray-200 dark:bg-gray-700 p-2 text-gray-900 dark:text-white rounded focus:outline-none focus:ring focus:ring-blue-500"
              required
            />
          </div>
          {/* Category Field */}
          <div>
            <label className="block font-medium text-gray-700 dark:text-gray-300">
              Category
            </label>
            <input
              type="text"
              name="category"
              value={formData.category}
              onChange={handleChange}
              className="w-full bg-gray-200 dark:bg-gray-700 p-2 text-gray-900 dark:text-white rounded focus:outline-none focus:ring focus:ring-blue-500"
              required
            />
          </div>
          {/* Supplier Field */}
          <div>
            <label className="block font-medium text-gray-700 dark:text-gray-300">
              Supplier
            </label>
            <input
              type="text"
              name="supplier"
              value={formData.supplier}
              onChange={handleChange}
              className="w-full bg-gray-200 dark:bg-gray-700 p-2 text-gray-900 dark:text-white rounded focus:outline-none focus:ring focus:ring-blue-500"
              required
            />
          </div>
          {/* Quantity Field */}
          <div>
            <label className="block font-medium text-gray-700 dark:text-gray-300">
              Quantity
            </label>
            <input
              type="number"
              name="quantity"
              value={formData.quantity}
              onChange={handleChange}
              className="w-full bg-gray-200 dark:bg-gray-700 p-2 text-gray-900 dark:text-white rounded focus:outline-none focus:ring focus:ring-blue-500"
              required
            />
          </div>
          {/* Unit Price Field */}
          <div>
            <label className="block font-medium text-gray-700 dark:text-gray-300">
              Unit Price
            </label>
            <input
              type="number"
              name="unitPrice"
              value={formData.unitPrice}
              onChange={handleChange}
              className="w-full bg-gray-200 dark:bg-gray-700 p-2 text-gray-900 dark:text-white rounded focus:outline-none focus:ring focus:ring-blue-500"
              required
            />
          </div>
          {/* Buttons */}
          <div className="flex justify-end space-x-2">
            <Button variant="secondary" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" variant="default">
              Add Item
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
