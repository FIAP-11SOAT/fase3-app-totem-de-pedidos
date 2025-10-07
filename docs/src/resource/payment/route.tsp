import "@typespec/http";
import "@typespec/openapi";

import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

namespace main.payments;

@route("/payments")
@tag("Payments")
interface Payments {

    @summary("Get By ID")
    @doc("Get a payment by its unique identifier")
    @get getPayment(@path id: integer): PaymentResponse;

    @summary("Create a new payment")
    @doc("Create a new payment with the specified details")
    @post createPayment(@body payment: PaymentInput): CreatePaymentResponse;

    @summary("Webhook for payment notifications")
    @doc("Receive notifications for payment status updates")
    @route("webhook") @post webhook(@body notification: PaymentNotification): void;
}