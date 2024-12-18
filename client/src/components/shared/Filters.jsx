import React, { useState } from 'react';
import { Button } from '../ui/button';

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
    <div className="flex flex-row w-fit sm:flex-row gap-4 mb-6 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
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
              className="block w-44 px-4 py-2 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-black focus:border-black dark:bg-gray-800 dark:text-gray-200 dark:border-gray-600 transition duration-200"
            >
              <option value="" disabled>
                {field.placeholder}
              </option>
              {field.options.map((option) => (
                <option
                  key={option}
                  value={option}
                  className="text-gray-700 dark:text-gray-300"
                >
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
              className="rounded bg-gray-200 dark:bg-gray-700 p-2 w-28"
            />
          );
        }
        return null;
      })}
      <Button variant="blue" onClick={handleApplyFilters}>
        Apply Filters
      </Button>
    </div>
  );
}
