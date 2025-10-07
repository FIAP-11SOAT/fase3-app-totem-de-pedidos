import "@typespec/http";
import "@typespec/openapi";

import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

namespace main.orders;

@route("/orders")
@tag("Orders")
interface Orders {

    @summary("Get all orders")
    @doc("Get all orders with optional filters")
    @get getOrders(
        @query id?: integer,
        @query status?: string,
        @query customer_id?: integer,
        @query notification_attempts?: integer
    ): Order[];

    @summary("Get an order by ID")
    @doc("Get an order by its unique identifier")
    @get getOrder(@path id: integer): Order;

    @summary("Create a new order")
    @doc("Create a new order with associated items")
    @post createOrder(@body order: OrderInput): OrderCreationResponse;

    @summary("Update an order by ID")
    @doc("Update an order's status by its unique identifier")
    @put updateOrder(@path id: integer, @body update: UpdateOrderInput): OrderUpdateResponse;
}
