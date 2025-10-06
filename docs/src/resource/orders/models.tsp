model Order {
    id: integer;
    order_date: string; 
    notification_attempts: integer;
    status: string;
    total_amount: float64;
    created_at: string; 
    customer_id: integer;
    items: OrderItem[];
}

model OrderItem {
    id: integer;
    quantity: integer;
    price: float64;
    created_at: string; 
    order_id: integer;
    product_id: integer;
}

model OrderInput {
    customer_id: integer;
    items: OrderItemInput[];
}

model OrderItemInput {
    product_id: integer;
    quantity: integer;
}

model OrderCreationResponse {
    message: string;
    orderId: string;
}

model UpdateOrderInput {
    status: string; 
}

model OrderUpdateResponse {
    message: string;
}
