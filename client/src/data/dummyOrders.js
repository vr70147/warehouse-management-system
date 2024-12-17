const dummyOrders = [
    {
        id: 1,
        customerName: "John Doe",
        status: "Pending",
        orderDate: "2024-12-01",
        deliveryDate: "2024-12-05",
        totalAmount: 150.5,
    },
    {
        id: 2,
        customerName: "Jane Smith",
        status: "Shipped",
        orderDate: "2024-11-28",
        deliveryDate: "2024-12-02",
        totalAmount: 200.0,
    },
    {
        id: 3,
        customerName: "Mike Johnson",
        status: "Delivered",
        orderDate: "2024-11-25",
        deliveryDate: "2024-11-29",
        totalAmount: 75.0,
    },
    {
        id: 4,
        customerName: "Emily Davis",
        status: "Canceled",
        orderDate: "2024-11-20",
        deliveryDate: "2024-11-25",
        totalAmount: 300.0,
    },
    {
        id: 5,
        customerName: "David Brown",
        status: "Pending",
        orderDate: "2024-11-18",
        deliveryDate: "2024-11-22",
        totalAmount: 120.0,
    },
    {
        id: 6,
        customerName: "Sophia Wilson",
        status: "Delivered",
        orderDate: "2024-11-15",
        deliveryDate: "2024-11-19",
        totalAmount: 220.0,
    },
    {
        id: 7,
        customerName: "Daniel Anderson",
        status: "Shipped",
        orderDate: "2024-11-12",
        deliveryDate: "2024-11-16",
        totalAmount: 85.0,
    },
    {
        id: 8,
        customerName: "Olivia Martinez",
        status: "Pending",
        orderDate: "2024-11-10",
        deliveryDate: "2024-11-14",
        totalAmount: 60.0,
    },
    {
        id: 9,
        customerName: "Matthew Garcia",
        status: "Delivered",
        orderDate: "2024-11-08",
        deliveryDate: "2024-11-12",
        totalAmount: 110.0,
    },
    {
        id: 10,
        customerName: "Isabella Taylor",
        status: "Shipped",
        orderDate: "2024-11-05",
        deliveryDate: "2024-11-09",
        totalAmount: 95.0,
    },
];


for (let i = 11; i <= 50; i++) {
    dummyOrders.push({
        id: i,
        customerName: `Customer ${i}`,
        status: ["Pending", "Shipped", "Delivered", "Canceled"][
            Math.floor(Math.random() * 4)
        ],
        orderDate: `2024-11-${Math.ceil(Math.random() * 30)
            .toString()
            .padStart(2, "0")}`,
        deliveryDate: `2024-12-${Math.ceil(Math.random() * 15)
            .toString()
            .padStart(2, "0")}`,
        totalAmount: parseFloat((Math.random() * 300).toFixed(2)),
    });
}

export default dummyOrders;
