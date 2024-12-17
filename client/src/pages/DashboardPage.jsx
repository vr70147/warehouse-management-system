import React from 'react';
import InfoCard from '../components/features/dashboard/InfoCard';
import GraphCard from '../components/features/dashboard/GraphCard';

export default function DashboardPage() {
  const salesData = {
    labels: ['January', 'February', 'March', 'April', 'May', 'June'],
    datasets: [
      {
        label: 'Sales ($)',
        data: [1200, 1900, 3000, 5000, 2300, 4000],
        backgroundColor: 'rgba(54, 162, 235, 0.5)',
        borderColor: 'rgba(54, 162, 235, 1)',
        borderWidth: 1,
      },
    ],
  };

  const salesOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Monthly Sales Data',
      },
    },
  };

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white">
          Dashboard Overview
        </h1>
        <p className="text-gray-600 dark:text-gray-400">
          Welcome back! Here's an overview of the warehouse.
        </p>
      </div>

      {/* InfoCards Section */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        <InfoCard
          title="Pending Orders"
          value="120"
          iconColor="bg-blue-100 dark:bg-blue-900"
          iconPath="M9.75 9.75h4.5v4.5h-4.5z"
        />
        <InfoCard
          title="Shipped Orders"
          value="98"
          iconColor="bg-green-100 dark:bg-green-900"
          iconPath="M3 3l18 18"
        />
        <InfoCard
          title="Low Inventory"
          value="15"
          iconColor="bg-yellow-100 dark:bg-yellow-900"
          iconPath="M12 8v8m0 0h8m-8 0H4"
        />
        <InfoCard
          title="Customers"
          value="250"
          iconColor="bg-purple-100 dark:bg-purple-900"
          iconPath="M17 9l4 4m0 0l-4 4m4-4H3"
        />
      </div>

      {/* GraphCards Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <GraphCard
          title="Sales Trend"
          chartData={salesData}
          options={salesOptions}
        />
        <GraphCard
          title="Inventory Breakdown"
          chartData={{
            labels: ['Electronics', 'Furniture', 'Clothing', 'Accessories'],
            datasets: [
              {
                label: 'Items',
                data: [50, 30, 20, 10],
                backgroundColor: [
                  'rgba(54, 162, 235, 0.5)',
                  'rgba(255, 206, 86, 0.5)',
                  'rgba(75, 192, 192, 0.5)',
                  'rgba(153, 102, 255, 0.5)',
                ],
                borderColor: [
                  'rgba(54, 162, 235, 1)',
                  'rgba(255, 206, 86, 1)',
                  'rgba(75, 192, 192, 1)',
                  'rgba(153, 102, 255, 1)',
                ],
                borderWidth: 1,
              },
            ],
          }}
          options={{
            responsive: true,
            plugins: {
              legend: {
                position: 'right',
              },
              title: {
                display: true,
                text: 'Inventory Distribution',
              },
            },
          }}
        />
      </div>

      {/* Table Section */}
      <div className="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4">
        <h3 className="text-lg font-bold text-gray-800 dark:text-white mb-4">
          Recent Orders
        </h3>
        <table className="w-full text-left text-sm text-gray-600 dark:text-gray-400">
          <thead>
            <tr>
              <th className="py-2 px-4">Order ID</th>
              <th className="py-2 px-4">Customer</th>
              <th className="py-2 px-4">Status</th>
              <th className="py-2 px-4">Total</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td className="py-2 px-4">#1001</td>
              <td className="py-2 px-4">John Doe</td>
              <td className="py-2 px-4">Pending</td>
              <td className="py-2 px-4">$250.00</td>
            </tr>
            <tr>
              <td className="py-2 px-4">#1002</td>
              <td className="py-2 px-4">Jane Smith</td>
              <td className="py-2 px-4">Shipped</td>
              <td className="py-2 px-4">$300.00</td>
            </tr>
            {/* Add more rows */}
          </tbody>
        </table>
      </div>
    </div>
  );
}
