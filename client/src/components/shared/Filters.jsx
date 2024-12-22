import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

export default function Filters({ fields, onFilter }) {
  const [filters, setFilters] = useState({});

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFilters((prev) => ({ ...prev, [name]: value }));
  };

  const handleSelectChange = (name, value) => {
    setFilters((prev) => ({ ...prev, [name]: value }));
  };

  const handleApplyFilters = () => {
    onFilter(filters);
  };

  const handleResetFilters = () => {
    setFilters({});
    onFilter({});
  };

  return (
    <div className="flex flex-row w-fit sm:flex-row gap-4 mb-6 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
      {fields.map((field) => {
        if (field.type === 'text') {
          return (
            <Input
              key={field.name}
              type="text"
              name={field.name}
              placeholder={field.placeholder}
              value={filters[field.name] || ''}
              onChange={handleInputChange}
              className="w-44"
            />
          );
        }
        if (field.type === 'select') {
          return (
            <Select
              key={field.name}
              onValueChange={(value) => handleSelectChange(field.name, value)}
              value={filters[field.name] || ''}
            >
              <SelectTrigger className="w-44 border border-gray-300 bg-white dark:bg-gray-800 dark:border-gray-700">
                <SelectValue placeholder={field.placeholder} />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700">
                <SelectGroup>
                  {field.options.map((option) => (
                    <SelectItem key={option} value={option}>
                      {option}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          );
        }
        if (field.type === 'number') {
          return (
            <Input
              key={field.name}
              type="number"
              name={field.name}
              placeholder={field.placeholder}
              value={filters[field.name] || ''}
              onChange={handleInputChange}
              className="w-36"
            />
          );
        }
        return null;
      })}
      <Button variant="blue" onClick={handleApplyFilters}>
        Apply Filters
      </Button>
      <Button variant="destructive" onClick={handleResetFilters}>
        Reset Filters
      </Button>
    </div>
  );
}
