import React from 'react';
import { useState } from 'react';

export default function InventoryFilters({ onFilter, onSearch, onSort }) {
  const [category, setCategory] = useState('');
  const [supplier, setSupplier] = useState('');
  const [priceRange, setPriceRange] = useState({ min: '', max: '' });
  const [sortOrder, setSortOrder] = useState('');

  const handleFilter = () => {
    onFilter({ category, supplier, priceRange });
  };

  const handleSearch = (e) => {
    const searchTerm = e.target.value;
    onSearch(searchTerm);
  };

  const handleSort = () => {
    const newOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    setSortOrder(newOrder);
    onSort(newOrder);
  };

  return (
    <div className="flex flex-row gap-4 mb-6 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
      <input
        type="text"
        placeholder="Search by product name..."
        onChange={handleSearch}
        className="rounded bg-gray-200 dark:bg-gray-700 p-2"
      />
      <button
        onClick={handleFilter}
        className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      >
        Apply Filter
      </button>
      <button
        onClick={handleSort}
        className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      >
        Sort by Quantity ({sortOrder === 'asc' ? 'asc' : 'desc'}){' '}
      </button>
    </div>
  );
}
