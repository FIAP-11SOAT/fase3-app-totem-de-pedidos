package response

type DynamicQRResponse struct {
	InStoreOrderID string `json:"in_store_order_id"`
	QRData         string `json:"qr_data"`
}

type PaymentIDResponse struct {
	ExternalReference string `json:"external_reference"`
	Status            string `json:"status"`
}
