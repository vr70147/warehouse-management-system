import React, { useEffect, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';

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
    const payload = {
      id: order?.id,
      ...formData,
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
      <div className="relative bg-white dark:bg-gray-800 p-6 rounded-lg shadow-lg w-full max-w-lg">
        <h2 className="text-xl font-bold mb-6 text-gray-800 dark:text-white">
          {mode === 'add' ? 'Add Order' : 'Edit Order'}
        </h2>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="customer_name">Customer Name</Label>
            <Input
              type="text"
              id="customer_name"
              name="customer_name"
              value={formData.customer_name}
              onChange={handleChange}
              placeholder="Enter customer name"
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <Input
              type="text"
              id="status"
              name="status"
              value={formData.status}
              onChange={handleChange}
              placeholder="Enter status"
              required
            />
          </div>
          <div className="space-y-2">
            <Label>Items</Label>
            <div className="space-y-3">
              {formData.items.map((item, index) => (
                <div key={index} className="flex gap-4 items-center">
                  <Input
                    type="text"
                    value={item.name}
                    onChange={(e) =>
                      handleItemsChange(index, 'name', e.target.value)
                    }
                    placeholder="Item Name"
                    className="flex-1"
                  />
                  <Input
                    type="number"
                    value={item.quantity}
                    onChange={(e) =>
                      handleItemsChange(index, 'quantity', e.target.value)
                    }
                    placeholder="Quantity"
                    min="1"
                    className="w-20"
                  />
                  <Button
                    variant="destructive"
                    size="icon"
                    onClick={() => handleRemoveItem(index)}
                  >
                    ×
                  </Button>
                </div>
              ))}
            </div>
            <Button variant="outline" className="mt-3" onClick={handleAddItem}>
              + Add Item
            </Button>
          </div>
          <div className="space-y-2">
            <Label htmlFor="total_price">Total Price</Label>
            <Input
              type="number"
              id="total_price"
              name="total_price"
              value={formData.total_price}
              onChange={handleChange}
              placeholder="Enter total price"
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="shipping_date">Shipping Date</Label>
            <Input
              type="date"
              id="shipping_date"
              name="shipping_date"
              value={formData.shipping_date}
              onChange={handleChange}
              required
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="notes">Notes</Label>
            <Textarea
              id="notes"
              name="notes"
              value={formData.notes}
              onChange={handleChange}
              placeholder="Enter any additional notes"
            />
          </div>
          <div className="flex justify-end gap-4 mt-4">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" variant="blue">
              {mode === 'add' ? 'Add' : 'Save'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
