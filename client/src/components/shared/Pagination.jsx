import React from 'react';

export default function Pagination({ currentPage, totalPages, onPageChange }) {
  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handlePageClick = (page) => {
    onPageChange(page);
  };

  return (
    <nav className="flex justify-between items-center mt-4">
      <button
        onClick={handlePrevious}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white rounded disabled:opacity-50"
        disabled={currentPage === 1}
      >
        Previous
      </button>
      <div className="flex space-x-2">
        {Array.from({ length: totalPages }, (_, index) => (
          <button
            key={index + 1}
            onClick={() => handlePageClick(index + 1)}
            className={`px-3 py-1 rounded ${
              currentPage === index + 1
                ? 'bg-blue-500 text-white'
                : 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white'
            }`}
          >
            {index + 1}
          </button>
        ))}
      </div>
      <button
        onClick={handleNext}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white rounded disabled:opacity-50"
        disabled={currentPage === totalPages}
      >
        Next
      </button>
    </nav>
  );
}
