import "@typespec/http";
import "@typespec/openapi";

import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

namespace main.customers;

@route("/customer")
@tag("Customer")
interface Customers {

    @summary("Get a customer by tax ID")
    @doc("Get a customer by tax ID")
    @get getCustomer(@query taxid?: string): Customer;

    @summary("Create a new customer")
    @doc("Create a new customer")
    @post createCustomer(@body customer: CustomerInput): Customer;
}