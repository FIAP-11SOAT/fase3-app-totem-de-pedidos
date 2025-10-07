import "@typespec/http";
import "@typespec/openapi";

import "./models.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

namespace main.products;

@route("/products")
@tag("Products")
interface Products {

    @summary("Get all products")
    @doc("Get all products")
    @get ListAllProducts(): Product[];

    @summary("Get a product by ID")
    @doc("Get a product by ID")
    @get FindProductByID(@path id: string): Product;

    @summary("Update a product")
    @doc("Update a product")
    @put UpdateProduct(@path id: string, @body product: ProductInput): Product;

    @summary("Delete a product")
    @doc("Delete a product")
    @delete DeleteProduct(@path id: string): void;

    @summary("Create a new product")
    @doc("Create a new product")
    @post CreateProduct(@body product: ProductInput): Product;
}