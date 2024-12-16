import React from 'react';

export default function InfoCard({ title, value, iconColor, iconPath }) {
  return (
    <div className="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4 flex items-center">
      <div className={`p-3 ${iconColor} rounded-full`}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-6 w-6 text-blue-500 dark:text-blue-300"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d={iconPath}
          />
        </svg>
      </div>
    </div>
  );
}
