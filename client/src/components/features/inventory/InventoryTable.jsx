import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { useDispatch } from 'react-redux';
import { deleteItem, updateItem } from '@/redux/slices/inventorySlice';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from '@/components/ui/pagination';
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
import UnifiedItemModal from '@/components/features/inventory/UnifiedItemModal';

export default function InventoryTable({ items, loading }) {
  const [currentPage, setCurrentPage] = useState(1);
  const pageSize = 10;

  const totalPages = Math.ceil(items.length / pageSize);
  const PaginatedItems = paginate(items, currentPage, pageSize);

  const dispatch = useDispatch();

  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [currentItem, setCurrentItem] = useState(null);

  const handleEdit = (item) => {
    setCurrentItem(item);
    setIsEditModalOpen(true);
  };

  const handleUpdateItem = (updatedItem) => {
    dispatch(updateItem(updatedItem));
  };

  const handleDelete = (id) => {
    dispatch(deleteItem(id));
  };

  if (loading) {
    return (
      <div className="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4">
        <table className="w-full text-left text-md text-gray-600 dark:text-gray-400">
          <thead>
            <tr>
              <th className="py-2 px-4">Name</th>
              <th className="py-2 px-4">Category</th>
              <th className="py-2 px-4">Supplier</th>
              <th className="py-2 px-4">Quantity</th>
              <th className="py-2 px-4">Unit Price</th>
              <th className="py-2 px-4">Actions</th>
            </tr>
          </thead>
          <tbody>
            {Array.from({ length: 5 }).map((_, index) => (
              <tr key={index}>
                <td className="py-3 px-4">
                  <Skeleton className="h-4 w-24" />
                </td>
                <td className="py-3 px-4">
                  <Skeleton className="h-4 w-32" />
                </td>
                <td className="py-3 px-4">
                  <Skeleton className="h-4 w-28" />
                </td>
                <td className="py-3 px-4">
                  <Skeleton className="h-4 w-16" />
                </td>
                <td className="py-3 px-4">
                  <Skeleton className="h-4 w-20" />
                </td>
                <td className="py-3 px-4">
                  <Skeleton className="h-8 w-16" />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4">
      <table className="w-full text-left text-md text-gray-600 dark:text-gray-400">
        <thead>
          <tr>
            <th className="py-2 px-4">Name</th>
            <th className="py-2 px-4">Category</th>
            <th className="py-2 px-4">Supplier</th>
            <th className="py-2 px-4">Quantity</th>
            <th className="py-2 px-4">Unit Price</th>
            <th className="py-2 px-4">Actions</th>
          </tr>
        </thead>
        <tbody>
          {PaginatedItems.map((item) => (
            <tr
              key={item.id}
              className={` hover:scale-102 hover:shadow transition-transform duration-300 ease-out ${
                item.quantity <= 20
                  ? 'even:bg-red-200 dark:even:bg-red-950 odd:bg-red-100 dark:odd:bg-red-900'
                  : 'even:bg-gray-50 odd:bg-white dark:even:bg-gray-700 dark:odd:bg-gray-800'
              }`}
            >
              <td className="py-3 px-4 dark:border-gray-700">{item.name}</td>
              <td className="py-3 px-4 dark:border-gray-700">
                {item.category}
              </td>
              <td className="py-3 px-4 dark:border-gray-700">
                {item.supplier}
              </td>
              <td
                className={`py-3 px-4 dark:border-gray-700 ${
                  item.quantity <= 20 ? 'text-red-600 font-bold' : ''
                }`}
              >
                {item.quantity}
              </td>
              <td className="py-3 px-4 dark:border-gray-700">
                $
                {isNaN(item.unitPrice)
                  ? '0.00'
                  : parseFloat(item.unitPrice).toFixed(2)}
              </td>
              <td className="py-3 px-4 dark:border-gray-700 flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleEdit(item)}
                >
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
                      <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                      <AlertDialogDescription>
                        This action cannot be undone. It will permanently delete
                        <strong> {item.name}</strong>
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction onClick={() => handleDelete(item.id)}>
                        Delete
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
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
                href="#"
                isActive={currentPage === index + 1}
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
      <UnifiedItemModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        mode="edit"
        item={currentItem}
        onSubmit={handleUpdateItem}
      />
    </div>
  );
}

function paginate(items, currentPage, pageSize) {
  if (!items || items.length === 0) return [];
  const startIndex = (currentPage - 1) * pageSize;
  return items.slice(startIndex, startIndex + pageSize);
}
