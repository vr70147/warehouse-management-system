import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Button } from '@/components/ui/button';
import UnifiedItemModal from '@/components/shared/UnifiedItemModal';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import {
  addItem,
  updateItem,
  deleteItem,
  filterItems,
  setLoading,
} from '@/redux/slices/inventorySlice';
import Filter from '@/components/shared/Filters';
import Table from '@/components/shared/Table';
import { v4 as uuidv4 } from 'uuid';

export default function InventoryPage() {
  const dispatch = useDispatch();
  const filteredItems = useSelector((state) => state.inventory.filteredItems);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const loading = useSelector((state) => state.inventory.loading);
  const [currentItem, setCurrentItem] = useState(null);
  const [modalMode, setModalMode] = useState('add');

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

  const handleAddItem = (newItem) => {
    dispatch(addItem({ ...newItem, id: uuidv4() }));
    setIsModalOpen(false);
  };

  const handleFilter = (filters) => {
    dispatch(filterItems(filters));
  };

  const handleEditItem = (item) => {
    setModalMode('edit');
    setCurrentItem(item);
    setIsModalOpen(true);
  };

  const handleOpenAddModal = () => {
    setModalMode('add');
    setCurrentItem(null);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setCurrentItem(null);
  };

  const handleUpdateItem = (updatedItem) => {
    dispatch(updateItem(updatedItem));
    setIsModalOpen(false);
  };

  const handleDeleteItem = (id) => {
    dispatch(deleteItem(id));
  };

  const renderActions = (row) => (
    <>
      <Button variant="outline" size="sm" onClick={() => handleEditItem(row)}>
        Edit
      </Button>
      <AlertDialog>
        <AlertDialogTrigger asChild>
          <Button variant="destructive" size="sm">
            Delete
          </Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Confir`m Deletion</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this order? This action cannot be
              undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                handleDeleteItem(row.id);
              }}
            >
              Confirm
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );

  const filterFields = [
    {
      name: 'category',
      type: 'select',
      placeholder: 'Select Category',
      options: Array.from(new Set(filteredItems.map((item) => item.category))),
    },
    {
      name: 'name',
      type: 'text',
      placeholder: 'Search by Name',
    },
    {
      name: 'supplier',
      type: 'text',
      placeholder: 'Search by Supplier',
    },
    {
      name: 'priceMin',
      type: 'number',
      placeholder: 'Min Price',
    },
    {
      name: 'priceMax',
      type: 'number',
      placeholder: 'Max Price',
    },
  ];

  const columns = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Category', key: 'category', sortable: true },
    { title: 'Supplier', key: 'supplier', sortable: true },
    { title: 'Quantity', key: 'quantity', sortable: true },
    { title: 'Unit Price', key: 'unitPrice', sortable: true },
    { title: 'Last Updated', key: 'lastUpdated', sortable: true },
  ];

  return (
    <div className="p-6 bg-gray-100 dark:bg-gray-900 min-h-screen">
      <div className="mb-6 flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-800 dark:text-white">
          Inventory Management
        </h1>
        <Button variant="blue" onClick={handleOpenAddModal}>
          Add New Item
        </Button>
      </div>
      <Filter onFilter={handleFilter} fields={filterFields} />
      <Table
        columns={columns}
        actions={renderActions}
        data={filteredItems}
        loading={loading}
        pageSize={DEFAULT_PAGE_SIZE}
        rowHighlightCondition={(row) => row.quantity <= 20}
        rowHighlightClasses={{
          even: 'even:bg-gray-50 dark:even:bg-gray-700',
          odd: 'odd:bg-white dark:odd:bg-gray-800',
        }}
      />
      <UnifiedItemModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        mode={currentItem ? 'edit' : 'add'}
        item={currentItem}
        onSubmit={modalMode === 'add' ? handleAddItem : handleUpdateItem}
      />
    </div>
  );
}
