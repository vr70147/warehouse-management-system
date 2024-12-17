import React, { useState } from 'react';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
import { paginate } from '@/utils/pagination';
import { ChevronUp, ChevronDown } from 'lucide-react';

export default function Table({
  columns,
  data,
  actions,
  loading,
  pageSize = 10,
}) {
  const [currentPage, setCurrentPage] = useState(1);
  const [sortConfig, setSortConfig] = useState({ key: null, direction: 'asc' });

  const sortedData = React.useMemo(() => {
    if (!sortConfig.key) return data;
    return [...data].sort((a, b) => {
      const valueA = a[sortConfig.key] ?? '';
      const valueB = b[sortConfig.key] ?? '';
      return sortConfig.direction === 'asc'
        ? valueA > valueB
          ? 1
          : -1
        : valueA < valueB
        ? 1
        : -1;
    });
  }, [data, sortConfig]);

  const paginatedData = paginate(sortedData, currentPage, pageSize);
  const totalPages = Math.ceil(sortedData.length / pageSize);

  const handleSort = (key) => {
    setSortConfig((prev) => ({
      key,
      direction: prev.key === key && prev.direction === 'asc' ? 'desc' : 'asc',
    }));
  };

  return (
    <div className="bg-white dark:bg-gray-800 p-4 overflow-x-auto rounded-lg shadow-lg">
      <table className="w-full text-left text-md text-gray-600 dark:text-gray-400">
        <thead>
          <tr>
            {columns.map((column) => (
              <th
                key={column.key}
                onClick={() => handleSort(column.key)}
                className="py-2 px-4 cursor-pointer select-none hover:text-blue-500"
              >
                <div className="flex items-center justify-between">
                  {column.title}
                  {sortConfig.key === column.key &&
                    (sortConfig.direction === 'asc' ? (
                      <ChevronUp className="w-4 h-4 inline-block" />
                    ) : (
                      <ChevronDown className="w-4 h-4 inline-block" />
                    ))}
                </div>
              </th>
            ))}
            {actions && <th className="py-2 px-4">Actions</th>}
          </tr>
        </thead>
        <tbody>
          {loading
            ? Array.from({ length: pageSize }).map((_, rowIndex) => (
                <tr key={rowIndex}>
                  {columns.map((column) => (
                    <td key={column.key} className="py-3 px-4">
                      <Skeleton className="h-6 w-full" />
                    </td>
                  ))}
                  {actions && (
                    <td className="py-3 px-4">
                      <Skeleton className="h-6 w-full" />
                    </td>
                  )}
                </tr>
              ))
            : paginatedData.map((row, rowIndex) => (
                <tr
                  key={rowIndex}
                  className={`hover:scale-102 hover:shadow transition-transform duration-300 ease-out ${
                    row.quantity <= 20
                      ? 'even:bg-red-200 dark:even:bg-red-950 odd:bg-red-100 dark:odd:bg-red-900'
                      : 'even:bg-gray-50 odd:bg-white dark:even:bg-gray-700 dark:odd:bg-gray-800'
                  }`}
                >
                  {columns.map((column) => (
                    <td key={column.key} className="py-3 px-4">
                      {row[column.key] || '-'}
                    </td>
                  ))}
                  {actions && (
                    <td className="py-3 px-4 flex gap-2">{actions(row)}</td>
                  )}
                </tr>
              ))}
        </tbody>
      </table>
      {totalPages > 1 && (
        <Pagination className="mt-4">
          <PaginationContent>
            <PaginationItem>
              <PaginationPrevious
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  if (currentPage > 1) setCurrentPage(currentPage - 1);
                }}
              />
            </PaginationItem>
            {Array.from({ length: totalPages }, (_, index) => (
              <PaginationItem key={index}>
                <PaginationLink
                  className={
                    currentPage === index + 1
                      ? 'bg-gray-200 dark:bg-gray-700'
                      : ''
                  }
                  href="#"
                  onClick={(e) => {
                    e.preventDefault();
                    setCurrentPage(index + 1);
                  }}
                >
                  {index + 1}
                </PaginationLink>
              </PaginationItem>
            ))}
            <PaginationItem>
              <PaginationNext
                href="#"
                onClick={(e) => {
                  e.preventDefault();
                  if (currentPage < totalPages) setCurrentPage(currentPage + 1);
                }}
              />
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      )}
    </div>
  );
}
