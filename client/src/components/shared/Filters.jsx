import React, { useState } from 'react';

export default function Filters({ fields, onFilter }) {
  const [filters, setFilters] = useState({});

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFilters((prev) => ({ ...prev, [name]: value }));
  };

  const handleApplyFilters = () => {
    onFilter(filters);
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
      <button
        onClick={handleApplyFilters}
        className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
      >
        Apply Filters
      </button>
    </div>
  );
}
