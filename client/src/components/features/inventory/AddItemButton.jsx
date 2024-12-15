import { Button } from '@/components/ui/button';
import { useDispatch } from 'react-redux';
import { addItem } from '@/redux/slices/inventorySlice';
import { Item } from '@radix-ui/react-dropdown-menu';

const AddItemButton = () => {
  const dispatch = useDispatch();

  const handleAddItem = () => {
    const newItem = {
      id: Math.random(),
      name: 'New Product',
      category: 'Misc',
      quantity: 10,
      unitPrice: 25.0,
      supplier: 'Default Supplier',
      lastUpdated: new Date().toISOString(),
    };

    dispatch(addItem(newItem));
  };
  return (
    <Button variant="outline" onClick={handleAddItem}>
      Add New Item
    </Button>
  );
};

export default AddItemButton;
