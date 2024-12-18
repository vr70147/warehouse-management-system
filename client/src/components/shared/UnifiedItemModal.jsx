import React, { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';

export default function UnifiedItemModal({
  isOpen,
  onClose,
  mode,
  item,
  onSubmit,
}) {
  const [formData, setFormData] = useState({
    name: '',
    category: '',
    quantity: 0,
    unitPrice: 0.0,
    supplier: '',
    lastUpdated: new Date().toISOString().split('T')[0],
  });

  useEffect(() => {
    if (mode === 'edit' && item) {
      setFormData({
        name: item.name || '',
        category: item.category || '',
        supplier: item.supplier || '',
        quantity: item.quantity || 0,
        unitPrice: item.unitPrice || 0.0,
        lastUpdated: item.lastUpdated || new Date().toISOString().split('T')[0],
      });
    } else {
      setFormData({
        name: '',
        category: '',
        supplier: '',
        quantity: 0,
        unitPrice: 0.0,
        lastUpdated: new Date().toISOString().split('T')[0],
      });
    }
  }, [mode, item]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const payload = {
      id: item?.id,
      ...formData,
    };
    onSubmit(payload);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50">
      <div
        className="absolute inset-0 bg-black bg-opacity-50"
        onClick={onClose}
      ></div>
      <div className="relative bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg w-1/3 z-10">
        <h2 className="text-xl font-bold text-gray-800 dark:text-white mb-4">
          {mode === 'add' ? 'Add New Item' : 'Edit Item'}
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
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
          <div className="flex justify-end space-x-2">
            <Button variant="secondary" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" variant="blue">
              {mode === 'add' ? 'Add Item' : 'Save Changes'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
