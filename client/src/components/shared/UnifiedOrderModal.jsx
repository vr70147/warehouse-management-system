import React, { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';

export default function UnifiedOrderModal({
  isOpen,
  onClose,
  mode,
  order,
  onSubmit,
}) {
  const [formData, setFormData] = useState({
    customer_name: '',
    status: '',
    items: [],
    total_price: 0,
    shipping_date: '',
    notes: '',
  });

  useEffect(() => {
    if (mode === 'edit' && order) {
      setFormData({
        customer_name: order.customer_name || '',
        status: order.status || '',
        items: order.items || [],
        total_price: order.total_price || 0,
        shipping_date: order.shipping_date || '',
        notes: order.notes || '',
      });
    } else if (mode === 'add') {
      setFormData({
        customer_name: '',
        status: '',
        items: [],
        total_price: 0,
        shipping_date: '',
        notes: '',
      });
    }
  }, [mode, order]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleItemsChange = (index, field, value) => {
    const updatedItems = [...formData.items];
    updatedItems[index] = { ...updatedItems[index], [field]: value };
    setFormData((prev) => ({ ...prev, items: updatedItems }));
  };

  const handleAddItem = () => {
    setFormData((prev) => ({
      ...prev,
      items: [...prev.items, { name: '', quantity: 1 }],
    }));
  };

  const handleRemoveItem = (index) => {
    const updatedItems = formData.items.filter((_, i) => i !== index);
    setFormData((prev) => ({ ...prev, items: updatedItems }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!Array.isArray(formData.items)) {
      console.error('Items is not a valid array:', formData.items);
      return;
    }
    const payload = {
      id: order?.id,
      ...formData,
      items: formData.items || [],
    };
    onSubmit(payload);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center z-50">
      <div
        className="absolute inset-0 bg-black bg-opacity-50"
        onClick={onClose}
      ></div>
      <div className="relative bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg w-1/3">
        <h2 className="text-xl font-bold mb-4">
          {mode === 'add' ? 'Add Order' : 'Edit Order'}
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label>Customer Name</label>
            <input
              type="text"
              name="customer_name"
              value={formData.customer_name}
              onChange={handleChange}
              required
            />
          </div>
          <div>
            <label>Status</label>
            <input
              type="text"
              name="status"
              value={formData.status}
              onChange={handleChange}
              required
            />
          </div>
          <div>
            <label>Items</label>
            {formData.items.map((item, index) => (
              <div key={index} className="flex gap-2">
                <input
                  type="text"
                  value={item.name}
                  onChange={(e) =>
                    handleItemsChange(index, 'name', e.target.value)
                  }
                  placeholder="Item Name"
                />
                <input
                  type="number"
                  value={item.quantity}
                  onChange={(e) =>
                    handleItemsChange(index, 'quantity', e.target.value)
                  }
                  placeholder="Quantity"
                  min="1"
                />
                <Button onClick={() => handleRemoveItem(index)}>Remove</Button>
              </div>
            ))}
            <Button onClick={handleAddItem}>Add Item</Button>
          </div>
          <div>
            <label>Total Price</label>
            <input
              type="number"
              name="total_price"
              value={formData.total_price}
              onChange={handleChange}
              required
            />
          </div>
          <div>
            <label>Shipping Date</label>
            <input
              type="date"
              name="shipping_date"
              value={formData.shipping_date}
              onChange={handleChange}
              required
            />
          </div>
          <div>
            <label>Notes</label>
            <textarea
              name="notes"
              value={formData.notes}
              onChange={handleChange}
            />
          </div>
          <div className="flex justify-end gap-2">
            <Button type="button" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit">{mode === 'add' ? 'Add' : 'Save'}</Button>
          </div>
        </form>
      </div>
    </div>
  );
}
