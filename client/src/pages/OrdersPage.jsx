import OrderTable from '@/components/features/orders/OrderTable';
import React from 'react';

export default function OrdersPage() {
  return (
    <div className="p-6 min-h-screen bg-gray-100 dark:bg-gray-900 dark:text-white">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-3xl font-bold">Orders Management</h1>
      </div>
      <OrderTable />
    </div>
  );
}
