import Filters from '@/components/shared/Filters';
import Table from '@/components/shared/Table';
import { Button } from '@/components/ui/button';
import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader } from '@/components/ui/card';

export default function ReportsAnalyticsPage() {
  const [filters, setFilters] = useState({});
  const [data, setData] = useState([]);
  const [filteredData, setFilteredData] = useState([]);
  const [loading, setLoading] = useState(false);

  const kpiData = [
    { title: 'Total Revenue', value: '$10,000' },
    { title: 'Orders', value: '150' },
    { title: 'Low Stock', value: '12 Items' },
    { title: 'Average Shipping Time', value: '3 Days' },
  ];

  const tableColumns = [
    { title: 'Date', key: 'date', sortable: true },
    { title: 'Category', key: 'category', sortable: true },
    { title: 'Metric', key: 'metric', sortable: true },
    { title: 'Value', key: 'value', sortable: true },
  ];

  useEffect(() => {
    setLoading(true);
    // Fetch data from API
    const timeout = setTimeout(() => {
      const fetchedData = Array.from({ length: 50 }, (_, i) => ({
        id: i + 1,
        date: `2024-${String(Math.ceil(Math.random() * 12)).padStart(
          2,
          '0'
        )}-${String(Math.ceil(Math.random() * 28)).padStart(2, '0')}`,
        category: ['Inventory', 'Orders', 'Shipping'][
          Math.floor(Math.random() * 3)
        ],
        metric: ['Revenue', 'Orders Processed', 'Items Sold'][
          Math.floor(Math.random() * 3)
        ],
        value: (Math.random() * 1000).toFixed(2),
      }));
      setData(fetchedData);
      setFilteredData(fetchedData);
      setLoading(false);
    }, 1000);

    return () => clearTimeout(timeout);
  }, []);

  const handleFilter = (appliedFilters) => {
    setFilters(appliedFilters);
    const filtered = data.filter((row) => {
      const matchCategory = appliedFilters.category
        ? row.category.toLowerCase() === appliedFilters.category.toLowerCase()
        : true;
      const matchDate =
        appliedFilters.startDate && appliedFilters.endDate
          ? row.date >= appliedFilters.startDate &&
            row.date <= appliedFilters.endDate
          : true;

      return matchCategory && matchDate;
    });
    setFilteredData(filtered);
  };

  const filterFields = [
    {
      name: 'category',
      type: 'select',
      placeholder: 'Select Category',
      options: ['All', 'Inventory', 'Orders', 'Shipping'],
    },
    {
      name: 'startDate',
      type: 'date',
      placeholder: 'Start Date',
    },
    {
      name: 'endDate',
      type: 'date',
      placeholder: 'End Date',
    },
  ];

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      <h1 className="text-3xk font-bold mb-6 text-gray-800 dark:text-white">
        Reports and Analytics
      </h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        {kpiData.map((kpi, index) => (
          <Card key={index} className="shadow-lg">
            <CardHeader className="font-bold text-lg">{kpi.title}</CardHeader>
            <CardContent className="text-2xl font-semibold">
              {kpi.value}
            </CardContent>
          </Card>
        ))}
      </div>
      <Filters fields={filterFields} onFilter={handleFilter} />
      <Table
        columns={tableColumns}
        data={filteredData}
        loading={loading}
        pageSize={10}
        actions={(row) => (
          <Button variant="link" onClick={() => console.log(row)}>
            View Details
          </Button>
        )}
      />
    </div>
  );
}
