model PaymentInput {
    order_id: string;
    amount: float;
    title: string;
}

model CreatePaymentResponse {
    payment_id: integer;
    order_id: integer;
    qrcode: string;
    qrcode_b64: string;
}

model PaymentResponse {
    id: integer;
    amount: float;
    payment_date: string;
    status: string;
    provider: string;
}

model PaymentNotification {}

