const companyNames = [
    'Teva Pharmaceutical', 'Check Point Software', 'Amdocs', 'NICE Systems',
    'Elbit Systems', 'Bank Leumi', 'Bank Hapoalim', 'Israel Aerospace Industries',
    'Mobileye', 'Wix', 'Monday.com', 'Fiverr', 'Taboola', 'IronSource',
    'Playtika', 'Tower Semiconductor', 'Strauss Group', 'El Al', 'ZIM Shipping', 'ICL Group'
];

const generateRandomDate = (isFuture) => {
    const today = new Date();
    const offset = Math.floor(Math.random() * 365);
    const newDate = new Date(today);
    newDate.setDate(isFuture ? today.getDate() + offset : today.getDate() - offset);
    return newDate.toISOString().split('T')[0];
};

const dummyOrders = Array.from({ length: 50 }, (_, i) => ({
    id: `ORD${String(i + 1).padStart(3, '0')}`,
    customer_name: companyNames[Math.floor(Math.random() * companyNames.length)],
    status: ['Pending', 'Shipped', 'Delivered', 'Cancelled'][Math.floor(Math.random() * 4)],
    items: Array.from({ length: Math.ceil(Math.random() * 3) }, () => ({
        name: [
            'Laptop', 'Mouse', 'Keyboard', 'Monitor', 'Desk Chair',
            'Headphones', 'Smartphone', 'Charger'
        ][Math.floor(Math.random() * 8)],
        quantity: Math.ceil(Math.random() * 5),
    })),
    total_price: (Math.random() * 1000 + 50).toFixed(2),
    created_at: generateRandomDate(Math.random() > 0.5),
    updated_at: generateRandomDate(Math.random() > 0.5),
    shipping_date: generateRandomDate(Math.random() > 0.5),
    notes: Math.random() > 0.5 ? 'Special delivery instructions' : '',
}));

export default dummyOrders;
