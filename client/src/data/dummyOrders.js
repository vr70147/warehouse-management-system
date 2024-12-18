const dummyOrders = Array.from({ length: 50 }, (_, i) => ({
    id: `ORD${String(i + 1).padStart(3, '0')}`,
    customer_name: `Customer ${i + 1}`,
    status: ['Pending', 'Shipped', 'Delivered', 'Cancelled'][Math.floor(Math.random() * 4)],
    items: Array.from({ length: Math.ceil(Math.random() * 3) }, (_, j) => ({
        name: ['Laptop', 'Mouse', 'Keyboard', 'Monitor', 'Desk Chair', 'Headphones', 'Smartphone', 'Charger'][Math.floor(Math.random() * 8)],
        quantity: Math.ceil(Math.random() * 5),
    })),
    total_price: (Math.random() * 1000 + 50).toFixed(2),
    created_at: `2024-${String(Math.ceil(Math.random() * 12)).padStart(2, '0')}-${String(Math.ceil(Math.random() * 28)).padStart(2, '0')}`,
    updated_at: `2024-${String(Math.ceil(Math.random() * 12)).padStart(2, '0')}-${String(Math.ceil(Math.random() * 28)).padStart(2, '0')}`,
    shipping_date: `2024-${String(Math.ceil(Math.random() * 12)).padStart(2, '0')}-${String(Math.ceil(Math.random() * 28)).padStart(2, '0')}`,
    notes: Math.random() > 0.5 ? 'Special delivery instructions' : '',
}));

export default dummyOrders;
