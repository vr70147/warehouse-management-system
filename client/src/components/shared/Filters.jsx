import React, { useState } from 'react';
import { Button } from '../ui/button';
import { useDispatch } from 'react-redux';
import { sortItem } from '@/redux/slices/inventorySlice';

export default function Filters({ fields, onFilter }) {
  const [filters, setFilters] = useState({});
  const dispatch = useDispatch();

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFilters((prev) => ({ ...prev, [name]: value }));
  };

  const handleApplyFilters = () => {
    onFilter(filters);
  };
  const [sortOrders, setSortOrders] = useState({
    quantity: 'asc',
    price: 'asc',
  });

  const handleSortToggle = (field) => {
    const currentOIrder = sortOrders[field];
    const newOrder = currentOIrder === 'asc' ? 'desc' : 'asc';
    setSortOrders((prev) => ({ ...prev, [field]: newOrder }));
    dispatch(sortItem({ field, order: newOrder }));
  };

  return (
    <div className="flex flex-col sm:flex-row gap-4 mb-6 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
      {fields.map((field) => {
        if (field.type === 'text') {
          return (
            <input
              key={field.name}
              type="text"
              name={field.name}
              placeholder={field.placeholder}
              onChange={handleInputChange}
              className="rounded bg-gray-200 dark:bg-gray-700 p-2"
            />
          );
        }
        if (field.type === 'select') {
          return (
            <select
              key={field.name}
              name={field.name}
              onChange={handleInputChange}
              className="p-2 border rounded-md bg-gray-200 dark:bg-gray-700 dark:text-white"
            >
              <option value="">{field.placeholder}</option>
              {field.options.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </select>
          );
        }
        if (field.type === 'number') {
          return (
            <input
              key={field.name}
              type="number"
              name={field.name}
              placeholder={field.placeholder}
              onChange={handleInputChange}
              className="rounded bg-gray-200 dark:bg-gray-700 p-2"
            />
          );
        }
        return null;
      })}
      <Button variant="blue" onClick={handleApplyFilters}>
        Apply Filters
      </Button>
      <Button onClick={() => handleSortToggle('quantity')}>
        Sort by Quantity (
        {sortOrders.quantity === 'asc' ? 'Low to High' : 'High to Low'})
      </Button>
      <Button onClick={() => handleSortToggle('unitPrice')}>
        Sort by Price (
        {sortOrders.price === 'asc' ? 'Low to High' : 'High to Low'})
      </Button>
    </div>
  );
}
