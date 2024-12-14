const { Button } = require('@/components/ui/button');
const { useDispatch } = require('react-redux');

const AddItemButton = () => {
  const dispatch = useDispatch();
};

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
  <Button variant="default" onClick={handleAddItem}>
    Add New Item
  </Button>
);

export default AddItemButton;
