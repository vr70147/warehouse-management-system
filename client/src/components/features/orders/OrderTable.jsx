import { Button } from '@/components/ui/button';
import React from 'react';

export default function OrderTable() {
  return (
    <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-lg">
      <table className="w-full text-left text-md text-gray-600 dark:text-gray-400">
        <thead>
          <th className="py-2 px-4">Id</th>
          <th className="py-2 px-4">Customer Name</th>
          <th className="py-2 px-4">Status</th>
          <th className="py-2 px-4">Order Date</th>
          <th className="py-2 px-4">Delivery Date</th>
          <th className="py-2 px-4">Total Amount</th>
          <th className="py-2 px-4">Action</th>
        </thead>
        <tbody>
          <tr>
            <td className="py-3 px-4">1</td>
            <td className="py-3 px-4">John Doe</td>
            <td className="py-3 px-4">Processing</td>
            <td className="py-3 px-4">2023-06-01</td>
            <td className="py-3 px-4">2023-06-05</td>
            <td className="py-3 px-4">$50.00</td>
            <td className="py-3 px-4">
              <Button variant="blue">View</Button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
