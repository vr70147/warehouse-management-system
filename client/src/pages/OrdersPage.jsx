import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button } from '@/components/ui/button';
import UnifiedOrderModal from '@/components/shared/UnifiedOrderModal';
import dummyOrders from '@/data/dummyOrders';
import {
  addOrder,
  deleteOrder,
  updateOrder,
  filterOrders,
  setLoading,
} from '@/redux/slices/ordersSlice';
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import Filter from '@/components/shared/Filters';
import Table from '@/components/shared/Table';
import { v4 as uuidv4 } from 'uuid';

export default function OrderPage() {
  const dispatch = useDispatch();
  const filteredOrders = useSelector((state) => state.orders.filteredOrders);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const loading = useSelector((state) => state.orders.loading);
  const [currentOrder, setCurrentOrder] = useState(null);
  const [modalMode, setModalMode] = useState('add');
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [sheetContent, setSheetContent] = useState([]);

  const DEFAULT_PAGE_SIZE = 10;

  useEffect(() => {
    dispatch(setLoading(true));
    const timer = setTimeout(() => {
      dispatch(setLoading(false));
    }, 1000);

    return () => {
      clearTimeout(timer);
    };
  }, [dispatch]);

  const handleAddOrder = (newItem) => {
    dispatch(addOrder({ ...newItem, id: uuidv4() }));
    setIsModalOpen(false);
  };

  const openItemSheet = (items) => {
    console.log('Opening sheet with items:', items);
    if (Array.isArray(items)) {
      setSheetContent(items);
      setIsSheetOpen(true);
    } else {
      console.error('Invalid items data:', items);
    }
  };

  const handleFilter = (filters) => {
    dispatch(filterOrders(filters));
  };

  const handleEditOrder = (order) => {
    setModalMode('edit');
    setCurrentOrder(order);
    setIsModalOpen(true);
  };

  const handleOpenAddModal = () => {
    setModalMode('add');
    setCurrentOrder(null);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setCurrentOrder(null);
  };

  const handleUpdateOrder = (updatedOrder) => {
    dispatch(updateOrder(updatedOrder));
    setIsModalOpen(false);
  };

  const handleDeleteOrder = (id) => {
    dispatch(deleteOrder(id));
  };

  const renderActions = (row) => (
    <>
      <Button variant="outline" size="sm" onClick={() => handleEditOrder(row)}>
        Edit
      </Button>
      <Button
        variant="destructive"
        size="sm"
        onClick={() => handleDeleteOrder(row.id)}
      >
        Delete
      </Button>
    </>
  );

  const filterFields = [
    {
      name: 'status',
      type: 'select',
      placeholder: 'Select Status',
      options: [
        'Select Status',
        ...new Set(dummyOrders.map((order) => order.status)),
      ],
    },
    {
      name: 'customer_name',
      type: 'text',
      placeholder: 'Search by Customer Name',
    },
    {
      name: 'totalPriceMin',
      type: 'number',
      placeholder: 'Min Total Price',
    },
    {
      name: 'totalPriceMax',
      type: 'number',
      placeholder: 'Max Total Price',
    },
    {
      name: 'shippingDateFrom',
      type: 'date',
      placeholder: 'Shipping Date From',
    },
    {
      name: 'shippingDateTo',
      type: 'date',
      placeholder: 'Shipping Date To',
    },
  ];

  const columns = [
    { title: 'Order ID', key: 'id' },
    { title: 'Customer Name', key: 'customer_name', sortable: true },
    { title: 'Order Status', key: 'status' },
    {
      title: 'Items',
      key: 'items',
      render: (items) =>
        Array.isArray(items) && items.length > 0 ? (
          <Button variant="outline" onClick={() => openItemSheet(items)}>
            View Items
          </Button>
        ) : (
          'No Items'
        ),
    },
    { title: 'Total Price', key: 'total_price', sortable: true },
    { title: 'Created At', key: 'created_at', sortable: true },
    { title: 'Updated At', key: 'updated_at', sortable: true },
    { title: 'Shipping Date', key: 'shipping_date', sortable: true },
    { title: 'Notes', key: 'notes' },
  ];

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white">
          Orders Management
        </h1>
        <Button variant="blue" onClick={handleOpenAddModal}>
          Add New Item
        </Button>
      </div>
      <Filter onFilter={handleFilter} fields={filterFields} />
      <Table
        columns={columns}
        actions={renderActions}
        data={filteredOrders}
        loading={loading}
        pageSize={DEFAULT_PAGE_SIZE}
        rowHighlightCondition={(row) =>
          new Date(row.shipping_date) < new Date()
        }
        rowHighlightClasses={{
          even: 'even:bg-blue-50 dark:even:bg-blue-900',
          odd: 'odd:bg-white dark:odd:bg-blue-800',
        }}
      />
      <UnifiedOrderModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        mode={currentOrder ? 'edit' : 'add'}
        order={currentOrder}
        onSubmit={modalMode === 'add' ? handleAddOrder : handleUpdateOrder}
      />
      <Sheet open={isSheetOpen} onOpenChange={setIsSheetOpen}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Order Items</SheetTitle>
            <SheetDescription>
              Below is the list of items in this order.
            </SheetDescription>
          </SheetHeader>
          <div className="grid gap-4 py-4">
            {sheetContent.map((item, index) => (
              <div key={index} className="grid grid-cols-2 gap-4">
                <span className="font-medium">{item.name}</span>
                <span>Quantity: {item.quantity}</span>
              </div>
            ))}
          </div>
          <SheetFooter>
            <SheetClose asChild>
              <Button variant="outline">Close</Button>
            </SheetClose>
          </SheetFooter>
        </SheetContent>
      </Sheet>
    </div>
  );
}
